package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
)

type AmUsecase struct {
	log                *log.Helper
	conf               *conf.Data
	CommonUsecase      *CommonUsecase
	TUsecase           *TUsecase
	MapUsecase         *MapUsecase
	TaskCreateUsecase  *TaskCreateUsecase
	BoxUsecase         *BoxUsecase
	AttorneybuzUsecase *AttorneybuzUsecase
	DataEntryUsecase   *DataEntryUsecase
	AttorneyUsecase    *AttorneyUsecase
}

func NewAmUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	BoxUsecase *BoxUsecase,
	AttorneybuzUsecase *AttorneybuzUsecase,
	DataEntryUsecase *DataEntryUsecase,
	AttorneyUsecase *AttorneyUsecase,
) *AmUsecase {
	uc := &AmUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		TUsecase:           TUsecase,
		MapUsecase:         MapUsecase,
		TaskCreateUsecase:  TaskCreateUsecase,
		BoxUsecase:         BoxUsecase,
		AttorneybuzUsecase: AttorneybuzUsecase,
		DataEntryUsecase:   DataEntryUsecase,
		AttorneyUsecase:    AttorneyUsecase,
	}
	return uc
}

func (c *AmUsecase) DoHandleAmContractPending(tCase TData) error {

	return c.TaskCreateUsecase.CreateTask(tCase.Id(), nil, Task_Dag_CreateEnvelopeAndSentFromBoxAm, 0, "", "")

	return nil
}

func (c *AmUsecase) HandleAmContractPending(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	contractSource := tCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource)
	if contractSource != ContractSource_AM {
		return errors.New("The source of the contract is not AM")
	}
	key := MapKeyAmContractPending(tCase.Id())
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err = c.DoHandleAmContractPending(*tCase)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *AmUsecase) HandleAmInformationIntake(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	contractSource := tCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource)
	if contractSource != ContractSource_AM {
		return errors.New("The source of the contract is not AM")
	}

	key := MapKeyAmInformationIntake(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err = c.DoHandleAmInformationIntake(*tCase)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}

	return nil
}

func (c *AmUsecase) DoAttorney(tCase TData) (*AttorneyEntity, error) {

	attorney, err := c.AttorneyUsecase.GetByGid(tCase.CustomFields.TextValueByNameBasic(FieldName_attorney_uniqid))
	if err != nil {
		return nil, err
	}
	if attorney == nil {
		attorney, err = c.AttorneybuzUsecase.GetAnAttorney()
		if err != nil {
			return nil, err
		}
		if attorney == nil {
			return nil, errors.New("attorney is nil")
		}
		dataEntry := make(TypeDataEntry)
		dataEntry[DataEntry_gid] = tCase.Gid()
		dataEntry[FieldName_attorney_uniqid] = attorney.Gid
		c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
	}
	return attorney, nil
}

func (c *AmUsecase) DoHandleAmInformationIntake(tCase TData) error {

	email := tCase.CustomFields.TextValueByNameBasic(FieldName_email)
	if email == "" {
		return errors.New("The email is empty")
	}

	attorney, err := c.AttorneybuzUsecase.GetAnAttorney()
	if err != nil {
		return err
	}
	if attorney == nil {
		return errors.New("attorney is nil")
	}
	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_gid] = tCase.Gid()
	dataEntry[FieldName_attorney_uniqid] = attorney.Gid
	c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)

	typeMap := make(lib.TypeMap)
	typeMap.Set("Genre", MailGenre_StartYourVADisabilityClaimRepresentation)
	typeMap.Set("Email", email)
	err = c.TaskCreateUsecase.CreateTask(tCase.Id(), typeMap, Task_Dag_BuzEmail, 0, "", "")
	if err != nil {
		return err
	}
	return nil
}
