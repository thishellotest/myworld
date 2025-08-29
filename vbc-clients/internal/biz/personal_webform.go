package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	. "vbc/lib/builder"
)

type PersonalWebformUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	TUsecase         *TUsecase
	MapUsecase       *MapUsecase
	FeeUsecase       *FeeUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewPersonalWebformUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase,
	FeeUsecase *FeeUsecase,
	DataEntryUsecase *DataEntryUsecase,
) *PersonalWebformUsecase {
	uc := &PersonalWebformUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		MapUsecase:       MapUsecase,
		FeeUsecase:       FeeUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

func (c *PersonalWebformUsecase) IsUseNewPersonalWebForm(caseId int32) (bool, error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return false, err
	}
	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	if tCase.CustomFields.TextValueByNameBasic(FieldName_personal_statement_type) == Personal_statement_type_Webform {
		return true, nil
	}
	return false, nil

	//key := MapKeyPersonalStatementUsePW(caseId)
	//val, err := c.MapUsecase.GetForString(key)
	//if err != nil {
	//	return false, err
	//}
	//if val == "1" {
	//	return true, nil
	//}
	//return false, nil
}

func (c *PersonalWebformUsecase) HandleUseNewPersonalWebForm(caseId int32) error {

	return nil
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	personalStatementType := tCase.CustomFields.TextValueByNameBasic(FieldName_personal_statement_type)

	if personalStatementType != "" {
		return nil
	}
	needUseNewPersonalWebForm, err := c.NeedUseNewPersonalWebForm(caseId)
	if err != nil {
		return err
	}
	if needUseNewPersonalWebForm {
		//key := MapKeyPersonalStatementUsePW(caseId)
		//err := c.MapUsecase.Set(key, "1")
		//if err != nil {
		//	return err
		//}

		data := make(TypeDataEntry)
		data[DataEntry_gid] = tCase.Gid()
		data[FieldName_personal_statement_type] = Personal_statement_type_Webform
		c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
	}
	return nil
}

// NeedUseNewPersonalWebForm 判断是否可以使用新的PW系统
func (c *PersonalWebformUsecase) NeedUseNewPersonalWebForm(caseId int32) (bool, error) {
	return true, nil
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return false, err
	}
	if tCase == nil {
		return false, errors.New("The Case does not exist.")
	}
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		return false, err
	}
	currentTime := time.Now()
	if isPrimaryCase {
		if currentTime.Unix() > 1753931763 {
			return true, nil
		}
	} else {
		if primaryCase == nil {
			return false, errors.New("The Primary Case does not exist.")
		}
		useNewPW, err := c.IsUseNewPersonalWebForm(primaryCase.Id())
		if err != nil {
			return false, err
		}
		if useNewPW {
			return true, nil
		}
	}
	return false, nil
}

func (c *PersonalWebformUsecase) ManualHistoryData() error {

	records, err := c.TUsecase.ListByCond(Kind_client_cases, And(In(FieldName_stages,
		config_vbc.Stages_StatementsFinalized,
		config_vbc.Stages_CurrentTreatment,
		config_vbc.Stages_CurrentTreatmentReview,
		config_vbc.Stages_StatementUpdates,
		config_vbc.Stages_PreparingDocumentsTinnitusLetter,
		config_vbc.Stages_AwaitingNexusLetter,
		config_vbc.Stages_MiniDBQs_Draft,
		config_vbc.Stages_MiniDBQs,
		config_vbc.Stages_MiniDBQ_Forms,
		config_vbc.Stages_MedicalTeamFormsSigned,
		config_vbc.Stages_MedicalTeam,
		config_vbc.Stages_MedicalTeamPaymentCollected,
		config_vbc.Stages_MedicalTeamExamsScheduled,
		config_vbc.Stages_MedicalTeamCallVet,
		config_vbc.Stages_MedicalTeamPrefilledFormsReview,
		config_vbc.Stages_DBQ_Completed,
		config_vbc.Stages_StatementFinalChanges,
		config_vbc.Stages_FileClaims,
		config_vbc.Stages_FileClaims_Draft,
		config_vbc.Stages_FileHLRDraft,
		config_vbc.Stages_FileHLRWithClient,
		config_vbc.Stages_VerifyEvidenceReceived,
		config_vbc.Stages_AwaitingDecision,
		config_vbc.Stages_AwaitingPayment,
		config_vbc.Stages_27_AwaitingBankReconciliation,
		config_vbc.Stages_Completed,
		config_vbc.Stages_AmAwaitingPayment,
	)))
	if err != nil {
		return err
	}
	for _, v := range records {
		data := make(TypeDataEntry)
		data[DataEntry_gid] = v.Gid()
		data[FieldName_personal_statement_type] = Personal_statement_type_WorddocumentWithinBox
		c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
	}

	return nil
}
