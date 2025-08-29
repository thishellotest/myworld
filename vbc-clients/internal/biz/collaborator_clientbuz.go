package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
)

type CollaboratorClientbuzUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	UserUsecase      *UserUsecase
	TUsecase         *TUsecase
	DataEntryUsecase *DataEntryUsecase
	DataComboUsecase *DataComboUsecase
}

func NewCollaboratorClientbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	UserUsecase *UserUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
	DataComboUsecase *DataComboUsecase,
) *CollaboratorClientbuzUsecase {
	uc := &CollaboratorClientbuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		UserUsecase:      UserUsecase,
		TUsecase:         TUsecase,
		DataEntryUsecase: DataEntryUsecase,
		DataComboUsecase: DataComboUsecase,
	}

	return uc
}

func (c *CollaboratorClientbuzUsecase) DoCollaboratorByChangeHistory(entity ChangeHistoryEntity, tCase TData) error {

	if entity.Kind == Kind_client_cases && (entity.FieldName == FieldName_primary_vs ||
		entity.FieldName == FieldName_lead_co ||
		entity.FieldName == FieldName_stages) {
		if entity.FieldName == FieldName_primary_vs {
			stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
			if IsForLeadCOStages(stages) {
				return c.OperationCollaboratorByFullName(tCase, entity.OldValue)
			}
			return nil
		} else if entity.FieldName == FieldName_stages {
			stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
			if !IsForLeadCOStages(stages) {
				leadCo := tCase.CustomFields.TextValueByNameBasic(FieldName_lead_co)
				if leadCo != "" {
					dataEntry := make(TypeDataEntry)
					dataEntry[DataEntry_gid] = tCase.Gid()
					dataEntry[FieldName_lead_co] = ""
					_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
					if err != nil {
						c.log.Error(err, " Gid: ", tCase.Gid())
					}
				}
			}
		}
	}
	return nil
}

func (c *CollaboratorClientbuzUsecase) OperationCollaboratorByFullName(tCase TData, removeFullName string) error {

	c.log.Debug("OperationCollaboratorByFullName:", removeFullName)
	var removeUserGid string
	if removeFullName != "" {
		old, err := c.UserUsecase.GetByFullName(removeFullName)
		if err != nil {
			c.log.Error(err, " caseId: ", tCase.Id(), " removeFullName:", removeFullName)
		}
		if old != nil {
			removeUserGid = old.Gid()
		}
	}
	tClient, _, err := c.DataComboUsecase.ClientWithCase(tCase)
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	return c.OperationCollaborator(tCase, *tClient, removeUserGid)
}
func (c *CollaboratorClientbuzUsecase) OperationCollaborator(tCase TData, tClient TData, removeUserGid string) error {

	userGids, err := c.GetRequiredCollaborators(tCase)
	if err != nil {
		return err
	}
	dbCollaborators := tClient.CustomFields.TextValueByNameBasic(FieldName_collaborators)
	if removeUserGid != "" {
		dbCollaborators = strings.ReplaceAll(dbCollaborators, removeUserGid+",", "")
		if dbCollaborators == "," {
			dbCollaborators = ""
		}
	}
	for _, v := range userGids {
		if !strings.Contains(dbCollaborators, ","+v+",") {
			if dbCollaborators == "" {
				dbCollaborators = "," + v + ","
			} else {
				dbCollaborators += v + ","
			}
		}
	}
	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_gid] = tClient.Gid()
	dataEntry[FieldName_collaborators] = dbCollaborators
	c.DataEntryUsecase.HandleOne(Kind_clients, dataEntry, DataEntry_gid, nil)
	return nil
}

func (c *CollaboratorClientbuzUsecase) GetRequiredCollaborators(tCase TData) ([]string, error) {

	var userGids []string
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if IsForLeadCOStages(stages) {
		leadVS, _ := c.UserUsecase.GetUserByLeadVS(&tCase)
		if leadVS != nil {
			if !lib.InArray(leadVS.Gid(), userGids) {
				userGids = append(userGids, leadVS.Gid())
			}
		}
	}
	return userGids, nil
}
