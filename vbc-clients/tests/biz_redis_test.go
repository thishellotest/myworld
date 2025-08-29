package tests

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_RedisUsecase_test(t *testing.T) {

	a := UT.CommonUsecase.RedisClient().Get(context.TODO(), "foo")
	fmt.Println(a.Val())
	//UT.RedisUsecase.Test()
}

func Test_RedisUsecase_CustomTaskParams11(t *testing.T) {

	cc := []biz.CustomTaskParams{{
		UniqueKey: "u1",
		Params:    "afdsafsaf",
	}}

	err := UT.CommonUsecase.RedisClient().LPush(context.TODO(), "lis", cc).Err()
	lib.DPrintln(err)
}

func Test_RedisUsecase_CustomTaskParams1(t *testing.T) {

	cc := biz.CustomTaskParams{
		UniqueKey: "u1",
		Params:    "afdsafsaf",
	}

	err := UT.CommonUsecase.RedisClient().LPush(context.TODO(), "lis", cc, cc).Err()
	lib.DPrintln(err)
}

func Test_RedisUsecase_CustomTaskParams2(t *testing.T) {

	a := UT.CommonUsecase.RedisClient().LPop(context.TODO(), "lis")
	//lib.DPrintln(a.Val())
	if a != nil {
		entity := &biz.CustomTaskParams{}
		err := a.Scan(entity)
		lib.DPrintln(err, entity)
	}
}

func Test_RedisUsecase_CustomTaskParams3(t *testing.T) {

	a := UT.CommonUsecase.RedisClient().LRange(context.TODO(), "lis", 0, -1)

	var res []biz.CustomTaskParams
	err := a.ScanSlice(&res)
	lib.DPrintln(err, res)

}

func Test_redis(t *testing.T) {
	redis := redisInit()
	redis.Set(context.TODO(), "a", "b", time.Second*30)
	c := redis.Get(context.TODO(), "a")
	lib.DPrintln(c.String())
}

func redisInit() *redis.Client {
	url := "redis://default:cb4Qb9LfydHXUcE4Xz7o5XvWgQrbsK40@redis-18626.c17.us-east-1-4.ec2.cloud.redislabs.com:18626?dial_timeout=3&db=0&read_timeout=6s&max_retries=2"
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	opt.MaxActiveConns = 5
	opt.PoolSize = 5
	client := redis.NewClient(opt)
	ping := client.Ping(context.TODO())
	if ping.Err() != nil {
		if ping.Err() != redis.Nil {
			panic(ping.Err())
		}
	}
	return client
}
