package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ConditionbuzUsecase struct {
	log                      *log.Helper
	CommonUsecase            *CommonUsecase
	conf                     *conf.Data
	TUsecase                 *TUsecase
	ConditionUsecase         *ConditionUsecase
	AiUsecase                *AiUsecase
	ConditionLogAiUsecase    *ConditionLogAiUsecase
	ConditionRelaAiUsecase   *ConditionRelaAiUsecase
	ConditionCategoryUsecase *ConditionCategoryUsecase
}

func NewConditionbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data, TUsecase *TUsecase,
	ConditionUsecase *ConditionUsecase,
	AiUsecase *AiUsecase,
	ConditionLogAiUsecase *ConditionLogAiUsecase,
	ConditionRelaAiUsecase *ConditionRelaAiUsecase,
	ConditionCategoryUsecase *ConditionCategoryUsecase) *ConditionbuzUsecase {
	uc := &ConditionbuzUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		TUsecase:                 TUsecase,
		ConditionUsecase:         ConditionUsecase,
		AiUsecase:                AiUsecase,
		ConditionLogAiUsecase:    ConditionLogAiUsecase,
		ConditionRelaAiUsecase:   ConditionRelaAiUsecase,
		ConditionCategoryUsecase: ConditionCategoryUsecase,
	}

	return uc
}

func (c *ConditionbuzUsecase) HandleAllCondition() error {
	c.log.Info("HandleAllCondition begin")
	records, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{"deleted_at": 0, "biz_deleted_at": 0})
	if err != nil {
		return err
	}
	for _, v := range records {
		serviceConnections := v.CustomFields.TextValueByNameBasic("service_connections")
		previousDenials := v.CustomFields.TextValueByNameBasic("previous_denials")
		claimsOnline := v.CustomFields.TextValueByNameBasic("claims_online")
		claimsNextRound := v.CustomFields.TextValueByNameBasic("claims_next_round")
		claimsSupplemental := v.CustomFields.TextValueByNameBasic("claims_supplemental")

		var strs []string
		strs = append(strs, serviceConnections)
		strs = append(strs, previousDenials)
		strs = append(strs, claimsOnline)
		strs = append(strs, claimsNextRound)
		strs = append(strs, claimsSupplemental)

		for _, v1 := range strs {
			res := CaseClaimsDivideV2(v1)
			for _, v2 := range res {
				err = c.ConditionUsecase.Upsert(v2.Condition)
				if err != nil {
					c.log.Error()
				}
			}
		}
	}
	c.log.Info("HandleAllCondition end")
	return nil
}

func (c *ConditionbuzUsecase) HandleOnceConditionSourceWithAi() error {

	promptKey := "multiline_condition_parser"
	res, err := c.ConditionUsecase.ListByCondWithPaging(Eq{"type": Condition_Type_Source,
		"source_status": Condition_SourceStatus_NotHandled}, "", 1, 8)
	c.log.Info("HandleOnceConditionSourceWithAi res length: ", len(res))
	if err != nil {
		return err
	}
	var conditions []string
	for _, v := range res {
		conditions = append(conditions, v.ConditionName)

	}

	if len(conditions) > 0 {
		conditionStr := strings.Join(conditions, "\n")
		dynamicParamsExample := lib.TypeMap{
			"text": conditionStr,
		}
		parseResult, _, err := c.AiUsecase.ExecuteWithClaude3(promptKey, dynamicParamsExample)
		if err != nil {
			c.log.Warn(err)
		}
		err = c.HandleParseConditionResultFromAi(parseResult, promptKey)
		if err != nil {
			return err
		}

		for k, _ := range res {
			res[k].SourceStatus = Condition_SourceStatus_Handled
			res[k].SourceHandleCount += 1
			res[k].UpdatedAt = time.Now().Unix()
			err = c.CommonUsecase.DB().Save(&res[k]).Error
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}

// HandleParseConditionResultFromAi [{"PrimaryCondition":"Left ankle pain","DirectSecondaryConditions":[],"AggravationConditions":[],"SourceData":"Left ankle pain (str, opinion)"},{"PrimaryCondition":"Depression","DirectSecondaryConditions":["tinnitus","hearing loss"],"AggravationConditions":[],"SourceData":"70 - Depression secondary to tinnitus and hearing loss"},{"PrimaryCondition":"Vertigo","DirectSecondaryConditions":["tinnitus","hearing loss"],"AggravationConditions":[],"SourceData":"60 - Vertigo secondary to tinnitus and hearing loss (opinion)"},{"PrimaryCondition":"Headaches","DirectSecondaryConditions":["tinnitus"],"AggravationConditions":[],"SourceData":"50 - Headaches secondary to tinnitus (opinion)"},{"PrimaryCondition":"Sleep apnea","DirectSecondaryConditions":["tinnitus","hearing loss"],"AggravationConditions":[],"SourceData":"50 - Sleep apnea secondary to tinnitus and hearing loss (opinion)"},{"PrimaryCondition":"Vertigo","DirectSecondaryConditions":[],"AggravationConditions":[],"SourceData":"30 - Vertigo"},{"PrimaryCondition":"Right knee patellar tendonitis","DirectSecondaryConditions":[],"AggravationConditions":["limitation of flexion","limitation of extension"],"SourceData":"20x - Right knee patellar tendonitis with limitation of flexion and extension (str, opinion)"}]
func (c *ConditionbuzUsecase) HandleParseConditionResultFromAi(str string, promptKey string) error {
	res := lib.ToTypeListByString(str)
	//c.log.Info("HandleOnceConditionSourceWithAi HandleParseConditionResultFromAi res:", str)
	if len(res) == 0 {
		return errors.New("Ai Response error")
	}

	for _, v := range res {
		PrimaryCondition := v.GetString("PrimaryCondition")
		AggravationConditions := v.GetTypeListInterface("AggravationConditions")
		DirectSecondaryConditions := v.GetTypeListInterface("DirectSecondaryConditions")
		SourceData := v.GetString("SourceData")

		sourceConditionEntity, _, err := c.ConditionUsecase.UpsertSourceCondition(SourceData, Condition_Type_Source_From_Ai)
		if err != nil {
			return err
		}
		err = c.ConditionLogAiUsecase.AddLogPromptAiCondition(sourceConditionEntity.ID, 0, promptKey)
		if err != nil {
			c.log.Error(err)
		}

		primaryConditionEntity, _, err := c.ConditionUsecase.UpsertPrimaryCondition(PrimaryCondition, 0)
		if err != nil {
			return err
		}
		err = c.ConditionLogAiUsecase.AddLogPromptAiCondition(primaryConditionEntity.ID, sourceConditionEntity.ID, promptKey)
		if err != nil {
			c.log.Error(err)
		}
		for _, v1 := range AggravationConditions {
			secondaryConditionEntity, _, err := c.ConditionUsecase.UpsertSecondaryCondition(InterfaceToString(v1), Condition_SecondaryType_Aggravation)
			if err != nil {
				return err
			}
			err = c.ConditionRelaAiUsecase.Upsert(primaryConditionEntity.ID, secondaryConditionEntity.ID, promptKey)
			if err != nil {
				return err
			}
			err = c.ConditionLogAiUsecase.AddLogPromptAiCondition(secondaryConditionEntity.ID, sourceConditionEntity.ID, promptKey)
			if err != nil {
				return err
			}

		}
		for _, v1 := range DirectSecondaryConditions {
			secondaryConditionEntity, _, err := c.ConditionUsecase.UpsertSecondaryCondition(InterfaceToString(v1), Condition_SecondaryType_DirectSecondary)
			if err != nil {
				return err
			}
			err = c.ConditionRelaAiUsecase.Upsert(primaryConditionEntity.ID, secondaryConditionEntity.ID, promptKey)
			if err != nil {
				return err
			}
			err = c.ConditionLogAiUsecase.AddLogPromptAiCondition(secondaryConditionEntity.ID, sourceConditionEntity.ID, promptKey)
			if err != nil {
				return err
			}

		}
		break
		//lib.DPrintln("PrimaryCondition:", PrimaryCondition)
		//lib.DPrintln("AggravationConditions:", AggravationConditions)
		//lib.DPrintln("DirectSecondaryConditions:", DirectSecondaryConditions)
		//lib.DPrintln("SourceData:", SourceData)
		//break
	}

	return nil
}

func (c *ConditionbuzUsecase) HandleTCondition(tCondition *TData) {
	if tCondition != nil {
		for k, v := range tCondition.CustomFields {
			if v.Type == FieldType_lookup && v.Name == ConditionFieldName_condition_category_id {
				if v.TextValue == nil || *v.TextValue == "0" {
					tCondition.CustomFields[k].DisplayValue = &UndefinedConditionCategory.ConditionCategoryName
				} else {
					categoryEntity, _ := c.ConditionCategoryUsecase.GetByCond(Eq{"id": *v.TextValue})
					if categoryEntity != nil {
						tCondition.CustomFields[k].DisplayValue = &categoryEntity.CategoryName
					}

				}
			}
		}
	}
}

func (c *ConditionbuzUsecase) GetAllSecondariesByConditionIds(conditionIds []int32) ([]*ConditionEntity, error) {
	builder := Dialect(MYSQL).Select("c.id").From("conditions", "c")
	builder.InnerJoin("condition_relas_ai r", "r.secondary_condition_id=c.id and r.deleted_at=0")
	builder.Where(And(In("r.primary_condition_id", conditionIds), Eq{"c.deleted_at": 0, "c.biz_deleted_at": 0}))
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	newSql := fmt.Sprintf("select * from conditions  where id in (%s)", sql)
	allSecondaries, err := c.ConditionUsecase.AllByRawSql(newSql)
	return allSecondaries, err
}

func (c *ConditionbuzUsecase) GetAllConditionRelasByConditionIds(conditionIds []int32) ([]*ConditionRelaAiEntity, error) {

	// 获取关系
	builder1 := Dialect(MYSQL).Select("r.*").From("conditions", "c")
	builder1.InnerJoin("condition_relas_ai r", "r.secondary_condition_id=c.id and r.deleted_at=0")
	builder1.Where(And(In("r.primary_condition_id", conditionIds), Eq{"c.deleted_at": 0}))
	sql1, err := builder1.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	allConditionRelas, err := c.ConditionRelaAiUsecase.AllByRawSql(sql1)
	if err != nil {
		return nil, err
	}
	return allConditionRelas, nil
}

func (c *ConditionbuzUsecase) GetAllCategoriesByConditionIds(conditionIds []int32) ([]*ConditionCategoryEntity, error) {
	builder2 := Dialect(MYSQL).Select("r.id").From("conditions", "c")
	builder2.InnerJoin("condition_categories r", "r.id=c.condition_category_id and r.biz_deleted_at=0")
	builder2.Where(And(In("c.id", conditionIds), Eq{"c.deleted_at": 0}))
	sql2, err := builder2.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	newSql2 := fmt.Sprintf("select * from condition_categories  where id in (%s)", sql2)
	AllCategories, err := c.ConditionCategoryUsecase.AllByRawSql(newSql2)
	return AllCategories, err
}
