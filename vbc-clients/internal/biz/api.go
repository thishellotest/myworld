package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type ApiUsecase struct {
	log *log.Helper
}

func NewApiUsecase() *ApiUsecase {
	return &ApiUsecase{}
}

func (u *ApiUsecase) SayHi(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "hello world",
	})
}
