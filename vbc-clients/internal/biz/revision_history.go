package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
)

const (
	RevisionHistory_BizType_contract = "contract"
)

type RevisionHistoryEntity struct {
	ID        int32 `gorm:"primaryKey"`
	BizType   string
	Uniqid    string
	OldValue  string
	NewValue  string
	ChangedBy string
	CreatedAt int64
	UpdatedAt int64
}

func (RevisionHistoryEntity) TableName() string {
	return "revision_history"
}

type RevisionHistoryToContractApi struct {
	Uniqid         string      `json:"uniqid"`
	ClientCaseName string      `json:"client_case_name"`
	OldValue       lib.TypeMap `json:"old_value"`
	NewValue       lib.TypeMap `json:"new_value"`
	ChangedByName  string      `json:"changed_by_name"`
	UpdatedAt      int32       `json:"updated_at"`
}

func (c *RevisionHistoryEntity) ToContractApi(cases map[string]*TData, users map[string]*TData) (vo RevisionHistoryToContractApi) {

	if v, ok := users[c.ChangedBy]; ok {
		vo.ChangedByName = v.CustomFields.TextValueByNameBasic(UserFieldName_fullname)
	}
	if v, ok := cases[c.Uniqid]; ok {
		vo.ClientCaseName = v.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	}
	vo.Uniqid = c.Uniqid
	vo.OldValue = lib.ToTypeMapByString(c.OldValue)
	vo.NewValue = lib.ToTypeMapByString(c.NewValue)
	vo.UpdatedAt = int32(c.UpdatedAt)
	return
}

type RevisionHistoryUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[RevisionHistoryEntity]
}

func NewRevisionHistoryUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *RevisionHistoryUsecase {
	uc := &RevisionHistoryUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *RevisionHistoryUsecase) Add(bizType, uniqid, oldValue, newValue string, userGid string) (entity RevisionHistoryEntity, err error) {

	entity = RevisionHistoryEntity{
		BizType:   bizType,
		Uniqid:    uniqid,
		OldValue:  oldValue,
		NewValue:  newValue,
		ChangedBy: userGid,
		UpdatedAt: time.Now().Unix(),
		CreatedAt: time.Now().Unix(),
	}
	err = c.CommonUsecase.DB().Save(&entity).Error
	return
}
