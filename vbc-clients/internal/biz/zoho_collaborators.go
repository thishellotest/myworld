package biz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

type ZohoCollaboratorUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	ZohoUsecase      *ZohoUsecase
	DataEntryUsecase *DataEntryUsecase
	TUsecase         *TUsecase
}

func NewZohoCollaboratorUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	ZohoUsecase *ZohoUsecase,
	DataEntryUsecase *DataEntryUsecase,
	TUsecase *TUsecase,
) *ZohoCollaboratorUsecase {
	uc := &ZohoCollaboratorUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		ZohoUsecase:      ZohoUsecase,
		DataEntryUsecase: DataEntryUsecase,
		TUsecase:         TUsecase,
	}

	return uc
}

func (c *ZohoCollaboratorUsecase) BizHandleClientCases(gid string) error {

	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	dbCollaborators := tCase.CustomFields.TextValueByNameBasic(FieldName_collaborators)
	sourceDbCollaborators := dbCollaborators
	if gid != "" {
		collaborators, err := c.GetCollaboratorsFromZoho(config_zoho.Deals, gid)
		if err != nil {
			c.log.Error(err)
			return err
		}
		var zohoCollaborators []string
		for _, v := range collaborators {
			zohoCollaborators = append(zohoCollaborators, v.GetString("Collaborators.id"))
		}
		zohoCollaborators = lib.ArrayReverse(zohoCollaborators)

		for _, v := range zohoCollaborators {
			if strings.Index(dbCollaborators, v) < 0 {
				if dbCollaborators == "" {
					dbCollaborators = fmt.Sprintf(",%s,", v)
				} else {
					dbCollaborators += fmt.Sprintf("%s,", v)
				}
			}
		}
		if sourceDbCollaborators != dbCollaborators {
			destData := make(TypeDataEntry)
			destData[DataEntry_gid] = gid
			destData[FieldName_collaborators] = dbCollaborators
			_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, destData, DataEntry_gid, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ZohoCollaboratorUsecase) HandleClientCases(ctx context.Context) error {
	c.log.Info("HandleClientCases")
	sqlRows, err := c.WaitingClientCases(ctx)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	for sqlRows.Next() {
		_, row, err := lib.SqlRowsToMap(sqlRows)
		if err != nil {
			return err
		}
		err = c.BizHandleClientCases(row.GetString("gid"))
		if err != nil {
			c.log.Error(err)
		}
	}
	return nil
}

func (c *ZohoCollaboratorUsecase) BizHandleClients(gid string) error {

	tClient, err := c.TUsecase.DataByGid(Kind_clients, gid)
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	dbCollaborators := tClient.CustomFields.TextValueByNameBasic(FieldName_collaborators)
	sourceDbCollaborators := dbCollaborators

	if gid != "" {
		collaborators, err := c.GetCollaboratorsFromZoho(config_zoho.Contacts, gid)
		if err != nil {
			c.log.Error(err)
			return err
		}
		var zohoCollaborators []string
		for _, v := range collaborators {
			zohoCollaborators = append(zohoCollaborators, v.GetString("Collaborators.id"))
		}
		zohoCollaborators = lib.ArrayReverse(zohoCollaborators)

		for _, v := range zohoCollaborators {
			if strings.Index(dbCollaborators, v) < 0 {
				if dbCollaborators == "" {
					dbCollaborators = fmt.Sprintf(",%s,", v)
				} else {
					dbCollaborators += fmt.Sprintf("%s,", v)
				}
			}
		}
		if sourceDbCollaborators != dbCollaborators {
			destData := make(TypeDataEntry)
			destData[DataEntry_gid] = gid
			destData[FieldName_collaborators] = dbCollaborators
			_, err = c.DataEntryUsecase.HandleOne(Kind_clients, destData, DataEntry_gid, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ZohoCollaboratorUsecase) HandleClients(ctx context.Context) error {

	c.log.Info("HandleClients")
	sqlRows, err := c.WaitingClients(ctx)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	for sqlRows.Next() {
		_, row, err := lib.SqlRowsToMap(sqlRows)
		if err != nil {
			return err
		}
		gid := row.GetString("gid")
		err = c.BizHandleClients(gid)
		if err != nil {
			c.log.Error(err)
		}
	}
	return nil
}

// GetCollaboratorsFromZoho collaborators: [{"Collaborators":{"id":"6159272000005147005","name":"Andrea Ladd"},"id":"6159272000005174075"},{"Collaborators":{"id":"6159272000001027094","name":"Lili Wang"},"id":"6159272000005174074"},{"Collaborators":{"id":"6159272000000453640","name":"Edward Bunting"},"id":"6159272000005174073"},{"Collaborators":{"id":"6159272000001027142","name":"Donald Pratko"},"id":"6159272000001381123"},{"Collaborators":{"id":"6159272000001027129","name":"Victoria Enriquez"},"id":"6159272000001381122"}],"id":"6159272000001017001"}]
func (c *ZohoCollaboratorUsecase) GetCollaboratorsFromZoho(zohoModuleName string, gid string) (collaborators lib.TypeList, err error) {
	libMap, err := c.ZohoUsecase.GetASpecificRecord(zohoModuleName, []string{config_zoho.Zoho_Collaborators}, gid)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	data := libMap.GetTypeList("data")
	for _, v := range data {
		return v.GetTypeList(config_zoho.Zoho_Collaborators), nil
	}
	return nil, nil
}

func (c *ZohoCollaboratorUsecase) WaitingClientCases(ctx context.Context) (*sql.Rows, error) {

	sql := `select gid from client_cases
  where  client_cases.deleted_at=0 and client_cases.biz_deleted_at=0  order by id desc`

	return c.CommonUsecase.DB().Raw(sql).Rows()
}

func (c *ZohoCollaboratorUsecase) WaitingClients(ctx context.Context) (*sql.Rows, error) {

	sql := `select gid from clients
  where  clients.deleted_at=0 and clients.biz_deleted_at=0 order by id desc`
	return c.CommonUsecase.DB().Raw(sql).Rows()
}
