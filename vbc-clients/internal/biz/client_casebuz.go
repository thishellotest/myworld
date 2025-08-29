package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib/uuid"
)

// needClientSyncCaseFields 需要互相同步的字段
var needClientSyncCaseFields = []string{"email", "phone", "ssn", "dob", "state", "city", "address", "zip_code",
	"place_of_birth_city",
	"place_of_birth_state_province",
	"place_of_birth_country", FieldName_current_occupation}

// needClientSyncCaseFieldsForOnce 需要在初始化同步一次的字段 key：client value: case
var needClientSyncCaseFieldsForOnce = map[string]string{
	"current_rating":           "current_rating",
	"effective_current_rating": "effective_current_rating",
	"retired":                  "retired",
	"active_duty":              "active_duty",
	"branch":                   "branch",
	"source":                   "source",
	"referrer":                 "referrer",
	"collaborators":            "collaborators",
}

type ClientCasebuzUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	TUsecase         *TUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewClientCasebuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
) *ClientCasebuzUsecase {
	uc := &ClientCasebuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

func (c *ClientCasebuzUsecase) CreateACaseByClientGid(clientGid string, operUser *TData) error {

	tClient, _ := c.TUsecase.DataByGid(Kind_clients, clientGid)
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	mapping := needClientSyncCaseFieldsForOnce
	for _, v := range needClientSyncCaseFields {
		mapping[v] = v
	}
	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_gid] = uuid.UuidWithoutStrike()
	if operUser != nil {
		dataEntry[FieldName_user_gid] = operUser.Gid()
	}
	dataEntry[FieldName_client_gid] = clientGid
	dataEntry[FieldName_deal_name] = tClient.CustomFields.TextValueByNameBasic(FieldName_full_name)

	dataEntry[FieldName_ContractSource] = NewCaseDefaultContractSource
	if NewCaseDefaultContractSource == ContractSource_VBC {
		dataEntry[FieldName_stages] = config_vbc.Stages_IncomingRequest
	} else {
		dataEntry[FieldName_stages] = config_vbc.Stages_AmIncomingRequest
	}

	clientMaps := tClient.CustomFields.ToMaps()

	for k, v := range mapping {
		if _, ok := clientMaps[k]; ok {
			val := InterfaceToString(clientMaps[k])
			if val != "" {
				dataEntry[v] = val
			}
		}
	}

	_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, operUser)

	return err
}
