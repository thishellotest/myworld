package data

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"vbc/internal/biz"
)

type CommonRepo struct {
	data *Data
}

func (c CommonRepo) RedisClient() *redis.Client {
	return c.data.RedisClient
}

func (c CommonRepo) DB() *gorm.DB {

	return c.data.Db
}

func NewCommonRepo(data *Data) biz.CommonRepo {

	return &CommonRepo{
		data: data,
	}
}
