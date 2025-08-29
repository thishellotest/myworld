package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type HttpUserUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	JWTUsecase       *JWTUsecase
	UserUsecase      *UserUsecase
	AppTokenUsecase  *AppTokenUsecase
	TimezonesUsecase *TimezonesUsecase
	RoleUsecase      *RoleUsecase
	MenuUsecase      *MenuUsecase
}

func NewHttpUserUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	UserUsecase *UserUsecase,
	AppTokenUsecase *AppTokenUsecase,
	TimezonesUsecase *TimezonesUsecase,
	RoleUsecase *RoleUsecase,
	MenuUsecase *MenuUsecase) *HttpUserUsecase {
	uc := &HttpUserUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		JWTUsecase:       JWTUsecase,
		UserUsecase:      UserUsecase,
		AppTokenUsecase:  AppTokenUsecase,
		TimezonesUsecase: TimezonesUsecase,
		RoleUsecase:      RoleUsecase,
		MenuUsecase:      MenuUsecase,
	}

	return uc
}

func (c *HttpUserUsecase) Info(ctx *gin.Context) {
	reply := CreateReply()

	userFacade, err := c.JWTUsecase.JWTUserFacade(ctx)
	if err != nil {
		reply.CommonError(err)
	} else {
		userProfile, err := c.UserUsecase.GetProfile(&userFacade.TData)
		if err != nil {
			reply.CommonError(err)
		} else {
			if userProfile == nil {
				reply.CommonError(errors.New("Profile is nil"))
			} else {

				userRole, _ := c.RoleUsecase.GetRole(userFacade.CustomFields.TextValueByNameBasic(UserFieldName_role_gid))
				if userRole == nil {
					reply.CommonError(errors.New("Role is nil"))
				} else {

					data, err := c.BizUserToInfoApi(&userFacade, *userProfile, *userRole)
					if err != nil {
						reply.CommonError(err)
					} else {
						reply["data"] = data
						reply.Success()
					}
					//reply["data"] = UserToInfoApi(&userFacade, *userProfile, *userRole, c.TimezonesUsecase)

				}
			}
		}
	}
	ctx.JSON(200, reply)
}

func (c *HttpUserUsecase) BizUserToInfoApi(uf *UserFacade, userProfile TData, userRole TData) (lib.TypeMap, error) {

	menu, err := c.MenuUsecase.GetMenu(*uf)
	if err != nil {
		return nil, err
	}
	data := UserToInfoApi(uf, userProfile, userRole, c.TimezonesUsecase)
	data.Set("is_mgmt", 0)
	data.Set("mgmt_dashboard_url", "")
	isOk := false
	for _, v := range menu {
		for _, v1 := range v.Children {
			data.Set("is_mgmt", 1)
			data.Set("mgmt_dashboard_url", v1.Url)
			isOk = true
			break
		}
		if isOk {
			break
		}
	}
	return data, nil
}

func (c *HttpUserUsecase) SignIn(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizSignIn(body.GetString("email"), body.GetString("password"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpUserUsecase) BizSignIn(email string, password string) (lib.TypeMap, error) {

	if configs.IsProd() {
		return nil, errors.New("Nonsupport")
	}
	data := make(lib.TypeMap)
	user, err := c.UserUsecase.GetByCond(Eq{"email": email})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("The user does not exist or the password is incorrect")
	}
	if password == "" {
		return nil, errors.New("Please enter password")
	}
	if email == "" {
		return nil, errors.New("Please enter email")
	}
	ok := c.UserUsecase.VerifyPassword(password, user.Password)
	if !ok {
		return nil, errors.New("The user does not exist or the password is incorrect")
	}

	appToken, err := c.AppTokenUsecase.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}
	data.Set("Data.jwt", appToken.AccessToken)
	return data, nil
}
