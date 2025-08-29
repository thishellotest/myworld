package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
)

type MenuUsecase struct {
	log                   *log.Helper
	conf                  *conf.Data
	CommonUsecase         *CommonUsecase
	UserUsecase           *UserUsecase
	MgmtPermissionUsecase *MgmtPermissionUsecase
	RoleUsecase           *RoleUsecase
}

func NewMenuUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	UserUsecase *UserUsecase,
	MgmtPermissionUsecase *MgmtPermissionUsecase,
	RoleUsecase *RoleUsecase,
) *MenuUsecase {
	uc := &MenuUsecase{
		log:                   log.NewHelper(logger),
		CommonUsecase:         CommonUsecase,
		conf:                  conf,
		UserUsecase:           UserUsecase,
		MgmtPermissionUsecase: MgmtPermissionUsecase,
		RoleUsecase:           RoleUsecase,
	}

	return uc
}

type MenuItemList []MenuItem
type MenuItem struct {
	Uniqid    string       `json:"-"`
	Label     string       `json:"label"`
	Url       string       `json:"url,omitempty"`
	Children  MenuItemList `json:"children,omitempty"`
	OnlyAdmin bool         `json:"-"`
}

var ConfigMenuItemList = MenuItemList{
	{
		Label: "General",
		Children: MenuItemList{
			{
				Uniqid: "Users",
				Label:  "Users",
				Url:    "/settings/users",
			},
			{
				Uniqid: "Contract",
				Label:  "Contract",
				Url:    "/settings/contract",
			},
			{
				Uniqid: "Attorneys",
				Label:  "Attorneys",
				Url:    "/settings/attorneys",
			},
			//{
			//	Label:     "Roles",
			//	Url:       "/settings/roles",
			//	OnlyAdmin: true,
			//},
		},
	},
	{
		Label: "Data Administration",
		Children: MenuItemList{
			{
				Uniqid: "RecycleBin",
				Label:  "Recycle Bin",
				Url:    "/settings/recyclebin",
				//OnlyAdmin: true,
			},
		},
	},
}

func (c *MenuUsecase) GetMenu(userFacade UserFacade) (root MenuItemList, err error) {

	tProfile, err := c.UserUsecase.GetProfile(&userFacade.TData)
	if err != nil {
		return nil, err
	}
	if tProfile == nil {
		return nil, errors.New("tProfile is nil")
	}

	userRole, err := c.RoleUsecase.GetRole(userFacade.CustomFields.TextValueByNameBasic(UserFieldName_role_gid))
	if err != nil {
		return nil, err
	}
	if userRole == nil {
		return nil, errors.New("userRole is nil")
	}

	if IsAdminProfile(tProfile) {
		return ConfigMenuItemList, nil
	}
	for k, _ := range ConfigMenuItemList {
		HandleMenu(c.log, ConfigMenuItemList[k], &root, userFacade, *tProfile, c.MgmtPermissionUsecase, *userRole)
	}

	root = CleanMenu(root)
	return root, nil
}

func HandleMenu(log *log.Helper, menuItem MenuItem, result *MenuItemList, userFacade UserFacade, tProfile TData, MgmtPermissionUsecase *MgmtPermissionUsecase, userRole TData) {
	if menuItem.OnlyAdmin == false {
		newMenuItem := menuItem
		var children MenuItemList
		for k, v := range menuItem.Children {
			if v.Uniqid == "Contract" {
				hasPermission, err := MgmtPermissionUsecase.Verify(userFacade, MgmtPermission_ReviseContract)
				if err != nil {
					log.Error(err)
					continue
				}
				if !hasPermission {
					continue
				}
			} else if v.Uniqid == "RecycleBin" {
				hasPermission, err := MgmtPermissionUsecase.Verify(userFacade, MgmtPermission_DataManagementRecycleBin)
				if err != nil {
					log.Error(err)
					continue
				}
				if !hasPermission {
					continue
				}
			} else if v.Uniqid == "Attorneys" {
				hasPermission, err := MgmtPermissionUsecase.Verify(userFacade, MgmtPermission_ManagementAttorneys)
				if err != nil {
					log.Error(err)
					continue
				}
				if !hasPermission {
					continue
				}
			} else if v.Uniqid == "Users" {
				if userRole.CustomFields.NumberValueByNameBasic(RoleFieldName_is_user_manager) == 0 {
					continue
				}
			}

			HandleMenu(log, menuItem.Children[k], &children, userFacade, tProfile, MgmtPermissionUsecase, userRole)
		}
		newMenuItem.Children = children
		*result = append(*result, newMenuItem)
	}
}

func CleanMenu(menu MenuItemList) (r MenuItemList) {
	for k, v := range menu {
		if len(v.Children) > 0 {
			r = append(r, menu[k])
		}
	}
	return
}
