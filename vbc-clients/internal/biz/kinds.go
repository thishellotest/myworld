package biz

import (
	"fmt"
	"github.com/pkg/errors"
	"vbc/configs"
)

const NoChangeHistory_Yes = 1
const NoTimelines_Yes = 1

type KindEntity struct {
	ID               int32 `gorm:"primaryKey"`
	Kind             string
	Tablename        string
	Label            string // 例：用于说明这条记录叫什么
	TabLabel         string // 例：用于显示导航名称等
	PrimaryFieldName string
	NoChangeHistory  int
	NoTimelines      int
	DeletedAt        int64
}

func (c *KindEntity) KindTableName() string {
	return c.Tablename
}

func (c *KindEntity) ModuleName() string {
	return KindConvertToModule(c.Kind)
}

func (c *KindEntity) ModuleLabel() string {
	return c.Label
}

func (KindEntity) TableName() string {
	return "kinds"
}

type KindUsecase struct {
	CommonUsecase  *CommonUsecase
	GoCacheUsecase *GoCacheUsecase
}

func NewKindUsecase(CommonUsecase *CommonUsecase, GoCacheUsecase *GoCacheUsecase) *KindUsecase {

	return &KindUsecase{
		CommonUsecase:  CommonUsecase,
		GoCacheUsecase: GoCacheUsecase,
	}
}

// *KindEntity 改为： KindEntity 指针在cache时，有可能出现问题
type TypeKindList []KindEntity

func (c TypeKindList) TableNameByKind(kind string) (string, error) {
	for _, v := range c {
		if v.Kind == kind {
			return v.Tablename, nil
		}
	}
	return "", errors.New(kind + ":未找到对应表名")
}

func (c TypeKindList) GetByKind(kind string) *KindEntity {
	for _, v := range c {
		if v.Kind == kind {
			e := v
			return &e
		}
	}
	return nil
}

func (c *KindUsecase) List() (r TypeKindList, err error) {
	err = c.CommonUsecase.DB().Where("deleted_at=0").Find(&r).Error
	return
}
func (c *KindUsecase) CacheList() (r TypeKindList, err error) {
	key := fmt.Sprintf("%s", GOCACHE_PREFIX_kind)
	r, found := GoCacheGet[TypeKindList](c.GoCacheUsecase, key)
	if found {
		return r, nil
	}
	r, err = c.List()
	if err != nil {
		return nil, err
	}
	GoCacheSet(c.GoCacheUsecase, key, r, configs.CacheExpiredDurationDefault)
	return r, nil
}

func (c *KindUsecase) CacheTableNameByKind(kind string) (tableName string, err error) {
	r, err := c.CacheList()
	if err != nil {
		return "", err
	}
	return r.TableNameByKind(kind)
}

func (c *KindUsecase) GetByKind(kind string) (*KindEntity, error) {
	r, err := c.CacheList()
	if err != nil {
		return nil, err
	}
	return r.GetByKind(kind), nil
}
