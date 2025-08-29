package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ContractHttpUsecase struct {
	log                    *log.Helper
	conf                   *conf.Data
	JWTUsecase             *JWTUsecase
	ContractbuzUsecase     *ContractbuzUsecase
	RevisionHistoryUsecase *RevisionHistoryUsecase
	TUsecase               *TUsecase
	MgmtPermissionUsecase  *MgmtPermissionUsecase
}

func NewContractHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	ContractbuzUsecase *ContractbuzUsecase,
	RevisionHistoryUsecase *RevisionHistoryUsecase,
	TUsecase *TUsecase,
	MgmtPermissionUsecase *MgmtPermissionUsecase) *ContractHttpUsecase {
	return &ContractHttpUsecase{
		log:                    log.NewHelper(logger),
		conf:                   conf,
		JWTUsecase:             JWTUsecase,
		ContractbuzUsecase:     ContractbuzUsecase,
		RevisionHistoryUsecase: RevisionHistoryUsecase,
		TUsecase:               TUsecase,
		MgmtPermissionUsecase:  MgmtPermissionUsecase,
	}
}

// Get 获取合同信息并确定是否可以修改
func (c *ContractHttpUsecase) Get(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))

	// 通过路由获取的
	//moduleName := ctx.Param("module_name")
	// tUser, _ := c.JWTUsecase.JWTUser(ctx)
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.ContractbuzUsecase.BizGet(body.GetInt("case_id"), userFacade)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type ContractHttpSaveVo struct {
	CaseId     int32 `json:"case_id"`
	ActiveDuty bool  `json:"active_duty"`
	Rating     int32 `json:"rating"`
}

func VerifyCurrentRating(rating int32) bool {
	if rating == 0 || rating == 10 || rating == 20 || rating == 30 ||
		rating == 40 || rating == 50 || rating == 60 || rating == 70 || rating == 80 || rating == 90 {
		return true
	}
	return false
}

// Save 获取合同信息并确定是否可以修改
func (c *ContractHttpUsecase) Save(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	var contractHttpSaveVo ContractHttpSaveVo
	contractHttpSaveVo = lib.BytesToTDef(rawData, contractHttpSaveVo)
	// 通过路由获取的
	//moduleName := ctx.Param("module_name")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.ContractbuzUsecase.BizSave(contractHttpSaveVo, userFacade)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ContractHttpUsecase) List(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	page := HandlePage(ctx.Query("page"))
	pageSize := HandlePageSize(ctx.Query("page_size"))
	data, err := c.BizList(userFacade, page, pageSize)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ContractHttpUsecase) BizList(userFacade UserFacade, page int, pageSize int) (lib.TypeMap, error) {

	hasPermission, err := c.MgmtPermissionUsecase.Verify(userFacade, MgmtPermission_ReviseContract)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New(Error_UnauthorizedOperation)
	}

	cond := Eq{"biz_type": RevisionHistory_BizType_contract}
	records, err := c.RevisionHistoryUsecase.ListByCondWithPaging(cond, "id desc", page, pageSize)
	if err != nil {
		return nil, err
	}

	changedBys := make(lib.TypeMap)
	caseGids := make(lib.TypeMap)
	for _, v := range records {
		changedBys.Set(v.ChangedBy, 1)
		caseGids.Set(v.Uniqid, 1)
	}
	var changedByList []string
	var caseGidList []string
	for k, _ := range changedBys {
		changedByList = append(changedByList, k)
	}

	for k, _ := range caseGids {
		caseGidList = append(caseGidList, k)
	}

	cases := make(map[string]*TData)
	users := make(map[string]*TData)
	if len(changedByList) > 0 {
		res1, err := c.TUsecase.ListByCond(Kind_client_cases, And(In(FieldName_gid, caseGidList), Eq{FieldName_biz_deleted_at: 0}))
		if err != nil {
			return nil, err
		}
		for k, v := range res1 {
			cases[v.Gid()] = res1[k]
		}

		res2, err := c.TUsecase.ListByCond(Kind_users, And(In(FieldName_gid, changedByList), Eq{FieldName_biz_deleted_at: 0}))
		if err != nil {
			return nil, err
		}
		for k, v := range res2 {
			users[v.Gid()] = res2[k]
		}
	}

	var destRecords []RevisionHistoryToContractApi
	for _, v := range records {
		destRecords = append(destRecords, v.ToContractApi(cases, users))
	}

	total, err := c.RevisionHistoryUsecase.Total(cond)
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
