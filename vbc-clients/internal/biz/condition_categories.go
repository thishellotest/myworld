package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type ConditionCategoryEntity struct {
	ID           int32 `gorm:"primaryKey"`
	CategoryName string
	BizDeletedAt int64
	CreatedAt    int64
	UpdatedAt    int64
}

func (c *ConditionCategoryEntity) ToFabFieldOption() FabFieldOption {
	return FabFieldOption{
		OptionLabel: c.CategoryName,
		OptionValue: InterfaceToString(c.ID),
	}
}

func (ConditionCategoryEntity) TableName() string {
	return "condition_categories"
}

type ConditionCategoryVo struct {
	ConditionCategoryId   int32  `json:"condition_category_id"`
	ConditionCategoryName string `json:"condition_category_name"`
}

var UndefinedConditionCategory = ConditionCategoryVo{
	ConditionCategoryId:   0,
	ConditionCategoryName: "Undefined Category",
}

func GetConditionCategoryVo(entity *ConditionCategoryEntity) ConditionCategoryVo {
	if entity == nil {
		return UndefinedConditionCategory
	}
	return entity.ToConditionCategoryVo()
}

func (c *ConditionCategoryEntity) ToConditionCategoryVo() ConditionCategoryVo {
	return ConditionCategoryVo{
		ConditionCategoryId:   c.ID,
		ConditionCategoryName: c.CategoryName,
	}
}

type ConditionCategoryUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[ConditionCategoryEntity]
}

func NewConditionCategoryUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *ConditionCategoryUsecase {
	uc := &ConditionCategoryUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *ConditionCategoryUsecase) Upsert(categoryName string) (entity *ConditionCategoryEntity, err error) {
	entity, err = c.GetByCond(Eq{"category_name": categoryName})
	if err != nil {
		return nil, err
	}
	if entity != nil {
		return entity, nil
	}
	entity = &ConditionCategoryEntity{
		CategoryName: categoryName,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = c.CommonUsecase.DB().Save(&entity).Error
	return
}
