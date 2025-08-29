package biz

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
)

type LeadVSChangeUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	TUsecase         *TUsecase
	MapUsecase       *MapUsecase
	UserUsecase      *UserUsecase
	DataComboUsecase *DataComboUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewLeadVSChangeUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase,
	UserUsecase *UserUsecase,
	DataComboUsecase *DataComboUsecase,
	DataEntryUsecase *DataEntryUsecase,
) *LeadVSChangeUsecase {
	uc := &LeadVSChangeUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		MapUsecase:       MapUsecase,
		UserUsecase:      UserUsecase,
		DataComboUsecase: DataComboUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

func (c *LeadVSChangeUsecase) GetLeadVSChangeLogVo(caseId int32) *LeadVSChangeLogVo {

	var leadVSChangeLogVo LeadVSChangeLogVo
	key := MapKeyLeadVSChangeLog(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		c.log.Error(err, " caseId: ", caseId)
		return nil
	}
	if val == "" {
		return nil
	}
	err = json.Unmarshal([]byte(val), &leadVSChangeLogVo)
	if err != nil {
		c.log.Error(err, " caseId: ", caseId)
		return nil
	}
	return &leadVSChangeLogVo
}

func (c *LeadVSChangeUsecase) HandleLeadVSChangeForClaimAnalysisToScheduleCall(changeHistoryEntity ChangeHistoryEntity) error {

	if changeHistoryEntity.Kind != Kind_client_cases {
		return errors.New("changeHistoryEntity.Kind is wrong")
	}
	if changeHistoryEntity.FieldName != FieldName_primary_vs {
		return errors.New("changeHistoryEntity.FieldName is wrong")
	}
	if changeHistoryEntity.OldValue == "" {
		return nil
	}
	if changeHistoryEntity.NewValue == "" {
		return nil
	}

	tCase, err := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	stage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stage == config_vbc.Stages_ClaimAnalysis ||
		stage == config_vbc.Stages_ClaimAnalysisReview ||
		stage == config_vbc.Stages_ScheduleCall ||
		stage == config_vbc.Stages_AmClaimAnalysis ||
		stage == config_vbc.Stages_AmClaimAnalysisReview ||
		stage == config_vbc.Stages_AmScheduleCall {
		key := MapKeyLeadVSChangeLog(tCase.Id())

		prevUser, _ := c.UserUsecase.GetByFullName(changeHistoryEntity.OldValue)
		if prevUser == nil {
			return errors.New("prevUser is nil")
		}
		newUser, _ := c.UserUsecase.GetByFullName(changeHistoryEntity.NewValue)
		if newUser == nil {
			return errors.New("newUser is nil")
		}
		leadVSChangeLogVo := LeadVSChangeLogVo{
			PreviousVSUserGid: prevUser.Gid(),
			NewVSUserGid:      newUser.Gid(),
		}
		err = c.MapUsecase.Set(key, InterfaceToString(leadVSChangeLogVo))
		if err != nil {
			return err
		}
	}
	return nil
}

func IsForLeadCOStages(stages string) bool {
	if stages == config_vbc.Stages_AmIncomingRequest ||
		stages == config_vbc.Stages_AmInformationIntake ||
		stages == config_vbc.Stages_AmContractPending ||
		stages == config_vbc.Stages_AmAwaitingClientRecords ||
		stages == config_vbc.Stages_AmSTRRequestPending ||
		stages == config_vbc.Stages_IncomingRequest ||
		stages == config_vbc.Stages_FeeScheduleandContract ||
		stages == config_vbc.Stages_GettingStartedEmail ||
		stages == config_vbc.Stages_AwaitingClientRecords ||
		stages == config_vbc.Stages_STRRequestPending {
		return true
	}
	return false
}

func (c *LeadVSChangeUsecase) HandleLeadVSSyncClient(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	return c.DoHandleLeadVSSyncClient(*tCase)
}

func (c *LeadVSChangeUsecase) DoHandleLeadVSSyncClient(tCase TData) error {

	leadVS := tCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs)
	if leadVS == "" {
		return nil
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if IsForLeadCOStages(stages) {
		return nil
	}

	tUser, err := c.UserUsecase.GetByFullName(leadVS)
	if err != nil {
		return err
	}
	if tUser == nil {
		return errors.New("LeadVS User is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	userGid := tClient.CustomFields.TextValueByNameBasic(FieldName_user_gid)
	if userGid != tUser.Gid() {
		dataEntry := make(TypeDataEntry)
		dataEntry[DataEntry_gid] = tClient.Gid()
		dataEntry[FieldName_user_gid] = tUser.Gid()
		_, err = c.DataEntryUsecase.HandleOne(Kind_clients, dataEntry, DataEntry_gid, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
