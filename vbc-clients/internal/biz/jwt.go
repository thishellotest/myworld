package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"vbc/internal/conf"
)

const CTX_USER_KEY = "__user"

type JWTUsecase struct {
	log             *log.Helper
	CommonUsecase   *CommonUsecase
	conf            *conf.Data
	AppTokenUsecase *AppTokenUsecase
	UserUsecase     *UserUsecase
	TUsecase        *TUsecase
}

func NewJWTUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AppTokenUsecase *AppTokenUsecase,
	UserUsecase *UserUsecase,
	TUsecase *TUsecase) *JWTUsecase {
	uc := &JWTUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		AppTokenUsecase: AppTokenUsecase,
		UserUsecase:     UserUsecase,
		TUsecase:        TUsecase,
	}

	return uc
}

// JWTAuth tokenOk=true时 json没有数据，=false：json有具体错误信息
func (c *JWTUsecase) JWTAuth(authHeader string, ctx *gin.Context) (tokenOk bool, tUser *TData, json gin.H) {
	defer func() {
		if r := recover(); r != nil {
			tokenOk = false
			json = gin.H{
				"code": Reply_code_BadRequest,
				"msg":  "recover:" + InterfaceToString(r),
			}
			return
		}
	}()
	//CreateReply()
	//authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		json = gin.H{
			"code": Reply_code_BadRequest,
			"msg":  "Authorization is empty",
		}
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		json = gin.H{
			"code": Reply_code_jwt_expired,
			"msg":  "Authorization Bearer is wrong",
		}
		return
	}
	userId, isExp, err := c.AppTokenUsecase.Parse(parts[1])
	if err != nil {
		json = gin.H{
			"code": Reply_code_jwt_expired,
			"msg":  err.Error(),
		}
		return
	}
	if isExp {
		json = gin.H{
			"code": Reply_code_jwt_expired,
			"msg":  "JWT has expired",
		}
		return
	}

	tUser, err = c.TUsecase.DataById(Kind_users, userId)
	if err != nil {
		json = gin.H{
			"code": Reply_code_internal_error,
			"msg":  err.Error(),
		}
		return
	}
	if tUser == nil {
		json = gin.H{
			"code": Reply_code_BadRequest,
			"msg":  "User does not exist",
		}
		return
	}
	if tUser.CustomFields.NumberValueByNameBasic(UserFieldName_status) == 0 {
		json = gin.H{
			"code": Reply_code_jwt_expired,
			"msg":  "JWT has expired",
		}
		return
	}
	tokenOk = true
	if ctx != nil {
		ctx.Set(CTX_USER_KEY, tUser.Id())
		ctx.Set(CTX_USER_KEY, *tUser)
	}
	return
}

// HandleJWTAuth authHeader: Bearer aaa
func (c *JWTUsecase) HandleJWTAuth(authHeader string, ctx *gin.Context) {
	tokenOk, _, json := c.JWTAuth(authHeader, ctx)
	if !tokenOk {
		ctx.JSON(http.StatusOK, json)
		ctx.Abort()
		return
	}
	ctx.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
}

// JWTAuthMiddleware 基于JWT的认证中间件
func (c *JWTUsecase) JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		c.HandleJWTAuth(authHeader, ctx)
	}
}

func (c *JWTUsecase) JWTUser(ctx *gin.Context) (user TData, err error) {
	entity, exists := ctx.Get(CTX_USER_KEY)
	if !exists {
		err = errors.New("CTX_USER_KEY does not exist")
		return
	}
	return entity.(TData), nil
}

func (c *JWTUsecase) JWTUserFacade(ctx *gin.Context) (user UserFacade, err error) {
	entity, exists := ctx.Get(CTX_USER_KEY)
	if !exists {
		err = errors.New("CTX_USER_KEY does not exist")
		return
	}
	userFacade := UserFacade{
		TData: entity.(TData),
	}
	return userFacade, nil
}
