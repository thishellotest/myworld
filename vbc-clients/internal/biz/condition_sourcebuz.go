package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ConditionSourcebuzUsecase struct {
	log                    *log.Helper
	conf                   *conf.Data
	CommonUsecase          *CommonUsecase
	ConditionSourceUsecase *ConditionSourceUsecase
	TUsecase               *TUsecase
}

func NewConditionSourcebuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	ConditionSourceUsecase *ConditionSourceUsecase,
	TUsecase *TUsecase,
) *ConditionSourcebuzUsecase {
	uc := &ConditionSourcebuzUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		ConditionSourceUsecase: ConditionSourceUsecase,
		TUsecase:               TUsecase,
	}

	return uc
}

func (c *ConditionSourcebuzUsecase) Handle() error {
	cases, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{FieldName_biz_deleted_at: 0})
	if err != nil {
		return err
	}
	for k, _ := range cases {
		err = c.HandleField("client_cases", FieldName_description, cases[k])
		if err != nil {
			c.log.Warn(err)
		}
		err = c.HandleField("client_cases", "service_connections", cases[k])
		if err != nil {
			c.log.Warn(err)
		}
		err = c.HandleField("client_cases", "previous_denials", cases[k])
		if err != nil {
			c.log.Warn(err)
		}
		err = c.HandleField("client_cases", "claims_next_round", cases[k])
		if err != nil {
			c.log.Warn(err)
		}
		err = c.HandleField("client_cases", "claims_supplemental", cases[k])
		if err != nil {
			c.log.Warn(err)
		}
		err = c.HandleField("client_cases", "claims_online", cases[k])
		if err != nil {
			c.log.Warn(err)
		}
	}

	return nil
}

func (c *ConditionSourcebuzUsecase) HandleField(tableName string, fieldName string, tData *TData) error {
	from := fmt.Sprintf("%s:%s", tableName, fieldName)
	content := tData.CustomFields.TextValueByNameBasic(fieldName)
	return c.Upsert(content, from, tData.Id())
}

func (c *ConditionSourcebuzUsecase) Upsert(content string, from string, caseId int32) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	md5 := lib.MD5Hash(content)
	entity, err := c.ConditionSourceUsecase.GetByCond(Eq{"content_md5": md5})
	if err != nil {
		return err
	}
	if entity != nil {
		entity.CaseId = caseId
		entity.UpdatedAt = time.Now().Unix()
	} else {
		entity = &ConditionSourceEntity{
			From:       from,
			Content:    content,
			ContentMd5: md5,
			CaseId:     caseId,
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
		}
	}
	return c.CommonUsecase.DB().Save(&entity).Error
}
