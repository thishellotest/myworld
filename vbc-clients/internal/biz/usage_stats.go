package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	UsageType_HandleAmount         = "HandleAmount"
	UsageType_HandleClientCaseName = "HandleClientCaseName"
	UsageType_GetContactRecords    = "GetContactRecords"
	UsageType_GetDealRecords       = "GetDealRecords"
	UsageType_GetTaskRecords       = "GetTaskRecords"
	UsageType_GetNoteRecords       = "GetNoteRecords"

	UsageType_PREFIX_BOX = "BOX:"
)

func UsageTypeValue(prefix string, val string) string {
	return UsageType_PREFIX_BOX + val
}

type UsageStatsEntity struct {
	ID        int32 `gorm:"primaryKey"`
	UsageType string
	Day       string
	Count     int
	CreatedAt int64
	UpdatedAt int64
}

func (UsageStatsEntity) TableName() string {
	return "usage_stats"
}

type UsageStatsUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[UsageStatsEntity]
}

func NewUsageStatsUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *UsageStatsUsecase {
	uc := &UsageStatsUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *UsageStatsUsecase) Stat(usageType string, day time.Time, count int) error {

	d := day.Format("2006-01-02")
	r, er := c.GetByCond(Eq{"usage_type": usageType, "day": d})
	if er != nil {
		return er
	}

	if r == nil {
		sql := fmt.Sprintf("INSERT INTO %s (usage_type, `day`, `count` ,created_at, updated_at) VALUES(\"%s\", \"%s\", %d, %d, %d) "+
			"ON DUPLICATE KEY UPDATE count=count+values(count),updated_at=values(updated_at),created_at=values(created_at)",
			UsageStatsEntity{}.TableName(), usageType, d, count, time.Now().Unix(), time.Now().Unix())

		return c.CommonUsecase.DB().Exec(sql).Error
	} else {
		c.CommonUsecase.DB().Model(&r).Updates(map[string]interface{}{
			"count": gorm.Expr("count+?", count),
		})
	}
	return nil
	//return nil
}
