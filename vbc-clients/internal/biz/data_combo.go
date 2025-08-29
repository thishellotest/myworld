package biz

import (
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type DataComboUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	TUsecase      *TUsecase
}

func NewDataComboUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase) *DataComboUsecase {
	uc := &DataComboUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
	}

	return uc
}

func (c *DataComboUsecase) Client(gid string) (*TData, TFields, error) {
	if gid == "" {
		return nil, nil, errors.New("Client gid is empty.")
	}
	tClient, err := c.TUsecase.Data(Kind_clients, Eq{"gid": gid})
	if err != nil {
		return nil, nil, err
	}
	if tClient != nil {

		return tClient, tClient.CustomFields, nil
	}
	return nil, nil, err
}

func (c *DataComboUsecase) ClientWithCase(tCase TData) (*TData, TFields, error) {
	return c.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
}
