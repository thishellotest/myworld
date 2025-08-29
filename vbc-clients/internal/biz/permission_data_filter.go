package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type PermissionDataFilterUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	FieldUsecase  *FieldUsecase
	TUsecase      *TUsecase
	RoleUsecase   *RoleUsecase
}

func NewPermissionDataFilterUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldUsecase *FieldUsecase,
	TUsecase *TUsecase,
	RoleUsecase *RoleUsecase) *PermissionDataFilterUsecase {
	uc := &PermissionDataFilterUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		FieldUsecase:  FieldUsecase,
		TUsecase:      TUsecase,
		RoleUsecase:   RoleUsecase,
	}

	return uc
}

func (c *PermissionDataFilterUsecase) Filter(kindEntity KindEntity, tProfile *TData, userFacade *UserFacade, aliasTableName string, normalUserOnlyOwner bool) (Cond, error) {
	if tProfile == nil {
		return nil, errors.New("tProfile is nil")
	}
	if IsAdminProfile(tProfile) {
		return nil, nil
	} else {
		fieldStruct, err := c.FieldUsecase.CacheStructByKind(kindEntity.Kind)
		if err != nil {
			return nil, err
		}
		if fieldStruct == nil {
			return nil, errors.New("fieldStruct is nil")
		}

		var conds []Cond
		if kindEntity.Kind == Kind_users {

			tRole, err := c.TUsecase.DataByGid(Kind_roles, userFacade.CustomFields.TextValueByNameBasic(UserFieldName_role_gid))
			if err != nil {
				return nil, err
			}
			if tRole == nil {
				return nil, errors.New("tRole is nil")
			}
			if tRole.CustomFields.NumberValueByNameBasic(RoleFieldName_is_user_manager) == Role_is_user_manager_Yes {
				childrenGids, err := c.RoleUsecase.ChildrenRolesGids(tRole)
				if err != nil {
					return nil, err
				}
				conds = append(conds, Or(Eq{TidyTableFieldForSql(DataEntry_gid, aliasTableName): userFacade.Gid()}, In(UserFieldName_role_gid, childrenGids)))
			} else {
				conds = append(conds, Eq{TidyTableFieldForSql(DataEntry_gid, aliasTableName): userFacade.Gid()})
			}
		} else {
			conds = append(conds, Eq{TidyTableFieldForSql(DataEntry_user_gid, aliasTableName): userFacade.Gid()})

			if !normalUserOnlyOwner {
				collaboratorField := fieldStruct.GetCollaborator()
				if collaboratorField != nil {
					conds = append(conds, Like{TidyTableFieldForSql(collaboratorField.FieldName, aliasTableName), fmt.Sprintf(",%s,", userFacade.Gid())})
				}
			}
		}

		return Or(conds...), nil
	}
}
