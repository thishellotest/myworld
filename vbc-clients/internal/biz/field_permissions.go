package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
)

type TypePermission string

const (
	Permission_ReadAndWrite = TypePermission("")
	Permission_ReadOnly     = TypePermission("readonly")
	Permission_DoNotShow    = TypePermission("donotshow")
)

type FieldPermissionEntity struct {
	ID         int32 `gorm:"primaryKey"`
	Gid        string
	ProfileGid string
	FieldKind  string
	FieldName  string
	Permission TypePermission
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  int64
}

func (FieldPermissionEntity) TableName() string {
	return "field_permissions"
}

type FieldPermissionUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[FieldPermissionEntity]
	GoCacheUsecase *GoCacheUsecase
	FieldUsecase   *FieldUsecase
}

func NewFieldPermissionUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	GoCacheUsecase *GoCacheUsecase,
	FieldUsecase *FieldUsecase) *FieldPermissionUsecase {
	uc := &FieldPermissionUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		conf:           conf,
		GoCacheUsecase: GoCacheUsecase,
		FieldUsecase:   FieldUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

type TypeFieldPermissionList []FieldPermissionEntity

type TypeFieldPermissionStruct struct {
	Kind         string
	List         TypeFieldPermissionList
	FieldNameIdx map[string]FieldPermissionEntity
}

func (c *TypeFieldPermissionStruct) GetByFieldName(fieldName string) *FieldPermissionEntity {
	if c.FieldNameIdx != nil {
		if _, ok := c.FieldNameIdx[fieldName]; ok {
			entity := c.FieldNameIdx[fieldName]
			return &entity
		}
	}
	return nil
}

func (c *TypeFieldPermissionStruct) Init(kind string, list TypeFieldPermissionList) {
	c.FieldNameIdx = make(map[string]FieldPermissionEntity)
	c.Kind = kind
	c.List = list
	for k, v := range list {
		c.FieldNameIdx[v.FieldName] = list[k]
	}
}

func (c *FieldPermissionUsecase) ListByKind(kind string, profileGid string) (list TypeFieldPermissionList, err error) {
	err = c.CommonUsecase.DB().Where("field_kind=? and deleted_at=0 and profile_gid=?", kind, profileGid).Find(&list).Error
	return
}

func (c *FieldPermissionUsecase) StructByKind(kind string, profileGid string) (*TypeFieldPermissionStruct, error) {
	list, err := c.ListByKind(kind, profileGid)
	if err != nil {
		return nil, err
	}
	res := &TypeFieldPermissionStruct{}
	res.Init(kind, list)
	return res, nil
}

func (c *FieldPermissionUsecase) CacheStructByKind(kind string, profileGid string) (*TypeFieldPermissionStruct, error) {
	key := fmt.Sprintf("%s%s:%s", GOCACHE_PREFIX_field_permission, kind, profileGid)
	res, found := GoCacheGet[*TypeFieldPermissionStruct](c.GoCacheUsecase, key)
	if found {
		return res, nil
	}
	var err error
	res, err = c.StructByKind(kind, profileGid)
	if err != nil {
		return nil, err
	}
	GoCacheSet[*TypeFieldPermissionStruct](c.GoCacheUsecase, key, res, configs.CacheExpiredDuration5Seconds)
	return res, nil
}

func (c *FieldPermissionUsecase) FieldPermissionCenter(kind string, profileGid string) (fieldPermissionCenter FieldPermissionCenter, err error) {
	fieldStruct, err := c.FieldUsecase.StructByKind(kind)
	if err != nil {
		return fieldPermissionCenter, err
	}
	return c.GetFieldPermissionCenter(kind, profileGid, fieldStruct)
	//if fieldStruct == nil {
	//	return fieldPermissionCenter, errors.New("fieldStruct is nil")
	//}
	//fieldPermissionStruct, err := c.StructByKind(kind, profileGid)
	//if err != nil {
	//	return fieldPermissionCenter, err
	//}
	//if fieldPermissionStruct == nil {
	//	return fieldPermissionCenter, errors.New("fieldPermissionStruct is nil")
	//}
	//fieldPermissionCenter.FieldStruct = *fieldStruct
	//fieldPermissionCenter.FieldPermissionStruct = *fieldPermissionStruct
	//return fieldPermissionCenter, nil
}

func (c *FieldPermissionUsecase) GetFieldPermissionCenter(kind string, profileGid string, fieldStruct *TypeFieldStruct) (fieldPermissionCenter FieldPermissionCenter, err error) {
	if err != nil {
		return fieldPermissionCenter, err
	}
	if fieldStruct == nil {
		return fieldPermissionCenter, errors.New("fieldStruct is nil")
	}
	fieldPermissionStruct, err := c.StructByKind(kind, profileGid)
	if err != nil {
		return fieldPermissionCenter, err
	}
	if fieldPermissionStruct == nil {
		return fieldPermissionCenter, errors.New("fieldPermissionStruct is nil")
	}
	fieldPermissionCenter.FieldStruct = *fieldStruct
	fieldPermissionCenter.FieldPermissionStruct = *fieldPermissionStruct
	return fieldPermissionCenter, nil
}

func (c *FieldPermissionUsecase) CacheFieldPermissionCenter(kind string, profileGid string) (fieldPermissionCenter FieldPermissionCenter, err error) {
	fieldStruct, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return fieldPermissionCenter, err
	}
	return c.GetFieldPermissionCenter(kind, profileGid, fieldStruct)
	//if fieldStruct == nil {
	//	return fieldPermissionCenter, errors.New("fieldStruct is nil")
	//}
	//fieldPermissionStruct, err := c.CacheStructByKind(kind, profileGid)
	//if err != nil {
	//	return fieldPermissionCenter, err
	//}
	//if fieldPermissionStruct == nil {
	//	return fieldPermissionCenter, errors.New("fieldPermissionStruct is nil")
	//}
	//fieldPermissionCenter.FieldStruct = *fieldStruct
	//fieldPermissionCenter.FieldPermissionStruct = *fieldPermissionStruct
	//return fieldPermissionCenter, nil
}

type FieldPermissionVo struct {
	FieldName  string
	Permission TypePermission
}

// CanShow 字段是否可以显示
func (c *FieldPermissionVo) CanShow() bool {
	if c.Permission == Permission_ReadOnly || c.Permission == Permission_ReadAndWrite {
		return true
	}
	return false
}

var FieldsOnlyWriteFromSystem = []string{DataEntry_created_at, DataEntry_updated_at, DataEntry_created_by, DataEntry_modified_by, DataEntry_updated_at, FieldName_full_name}

// CanWrite 字段是否可以写入（系统保留字段，只允许内部写入，此处不允许）
func (c *FieldPermissionVo) CanWrite() bool {
	if lib.InArray(c.FieldName, FieldsOnlyWriteFromSystem) {
		return false
	}
	if c.Permission == Permission_ReadAndWrite {
		return true
	}
	return false
}

type FieldPermissionCenter struct {
	FieldStruct           TypeFieldStruct
	FieldPermissionStruct TypeFieldPermissionStruct
}

func (c *FieldPermissionCenter) PermissionByFieldName(fieldName string) (fieldPermissionVo FieldPermissionVo, err error) {
	fieldEntity := c.FieldStruct.GetByFieldName(fieldName)
	if fieldEntity == nil {
		return fieldPermissionVo, errors.New(fieldName + " PermissionByFieldName: fieldEntity is nil")
	}
	fieldPermissionVo.FieldName = fieldName
	if fieldEntity.FieldPermissionLevel == FieldPermissionLevel_OnlyRead {
		fieldPermissionVo.Permission = Permission_ReadOnly
	} else if fieldEntity.FieldPermissionLevel == FieldPermissionLevel_OnlyReadAndWrite {
		fieldPermissionVo.Permission = Permission_ReadAndWrite
	} else {
		fieldPermissionEntity := c.FieldPermissionStruct.GetByFieldName(fieldName)
		if fieldPermissionEntity == nil {
			fieldPermissionVo.Permission = Permission_ReadAndWrite
		} else {
			fieldPermissionVo.Permission = fieldPermissionEntity.Permission
		}
	}
	return fieldPermissionVo, nil
}
