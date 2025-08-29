package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
)

type TypeMgmtPermission string

const (
	MgmtPermission_ManagementAttorneys      = TypeMgmtPermission("Management Attorneys")
	MgmtPermission_ReviseContract           = TypeMgmtPermission("Revise contract")
	MgmtPermission_DataManagementRecycleBin = TypeMgmtPermission("Data Management Recycle Bin")
)

type MgmtPermissionUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[TTemplateEntity]
	UserUsecase *UserUsecase
}

func NewMgmtPermissionUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	UserUsecase *UserUsecase,
) *MgmtPermissionUsecase {
	uc := &MgmtPermissionUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		UserUsecase:   UserUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *MgmtPermissionUsecase) Verify(userFacade UserFacade, mgmtPermission TypeMgmtPermission) (bool, error) {

	tProfile, err := c.UserUsecase.GetProfile(&userFacade.TData)
	if err != nil {
		return false, err
	}
	if tProfile == nil {
		return false, errors.New("tProfile is nil")
	}
	if IsAdminProfile(tProfile) {
		return true, nil
	}
	if mgmtPermission == MgmtPermission_ReviseContract { // 关闭此功能
		return false, nil
	}
	values := userFacade.CustomFields.TFieldMultiValuesByName(UserFieldName_permissions)
	for _, v := range values {
		if TypeMgmtPermission(v.Value) == mgmtPermission {
			return true, nil
		}
	}
	return false, nil
}
