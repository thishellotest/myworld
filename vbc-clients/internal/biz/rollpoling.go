package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	Rollpoing_Vendor_boxsign = "boxsign"
)

type RollpoingEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	NextAt             int64
	Timeout            int32
	HandleResultDetail string
	Vendor             string
	VendorUniqId       string
	ResponseText       string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

func (RollpoingEntity) TableName() string {
	return "rollpoling"
}

type RollpoingUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[RollpoingEntity]
}

func NewRollpoingUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *RollpoingUsecase {
	uc := &RollpoingUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

// Upsert 添加任务
func (c *RollpoingUsecase) Upsert(vendor string, vendorUniqId string) error {

	entity, err := c.GetByCond(Eq{"vendor": vendor, "vendor_uniq_id": vendorUniqId})
	if err != nil {
		return err
	}
	if entity == nil {
		entity = &RollpoingEntity{
			Vendor:       vendor,
			VendorUniqId: vendorUniqId,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		return c.CommonUsecase.DB().Save(&entity).Error
	}
	return nil
}
