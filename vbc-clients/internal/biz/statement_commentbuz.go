package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type StatementCommentBuzUsecase struct {
	log                       *log.Helper
	conf                      *conf.Data
	CommonUsecase             *CommonUsecase
	TUsecase                  *TUsecase
	DataComboUsecase          *DataComboUsecase
	StatementCommentUsecase   *StatementCommentUsecase
	StatementConditionUsecase *StatementConditionUsecase
	DataEntryUsecase          *DataEntryUsecase
	MapUsecase                *MapUsecase
	NotesbuzUsecase           *NotesbuzUsecase
	RecordbuzUsecase          *RecordbuzUsecase
}

func NewStatementCommentBuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	DataComboUsecase *DataComboUsecase,
	StatementCommentUsecase *StatementCommentUsecase,
	StatementConditionUsecase *StatementConditionUsecase,
	DataEntryUsecase *DataEntryUsecase,
	MapUsecase *MapUsecase,
	NotesbuzUsecase *NotesbuzUsecase,
	RecordbuzUsecase *RecordbuzUsecase,
) *StatementCommentBuzUsecase {
	uc := &StatementCommentBuzUsecase{
		log:                       log.NewHelper(logger),
		CommonUsecase:             CommonUsecase,
		conf:                      conf,
		TUsecase:                  TUsecase,
		DataComboUsecase:          DataComboUsecase,
		StatementCommentUsecase:   StatementCommentUsecase,
		StatementConditionUsecase: StatementConditionUsecase,
		DataEntryUsecase:          DataEntryUsecase,
		MapUsecase:                MapUsecase,
		NotesbuzUsecase:           NotesbuzUsecase,
		RecordbuzUsecase:          RecordbuzUsecase,
	}

	return uc
}

type BizStatementCommentSaveVo struct {
	StatementConditionId int32  `json:"statement_condition_id"`
	StatementSection     string `json:"statement_section"`
	Text                 string `json:"text"`
}

func (c *StatementCommentBuzUsecase) BizStatementCommentSave(usePasswordAccess bool, tUser *TData, caseGid string, commentId int32, rawData []byte) (lib.TypeMap, error) {

	if caseGid == "" {
		return nil, errors.New("Parameters Incorrect")
	}
	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("BizStatementSave: tClient is nil")
	}

	var bizStatementCommentSaveVo BizStatementCommentSaveVo
	bizStatementCommentSaveVo = lib.BytesToTDef(rawData, bizStatementCommentSaveVo)
	bizStatementCommentSaveVo.Text = strings.TrimSpace(bizStatementCommentSaveVo.Text)
	if len(bizStatementCommentSaveVo.Text) <= 0 {
		return nil, errors.New("Please enter the content.")
	}
	if bizStatementCommentSaveVo.StatementConditionId <= 0 {
		//return nil, errors.New("The parameters are incorrect.")
	}
	if len(bizStatementCommentSaveVo.StatementSection) == 0 {
		//return nil, errors.New("The parameters are incorrect.")
	}

	var commentEntity *StatementCommentEntity
	if commentId > 0 {
		commentEntity, _ = c.StatementCommentUsecase.GetByCond(Eq{"id": commentId,
			"case_id":                tCase.Id(),
			"statement_condition_id": bizStatementCommentSaveVo.StatementConditionId,
			"statement_section":      bizStatementCommentSaveVo.StatementSection,
		})
		if commentEntity == nil {
			return nil, errors.New("The comment does not exist.")
		}
		if usePasswordAccess && commentEntity.ModifiedBy != "" {
			return nil, errors.New("Cannot be modified")
		}
		if !usePasswordAccess && commentEntity.ModifiedBy == "" {
			return nil, errors.New("Cannot be modified")
		}

	} else {
		commentEntity = &StatementCommentEntity{
			CaseId:               tCase.Id(),
			StatementConditionId: bizStatementCommentSaveVo.StatementConditionId,
			StatementSection:     bizStatementCommentSaveVo.StatementSection,
			CreatedAt:            time.Now().Unix(),
		}
	}
	commentEntity.Text = bizStatementCommentSaveVo.Text
	commentEntity.UpdatedAt = time.Now().Unix()
	if tUser != nil {
		commentEntity.ModifiedBy = tUser.Gid()
	}
	err = c.CommonUsecase.DB().Save(&commentEntity).Error
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	statementConditions := make(map[int32]*StatementConditionEntity)
	res, err := c.StatementConditionUsecase.AllByCond(In("id", []int32{commentEntity.StatementConditionId}))
	if err != nil {
		c.log.Error(err)
	}
	for k, v := range res {
		statementConditions[v.ID] = res[k]
	}

	data.Set("data", commentEntity.ToStatementCommentVo(statementConditions))
	return data, nil
}

func (c *StatementCommentBuzUsecase) BizStatementCommentDelete(usePasswordAccess bool, tUser *TData, caseGid string, commentId int32) (lib.TypeMap, error) {

	if caseGid == "" {
		return nil, errors.New("Parameters Incorrect")
	}
	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("BizStatementSave: tClient is nil")
	}

	commentEntity, _ := c.StatementCommentUsecase.GetByCond(Eq{"id": commentId,
		"case_id": tCase.Id(),
	})
	if commentEntity == nil {
		return nil, errors.New("The record does not exist.")
	}
	if usePasswordAccess && commentEntity.ModifiedBy != "" {
		return nil, errors.New("Deletion is not allowed.")
	}
	if !usePasswordAccess && commentEntity.ModifiedBy == "" {
		return nil, errors.New("Deletion is not allowed.")
	}

	commentEntity.DeletedAt = time.Now().Unix()
	err = c.CommonUsecase.DB().Save(&commentEntity).Error
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)

	return data, nil
}

func (c *StatementCommentBuzUsecase) HasSubmitToReview(caseId int32) (hasSubmitToReview bool, err error) {

	key := MapKeyPersonalStatementSubmitForReview(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return false, err
	}
	if val == "1" {
		return true, nil
	}
	return false, nil
}

func (c *StatementCommentBuzUsecase) BizStatementCommentList(usePasswordAccess bool, tUser *TData, caseGid string, rawData []byte) (lib.TypeMap, error) {

	if caseGid == "" {
		return nil, errors.New("Parameters Incorrect")
	}
	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("BizStatementSave: tClient is nil")
	}

	comments, err := c.StatementCommentUsecase.AllByCondWithOrderBy(Eq{"case_id": tCase.Id(), "deleted_at": 0}, "id desc", 2000)
	if err != nil {
		return nil, err
	}

	var listStatementCommentVo ListStatementCommentVo

	var conditionIds []int32
	for _, v := range comments {
		conditionIds = append(conditionIds, v.StatementConditionId)
	}
	statementConditions := make(map[int32]*StatementConditionEntity)

	if len(conditionIds) > 0 {
		res, err := c.StatementConditionUsecase.AllByCond(In("id", conditionIds))
		if err != nil {
			c.log.Error(err)
		}
		for k, v := range res {
			statementConditions[v.ID] = res[k]
		}
	}

	for _, v := range comments {
		listStatementCommentVo = append(listStatementCommentVo, v.ToStatementCommentVo(statementConditions))
	}
	HasSubmitToReview, err := c.HasSubmitToReview(tCase.Id())
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)
	data.Set("records", listStatementCommentVo)
	data.Set("has_submit_to_review", HasSubmitToReview)
	return data, nil
}

func (c *StatementCommentBuzUsecase) BizStatementCommentSubmitForReview(caseGid string) (lib.TypeMap, error) {
	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("The parameter is incorrect.")
	}
	key := MapKeyPersonalStatementSubmitForReview(tCase.Id())
	err = c.MapUsecase.Set(key, "1")
	if err != nil {
		return nil, err
	}

	err = c.NotesbuzUsecase.HandlePWNotification(caseGid)
	if err != nil {
		c.log.Error(err)
	}
	data := make(lib.TypeMap)
	data.Set("has_submit_to_review", true)
	return data, nil
}

func (c *StatementCommentBuzUsecase) BizStatementCommentMarkComplete(userFacade UserFacade, caseGid string, commentId int32, action string) (lib.TypeMap, error) {

	tCase, err := c.RecordbuzUsecase.VerifyDataPermission(caseGid, userFacade.TData)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("The data does not exist or there is no permission")
	}

	commentEntity, err := c.StatementCommentUsecase.GetByCond(Eq{"id": commentId, "case_id": tCase.Id(), "deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if commentEntity == nil {
		return nil, errors.New("commentEntity is nil")
	}

	isComplete := StatementComment_IsComplete_Yes
	if action == StatementComment_Action_Unmark {
		isComplete = StatementComment_IsComplete_No
	}
	commentEntity.IsComplete = isComplete

	err = c.StatementCommentUsecase.UpdatesByCond(map[string]interface{}{
		"is_complete": isComplete,
	}, Eq{"id": commentId})
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)

	statementConditions := make(map[int32]*StatementConditionEntity)
	res, err := c.StatementConditionUsecase.AllByCond(In("id", []int32{commentEntity.StatementConditionId}))
	if err != nil {
		c.log.Error(err)
	}
	for k, v := range res {
		statementConditions[v.ID] = res[k]
	}

	data.Set("data", commentEntity.ToStatementCommentVo(statementConditions))

	return data, nil
}
