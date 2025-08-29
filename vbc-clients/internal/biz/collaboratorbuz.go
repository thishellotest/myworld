package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
)

type CollaboratorbuzUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	UserUsecase      *UserUsecase
	TUsecase         *TUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewCollaboratorbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	UserUsecase *UserUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
) *CollaboratorbuzUsecase {
	uc := &CollaboratorbuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		UserUsecase:      UserUsecase,
		TUsecase:         TUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

func (c *CollaboratorbuzUsecase) DoCollaboratorByChangeHistory(entity ChangeHistoryEntity) error {

	if entity.Kind == Kind_client_cases && (entity.FieldName == FieldName_user_gid ||
		entity.FieldName == FieldName_primary_vs ||
		entity.FieldName == FieldName_primary_cp) {

		tCase, err := c.TUsecase.DataById(Kind_client_cases, entity.IncrId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		if entity.FieldName == FieldName_user_gid {
			return c.OperationCollaborator(*tCase, entity.OldValue)
		} else if entity.FieldName == FieldName_primary_vs {
			removeUserGid := ""
			if entity.OldValue != "" {
				old, err := c.UserUsecase.GetByFullName(entity.OldValue)
				if err != nil {
					return err
				}
				if old != nil {
					removeUserGid = old.Gid()
				}
			}
			return c.OperationCollaborator(*tCase, removeUserGid)
		} else if entity.FieldName == FieldName_primary_cp {
			removeUserGid := ""
			if entity.OldValue != "" {
				old, err := c.UserUsecase.GetByFullName(entity.OldValue)
				if err != nil {
					return err
				}
				if old != nil {
					removeUserGid = old.Gid()
				}
			}
			return c.OperationCollaborator(*tCase, removeUserGid)
		} else if entity.FieldName == FieldName_lead_co {
			removeUserGid := ""
			if entity.OldValue != "" {
				old, err := c.UserUsecase.GetByFullName(entity.OldValue)
				if err != nil {
					return err
				}
				if old != nil {
					removeUserGid = old.Gid()
				}
			}
			return c.OperationCollaborator(*tCase, removeUserGid)
		}
	}
	return nil
}

func (c *CollaboratorbuzUsecase) OperationCollaboratorByFullName(tCase TData, removeFullName string) error {

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
	return c.OperationCollaborator(tCase, removeUserGid)
}
func (c *CollaboratorbuzUsecase) OperationCollaborator(tCase TData, removeUserGid string) error {

	userGids, err := c.GetRequiredCollaborators(tCase)
	if err != nil {
		return err
	}
	dbCollaborators := tCase.CustomFields.TextValueByNameBasic(FieldName_collaborators)
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
	dataEntry[DataEntry_gid] = tCase.Gid()
	dataEntry[FieldName_collaborators] = dbCollaborators
	c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
	return nil
}

func (c *CollaboratorbuzUsecase) GetRequiredCollaborators(tCase TData) ([]string, error) {

	var userGids []string
	userGid := tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid)
	if userGid != "" {
		userGids = append(userGids, userGid)
	}
	leadVS, _ := c.UserUsecase.GetUserByLeadVS(&tCase)
	if leadVS != nil {
		if !lib.InArray(leadVS.Gid(), userGids) {
			userGids = append(userGids, leadVS.Gid())
		}
	}
	leadCP, _ := c.UserUsecase.GetUserByLeadCP(&tCase)
	if leadCP != nil {
		if !lib.InArray(leadCP.Gid(), userGids) {
			userGids = append(userGids, leadCP.Gid())
		}
	}

	supportCP, _ := c.UserUsecase.GetUserBySupportCP(&tCase)
	if supportCP != nil {
		if !lib.InArray(supportCP.Gid(), userGids) {
			userGids = append(userGids, supportCP.Gid())
		}
	}
	leadCO, _ := c.UserUsecase.GetUserByLeadCO(&tCase)
	if leadCO != nil {
		if !lib.InArray(leadCO.Gid(), userGids) {
			userGids = append(userGids, leadCO.Gid())
		}
	}
	return userGids, nil
}
