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
	CounterKey_ClientUploadFiles = "ClientUploadFiles:"
)

func CounterKeyClientUploadFiles(caseId int32) string {
	return fmt.Sprintf("%s%d", CounterKey_ClientUploadFiles, caseId)
}

type CounterEntity struct {
	ID         int32 `gorm:"primaryKey"`
	CounterKey string
	Count      int
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  int64
}

func (CounterEntity) TableName() string {
	return "counters"
}

type CounterUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[CounterEntity]
}

func NewCounterUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *CounterUsecase {
	uc := &CounterUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *CounterUsecase) HasLimit(counterKey string, maxCount int) (hasLimit bool, err error) {
	counter, err := c.GetByCond(Eq{"counter_key": counterKey, "deleted_at": 0})
	if err != nil {
		return false, err
	}
	if counter == nil {
		return false, nil
	}
	if counter.Count >= maxCount {
		return true, nil
	}
	return false, nil
}

func (c *CounterUsecase) Stat(counterKey string, count int) error {

	r, er := c.GetByCond(Eq{"counter_key": counterKey, "deleted_at": 0})
	if er != nil {
		return er
	}

	if r == nil {
		sql := fmt.Sprintf("INSERT INTO %s (counter_key, `count` ,created_at, updated_at) VALUES(\"%s\",  %d, %d, %d) "+
			"ON DUPLICATE KEY UPDATE count=count+values(count),updated_at=values(updated_at),created_at=values(created_at)",
			CounterEntity{}.TableName(), counterKey, count, time.Now().Unix(), time.Now().Unix())

		return c.CommonUsecase.DB().Exec(sql).Error
	} else {
		return c.CommonUsecase.DB().Model(&r).Updates(map[string]interface{}{
			"count": gorm.Expr("count+?", count),
		}).Error
	}
}
