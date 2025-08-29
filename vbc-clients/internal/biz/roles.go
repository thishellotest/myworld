package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/configs"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	Role_is_user_manager_No  = 0
	Role_is_user_manager_Yes = 1

	RoleFieldName_is_user_manager = "is_user_manager"
)

type RoleUsecase struct {
	log            *log.Helper
	conf           *conf.Data
	CommonUsecase  *CommonUsecase
	TUsecase       *TUsecase
	GoCacheUsecase *GoCacheUsecase
}

func NewRoleUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	GoCacheUsecase *GoCacheUsecase,
) *RoleUsecase {
	uc := &RoleUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		conf:           conf,
		TUsecase:       TUsecase,
		GoCacheUsecase: GoCacheUsecase,
	}

	return uc
}

func (c *RoleUsecase) AllRoles() ([]*TData, error) {
	return c.TUsecase.ListByCond(Kind_roles, Eq{DataEntry_biz_deleted_at: 0})
}

func (c *RoleUsecase) CacheAllRoles() ([]*TData, error) {
	key := fmt.Sprintf("%s%s", "RoleUsecase:", "CacheAllRoles")
	res, found := GoCacheGet[[]*TData](c.GoCacheUsecase, key)
	if found {
		return res, nil
	}
	var err error
	res, err = c.AllRoles()
	if err != nil {
		return nil, err
	}
	GoCacheSet(c.GoCacheUsecase, key, res, configs.CacheExpiredDurationDefault)
	return res, err
}

func (c *RoleUsecase) HandleChildrenRoles(allRoles []*TData, currentRole *TData, childrenRoles *[]*TData) {
	for k, v := range allRoles {
		if v.CustomFields.TextValueByNameBasic("role_parent_gid") == currentRole.Gid() {
			*childrenRoles = append(*childrenRoles, allRoles[k])
			c.HandleChildrenRoles(allRoles, allRoles[k], childrenRoles)
		}
	}
}

func (c *RoleUsecase) ChildrenRoles(currentRole *TData) (childrenRoles []*TData, err error) {

	if currentRole == nil {
		return nil, errors.New("currentRole is nil")
	}
	allRoles, err := c.CacheAllRoles()
	if err != nil {
		return nil, err
	}
	c.HandleChildrenRoles(allRoles, currentRole, &childrenRoles)
	return
}

func (c *RoleUsecase) ChildrenRolesGids(currentRole *TData) (gids []string, err error) {
	dest, err := c.ChildrenRoles(currentRole)
	if err != nil {
		return nil, err
	}
	for _, v := range dest {
		gids = append(gids, v.Gid())
	}
	return gids, nil
}

func (c *RoleUsecase) GetRole(roleGid string) (role *TData, err error) {
	return c.TUsecase.DataByGid(Kind_roles, roleGid)
}
