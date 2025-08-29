package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ConditionHttpUsecase struct {
	log                      *log.Helper
	conf                     *conf.Data
	JWTUsecase               *JWTUsecase
	ConditionUsecase         *ConditionUsecase
	ConditionRelaAiUsecase   *ConditionRelaAiUsecase
	ConditionLogAiUsecase    *ConditionLogAiUsecase
	ConditionCategoryUsecase *ConditionCategoryUsecase
	ConditionbuzUsecase      *ConditionbuzUsecase
}

func NewConditionHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	ConditionUsecase *ConditionUsecase,
	ConditionRelaAiUsecase *ConditionRelaAiUsecase,
	ConditionLogAiUsecase *ConditionLogAiUsecase,
	ConditionCategoryUsecase *ConditionCategoryUsecase,
	ConditionbuzUsecase *ConditionbuzUsecase) *ConditionHttpUsecase {
	return &ConditionHttpUsecase{
		log:                      log.NewHelper(logger),
		conf:                     conf,
		JWTUsecase:               JWTUsecase,
		ConditionUsecase:         ConditionUsecase,
		ConditionRelaAiUsecase:   ConditionRelaAiUsecase,
		ConditionLogAiUsecase:    ConditionLogAiUsecase,
		ConditionCategoryUsecase: ConditionCategoryUsecase,
		ConditionbuzUsecase:      ConditionbuzUsecase,
	}
}

func (c *ConditionHttpUsecase) Sources(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	//
	//// 通过路由获取的
	//moduleName := ctx.Param("module_name")
	//lib.DPrintln(moduleName)

	data, err := c.BizSources(body.GetInt("id"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ConditionHttpUsecase) BizSources(id int32) (lib.TypeMap, error) {
	if id <= 0 {
		return nil, errors.New("Parameter Error")
	}
	records, err := c.ConditionLogAiUsecase.AllByCond(Eq{"condition_id": id})
	if err != nil {
		return nil, err
	}

	var destRecords []*ConditionLogAiVo
	for _, v := range records {
		row := v.ToApi(c.log, c.ConditionUsecase)
		if row != nil {
			destRecords = append(destRecords, row)
		}
	}
	data := make(lib.TypeMap)
	data.Set("records", destRecords)
	return data, nil
}

func (c *ConditionHttpUsecase) Delete(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizDelete(userFacade, body, ctx)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ConditionHttpUsecase) BizDelete(userFacade UserFacade, params lib.TypeMap, ctx *gin.Context) (lib.TypeMap, error) {

	conditionId := params.GetInt("condition_id")
	if conditionId <= 0 {
		return nil, errors.New("Parameter Error")
	}

	conditionEntity, _ := c.ConditionUsecase.GetByCond(Eq{"id": conditionId, "biz_deleted_at": 0})
	if conditionEntity == nil {
		c.log.Warn("conditionEntity is nil")
		return nil, errors.New("Parameter Error")
	}

	if conditionEntity.Type == Condition_Type_Primary {

		secondaries, _ := c.ConditionbuzUsecase.GetAllSecondariesByConditionIds([]int32{conditionEntity.ID})
		if len(secondaries) > 0 {
			return nil, errors.New("Delete all secondary conditions before performing this operation.")
		}
	} else if conditionEntity.Type == Condition_Type_SecondaryCondition {

	} else {
		c.log.Warn("conditionEntity.Type is wrong")
		return nil, errors.New("Parameter Error")
	}
	conditionEntity.BizDeletedAt = time.Now().Unix()
	err := c.ConditionUsecase.CommonUsecase.DB().Save(&conditionEntity).Error
	if err != nil {
		return nil, err
	}
	_, er := c.ConditionLogAiUsecase.UpsertCenter(ConditionLogAi_LogType_ConditionSource, conditionEntity.ID, "", ConditionLogAi_FromType_ManualDelete, "", userFacade.Gid(), "")
	if er != nil {
		c.log.Warn(er)
	}
	data := make(lib.TypeMap)

	if conditionEntity.Type == Condition_Type_SecondaryCondition {

		belongConditionId, _ := ctx.GetQuery("belong_condition_id")
		belongConditionIdInt, _ := strconv.ParseInt(belongConditionId, 10, 32)
		primaryConditionEntity, err := c.ConditionUsecase.GetByCond(Eq{"id": belongConditionIdInt})
		if err != nil {
			return nil, err
		}

		conditionIds := []int32{primaryConditionEntity.ID}
		AllSecondaries, _ := c.ConditionbuzUsecase.GetAllSecondariesByConditionIds(conditionIds)
		AllCategories, _ := c.ConditionbuzUsecase.GetAllCategoriesByConditionIds(conditionIds)
		conditionCategoryEntity := GetConditionCategoryById(primaryConditionEntity.ConditionCategoryId, AllCategories)
		primaryConditionVo := primaryConditionEntity.ToPrimaryCondition(AllSecondaries, conditionCategoryEntity)
		data.Set("primary_condition", primaryConditionVo)
	}

	return data, nil
}

func (c *ConditionHttpUsecase) List(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	//
	//// 通过路由获取的
	//moduleName := ctx.Param("module_name")
	//lib.DPrintln(moduleName)

	page := HandlePage(ctx.Query("page"))
	pageSize := HandlePageSize(ctx.Query("page_size"))

	data, err := c.BizList(page, pageSize, body)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func GetConditionCategoryById(conditionCategoryId int32, AllCategories []*ConditionCategoryEntity) *ConditionCategoryEntity {
	for k, v := range AllCategories {
		if v.ID == conditionCategoryId {
			return AllCategories[k]
		}
	}
	return nil
}

func GetSecondariesByPrimaryConditionId(primaryConditionId int32, allSecondaries []*ConditionEntity, allConditionRelas []*ConditionRelaAiEntity) (secondaries []*ConditionEntity) {
	var secondariesConditionIds []int32
	for _, v := range allConditionRelas {
		if v.PrimaryConditionId == primaryConditionId {
			secondariesConditionIds = append(secondariesConditionIds, v.SecondaryConditionId)
		}
	}
	secondariesConditionIds = lib.RemoveDuplicates(secondariesConditionIds)
	for _, v := range secondariesConditionIds {
		for k1, v1 := range allSecondaries {
			if v1.ID == v {
				secondaries = append(secondaries, allSecondaries[k1])
			}
		}
	}
	return secondaries
}

func (c *ConditionHttpUsecase) BizList(page, pageSize int, params lib.TypeMap) (lib.TypeMap, error) {

	categories := params.GetTypeList("category.value")
	var categoriesIds []int32
	for _, v := range categories {
		categoriesIds = append(categoriesIds, v.GetInt("value"))
	}
	condition := params.GetString("condition.value")
	var conds []Cond
	cond := Eq{"type": Condition_Type_Primary,
		"deleted_at": 0, "biz_deleted_at": 0}
	conds = append(conds, cond)
	if len(categoriesIds) > 0 {
		conds = append(conds, In("condition_category_id", categoriesIds))
	}
	if condition != "" {
		conds = append(conds, Like{"condition_name", condition})
	}

	records, err := c.ConditionUsecase.ListByCondWithPaging(And(conds...), "id desc", page, pageSize)
	if err != nil {
		return nil, err
	}
	var conditionIds []int32
	for _, v := range records {
		conditionIds = append(conditionIds, v.ID)
	}
	if len(conditionIds) == 0 {
		return nil, nil
	}
	//builder := Dialect(MYSQL).Select("c.id").From("conditions", "c")
	//builder.InnerJoin("condition_relas_ai r", "r.secondary_condition_id=c.id and r.deleted_at=0")
	//builder.Where(And(In("r.primary_condition_id", conditionIds), Eq{"c.deleted_at": 0}))
	//sql, err := builder.ToBoundSQL()
	//if err != nil {
	//	return nil, err
	//}
	//newSql := fmt.Sprintf("select * from conditions  where id in (%s)", sql)
	//allSecondaries, err := c.ConditionUsecase.AllByRawSql(newSql)
	allSecondaries, err := c.ConditionbuzUsecase.GetAllSecondariesByConditionIds(conditionIds)
	if err != nil {
		return nil, err
	}

	// 获取关系
	//builder1 := Dialect(MYSQL).Select("r.*").From("conditions", "c")
	//builder1.InnerJoin("condition_relas_ai r", "r.secondary_condition_id=c.id and r.deleted_at=0")
	//builder1.Where(And(In("r.primary_condition_id", conditionIds), Eq{"c.deleted_at": 0}))
	//sql1, err := builder1.ToBoundSQL()
	//if err != nil {
	//	return nil, err
	//}
	//allConditionRelas, err := c.ConditionRelaAiUsecase.AllByRawSql(sql1)
	allConditionRelas, err := c.ConditionbuzUsecase.GetAllConditionRelasByConditionIds(conditionIds)
	if err != nil {
		return nil, err
	}

	// 获取所有分类
	//builder2 := Dialect(MYSQL).Select("r.id").From("conditions", "c")
	//builder2.InnerJoin("condition_categories r", "r.id=c.condition_category_id and r.biz_deleted_at=0")
	//builder2.Where(And(In("c.id", conditionIds), Eq{"c.deleted_at": 0}))
	//sql2, err := builder2.ToBoundSQL()
	//if err != nil {
	//	return nil, err
	//}
	//newSql2 := fmt.Sprintf("select * from condition_categories  where id in (%s)", sql2)
	//AllCategories, err := c.ConditionCategoryUsecase.AllByRawSql(newSql2)
	AllCategories, err := c.ConditionbuzUsecase.GetAllCategoriesByConditionIds(conditionIds)

	var destRecords []PrimaryConditionVo
	for _, v := range records {
		conditionCategoryEntity := GetConditionCategoryById(v.ConditionCategoryId, AllCategories)
		secondaries := GetSecondariesByPrimaryConditionId(v.ID, allSecondaries, allConditionRelas)
		row := v.ToPrimaryCondition(secondaries, conditionCategoryEntity)
		destRecords = append(destRecords, row)
	}

	total, err := c.ConditionUsecase.Total(cond)
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)
	data.Set(Fab_TRecords, destRecords)
	data.Set(Fab_TTotal, int32(total))
	data.Set(Fab_TPage, page)
	data.Set(Fab_TPageSize, pageSize)

	return data, nil
}
