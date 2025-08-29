package biz

import (
	"github.com/patrickmn/go-cache"
	"time"
	"vbc/configs"
)

type GoCacheResult struct {
	Value interface{}
}

type GoCacheUsecase struct {
	Cache *cache.Cache
}

const (
	GOCACHE_PREFIX_field_validator  = "field_validator:"
	GOCACHE_PREFIX_field_permission = "field_permission:"
	GOCACHE_PREFIX_field_option     = "field_option:"
	GOCACHE_PREFIX_field            = "field:"
	GOCACHE_PREFIX_kind             = "kind"
)

func NewGoCacheUsecase() *GoCacheUsecase {

	c := cache.New(5*time.Minute, 10*time.Minute)

	return &GoCacheUsecase{
		Cache: c,
	}
}

func GoCacheGet[T any](goCacheUsecase *GoCacheUsecase, k string) (T, bool) {
	val, found := goCacheUsecase.Cache.Get(k)
	if !found {
		var a T
		return a, false
	}
	return val.(T), true
}

// GoCacheSet nil也可以cache
func GoCacheSet[T any](goCacheUsecase *GoCacheUsecase, k string, val T, duration time.Duration) {
	if duration == 0 {
		duration = configs.CacheExpiredDurationDefault
	}
	goCacheUsecase.Cache.Set(k, val, duration)
}
