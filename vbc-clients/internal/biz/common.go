package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"vbc/lib"
)

type CommonRepo interface {
	DB() *gorm.DB
	RedisClient() *redis.Client
}

type CommonUsecase struct {
	Repo CommonRepo
	Log  *log.Helper
}

func NewCommonUsecase(repo CommonRepo, logger log.Logger) *CommonUsecase {
	return &CommonUsecase{
		Repo: repo,
		Log:  log.NewHelper(logger),
	}
}

func (c *CommonUsecase) DB() *gorm.DB {
	return c.Repo.DB()
}

// Count 获取总数 select count(*) c from t1
func (c *CommonUsecase) Count(db *gorm.DB, sql string) (int64, error) {
	sqlRows, err := db.Raw(sql).Rows()
	if err != nil {
		return 0, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, list, err := lib.SqlRowsTrans(sqlRows)
	if err != nil {
		return 0, err
	}

	total, err := strconv.ParseInt(lib.InterfaceToString(list[0]["c"]), 10, 32)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (c *CommonUsecase) RedisClient() *redis.Client {
	return c.Repo.RedisClient()
}
