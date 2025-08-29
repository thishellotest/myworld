package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

func (c *StageTransUsecase) ClientCasesMappings(row lib.TypeMap) (typeMap lib.TypeMap, err error) {
	if row == nil {
		return nil, nil
	}
	res := make(lib.TypeMap)
	for k, v := range config_zoho.ClientCasesMappingConfigs {

		if k == "Stage" {
			zohoStage := row.GetString(k)
			dbStage, err := c.BizZohoStageToDBStage(zohoStage)
			if err != nil {
				c.log.Error(err)
			}
			res.Set(v, dbStage)
		} else {
			res.Set(v, row.GetString(k))
		}
	}
	return res, nil
}

func (c *StageTransUsecase) ClientCasesMappings2(row lib.TypeMap) (typeMap lib.TypeMap, err error) {
	if row == nil {
		return nil, nil
	}
	res := make(lib.TypeMap)
	for k, v := range config_zoho.ClientCasesMappingConfigs2 {

		if k == "Stage" {
			zohoStage := row.GetString(k)
			dbStage, err := c.BizZohoStageToDBStage(zohoStage)
			if err != nil {
				c.log.Error(err)
			}
			res.Set(v, dbStage)
		} else {
			res.Set(v, row.GetString(k))
		}
	}
	return res, nil
}

func DefaultStageTrans(stage string) string {
	// 26. Terminated
	aa := strings.Split(stage, ".")
	if len(aa) > 1 {
		r := ""
		for i := 1; i < len(aa); i++ {
			r += strings.TrimSpace(aa[i])
		}
		return r
	} else {
		return stage
	}
}

type StageTransUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	FieldOptionUsecase *FieldOptionUsecase
	FieldUsecase       *FieldUsecase
}

func NewStageTransUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldOptionUsecase *FieldOptionUsecase,
	FieldUsecase *FieldUsecase) *StageTransUsecase {
	uc := &StageTransUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		FieldOptionUsecase: FieldOptionUsecase,
		FieldUsecase:       FieldUsecase,
	}

	return uc
}

func (c *StageTransUsecase) BizZohoStageToDBStage(zohoStage string) (dbStage string, err error) {
	dbStage, err = c.ZohoStageToDBStage(zohoStage)
	if err != nil {
		return "", err
	}
	if dbStage == "" {
		dbStage = DefaultStageTrans(zohoStage)
	}
	return dbStage, nil
}
func (c *StageTransUsecase) ZohoStageToDBStage(zohoStage string) (dbStage string, err error) {

	fieldStructs, err := c.FieldOptionUsecase.CacheStructByKind(Kind_client_cases)
	if err != nil {
		return "", err
	}
	if fieldStructs == nil {
		return "", errors.New("fieldStructs is nil")
	}
	fieldStruct, err := c.FieldUsecase.CacheStructByKind(Kind_client_cases)
	if err != nil {
		return "", err
	}
	if fieldStruct == nil {
		return "", errors.New("ZohoStageToDBStage: fieldStruct is nil")
	}
	fieldEntity := fieldStruct.GetByFieldName(FieldName_stages)
	if fieldEntity == nil {
		return "", errors.New("ZohoStageToDBStage: fieldEntity is nil")
	}
	optionList := fieldStructs.AllByFieldName(*fieldEntity)
	option := optionList.GetByLabel(zohoStage)
	if option == nil {
		return "", nil
	} else {
		return option.OptionValue, nil
	}
}

func (c *StageTransUsecase) DBStageToZohoStage(dbStage string) (zohoStage string, err error) {
	fieldStructs, err := c.FieldOptionUsecase.CacheStructByKind(Kind_client_cases)
	if err != nil {
		return "", err
	}
	if fieldStructs == nil {
		return "", errors.New("fieldStructs is nil")
	}

	fieldStruct, err := c.FieldUsecase.CacheStructByKind(Kind_client_cases)
	if err != nil {
		return "", err
	}
	if fieldStruct == nil {
		return "", errors.New("ZohoStageToDBStage: fieldStruct is nil")
	}
	fieldEntity := fieldStruct.GetByFieldName(FieldName_stages)
	if fieldEntity == nil {
		return "", errors.New("ZohoStageToDBStage: fieldEntity is nil")
	}

	optionList := fieldStructs.AllByFieldName(*fieldEntity)
	option := optionList.GetByValue(dbStage)
	if option == nil {
		return "", errors.New("zohoStage or dbStage does not mapping")
	} else {
		return option.OptionLabel, nil
	}
}
