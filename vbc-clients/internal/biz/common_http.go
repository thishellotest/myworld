package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type CommonHttpUsecase struct {
	log                      *log.Helper
	conf                     *conf.Data
	JWTUsecase               *JWTUsecase
	ConditionUsecase         *ConditionUsecase
	QuestionnairesbuzUsecase *QuestionnairesbuzUsecase
	RelasLogUsecase          *RelasLogUsecase
	JotformSubmissionUsecase *JotformSubmissionUsecase
	TUsecase                 *TUsecase
}

func NewCommonHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	ConditionUsecase *ConditionUsecase,
	QuestionnairesbuzUsecase *QuestionnairesbuzUsecase,
	RelasLogUsecase *RelasLogUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	TUsecase *TUsecase) *CommonHttpUsecase {
	return &CommonHttpUsecase{
		log:                      log.NewHelper(logger),
		conf:                     conf,
		JWTUsecase:               JWTUsecase,
		ConditionUsecase:         ConditionUsecase,
		QuestionnairesbuzUsecase: QuestionnairesbuzUsecase,
		RelasLogUsecase:          RelasLogUsecase,
		JotformSubmissionUsecase: JotformSubmissionUsecase,
		TUsecase:                 TUsecase,
	}
}

const (
	CommonHttp_CommonType_ConditionQuestionnaires = "ConditionQuestionnaires"
)

func (c *CommonHttpUsecase) Save(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	// 通过路由获取的
	commonType := ctx.Param("common_type")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizSave(userFacade, commonType, body)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *CommonHttpUsecase) BizSave(userFacade UserFacade, commonType string, body lib.TypeMap) (lib.TypeMap, error) {
	data := make(lib.TypeMap)
	uniqid := body.GetString("uniqid")
	if commonType == CommonHttp_CommonType_ConditionQuestionnaires {

		if uniqid == "" {
			return nil, errors.New("uniqid is empty")
		}
		if configs.NewPSGen {

			jotformIdsList := body.GetTypeList(CommonFieldName_common_jotform_ids)
			var jotformIds []string
			sort := 10
			for _, v := range jotformIdsList {
				_, err := c.RelasLogUsecase.ConditionUpsert(InterfaceToString(uniqid),
					v.GetString("value"), &sort)
				sort++
				if err != nil {
					return nil, err
				}
				jotformIds = append(jotformIds, v.GetString("value"))
			}
			err := c.RelasLogUsecase.ConditionRemoveOtherTargetIds(InterfaceToString(uniqid), jotformIds)
			if err != nil {
				return nil, err
			}

		} else {
			conditionEntity, err := c.ConditionUsecase.ConditionUpsert(uniqid)
			if err != nil {
				return nil, err
			}
			if conditionEntity == nil {
				return nil, errors.New("conditionEntity is nil")
			}
			jotformIdsList := body.GetTypeList(CommonFieldName_common_jotform_ids)
			var jotformIds []string
			sort := 10
			for _, v := range jotformIdsList {
				_, err = c.RelasLogUsecase.ConditionUpsert(InterfaceToString(conditionEntity.ID),
					v.GetString("value"), &sort)
				sort++
				if err != nil {
					return nil, err
				}
				jotformIds = append(jotformIds, v.GetString("value"))
			}
			err = c.RelasLogUsecase.ConditionRemoveOtherTargetIds(InterfaceToString(conditionEntity.ID), jotformIds)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, errors.New("commonType is wrong")
	}
	return data, nil
}

func (c *CommonHttpUsecase) Get(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	// 通过路由获取的
	commonType := ctx.Param("common_type")
	uniqid := body.GetString("uniqid")
	caseGid := body.GetString("gid")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizGet(userFacade, commonType, uniqid, caseGid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *CommonHttpUsecase) BizGet(userFacade UserFacade, commonType, uniqid string, caseGid string) (lib.TypeMap, error) {
	data := make(lib.TypeMap)
	if commonType == CommonHttp_CommonType_ConditionQuestionnaires {
		if configs.NewPSGen {

			tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
			if tCase == nil {
				return nil, errors.New("tCase is nil")
			}

			items, err := c.JotformSubmissionUsecase.AllByUniqcodeAndConditionUniqid(uniqid)
			if err != nil {
				return nil, err
			}

			var fabFieldOptionList []FabFieldOptionNew
			for k, v := range items {

				newFileName, err := GenJotformNewFileNameForBox(items[k], tCase)
				if err != nil {
					return nil, err
				}

				fabFieldOptionList = append(fabFieldOptionList, FabFieldOptionNew{
					OptionValue: v.SubmissionId,
					OptionLabel: newFileName,
				})
			}
			data.Set("data."+CommonFieldName_common_jotform_ids, fabFieldOptionList)
			//sql := fmt.Sprintf("select q.* from questionnaires q inner join relas_log l on l.target_id=q.jotform_form_id and l.type=\"%s\" and l.source_id=\"%d\" where l.deleted_at=0 and q.deleted_at=0 order by q.base_title ",
			//	RelasLog_Type_condition_2_jotform_new, uniqid)

		} else {
			conditionEntity, err := c.ConditionUsecase.GetByCond(Eq{"condition_name": uniqid, "type": Condition_Type_Condition, "deleted_at": 0})
			if err != nil {
				return nil, err
			}
			if conditionEntity != nil {

				items, err := c.QuestionnairesbuzUsecase.AllForCondition(conditionEntity.ID)
				if err != nil {
					return nil, err
				}
				var fabFieldOptionList []FabFieldOptionNew
				for _, v := range items {
					fabFieldOptionList = append(fabFieldOptionList, FabFieldOptionNew{
						OptionValue: v.JotformFormId,
						OptionLabel: v.BaseTitle,
					})
				}
				data.Set("data."+CommonFieldName_common_jotform_ids, fabFieldOptionList)
			}
		}
	} else {
		return nil, errors.New("commonType is wrong")
	}
	return data, nil
}
