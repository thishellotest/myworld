package biz

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

type RecordHttpUsecase struct {
	log                         *log.Helper
	CommonUsecase               *CommonUsecase
	conf                        *conf.Data
	TUsecase                    *TUsecase
	JWTUsecase                  *JWTUsecase
	RecordbuzUsecase            *RecordbuzUsecase
	SettingSectionFieldUsecase  *SettingSectionFieldUsecase
	FieldUsecase                *FieldUsecase
	FieldOptionUsecase          *FieldOptionUsecase
	DataEntryUsecase            *DataEntryUsecase
	ZohoUsecase                 *ZohoUsecase
	TimelinesbuzUsecase         *TimelinesbuzUsecase
	BUsecase                    *BUsecase
	UserUsecase                 *UserUsecase
	FieldPermissionUsecase      *FieldPermissionUsecase
	FieldValidatorUsecase       *FieldValidatorUsecase
	TimezonesUsecase            *TimezonesUsecase
	PermissionDataFilterUsecase *PermissionDataFilterUsecase
	KindUsecase                 *KindUsecase
	RecordbuzSearchUsecase      *RecordbuzSearchUsecase
	QueueUsecase                *QueueUsecase
	ZohobuzUsecase              *ZohobuzUsecase
	FieldbuzUsecase             *FieldbuzUsecase
	ConditionbuzUsecase         *ConditionbuzUsecase
	ConditionLogAiUsecase       *ConditionLogAiUsecase
	ConditionRelaAiUsecase      *ConditionRelaAiUsecase
	ConditionUsecase            *ConditionUsecase
	FilterUsecase               *FilterUsecase
	RecordLogUsecase            *RecordLogUsecase
	MedicalDbqCostUsecase       *MedicalDbqCostUsecase
	LogUsecase                  *LogUsecase
	ClientCasebuzUsecase        *ClientCasebuzUsecase
	StatementUsecase            *StatementUsecase
	LeadsUsecase                *LeadsUsecase
	ClientReviewBuzUsecase      *ClientReviewBuzUsecase
	StatementCommentBuzUsecase  *StatementCommentBuzUsecase
	ClientUsecase               *ClientUsecase
}

func NewRecordHttpUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	JWTUsecase *JWTUsecase,
	RecordbuzUsecase *RecordbuzUsecase,
	SettingSectionFieldUsecase *SettingSectionFieldUsecase,
	FieldUsecase *FieldUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
	DataEntryUsecase *DataEntryUsecase,
	ZohoUsecase *ZohoUsecase,
	TimelinesbuzUsecase *TimelinesbuzUsecase,
	BUsecase *BUsecase,
	UserUsecase *UserUsecase,
	FieldPermissionUsecase *FieldPermissionUsecase,
	FieldValidatorUsecase *FieldValidatorUsecase,
	TimezonesUsecase *TimezonesUsecase,
	PermissionDataFilterUsecase *PermissionDataFilterUsecase,
	KindUsecase *KindUsecase,
	RecordbuzSearchUsecase *RecordbuzSearchUsecase,
	QueueUsecase *QueueUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
	FieldbuzUsecase *FieldbuzUsecase,
	ConditionbuzUsecase *ConditionbuzUsecase,
	ConditionLogAiUsecase *ConditionLogAiUsecase,
	ConditionRelaAiUsecase *ConditionRelaAiUsecase,
	ConditionUsecase *ConditionUsecase,
	FilterUsecase *FilterUsecase,
	RecordLogUsecase *RecordLogUsecase,
	MedicalDbqCostUsecase *MedicalDbqCostUsecase,
	LogUsecase *LogUsecase,
	ClientCasebuzUsecase *ClientCasebuzUsecase,
	StatementUsecase *StatementUsecase,
	LeadsUsecase *LeadsUsecase,
	ClientReviewBuzUsecase *ClientReviewBuzUsecase,
	StatementCommentBuzUsecase *StatementCommentBuzUsecase,
	ClientUsecase *ClientUsecase) *RecordHttpUsecase {
	uc := &RecordHttpUsecase{
		log:                         log.NewHelper(logger),
		CommonUsecase:               CommonUsecase,
		conf:                        conf,
		TUsecase:                    TUsecase,
		JWTUsecase:                  JWTUsecase,
		RecordbuzUsecase:            RecordbuzUsecase,
		SettingSectionFieldUsecase:  SettingSectionFieldUsecase,
		FieldUsecase:                FieldUsecase,
		FieldOptionUsecase:          FieldOptionUsecase,
		DataEntryUsecase:            DataEntryUsecase,
		ZohoUsecase:                 ZohoUsecase,
		TimelinesbuzUsecase:         TimelinesbuzUsecase,
		BUsecase:                    BUsecase,
		UserUsecase:                 UserUsecase,
		FieldPermissionUsecase:      FieldPermissionUsecase,
		FieldValidatorUsecase:       FieldValidatorUsecase,
		TimezonesUsecase:            TimezonesUsecase,
		PermissionDataFilterUsecase: PermissionDataFilterUsecase,
		KindUsecase:                 KindUsecase,
		RecordbuzSearchUsecase:      RecordbuzSearchUsecase,
		QueueUsecase:                QueueUsecase,
		ZohobuzUsecase:              ZohobuzUsecase,
		FieldbuzUsecase:             FieldbuzUsecase,
		ConditionbuzUsecase:         ConditionbuzUsecase,
		ConditionLogAiUsecase:       ConditionLogAiUsecase,
		ConditionRelaAiUsecase:      ConditionRelaAiUsecase,
		ConditionUsecase:            ConditionUsecase,
		FilterUsecase:               FilterUsecase,
		RecordLogUsecase:            RecordLogUsecase,
		MedicalDbqCostUsecase:       MedicalDbqCostUsecase,
		LogUsecase:                  LogUsecase,
		ClientCasebuzUsecase:        ClientCasebuzUsecase,
		StatementUsecase:            StatementUsecase,
		LeadsUsecase:                LeadsUsecase,
		ClientReviewBuzUsecase:      ClientReviewBuzUsecase,
		StatementCommentBuzUsecase:  StatementCommentBuzUsecase,
		ClientUsecase:               ClientUsecase,
	}
	return uc
}

// HandlePage 页码是从1开始的
func HandlePage(page string) int {
	v, _ := strconv.ParseInt(page, 10, 32)
	if v <= 0 {
		return 1
	}
	return int(v)
}

func HandlePageSize(pageSize string) int {
	v, _ := strconv.ParseInt(pageSize, 10, 32)
	if v <= 0 {
		return TDefaultPageSize
	}
	if v > TMaxPageSize {
		return TMaxPageSize
	}
	return int(v)
}

func HandleOffset(page int, pageSize int) int {
	return (page - 1) * pageSize
}

func (c *RecordHttpUsecase) StatementRevertVersion(ctx *gin.Context) {
	reply := CreateReply()
	caseGid := ctx.Param("caseGid")
	versionId := ctx.Query("versionId")
	commentIdInt, _ := strconv.ParseInt(versionId, 0, 32)
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.StatementUsecase.BizStatementRevertVersion(userFacade, caseGid, int32(commentIdInt))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) StatementDetailVersions(ctx *gin.Context) {

	reply := CreateReply()
	tokenOk, tUser, _ := c.JWTUsecase.JWTAuth(ctx.GetHeader("Authorization"), ctx)
	caseGid := ctx.Param("caseGid")
	usePasswordAccess := false
	if tokenOk { // 验证数据权限
		tCase, er := c.RecordbuzUsecase.VerifyDataPermission(caseGid, *tUser)
		if er != nil {
			c.log.Error(er)
		}
		if tCase == nil {
			tokenOk = false
		}
	}
	if !tokenOk { // 非登录base用户，需要使用密码访问
		password := ctx.Query("password")
		tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
		caseId := int32(0)
		if tCase != nil {
			caseId = tCase.Id()
		}

		isOk, _ := c.StatementUsecase.BizStatementVerifyPassword(caseId, password)
		if !isOk { // todo:lgl 验证密码
			reply["code"] = Reply_code_waiting_password
			reply["message"] = "Waiting Password"
			ctx.JSON(200, reply)
			return
		}
		usePasswordAccess = true
	}
	rawData, _ := ctx.GetRawData()
	data, err := c.StatementUsecase.BizStatementDetailVersions(usePasswordAccess, tUser, caseGid, rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) StatementDetailTest(ctx *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			c.log.Error("panic recovered: %v", r)
			ctx.JSON(500, gin.H{
				"code":    500,
				"message": "Internal Server Error",
			})
		}
	}()

	reply := CreateReply()

	data := make(lib.TypeMap)

	json.Unmarshal([]byte(TestStatementData2), &data)
	//lib.DPrintln(usePasswordAccess, tUser, caseGid)
	var err error
	//var data
	//data, err := c.StatementUsecase.BizStatementDetail(usePasswordAccess, tUser, caseGid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordbuzUsecase) VerifyDataPermission(caseGid string, tUser TData) (tCase *TData, err error) {

	kindEntity, err := c.KindUsecase.GetByKind(Kind_client_cases)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}

	userFacade := UserFacade{
		TData: tUser,
	}
	return c.GetRecordData(caseGid, *kindEntity, &userFacade, false)
}

func (c *RecordHttpUsecase) StatementDetail(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.log.Error("panic recovered: %v", r)
			ctx.JSON(500, gin.H{
				"code":    500,
				"message": "Internal Server Error",
			})
		}
	}()

	reply := CreateReply()
	tokenOk, tUser, _ := c.JWTUsecase.JWTAuth(ctx.GetHeader("Authorization"), ctx)
	caseGid := ctx.Param("caseGid")
	usePasswordAccess := false

	if tokenOk { // 验证数据权限
		tCase, er := c.RecordbuzUsecase.VerifyDataPermission(caseGid, *tUser)
		if er != nil {
			c.log.Error(er)
		}
		if tCase == nil {
			tokenOk = false
		}
	}

	if !tokenOk { // 非登录base用户，需要使用密码访问
		password := ctx.Query("password")
		tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
		caseId := int32(0)
		if tCase != nil {
			caseId = tCase.Id()
		}

		isOk, _ := c.StatementUsecase.BizStatementVerifyPassword(caseId, password)
		if !isOk { // todo:lgl 验证密码
			reply["code"] = Reply_code_waiting_password
			reply["message"] = "Waiting Password"
			ctx.JSON(200, reply)
			return
		}
		usePasswordAccess = true
	}
	//rawData, _ := ctx.GetRawData()

	//data := make(lib.TypeMap)

	//json.Unmarshal([]byte(TestStatementData2), &data)
	//lib.DPrintln(usePasswordAccess, tUser, caseGid)
	//var err error
	//var data
	data, err := c.StatementUsecase.BizStatementDetail(usePasswordAccess, tUser, caseGid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) StatementVerifyPassword(ctx *gin.Context) {

	password := ctx.Query("password")
	caseGid := ctx.Param("caseGid")
	reply := CreateReply()
	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	caseId := int32(0)
	if tCase != nil {
		caseId = tCase.Id()
	}
	isOk, err := c.StatementUsecase.BizStatementVerifyPassword(caseId, password)
	if err != nil {
		reply.CommonError(err)
	} else {
		if isOk {
			reply.Success()
		} else {
			reply.Update(Reply_code_password_error, "The password is incorrect")
		}
	}
	ctx.JSON(200, reply)

}

func (c *RecordHttpUsecase) StatementCommentSubmitForReview(ctx *gin.Context) {
	reply := CreateReply()
	caseGid := ctx.Param("caseGid")

	password := ctx.Query("password")
	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	caseId := int32(0)
	if tCase != nil {
		caseId = tCase.Id()
	}
	isOk, _ := c.StatementUsecase.BizStatementVerifyPassword(caseId, password)
	if !isOk { // todo:lgl 验证密码
		reply["code"] = Reply_code_waiting_password
		reply["message"] = "Waiting Password"
		ctx.JSON(200, reply)
		return
	}

	data, err := c.StatementCommentBuzUsecase.BizStatementCommentSubmitForReview(caseGid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) StatementCommentDelete(ctx *gin.Context) {
	reply := CreateReply()
	caseGid := ctx.Param("caseGid")
	commentId, _ := strconv.ParseInt(ctx.Query("comment_id"), 10, 32)

	tokenOk, tUser, _ := c.JWTUsecase.JWTAuth(ctx.GetHeader("Authorization"), ctx)
	usePasswordAccess := false
	if tokenOk { // 验证数据权限
		tCase, er := c.RecordbuzUsecase.VerifyDataPermission(caseGid, *tUser)
		if er != nil {
			c.log.Error(er)
		}
		if tCase == nil {
			tokenOk = false
		}
	}
	if !tokenOk { // 非登录base用户，需要使用密码访问
		password := ctx.Query("password")
		tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
		caseId := int32(0)
		if tCase != nil {
			caseId = tCase.Id()
		}
		isOk, _ := c.StatementUsecase.BizStatementVerifyPassword(caseId, password)
		if !isOk { // todo:lgl 验证密码
			reply["code"] = Reply_code_waiting_password
			reply["message"] = "Waiting Password"
			ctx.JSON(200, reply)
			return
		}
		usePasswordAccess = true
	}

	//rawData, _ := ctx.GetRawData()

	data, err := c.StatementCommentBuzUsecase.BizStatementCommentDelete(usePasswordAccess, tUser, caseGid, int32(commentId))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) StatementCommentSave(ctx *gin.Context) {
	reply := CreateReply()
	caseGid := ctx.Param("caseGid")
	commentId, _ := strconv.ParseInt(ctx.Query("comment_id"), 10, 32)

	tokenOk, tUser, _ := c.JWTUsecase.JWTAuth(ctx.GetHeader("Authorization"), ctx)
	usePasswordAccess := false
	if tokenOk { // 验证数据权限
		tCase, er := c.RecordbuzUsecase.VerifyDataPermission(caseGid, *tUser)
		if er != nil {
			c.log.Error(er)
		}
		if tCase == nil {
			tokenOk = false
		}
	}
	if !tokenOk { // 非登录base用户，需要使用密码访问
		password := ctx.Query("password")
		tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
		caseId := int32(0)
		if tCase != nil {
			caseId = tCase.Id()
		}
		isOk, _ := c.StatementUsecase.BizStatementVerifyPassword(caseId, password)
		if !isOk { // todo:lgl 验证密码
			reply["code"] = Reply_code_waiting_password
			reply["message"] = "Waiting Password"
			ctx.JSON(200, reply)
			return
		}
		usePasswordAccess = true
	}

	rawData, _ := ctx.GetRawData()

	data, err := c.StatementCommentBuzUsecase.BizStatementCommentSave(usePasswordAccess, tUser, caseGid, int32(commentId), rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) StatementCommentMarkComplete(ctx *gin.Context) {
	reply := CreateReply()
	caseGid := ctx.Param("caseGid")
	commentId := ctx.Query("commentId")
	commentIdInt, _ := strconv.ParseInt(commentId, 0, 32)

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.StatementCommentBuzUsecase.BizStatementCommentMarkComplete(userFacade, caseGid, int32(commentIdInt), StatementComment_Action_mark)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)

}

const (
	StatementComment_Action_mark   = "mark"
	StatementComment_Action_Unmark = "unmark"
)

func (c *RecordHttpUsecase) StatementCommentUnmarkComplete(ctx *gin.Context) {
	reply := CreateReply()
	caseGid := ctx.Param("caseGid")
	commentId := ctx.Query("commentId")
	commentIdInt, _ := strconv.ParseInt(commentId, 0, 32)

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.StatementCommentBuzUsecase.BizStatementCommentMarkComplete(userFacade, caseGid, int32(commentIdInt), StatementComment_Action_Unmark)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)

}

func (c *RecordHttpUsecase) StatementCommentList(ctx *gin.Context) {
	reply := CreateReply()
	caseGid := ctx.Param("caseGid")

	tokenOk, tUser, _ := c.JWTUsecase.JWTAuth(ctx.GetHeader("Authorization"), ctx)
	usePasswordAccess := false
	if tokenOk { // 验证数据权限
		tCase, er := c.RecordbuzUsecase.VerifyDataPermission(caseGid, *tUser)
		if er != nil {
			c.log.Error(er)
		}
		if tCase == nil {
			tokenOk = false
		}
	}
	if !tokenOk { // 非登录base用户，需要使用密码访问
		password := ctx.Query("password")
		tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
		caseId := int32(0)
		if tCase != nil {
			caseId = tCase.Id()
		}
		isOk, _ := c.StatementUsecase.BizStatementVerifyPassword(caseId, password)
		if !isOk { // todo:lgl 验证密码
			reply["code"] = Reply_code_waiting_password
			reply["message"] = "Waiting Password"
			ctx.JSON(200, reply)
			return
		}
		usePasswordAccess = true
	}

	rawData, _ := ctx.GetRawData()

	data, err := c.StatementCommentBuzUsecase.BizStatementCommentList(usePasswordAccess, tUser, caseGid, rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) StatementSave(ctx *gin.Context) {
	reply := CreateReply()
	caseGid := ctx.Param("caseGid")
	tokenOk, tUser, _ := c.JWTUsecase.JWTAuth(ctx.GetHeader("Authorization"), ctx)
	usePasswordAccess := false
	if tokenOk { // 验证数据权限
		tCase, er := c.RecordbuzUsecase.VerifyDataPermission(caseGid, *tUser)
		if er != nil {
			c.log.Error(er)
		}
		if tCase == nil {
			tokenOk = false
		}
	}
	if !tokenOk { // 非登录base用户，需要使用密码访问
		password := ctx.Query("password")
		tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
		caseId := int32(0)
		if tCase != nil {
			caseId = tCase.Id()
		}
		isOk, _ := c.StatementUsecase.BizStatementVerifyPassword(caseId, password)
		if !isOk { // todo:lgl 验证密码
			reply["code"] = Reply_code_waiting_password
			reply["message"] = "Waiting Password"
			ctx.JSON(200, reply)
			return
		}
		usePasswordAccess = true
	}

	rawData, _ := ctx.GetRawData()

	data, err := c.StatementUsecase.BizStatementSave(usePasswordAccess, tUser, caseGid, rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) LeadsSave(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()

	data, err := c.LeadsUsecase.BizLeadsSave(rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) ClientReviews(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()

	data, err := c.ClientReviewBuzUsecase.BizClientReviews()
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) MedicalCost(ctx *gin.Context) {
	reply := CreateReply()
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	caseGid := ctx.Param("caseGid")
	rawData, _ := ctx.GetRawData()

	data, err := c.BizMedicalCost(userFacade, caseGid, rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) BizMedicalCost(userFacade UserFacade, caseGid string, rawData []byte) (lib.TypeMap, error) {

	if caseGid == "" {
		return nil, errors.New("Parameters Incorrect")
	}
	medicalDbqCostUserInfo := lib.BytesToTDef(rawData, MedicalDbqCostUserInfo{})

	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	c.MedicalDbqCostUsecase.SaveMedicalDbqCostUserInfo(tCase.Id(), medicalDbqCostUserInfo)
	data := make(lib.TypeMap)
	data.Set("medical_dbq_cost", c.MedicalDbqCostUsecase.GetMedicalDbqCost(tCase.Id()))
	return data, nil
}

func (c *RecordHttpUsecase) RelatedClient(ctx *gin.Context) {
	reply := CreateReply()
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	clientGid := ctx.Param("clientGid")

	data, err := c.BizRelatedClient(userFacade, clientGid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) BizRelatedClient(userFacade UserFacade, clientGid string) (lib.TypeMap, error) {

	if clientGid == "" {
		return nil, errors.New("Parameters Incorrect")
	}
	specificFieldNames := []string{FieldName_deal_name, FieldName_amount, FieldName_stages}
	str := `{
	"filter": {
		"operator": "AND",
		"group": [{
			"comparator": "eq",
			"field": {
				"field_name": "client_gid"
			},
			"value": [{
				"value": "` + clientGid + `",
				"label": ""
			}]
		}]
	},
	"table_type": ""
}`
	tListRequest := lib.StringToTDef[*TListRequest](str, nil)
	tList, err := c.RecordbuzUsecase.List(Kind_client_cases, &userFacade, tListRequest, 1, 1000, specificFieldNames, false)
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	data.Set(Fab_TRecords, tList)
	return data, nil
}

func (c *RecordHttpUsecase) List(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	//lib.DPrintln(body)
	moduleName := ctx.Param("module_name")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	page := HandlePage(ctx.Query("page"))
	pageSize := HandlePageSize(ctx.Query("page_size"))
	keyword := strings.TrimSpace(ctx.Query("keyword"))

	data, err := c.BizList(ModuleConvertToKind(moduleName), userFacade, rawData, page, pageSize, keyword)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) BizList(kind string, userFacade UserFacade, rawData []byte, page int, pageSize int, searchKeyword string) (lib.TypeMap, error) {

	tList, total, err := c.DoBizList(kind, userFacade, rawData, page, pageSize, searchKeyword)
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)
	data.Set(Fab_TRecords, tList)
	data.Set(Fab_TTotal, int32(total))
	data.Set(Fab_TPage, page)
	data.Set(Fab_TPageSize, pageSize)

	return data, nil
}

func (c *RecordHttpUsecase) DoBizList(kind string, userFacade UserFacade, rawData []byte, page int, pageSize int, searchKeyword string) (tList TDataList, total int64, err error) {

	var tListRequest *TListRequest

	tProfile, _ := c.UserUsecase.GetProfile(&userFacade.TData)

	//if len(rawData) > 0 {
	tListRequest = &TListRequest{}
	json.Unmarshal(rawData, &tListRequest)
	if kind == Kind_users {
		//if !IsAdminProfile(tProfile) {
		//	return nil, 0, errors.New("No permission to access data")
		//}
	}

	if kind == Kind_client_cases {
		if tListRequest.TableType == SettingHttpCustomViewRequest_TableType_ongoing {
			tListRequest.Filter.Group = append(tListRequest.Filter.Group, TListCondition{
				Comparator: Comparator_neq,
				Field: TListConditionField{
					FieldName: FieldName_stages,
				},
				Value: []interface{}{
					map[string]interface{}{"value": config_vbc.Stages_Completed},
					map[string]interface{}{"value": config_vbc.Stages_Terminated},
					map[string]interface{}{"value": config_vbc.Stages_Dormant},
					map[string]interface{}{"value": config_vbc.Stages_AmCompleted},
					map[string]interface{}{"value": config_vbc.Stages_AmTerminated},
					map[string]interface{}{"value": config_vbc.Stages_AmDormant},
				},
			})
		} else if tListRequest.TableType == SettingHttpCustomViewRequest_TableType_overdue {
			tListRequest.Filter.Group = append(tListRequest.Filter.Group, TListCondition{
				Comparator: Comparator_neq,
				Field: TListConditionField{
					FieldName: FieldName_stages,
				},
				Value: []interface{}{
					map[string]interface{}{"value": config_vbc.Stages_Completed},
					map[string]interface{}{"value": config_vbc.Stages_Terminated},
					map[string]interface{}{"value": config_vbc.Stages_Dormant},
				},
			})
			currentTime := time.Now().In(configs.GetVBCDefaultLocation())
			tListRequest.Filter.Group = append(tListRequest.Filter.Group, TListCondition{
				Comparator: Comparator_lt,
				Field: TListConditionField{
					FieldName: DataEntry_sys__due_date,
				},
				Value: []interface{}{
					currentTime.Format(time.DateOnly),
				},
			})
		} else if tListRequest.TableType == SettingHttpCustomViewRequest_TableType_upcoming {
			tListRequest.Filter.Group = append(tListRequest.Filter.Group, TListCondition{
				Comparator: Comparator_neq,
				Field: TListConditionField{
					FieldName: FieldName_stages,
				},
				Value: []interface{}{
					map[string]interface{}{"value": config_vbc.Stages_Completed},
					map[string]interface{}{"value": config_vbc.Stages_Terminated},
					map[string]interface{}{"value": config_vbc.Stages_Dormant},
				},
			})
			currentTime := time.Now().In(configs.GetVBCDefaultLocation())
			currentTime = currentTime.AddDate(0, 0, -1)
			tListRequest.Filter.Group = append(tListRequest.Filter.Group, TListCondition{
				Comparator: Comparator_gt,
				Field: TListConditionField{
					FieldName: DataEntry_sys__due_date,
				},
				Value: []interface{}{
					currentTime.Format(time.DateOnly),
				},
			})
		} else {
			if !IsAdminProfile(tProfile) && !HaveAllDataPermissions(kind, tProfile.Gid()) {
				tListRequest.Filter.Group = append(tListRequest.Filter.Group, TListCondition{
					Comparator: Comparator_neq,
					Field: TListConditionField{
						FieldName: FieldName_stages,
					},
					Value: []interface{}{
						map[string]interface{}{"value": config_vbc.Stages_Completed},
						map[string]interface{}{"value": config_vbc.Stages_Terminated},
						map[string]interface{}{"value": config_vbc.Stages_AmCompleted},
						map[string]interface{}{"value": config_vbc.Stages_AmTerminated},
					},
				})
			}
		}
	}
	if tListRequest.TableType == SettingHttpCustomViewRequest_TableType_Search {
		if kind == Kind_clients || kind == Kind_client_cases {
			if searchKeyword == "" {
				return nil, 0, nil
			}
			tListRequest.Filter.Operator = TFilterVo_Operator_OR
			fieldStruct, err := c.FieldUsecase.CacheStructByKind(kind)
			if err != nil {
				return nil, 0, err
			}
			var filterGroup TListConditions
			for _, v := range fieldStruct.Records {
				if IsTextField(v.FieldType) {
					filterGroup = append(filterGroup, TListCondition{
						Comparator: Comparator_contains,
						Field: TListConditionField{
							FieldName: v.FieldName,
						},
						Value: []interface{}{
							searchKeyword,
						},
					})
				}
			}
			tListRequest.Filter.Group = filterGroup
		}
	}

	//if err != nil {
	//	return nil, 0, err
	//}
	//}

	var normalUserOnlyOwner bool
	if tListRequest.TableType == SettingHttpCustomViewRequest_TableType_ongoing ||
		tListRequest.TableType == SettingHttpCustomViewRequest_TableType_overdue ||
		tListRequest.TableType == SettingHttpCustomViewRequest_TableType_upcoming {
		normalUserOnlyOwner = true
	}

	tList, err = c.RecordbuzUsecase.List(kind, &userFacade, tListRequest, page, pageSize, nil, normalUserOnlyOwner)
	if err != nil {
		return nil, 0, err
	}
	fieldStruct, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, 0, err
	}
	if fieldStruct == nil {
		return nil, 0, errors.New("fieldStruct is nil")
	}

	// 格式化时间
	for k, _ := range tList {
		err := tList[k].HandleTDataTimezone(c.TimezonesUsecase, *fieldStruct, &userFacade)
		if err != nil {
			return nil, 0, err
		}
	}

	total, err = c.RecordbuzUsecase.Total(kind, &userFacade, tListRequest, normalUserOnlyOwner)
	if err != nil {
		return nil, 0, err
	}
	return tList, total, nil
}

type RecordHttpDetailVo struct {
	Gid string `json:"gid"`
}

func (c *RecordHttpUsecase) Related(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	moduleName := ctx.Param("module_name")
	gid := ctx.Param("gid")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.BizRelated(ModuleConvertToKind(moduleName), gid, userFacade, rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type ResponseRelated struct {
	NotesCount int32 `json:"notes_count"`
	CasesCount int32 `json:"cases_count"`
}

func (c *RecordHttpUsecase) BizRelated(kind string, gid string, userFacade UserFacade, rawData []byte) (lib.TypeMap, error) {

	if gid == "" {
		return nil, errors.New("Incorrect parameter")
	}

	var responseRelated ResponseRelated
	notesKindEntity, err := c.KindUsecase.GetByKind(Kind_notes)
	if err != nil {
		return nil, err
	}
	if notesKindEntity == nil {
		return nil, errors.New("notesKindEntity is nil")
	}

	notesCount, _ := c.TUsecase.TotalByCond(*notesKindEntity, Eq{
		Notes_FieldName_kind:     kind,
		Notes_FieldName_kind_gid: gid,
		DataEntry_deleted_at:     0,
	})
	responseRelated.NotesCount = int32(notesCount)
	if kind == Kind_clients {
		str := `{
	"filter": {
		"operator": "AND",
		"group": [{
			"comparator": "eq",
			"field": {
				"field_name": "client_gid"
			},
			"value": [{
				"value": "` + gid + `",
				"label": ""
			}]
		}]
	},
	"table_type": ""
}`
		tListRequest := lib.StringToTDef[*TListRequest](str, nil)
		total, err := c.RecordbuzUsecase.Total(Kind_client_cases, &userFacade, tListRequest, false)
		if err != nil {
			return nil, err
		}
		responseRelated.CasesCount = int32(total)
	}
	data := make(lib.TypeMap)
	data.Set("related", responseRelated)
	return data, nil
}

const (
	Record_From_create = "create"
	Record_From_edit   = "edit"
	Record_From_detail = "detail"
)

func (c *RecordHttpUsecase) Edit(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	//lib.DPrintln(body)
	moduleName := ctx.Param("module_name")
	gid := ctx.Param("gid")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.BizDetail(ModuleConvertToKind(moduleName), gid, userFacade, rawData, Record_From_edit)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) Delete(ctx *gin.Context) {
	reply := CreateReply()
	moduleName := ctx.Param("module_name")
	gid := ctx.Param("gid")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.BizDelete(ModuleConvertToKind(moduleName), gid, userFacade)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) BizCustomKindDelete(kind string, gid string, userFacade UserFacade) (lib.TypeMap, error) {

	if gid == "" {
		return nil, errors.New("Incorrect parameter")
	}

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}

	if kind == Kind_Custom_Filter {
		entity, err := c.FilterUsecase.GetByCond(Eq{"user_gid": userFacade.Gid(), "id": gid, "biz_deleted_at": 0, "deleted_at": 0})
		if err != nil {
			return nil, err
		}
		if entity == nil {
			return nil, errors.New("The record does not exist or has no permission to perform operations")
		}
	} else {
		return nil, errors.New("Nonsupport")
	}

	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_Incr_id_name] = gid
	dataEntry[DataEntry_biz_deleted_at] = time.Now().Unix()
	dataEntry[DataEntry_modified_by] = userFacade.Gid()
	_, err = c.DataEntryUsecase.HandleOne(kind, dataEntry, DataEntry_Incr_id_name, &userFacade.TData)
	return nil, nil
}

func (c *RecordHttpUsecase) BizDelete(kind string, gid string, userFacade UserFacade) (lib.TypeMap, error) {

	if IsCustomKind(kind) {
		return c.BizCustomKindDelete(kind, gid, userFacade)
	}

	if gid == "" {
		return nil, errors.New("Incorrect parameter")
	}

	tProfile, _ := c.UserUsecase.GetProfile(&userFacade.TData)
	if tProfile == nil {
		return nil, errors.New("tProfile is nil")
	}

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}

	searchCls := c.RecordbuzSearchUsecase.NewRecordbuzSearchCls(false)
	hasDeletePermission, err := searchCls.HasDeletePermission(userFacade, *tProfile, *kindEntity, gid)
	if err != nil {
		return nil, err
	}
	if !hasDeletePermission {
		return nil, errors.New("Records do not exist or have no permission to delete")
	}

	tData, err := c.TUsecase.Data(kind, Eq{"gid": gid, "deleted_at": 0, "biz_deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if tData == nil {
		return nil, errors.New("tData is nil")
	}
	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_gid] = gid
	dataEntry[DataEntry_biz_deleted_at] = time.Now().Unix()
	dataEntry[DataEntry_modified_by] = userFacade.Gid()

	_, err = c.DataEntryUsecase.HandleOne(kind, dataEntry, DataEntry_gid, &userFacade.TData)
	if err != nil {
		return nil, err
	}
	// 记录删除日志
	c.LogUsecase.SaveLog(tData.Id(), kind, dataEntry)
	if kind == Kind_clients {
		cases, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{FieldName_client_gid: tData.Gid(), DataEntry_biz_deleted_at: 0, DataEntry_deleted_at: 0})
		if err != nil {
			return nil, err
		}
		var dataEntryList TypeDataEntryList
		for _, v := range cases {
			entry := make(TypeDataEntry)
			entry[DataEntry_gid] = v.Gid()
			entry[DataEntry_biz_deleted_at] = time.Now().Unix()
			entry[DataEntry_modified_by] = userFacade.Gid()
			dataEntryList = append(dataEntryList, entry)
		}
		if len(dataEntryList) > 0 {
			c.DataEntryUsecase.Handle(Kind_client_cases, dataEntryList, DataEntry_gid, &userFacade.TData)
		}
	}

	if kind == Kind_client_tasks {
		tTask := tData
		whatGid := tTask.CustomFields.TextValueByNameBasic(TaskFieldName_what_id_gid)
		if whatGid != "" {
			c.QueueUsecase.PushClientTaskHandleWhatGidJobTasks(context.TODO(), []string{whatGid})
		}
		whoGid := tTask.CustomFields.TextValueByNameBasic(TaskFieldName_who_id_gid)
		if whoGid != "" {
			c.QueueUsecase.PushClientTaskHandleWhoGidJobTasks(context.TODO(), []string{whoGid})
		}
	}

	return nil, nil
}

func (c *RecordHttpUsecase) Detail(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	//lib.DPrintln(body)
	moduleName := ctx.Param("module_name")
	gid := ctx.Param("gid")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.BizDetail(ModuleConvertToKind(moduleName), gid, userFacade, rawData, Record_From_detail)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

const (
	RecordHttpFormColumnRequest_Type_Account                     = "account"
	RecordHttpFormColumnRequest_Type_FormCondition               = "form_condition"
	RecordHttpFormColumnRequest_Type_FormConditionSecondary      = "form_condition_secondary"
	RecordHttpFormColumnRequest_Type_FormFilter                  = "form_filter"
	RecordHttpFormColumnRequest_Type_FormConditionQuestionnaires = "form_condition_questionnaires"
	RecordHttpFormColumnRequest_Type_Form_mgmt_user              = "form_mgmt_user"
	RecordHttpFormColumnRequest_Type_Form_mgmt_attorney          = "form_mgmt_attorney"
)

var FormColumnMappingKind = map[string]string{
	RecordHttpFormColumnRequest_Type_Account:                     Kind_users,
	RecordHttpFormColumnRequest_Type_FormCondition:               Kind_Custom_Condition,
	RecordHttpFormColumnRequest_Type_FormConditionSecondary:      Kind_Custom_ConditionSecondary,
	RecordHttpFormColumnRequest_Type_FormFilter:                  Kind_Custom_Filter,
	RecordHttpFormColumnRequest_Type_FormConditionQuestionnaires: Kind_Custom_ConditionQuestionnaires,
	RecordHttpFormColumnRequest_Type_Form_mgmt_user:              Kind_users,
	RecordHttpFormColumnRequest_Type_Form_mgmt_attorney:          Kind_attorneys,
}

//type RecordHttpFormColumnRequest struct {
//	Type string `json:"type"`
//}

func (c *RecordHttpUsecase) FormColumn(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	params := lib.ToTypeMapByString(string(rawData))

	//recordHttpFormColumnRequest, _ := lib.StringToTE[RecordHttpFormColumnRequest](string(rawData), RecordHttpFormColumnRequest{})
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizFormColumn(userFacade, params)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) BizFormColumn(userFacade UserFacade, params lib.TypeMap) (lib.TypeMap, error) {

	//gid := ""

	kind := ""
	conditionId := params.GetInt("condition_id")
	formColumnType := params.GetString("type")
	uniqId := params.GetString("uniq_id") // 有可能是id：int 或gid：string

	if v, ok := FormColumnMappingKind[formColumnType]; ok {
		kind = v
	} else {
		return nil, errors.New("Parameters Incorrect")
	}

	if formColumnType == RecordHttpFormColumnRequest_Type_Account {
		uniqId = userFacade.Gid()
	}
	tProfile, err := c.UserUsecase.GetProfile(&userFacade.TData)
	if err != nil {
		return nil, err
	}
	if tProfile == nil {
		return nil, errors.New("tProfile is nil")
	}
	settingSectionVo, err := c.SettingSectionFieldUsecase.GetSettingSectionsByKindForFormColumn(kind,
		*tProfile,
		formColumnType)

	if err != nil {
		return nil, err
	}
	isUpdated := true
	if (kind == Kind_users || kind == Kind_attorneys) && uniqId == "" {
		isUpdated = false
	} else if IsCustomKind(kind) {
		if kind == Kind_Custom_Condition || kind == Kind_Custom_ConditionSecondary {
			if conditionId == 0 {
				isUpdated = false
			}
		} else if uniqId == "" {
			isUpdated = false
		}
	}
	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}
	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}

	var tRecord *TData
	var detailWhole DetailWhole
	if isUpdated {
		if IsCustomKind(kind) {
			if kind == Kind_Custom_Condition || kind == Kind_Custom_ConditionSecondary {
				tRecord, err = c.TUsecase.Data(kind, Eq{"id": conditionId, "deleted_at": 0})
				if err != nil {
					return nil, err
				}
			} else {
				tRecord, err = c.TUsecase.Data(kind, Eq{"id": uniqId, "deleted_at": 0})
				if err != nil {
					return nil, err
				}
			}
			if Kind_Custom_Condition == kind {
				c.ConditionbuzUsecase.HandleTCondition(tRecord)
			}
		} else {
			tRecord, err = c.TUsecase.Data(kind, Eq{"gid": uniqId, "deleted_at": 0, "biz_deleted_at": 0})
		}

		if err != nil {
			return nil, err
		}
		if tRecord == nil {
			return nil, errors.New("tRecord is nil")
		}
		err = tRecord.HandleTDataTimezone(c.TimezonesUsecase, *structField, &userFacade)
		if err != nil {
			return nil, err
		}
		detailWhole.PrimaryName = tRecord.CustomFields.TextValueByNameBasic(kindEntity.PrimaryFieldName)
		detailWhole.UpdatedAt = tRecord.CustomFields.NumberValueByNameBasic(DataEntry_updated_at)
		createdAt := tRecord.CustomFields.NumberValueByNameBasic(DataEntry_created_at)
		if createdAt > 0 {
			timeLocation, err := userFacade.GetTimeLocation(c.TimezonesUsecase)
			if err != nil {
				c.log.Error(err)
			} else {
				a := time.Unix(int64(createdAt), 0).In(timeLocation)
				detailWhole.CreatedTime = a.Format(configs.TimeFormatDate)
			}
		}
	}
	if kind == Kind_client_cases {
		fieldPermissionCenter, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, userFacade.ProfileGid())
		if err != nil {
			return nil, err
		}
		fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(FieldName_amount)
		if err != nil {
			return nil, err
		}
		if fieldPermissionVo.CanShow() {
			detailWhole.Amount = tRecord.CustomFields.DisplayValueByNameBasic(FieldName_amount)
		}
	}

	var detailSections DetailSections
	var fieldNames []string
	for _, v := range settingSectionVo.Sections {
		var detailSection DetailSection
		detailSection.SectionName = v.SectionName
		detailSection.SectionLabel = v.SectionLabel

		for _, v1 := range v.Left.Fields {

			fieldNames = append(fieldNames, v1.FieldName)
			detailField, err := v1.ToDetailField(structField, tRecord, c.StatementUsecase)
			if err != nil {
				c.log.Error(err)
				continue
			}
			if detailField == nil {
				continue
			}
			detailSection.Left.Fields = append(detailSection.Left.Fields, *detailField)
		}
		for _, v1 := range v.Right.Fields {
			fieldNames = append(fieldNames, v1.FieldName)
			detailField, err := v1.ToDetailField(structField, tRecord, c.StatementUsecase)
			if err != nil {
				c.log.Error(err)
				continue
			}
			if detailField == nil {
				continue
			}
			detailSection.Right.Fields = append(detailSection.Right.Fields, *detailField)
		}
		detailSections = append(detailSections, detailSection)
	}
	detailWhole.Sections = detailSections

	// fieldNames 经过了权限过虑
	fabFields, err := c.FieldbuzUsecase.FabFieldsForBasicdata(kind, &userFacade)
	if err != nil {
		return nil, err
	}
	lib.DPrintln("fieldNames:", fieldNames)
	lib.DPrintln("fabFields:", fabFields)
	var destFabFields []FabField
	for k, v := range fabFields {
		if lib.InArray(v.FieldName, fieldNames) {
			destFabFields = append(destFabFields, fabFields[k])
		}
	}
	data := make(lib.TypeMap)
	data.Set("data", detailWhole)
	data.Set("fields", destFabFields)
	return data, nil
}

func (c *RecordHttpUsecase) BizDetail(kind string, gid string, userFacade UserFacade, rawData []byte, from string) (lib.TypeMap, error) {

	//var recordHttpDetailVo RecordHttpDetailVo
	//err := json.Unmarshal(rawData, &recordHttpDetailVo)
	//if err != nil {
	//	return nil, err
	//}
	if gid == "" {
		return nil, errors.New("Incorrect parameter")
	}

	settingSectionVo, err := c.SettingSectionFieldUsecase.GetSettingSectionsByKind(kind,
		userFacade.CustomFields.TextValueByNameBasic(User_FieldName_profile_gid),
		from)
	if err != nil {
		return nil, err
	}

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}

	var tCase *TData

	tCase, err = c.RecordbuzUsecase.GetRecordData(gid, *kindEntity, &userFacade, false)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("The record does not exist or has no permission")
	}

	/*
		tCase, err := c.TUsecase.Data(kind, Eq{"gid": gid, "deleted_at": 0, "biz_deleted_at": 0})
		if err != nil {
			return nil, err
		}
		if tCase == nil {
			return nil, errors.New("tCase is nil")
		}*/

	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}

	err = tCase.HandleTDataTimezone(c.TimezonesUsecase, *structField, &userFacade)
	if err != nil {
		return nil, err
	}

	tProfile, err := c.UserUsecase.GetProfile(&userFacade.TData)
	if err != nil {
		return nil, err
	}
	if tProfile == nil {
		return nil, errors.New("tProfile is nil")
	}

	var detailWhole DetailWhole

	detailWhole.PrimaryName = tCase.CustomFields.TextValueByNameBasic(kindEntity.PrimaryFieldName)
	detailWhole.UpdatedAt = tCase.CustomFields.NumberValueByNameBasic(DataEntry_updated_at)
	createdAt := tCase.CustomFields.NumberValueByNameBasic(DataEntry_created_at)
	if createdAt > 0 {
		timeLocation, err := userFacade.GetTimeLocation(c.TimezonesUsecase)
		if err != nil {
			c.log.Error(err)
		} else {
			a := time.Unix(int64(createdAt), 0).In(timeLocation)
			detailWhole.CreatedTime = a.Format(time.DateOnly)
		}
	}

	if kind == Kind_client_cases {
		fieldPermissionCenter, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, userFacade.ProfileGid())
		if err != nil {
			return nil, err
		}
		fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(FieldName_amount)
		if err != nil {
			return nil, err
		}
		if fieldPermissionVo.CanShow() {
			detailWhole.Amount = tCase.CustomFields.DisplayValueByNameBasic(FieldName_amount)
		}
		var detailWholeCaseExtras DetailWholeCaseExtras
		stagesLog, _ := c.RecordLogUsecase.BizCrmStagesLatest(tCase.Gid())
		if stagesLog != nil {
			a, _ := GenTFieldExtendForSysDueDate(TimestampToDate(stagesLog.StartTime))
			if a != nil {
				detailWholeCaseExtras.StageStartDate = *a
			}
			a, _ = GenTFieldExtendForSysDueDate(TimestampToDate(stagesLog.EndTime))
			if a != nil {
				detailWholeCaseExtras.StageDueDate = *a
			}
		}
		c.TUsecase.DoFormula(*kindEntity, tCase)
		itfFormulaVal := tCase.CustomFields.TextValueByNameBasic(DataEntry_sys__itf_formula)
		itfFormula, _ := GenTFieldExtendForSysItfFormula(itfFormulaVal)
		if itfFormula != nil {
			detailWholeCaseExtras.ItfFormula = *itfFormula
		}
		itfExpiration := GenItfExpirationExtend(tCase)
		if itfExpiration != nil {
			detailWholeCaseExtras.ItfExpiration = *itfExpiration
		}
		detailWholeCaseExtras.MedicalDbqCost = c.MedicalDbqCostUsecase.GetMedicalDbqCost(tCase.Id())
		detailWholeCaseExtras.Pipelines = GetPipelinesByCase(tCase)
		detailWhole.Extras = detailWholeCaseExtras
		if userFacade.Id() == 4 { // todo:lgl 测试帐号先开启
			detailWhole.ShowRecordReview = true
		}
	} else if kind == Kind_clients {
		clientPipelines, err := c.ClientUsecase.GetOneClientPipeline(tCase.Gid())
		if err != nil {
			c.log.Error(err)
		}
		detailWhole.ClientPipelines = clientPipelines
	}

	var detailSections DetailSections
	for _, v := range settingSectionVo.Sections {

		var detailSection DetailSection
		detailSection.SectionName = v.SectionName
		detailSection.SectionLabel = v.SectionLabel

		for _, v1 := range v.Left.Fields {

			if !IsAdminProfile(tProfile) && IsAmContract(*tCase) {
				if v1.FieldName == FieldName_pricing_version {
					continue
				}
			}

			detailField, err := v1.ToDetailField(structField, tCase, c.StatementUsecase)
			if err != nil {
				c.log.Error(err)
				continue
			}
			if detailField == nil {
				continue
			}
			detailSection.Left.Fields = append(detailSection.Left.Fields, *detailField)
		}
		for _, v1 := range v.Right.Fields {

			if IsAmContract(*tCase) {

				if !IsAdminProfile(tProfile) && v1.FieldName == FieldName_pricing_version {
					continue
				}
			}

			detailField, err := v1.ToDetailField(structField, tCase, c.StatementUsecase)
			if err != nil {
				c.log.Error(err)
				continue
			}
			if detailField == nil {
				continue
			}
			detailSection.Right.Fields = append(detailSection.Right.Fields, *detailField)
		}
		detailSections = append(detailSections, detailSection)
	}
	detailWhole.Sections = detailSections
	//detailWhole.Pipelines = GetPipelinesByCase(tCase)

	data := make(lib.TypeMap)
	data.Set("detail", detailWhole)
	//data.Set("data.deal_name", dealName)
	return data, nil
}

func GetPipelinesByCase(tCase *TData) string {
	if tCase != nil {
		contractSource := tCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource)
		return ContractSourceToPipelines(contractSource)
	}
	return Pipelines_default
}

func PipelinesToContractSource(pipelines string) string {
	if pipelines == Pipelines_am {
		return ContractSource_AM
	} else if pipelines == Pipelines_vbc {
		return ContractSource_VBC
	}
	return ""
}

func ContractSourceToPipelines(contractSource string) string {
	if contractSource == ContractSource_AM {
		return Pipelines_am
	} else if contractSource == ContractSource_VBC {
		return Pipelines_vbc
	}
	return Pipelines_default
}

func (c *RecordHttpUsecase) Timelines(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	//lib.DPrintln(body)
	moduleName := ctx.Param("module_name")
	gid := ctx.Param("gid")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	page := HandlePage(ctx.Query("page"))
	pageSize := HandlePageSize(ctx.Query("page_size"))

	data, err := c.BizTimelines(ModuleConvertToKind(moduleName), gid, userFacade, rawData, page, pageSize)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) BizTimelines(kind string, gid string, userFacade UserFacade, rawData []byte, page int, pageSize int) (lib.TypeMap, error) {

	if gid == "" {
		return nil, errors.New("Incorrect parameter")
	}

	tData, err := c.TUsecase.Data(kind, Eq{"gid": gid, "deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if tData == nil {
		return nil, errors.New("tData is nil")
	}

	res, total, err := c.TimelinesbuzUsecase.List(kind, gid, &userFacade, page, pageSize)

	data := make(lib.TypeMap)
	data.Set("timelines", res)
	data.Set(Fab_TTotal, int32(total))
	data.Set(Fab_TPage, page)
	data.Set(Fab_TPageSize, pageSize)

	if int64(page*pageSize) >= total {
		data.Set(Fab_HasMore, false)
	} else {
		data.Set(Fab_HasMore, true)
	}

	//sql := fmt.Sprintf("")
	return data, err
}

func (c *RecordHttpUsecase) Layout(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	//lib.DPrintln(body)
	moduleName := ctx.Param("module_name")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.BizLayout(ModuleConvertToKind(moduleName), userFacade, rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *RecordHttpUsecase) InitUserGidValue(detailField *DetailField, structField *TypeFieldStruct, userFacade UserFacade) {
	if detailField == nil {
		return
	}
	if detailField.FieldName == FieldName_user_gid {
		fieldEntity := structField.GetByFieldName(FieldName_user_gid)
		if fieldEntity == nil {
			return
		}
		gid := userFacade.Gid()
		displayName := userFacade.CustomFields.DisplayValueByNameBasic(fieldEntity.RelaName)
		detailField.Value = ValueOp{
			Label: displayName,
			Value: gid,
		}
	}
}

func (c *RecordHttpUsecase) BizLayout(kind string, userFacade UserFacade, rawData []byte) (lib.TypeMap, error) {

	settingSectionVo, err := c.SettingSectionFieldUsecase.GetSettingSectionsByKind(kind, userFacade.ProfileGid(), Record_From_create)
	if err != nil {
		return nil, err
	}

	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}

	var detailWhole DetailWhole
	var detailSections DetailSections
	for _, v := range settingSectionVo.Sections {

		var detailSection DetailSection
		detailSection.SectionName = v.SectionName
		detailSection.SectionLabel = v.SectionLabel

		for _, v1 := range v.Left.Fields {

			detailField, err := v1.ToDetailField(structField, nil, c.StatementUsecase)
			if err != nil {
				c.log.Error(err)
				continue
			}
			if detailField == nil {
				continue
			}
			c.InitUserGidValue(detailField, structField, userFacade)
			detailSection.Left.Fields = append(detailSection.Left.Fields, *detailField)

		}
		for _, v1 := range v.Right.Fields {
			detailField, err := v1.ToDetailField(structField, nil, c.StatementUsecase)
			if err != nil {
				c.log.Error(err)
				continue
			}
			if detailField == nil {
				continue
			}
			detailSection.Right.Fields = append(detailSection.Right.Fields, *detailField)
		}
		detailSections = append(detailSections, detailSection)
	}
	detailWhole.Sections = detailSections

	data := make(lib.TypeMap)
	data.Set("layout", detailWhole)
	//data.Set("data.deal_name", dealName)
	return data, nil
}

func (c *RecordHttpUsecase) Save(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	moduleName := ctx.Param("module_name")
	gid := ctx.Param("gid")
	createACase := ctx.Query("create_a_case")
	isCreateACase := false
	if createACase == "true" {
		isCreateACase = true
	}
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	useZoho := false
	if configs.IsProd() {
		useZoho = true
	}
	if configs.StoppedZoho {
		useZoho = false
	}
	code, data, err := c.BizSave(ModuleConvertToKind(moduleName), gid, userFacade, body, useZoho, ctx, isCreateACase)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		if code != 0 {
			reply.Update(code, "")
		} else {
			reply.Success()
		}
	}
	ctx.JSON(200, reply)
}

func MultilookupTidyValuesString(values string) (r string) {
	news := strings.Split(values, ",")
	for _, v := range news {
		if v != "" {
			if r == "" {
				r = "," + v + ","
			} else {
				r += v + ","
			}
		}
	}
	return
}

const (
	Source_External_Referral = "External Referral"
	Source_Team_Referral     = "Team Referral"
)

func (c *RecordbuzUsecase) GetRecordData(gid string, kindEntity KindEntity, userFacade *UserFacade, normalUserOnlyOwner bool) (tData *TData, err error) {

	kind := kindEntity.Kind
	if IsCustomKind(kind) {
		tData, err = c.TUsecase.Data(kind, And(Eq{"id": gid, "deleted_at": 0}))
		return
	} else {
		cond, err := c.FilterDataPermissionCond(kindEntity, userFacade, normalUserOnlyOwner)
		if err != nil {
			return nil, err
		}
		//lib.DPrintln("__:", cond)
		// 此处判断有无此数据权限
		tData, err = c.TUsecase.Data(kind, And(Eq{"gid": gid, "biz_deleted_at": 0, "deleted_at": 0}, cond))
		if err != nil {
			return nil, err
		}
	}
	return tData, nil
}

func (c *RecordHttpUsecase) BizSave(kind string, gid string, userFacade UserFacade, body lib.TypeMap, useZoho bool, ctx *gin.Context, isCreateACase bool) (int, lib.TypeMap, error) {

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return 0, nil, err
	}
	if kindEntity == nil {
		return 0, nil, errors.New("kindEntity is nil")
	}
	isUpdate := false
	verifyDestData := make(lib.TypeMap)

	tProfile, err := c.UserUsecase.GetProfile(&userFacade.TData)
	if err != nil {
		return 0, nil, err
	}
	if tProfile == nil {
		return 0, nil, errors.New("tProfile is nil")
	}
	var tData *TData

	if kind == Kind_users {
		pwdValue := body.GetString("mail_password.value")
		if pwdValue != "" {
			pwdValue, err = EncryptSensitive(pwdValue)
			if err != nil {
				return 0, nil, err
			}
			body.Set("mail_password.value", pwdValue)
		}
	}

	if gid != "" { // 更新操作
		if kind == Kind_users {
			if !IsAdminProfile(tProfile) { // 非管理员禁止修改
				delete(body, User_FieldName_profile_gid)
			}
		}

		tData, err = c.RecordbuzUsecase.GetRecordData(gid, *kindEntity, &userFacade, false)
		if err != nil {
			return 0, nil, err
		}
		if tData == nil {
			return 0, nil, errors.New("The record does not exist or has no permission")
		}

		isUpdate = true
		dbData := tData.CustomFields.ToMaps()
		//ccc := tData.CustomFields.TextValueByNameBasic(FieldName_amount)
		//ccc1 := tData.CustomFields.DisplayValueByNameBasic(FieldName_amount)
		//lib.DPrintln(ccc)
		//lib.DPrintln(ccc1)
		verifyDestData = dbData
		for k, _ := range body {
			verifyDestData.Set(k, body.GetString(k+".value"))
		}

	} else {
		if kind == Kind_users {
			if !IsAdminProfile(tProfile) { // 非管理员，只能加标准用户
				body.Set(User_FieldName_profile_gid+".value", Profile_Standard_Gid)
			}
		}
		for k, _ := range body {
			verifyDestData.Set(k, body.GetString(k+".value"))
		}
	}
	c.log.Info("verifyDestData: ", verifyDestData, "gid:", gid, "isUpdated:", isUpdate)

	structField, err := c.FieldUsecase.StructByKind(kind)
	if err != nil {
		return 0, nil, err
	}

	fieldPermissionCenter, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, userFacade.ProfileGid())
	if err != nil {
		return 0, nil, errors.New("Permission configuration error")
	}

	fieldValidatorCenter, err := c.FieldValidatorUsecase.CacheFieldValidatorCenter(kind)
	if err != nil {
		return 0, nil, err
	}
	//lib.DPrintln(fieldValidatorCenter)
	destData := make(lib.TypeMap)

	var verifyFailureResultList VerifyFailureResultList

	if !isUpdate {

		if kind == Kind_clients {
			email := body.GetString("email.value")
			phone := body.GetString("phone.value")
			if email == "" && phone == "" {
				verifyFailureResultList = append(verifyFailureResultList, VerifyFailureResultItem{
					ModuleName: KindConvertToModule(Kind_clients),
					FieldName:  FieldName_email,
					Message:    "Email and Mobile cannot be empty at the same time",
				})
				verifyFailureResultList = append(verifyFailureResultList, VerifyFailureResultItem{
					ModuleName: KindConvertToModule(Kind_clients),
					FieldName:  FieldName_phone,
					Message:    "Email and Mobile cannot be empty at the same time",
				})
			}
		}
		if kind == Kind_clients || kind == Kind_client_cases {
			sourceValue := body.GetString(FieldName_source + ".value")
			if sourceValue == Source_External_Referral || sourceValue == Source_Team_Referral {
				referrerValue := body.GetString(FieldName_referrer + ".value")
				if referrerValue == "" {
					verifyFailureResultList = append(verifyFailureResultList, VerifyFailureResultItem{
						ModuleName: KindConvertToModule(kind),
						FieldName:  FieldName_referrer,
						Message:    "Referring Person cannot be empty",
					})
				}
			}

		}

		for _, v := range structField.Records {

			fieldPermission, err := fieldPermissionCenter.PermissionByFieldName(v.FieldName)
			if err != nil {
				return 0, nil, errors.New(v.FieldName + " Permission authentication error")
			}
			if fieldPermission.CanWrite() { // 只有允许写入的，才验证
				_, _, list, err := fieldValidatorCenter.Verify(v.FieldName, body.GetString(v.FieldName+".value"), TypeDataEntry(verifyDestData))
				if err != nil {
					return 0, nil, err
				}
				verifyFailureResultList = append(verifyFailureResultList, list...)
			}
		}
		verifyFailureResultList = verifyFailureResultList.RemoveDuplicateResult()
		if len(verifyFailureResultList) > 0 {
			data := make(lib.TypeMap)
			data.Set("verify_failure_result", verifyFailureResultList)
			return Reply_code_data_validation_failure, data, nil
		}
	}

	for fieldName, _ := range body {
		fieldEntity := structField.GetByFieldName(fieldName)
		if fieldEntity == nil {
			return 0, nil, errors.New(fieldName + " does not exists")
		}

		if fieldEntity.FieldType == FieldType_multilookup || fieldEntity.FieldType == FieldType_multidropdown {
			body.Set(fieldName+".value", MultilookupTidyValuesString(body.GetString(fieldName+".value")))
		}

		// 权限验证
		fieldPermission, err := fieldPermissionCenter.PermissionByFieldName(fieldName)
		if err != nil {
			return 0, nil, errors.New(fieldName + " Permission authentication error")
		}
		if !fieldPermission.CanWrite() {
			return 0, nil, errors.New(fieldName + " No write permission")
		}
		// todo:lgl 各类字段的值验证和过虑

		if isUpdate { // 更新
			if kind == Kind_client_cases {
				if fieldName == FieldName_statements {
					if tData == nil {
						return 0, nil, errors.New("No record was found.")
					}
					//useNewPersonalWebForm, err := c.StatementUsecase.IsUseNewPersonalWebForm(tData.Id())
					//if err != nil {
					//	return 0, nil, err
					//}
					//if useNewPersonalWebForm {
					//
					//	stageNumber, err := vbc_config.GetStageNumber(tData.CustomFields.DisplayValueByNameBasic(FieldName_stages))
					//	if err != nil {
					//		return 0, nil, err
					//	}
					//	if stageNumber >= vbc_config.Stages_StatementDrafts_Number {
					//		return 0, nil, errors.New("Clients with 'Client PW Active' status who need to update their conditions can do so in the Personal Statement Manager.")
					//	}
					//}
				} else if fieldName == FieldName_stages {
					pipelines := GetPipelinesByCase(tData)
					if pipelines != Pipelines_default {
						fieldOption, err := c.FieldOptionUsecase.GetByFieldName(Kind_client_cases, fieldName, body.GetString(fieldName+".value"))
						if err != nil {
							return 0, nil, err
						}
						if fieldOption == nil {
							return 0, nil, errors.New("stages: fieldOption is nil")
						}
						if fieldOption.Pipelines != pipelines {
							return 0, nil, errors.New("It cannot be switched to this stage")
						}
					}
				} else if fieldName == FieldName_ContractSource {
					stages := tData.CustomFields.TextValueByNameBasic(FieldName_stages)
					if stages != config_vbc.Stages_IncomingRequest && stages != config_vbc.Stages_AmIncomingRequest {
						return 0, nil, errors.New("The Contract Entity is not allowed to be modified at this stage")
					}
					newContractSource := body.GetString(fieldName + ".value")
					if newContractSource == "" {
						return 0, nil, errors.New("The Contract Entity cannot be empty")
					}
					if newContractSource == ContractSource_VBC {
						destData.Set(FieldName_stages, config_vbc.Stages_IncomingRequest)
					} else {
						destData.Set(FieldName_stages, config_vbc.Stages_AmIncomingRequest)
					}
				}
			}
		} else { // 创建
			if kind == Kind_client_cases {
				if fieldName == FieldName_stages {
					historyPipelines, err := c.ClientUsecase.GetOneClientPipeline(body.GetString(FieldName_client_gid + ".value"))
					if err != nil {
						c.log.Error(err)
					}
					if historyPipelines != Pipelines_default {
						fieldOption, err := c.FieldOptionUsecase.GetByFieldName(Kind_client_cases, fieldName, body.GetString(fieldName+".value"))
						if err != nil {
							return 0, nil, err
						}
						if fieldOption == nil {
							return 0, nil, errors.New("stages: fieldOption is nil")
						}
						if fieldOption.Pipelines != historyPipelines {
							return 0, nil, errors.New("Please select the correct value of Stage")
						}
					}
				}
			}
		}

		_, _, list, err := fieldValidatorCenter.Verify(fieldName, body.GetString(fieldName+".value"), TypeDataEntry(verifyDestData))
		if err != nil {
			return 0, nil, err
		}
		verifyFailureResultList = append(verifyFailureResultList, list...)
		val := body.GetTypeMap(fieldName)
		destData.Set(fieldName, val.GetString("value"))
	}

	verifyFailureResultList = verifyFailureResultList.RemoveDuplicateResult()

	if len(verifyFailureResultList) > 0 {
		data := make(lib.TypeMap)
		data.Set("verify_failure_result", verifyFailureResultList)
		return Reply_code_data_validation_failure, data, nil
	}

	//fieldOptionStruct, err := c.FieldOptionUsecase.CacheStructByKind(kind)
	//if err != nil {
	//	return 0, nil, err
	//}

	//lib.DPrintln("destData_:", destData)
	if len(destData) > 0 {
		destGid := gid
		if !isUpdate {
			destGid = uuid.UuidWithoutStrike()
			if kind == Kind_users {
				if destData.GetString(User_FieldName_profile_gid) == "" {
					destData.Set(User_FieldName_profile_gid, userFacade.CustomFields.TextValueByNameBasic(User_FieldName_profile_gid))
				}
			}
		}

		if len(destData) > 0 {
			recognizeFieldName := DataEntry_gid
			if IsCustomKind(kind) {
				recognizeFieldName = DataEntry_Incr_id_name
				if isUpdate {
					destData.Set(recognizeFieldName, gid)
				} else {
					if kind == Kind_Custom_Condition {
						destData.Set(ConditionFieldName_type, Condition_Type_Primary)
					} else if kind == Kind_Custom_ConditionSecondary {
						destData.Set(ConditionFieldName_type, Condition_Type_SecondaryCondition)
					} else if kind == Kind_Custom_Filter {
						destData.Set(Filter_FieldName_user_gid, userFacade.Gid())
						belongModuleName, _ := ctx.GetQuery("belong_module_name")
						destData.Set(Filter_FieldName_kind, ModuleConvertToKind(belongModuleName))
					}
				}

			} else {
				destData.Set("gid", destGid)
			}

			if kind == Kind_client_cases {
				if val, ok := destData[FieldName_statements]; ok {
					delete(destData, FieldName_statements)
					// [{"id":"new_1752662054701","rating":"3","condition":"2222","association":"increase","category":"General"},{"id":"new_1752662048782","rating":"2","condition":"22","association":"increase","category":"Supplemental"},{"id":"new_1752662043207","rating":"1","condition":"11","association":"new","category":"NO PRIVATE EXAMS"}]
					if isUpdate {
						newStatements, err := c.StatementUsecase.SaveCaseStatement(destGid, InterfaceToString(val))
						if err != nil {
							return 0, nil, err
						}
						destData[FieldName_statements] = newStatements
					}
				}
			}

			dataEntryOperResult, err := c.DataEntryUsecase.HandleOne(kind, TypeDataEntry(destData), recognizeFieldName, &userFacade.TData)
			if err != nil {
				return 0, nil, err
			}

			if !isUpdate && IsCustomKind(kind) {
				gid, _ = dataEntryOperResult.GetOne()

				if kind == Kind_Custom_Condition || kind == Kind_Custom_ConditionSecondary {
					gidInt, _ := strconv.ParseInt(gid, 10, 32)
					_, er := c.ConditionLogAiUsecase.AddConditionSource(int32(gidInt), ConditionLogAi_FromType_Manual, "", userFacade.Gid(), "")
					if er != nil {
						c.log.Error("AddConditionSource:", gidInt)
					}
					if kind == Kind_Custom_ConditionSecondary {
						belongConditionId, _ := ctx.GetQuery("belong_condition_id")
						belongConditionIdInt, _ := strconv.ParseInt(belongConditionId, 10, 32)
						er := c.ConditionRelaAiUsecase.Upsert(int32(belongConditionIdInt), int32(gidInt), "")
						if er != nil {
							c.log.Error("ConditionRelaAiUsecase Upsert:", gidInt, "belongConditionId:", belongConditionId)
						}
					}
				}
			}
			if kind == Kind_clients && isCreateACase {
				err = c.ClientCasebuzUsecase.CreateACaseByClientGid(destGid, &userFacade.TData)
				if err != nil {
					c.log.Error(err)
				}
			}
			//c.log.Info("HandleOne dataEntryOperResult:", dataEntryOperResult)
		}
		data := make(lib.TypeMap)
		responseData := make(lib.TypeMap)

		if IsCustomKind(kind) {
			if kind == Kind_Custom_Condition || kind == Kind_Custom_ConditionSecondary {
				var conditionEntity *ConditionEntity
				if kind == Kind_Custom_Condition {
					conditionEntity, err = c.ConditionUsecase.GetByCond(Eq{"id": gid})
					if err != nil {
						return 0, nil, err
					}
				} else if kind == Kind_Custom_ConditionSecondary {
					belongConditionId, _ := ctx.GetQuery("belong_condition_id")
					belongConditionIdInt, _ := strconv.ParseInt(belongConditionId, 10, 32)
					conditionEntity, err = c.ConditionUsecase.GetByCond(Eq{"id": belongConditionIdInt})
					if err != nil {
						return 0, nil, err
					}
				}
				if conditionEntity == nil {
					return 0, nil, errors.New("conditionEntity is nil")
				}
				conditionIds := []int32{conditionEntity.ID}
				AllSecondaries, _ := c.ConditionbuzUsecase.GetAllSecondariesByConditionIds(conditionIds)
				AllCategories, _ := c.ConditionbuzUsecase.GetAllCategoriesByConditionIds(conditionIds)
				conditionCategoryEntity := GetConditionCategoryById(conditionEntity.ConditionCategoryId, AllCategories)
				primaryConditionVo := conditionEntity.ToPrimaryCondition(AllSecondaries, conditionCategoryEntity)
				data.Set("primary_condition", primaryConditionVo)
			} else {
				responseData.Set("id", gid)
			}

		} else {
			var newTData *TData

			newTData, err = c.TUsecase.Data(kind, Eq{"gid": destGid})

			if err != nil {
				return 0, nil, err
			}
			if newTData == nil {
				return 0, nil, errors.New("newTData is nil")
			}

			for fieldName, _ := range body {
				val := newTData.CustomFields.ValueByName(fieldName)
				responseData.Set(fieldName, val)
				if kind == Kind_client_cases {
					if fieldName == FieldName_stages {
						dueDate := newTData.CustomFields.ValueByName(DataEntry_sys__due_date)
						responseData.Set(DataEntry_sys__due_date, dueDate)
						userGidValue := newTData.CustomFields.ValueByName(FieldName_user_gid)
						responseData.Set(FieldName_user_gid, userGidValue)
						leadCo := newTData.CustomFields.ValueByName(FieldName_lead_co)
						responseData.Set(FieldName_lead_co, leadCo)
					} else if fieldName == FieldName_statements {

						conditions, _ := c.StatementUsecase.GetCaseStatementExtend(*newTData)
						responseData.Set(FieldName_statements+"__extend", conditions)
					} else if fieldName == FieldName_ContractSource {
						responseData.Set(FieldName_stages, newTData.CustomFields.ValueByName(FieldName_stages))
					} else if fieldName == FieldName_user_gid ||
						fieldName == FieldName_primary_vs ||
						fieldName == FieldName_primary_cp ||
						fieldName == FieldName_support_cp ||
						fieldName == FieldName_lead_co {
						value := newTData.CustomFields.ValueByName(FieldName_collaborators)
						responseData.Set(FieldName_collaborators, value)
					}
				}
			}

			responseData.Set(DataEntry_gid, destGid)

			responseData.Set(DataEntry_updated_at, newTData.CustomFields.NumberValueByNameBasic(DataEntry_updated_at))
			responseData.Set("primary_name", newTData.CustomFields.TextValueByNameBasic(kindEntity.PrimaryFieldName))
			if kind == Kind_client_cases {
				fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(FieldName_amount)
				if err != nil {
					return 0, nil, err
				}
				if fieldPermissionVo.CanShow() {
					responseData.Set("amount", newTData.CustomFields.DisplayValueByNameBasic(FieldName_amount))
				}
				var detailWholeCaseExtras DetailWholeCaseExtras
				stagesLog, _ := c.RecordLogUsecase.BizCrmStagesLatest(newTData.Gid())
				if stagesLog != nil {
					a, _ := GenTFieldExtendForSysDueDate(TimestampToDate(stagesLog.StartTime))
					if a != nil {
						detailWholeCaseExtras.StageStartDate = *a
					}
					a, _ = GenTFieldExtendForSysDueDate(TimestampToDate(stagesLog.EndTime))
					if a != nil {
						detailWholeCaseExtras.StageDueDate = *a
					}
				}
				c.TUsecase.DoFormula(*kindEntity, newTData)
				itfFormulaVal := newTData.CustomFields.TextValueByNameBasic(DataEntry_sys__itf_formula)
				itfFormula, _ := GenTFieldExtendForSysItfFormula(itfFormulaVal)
				if itfFormula != nil {
					detailWholeCaseExtras.ItfFormula = *itfFormula
				}
				itfExpiration := GenItfExpirationExtend(newTData)
				if itfExpiration != nil {
					detailWholeCaseExtras.ItfExpiration = *itfExpiration
				}
				detailWholeCaseExtras.Pipelines = GetPipelinesByCase(newTData)

				responseData.Set("extras", detailWholeCaseExtras)
			}
		}

		data.Set("data", responseData)
		return 0, data, nil
	}
	//data := make(lib.TypeMap)
	//data.Set("data.val", "aaa")
	return 0, nil, nil
}
