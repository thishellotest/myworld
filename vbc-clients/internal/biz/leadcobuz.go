package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
)

type LeadcobuzUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	TUsecase         *TUsecase
	UserUsecase      *UserUsecase
	DataComboUsecase *DataComboUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewLeadcobuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	UserUsecase *UserUsecase,
	DataComboUsecase *DataComboUsecase,
	DataEntryUsecase *DataEntryUsecase,
) *LeadcobuzUsecase {
	uc := &LeadcobuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		UserUsecase:      UserUsecase,
		DataComboUsecase: DataComboUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}
	return uc
}

func (c *LeadcobuzUsecase) HandleLeadCOSyncClient(tCase TData) error {

	leadCO := tCase.CustomFields.TextValueByNameBasic(FieldName_lead_co)
	if leadCO == "" {
		return nil
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if !IsForLeadCOStages(stages) {
		return nil
	}

	tUser, err := c.UserUsecase.GetByFullName(leadCO)
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
