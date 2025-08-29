package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	. "vbc/lib/builder"
)

type ClientCaseSyncbuzUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	TUsecase         *TUsecase
	DataEntryUsecase *DataEntryUsecase
	StatementUsecase *StatementUsecase
}

func NewClientCaseSyncbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
	StatementUsecase *StatementUsecase,
) *ClientCaseSyncbuzUsecase {
	uc := &ClientCaseSyncbuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		DataEntryUsecase: DataEntryUsecase,
		StatementUsecase: StatementUsecase,
	}
	return uc
}

type ClientCaseSyncVo struct {
	FieldName  string
	FieldValue string
}

func (c *ClientCaseSyncbuzUsecase) ClientToCases(clientGid string, clientCaseSyncVo ClientCaseSyncVo, operUser *TData) error {

	cases, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{FieldName_biz_deleted_at: 0, FieldName_client_gid: clientGid})
	if err != nil {
		return err
	}

	caseFieldName := config_vbc.GetSyncFieldNameByClientForCase(clientCaseSyncVo.FieldName)
	if caseFieldName == "" {
		return nil
	}
	for _, v := range cases {
		caseValue := v.CustomFields.TextValueByNameBasic(caseFieldName)
		if caseValue != clientCaseSyncVo.FieldValue {
			dataEntry := make(TypeDataEntry)
			dataEntry[DataEntry_gid] = v.Gid()
			dataEntry[caseFieldName] = clientCaseSyncVo.FieldValue
			_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, operUser)
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}

func (c *ClientCaseSyncbuzUsecase) CaseToClient(caseId int32, clientCaseSyncVo ClientCaseSyncVo, operUser *TData) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New(InterfaceToString(caseId) + ":tCase is nil")
	}
	clientGid := tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid)
	if clientGid == "" {
		return nil
	}
	clientFieldName := config_vbc.GetSyncFieldNameByCaseForClient(clientCaseSyncVo.FieldName)
	if clientFieldName == "" {
		return nil
	}

	client, err := c.TUsecase.DataByGid(Kind_clients, clientGid)
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("client is nil")
	}
	clientValue := client.CustomFields.TextValueByNameBasic(clientFieldName)
	if clientValue == clientCaseSyncVo.FieldValue {
		return nil
	}

	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_gid] = clientGid
	dataEntry[clientFieldName] = clientCaseSyncVo.FieldValue
	_, err = c.DataEntryUsecase.HandleOne(Kind_clients, dataEntry, DataEntry_gid, operUser)

	return err
}

func PersonalStatementManagerUrl(caseGid string) (url string) {

	if configs.IsDev() {
		url = fmt.Sprintf("http://localhost:3000/ps/%s", caseGid)
	} else {
		url = fmt.Sprintf("%s/ps/%s", configs.Domain, caseGid)
	}
	return url
}

func (c *ClientCaseSyncbuzUsecase) UpdatePersonalStatementManagerUrl(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	val := tCase.CustomFields.TextValueByNameBasic(FieldName_personal_statement_manager)
	if val == "" {
		url := PersonalStatementManagerUrl(tCase.Gid())
		data := make(TypeDataEntry)
		data[DataEntry_gid] = tCase.Gid()
		data[FieldName_personal_statement_manager] = url
		password, err := c.StatementUsecase.PersonalStatementPassword(tCase.Id())
		if err != nil {
			c.log.Error(err, tCase.Id())
		}

		data[FieldName_personal_statement_password] = password
		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
		return err
	}

	return nil
}
