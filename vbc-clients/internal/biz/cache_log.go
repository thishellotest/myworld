package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	CacheType_AiResult = "AiResult"
)

type CacheLogEntity struct {
	ID        int32 `gorm:"primaryKey"`
	CacheType string
	CacheKey  string
	ResultId  string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

func (CacheLogEntity) TableName() string {
	return "cache_log"
}

type CacheLogUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[CacheLogEntity]
}

func NewCacheLogUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *CacheLogUsecase {
	uc := &CacheLogUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *CacheLogUsecase) GetForAiResult(cacheKey string) (*CacheLogEntity, error) {
	return c.Get(CacheType_AiResult, cacheKey)
}

func (c *CacheLogUsecase) AddForAiResult(cacheKey string, resultId string) error {
	return c.Add(CacheType_AiResult, cacheKey, resultId)
}

func (c *CacheLogUsecase) Get(cacheType, cacheKey string) (*CacheLogEntity, error) {
	return c.GetByCondWithOrderBy(Eq{"cache_type": cacheType, "cache_key": cacheKey, "deleted_at": 0}, "id desc")
}

func (c *CacheLogUsecase) Add(cacheType, cacheKey string, resultId string) error {
	entity := CacheLogEntity{
		CacheType: cacheType,
		CacheKey:  cacheKey,
		ResultId:  resultId,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	return c.CommonUsecase.DB().Save(&entity).Error
}
