package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	Condition_Type_Source             = 0
	Condition_Type_Source_From_Ai     = 4
	Condition_Type_Primary            = 1
	Condition_Type_SecondaryCondition = 2
	Condition_Type_Condition          = 10

	Condition_SecondaryType_DirectSecondary = "1"
	Condition_SecondaryType_Aggravation     = "2"

	Condition_SourceStatus_NotHandled = 0
	Condition_SourceStatus_Handled    = 1
)

const (
	ConditionFieldName_type                  = "type"
	ConditionFieldName_secondary_type        = "secondary_type"
	ConditionFieldName_condition_category_id = "condition_category_id"
	ConditionFieldName_condition_name        = "condition_name"
)

func ConditionSecondaryTypeNameById(id string) string {
	if id == Condition_SecondaryType_DirectSecondary {
		return "Direct Secondary"
	} else if id == Condition_SecondaryType_Aggravation {
		return "Aggravation"
	}
	return ""
}

type ConditionEntity struct {
	ID                  int32 `gorm:"primaryKey"`
	ConditionCategoryId int32
	Type                int
	SecondaryType       string
	ConditionName       string
	SourceStatus        int
	SourceHandleCount   int
	BizDeletedAt        int64
	CreatedAt           int64
	UpdatedAt           int64
	DeletedAt           int64
}

func (ConditionEntity) TableName() string {
	return "conditions"
}

func (c *ConditionEntity) ToConditionItem() ConditionItem {
	return ConditionItem{
		ID:   c.ID,
		Name: c.ConditionName,
	}
}

type ConditionItem struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
type ConditionSecondaries struct {
	SecondaryType     string          `json:"secondary_type"`
	SecondaryTypeName string          `json:"secondary_type_name"`
	List              []ConditionItem `json:"list"`
}

func InitConditionSecondaries() (DirectSecondary ConditionSecondaries, Aggravation ConditionSecondaries) {
	DirectSecondary.SecondaryType = Condition_SecondaryType_DirectSecondary
	DirectSecondary.SecondaryTypeName = ConditionSecondaryTypeNameById(Condition_SecondaryType_DirectSecondary)

	Aggravation.SecondaryType = Condition_SecondaryType_Aggravation
	Aggravation.SecondaryTypeName = ConditionSecondaryTypeNameById(Condition_SecondaryType_Aggravation)
	return
}

type PrimaryConditionVo struct {
	ID                  int32                  `json:"id"`
	Name                string                 `json:"name"`
	ConditionCategoryVo ConditionCategoryVo    `json:"category"`
	List                []ConditionSecondaries `json:"list"`
}

func (c *ConditionEntity) ToPrimaryCondition(secondaries []*ConditionEntity, conditionCategoryEntity *ConditionCategoryEntity) PrimaryConditionVo {
	var primaryConditionVo PrimaryConditionVo
	primaryConditionVo.ID = c.ID
	primaryConditionVo.Name = c.ConditionName
	primaryConditionVo.ConditionCategoryVo = GetConditionCategoryVo(conditionCategoryEntity)

	DirectSecondary, Aggravation := InitConditionSecondaries()

	for _, v := range secondaries {
		if v.SecondaryType == Condition_SecondaryType_DirectSecondary {
			DirectSecondary.List = append(DirectSecondary.List, v.ToConditionItem())
		} else if v.SecondaryType == Condition_SecondaryType_Aggravation {
			Aggravation.List = append(Aggravation.List, v.ToConditionItem())
		}
	}
	if len(DirectSecondary.List) > 0 {
		primaryConditionVo.List = append(primaryConditionVo.List, DirectSecondary)
	}
	if len(Aggravation.List) > 0 {
		primaryConditionVo.List = append(primaryConditionVo.List, Aggravation)
	}

	return primaryConditionVo
}

type ConditionUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[ConditionEntity]
}

func NewConditionUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *ConditionUsecase {
	uc := &ConditionUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *ConditionUsecase) UpsertSourceCondition(sourceConditionName string, sourceType int) (e *ConditionEntity, isNew bool, err error) {

	sourceEntity, err := c.GetByCond(And(Eq{"condition_name": sourceConditionName}, In("type", Condition_Type_Source, Condition_Type_Source_From_Ai)))
	if err != nil {
		return nil, false, err
	}
	if sourceEntity == nil {
		sourceEntity = &ConditionEntity{
			Type:          sourceType,
			ConditionName: sourceConditionName,
			CreatedAt:     time.Now().Unix(),
			UpdatedAt:     time.Now().Unix(),
		}
		err = c.CommonUsecase.DB().Save(&sourceEntity).Error
		if err != nil {
			return nil, false, err
		}
		isNew = true
	}
	return sourceEntity, isNew, nil
}

func (c *ConditionUsecase) UpsertPrimaryCondition(conditionName string, conditionCategoryId int32) (e *ConditionEntity, isNew bool, err error) {

	entity, err := c.GetByCond(And(Eq{"condition_name": conditionName}, In("type", Condition_Type_Primary)))
	if err != nil {
		return nil, false, err
	}
	if entity == nil {
		entity = &ConditionEntity{
			Type:                Condition_Type_Primary,
			ConditionName:       conditionName,
			CreatedAt:           time.Now().Unix(),
			UpdatedAt:           time.Now().Unix(),
			ConditionCategoryId: conditionCategoryId,
		}
		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return nil, false, err
		}
		isNew = true
	} else {
		if conditionCategoryId > 0 {
			if entity.ConditionCategoryId != conditionCategoryId {
				entity.ConditionCategoryId = conditionCategoryId
				entity.UpdatedAt = time.Now().Unix()
				err := c.CommonUsecase.DB().Save(&entity).Error
				if err != nil {
					return nil, false, err
				}
			}
		}
	}
	return entity, isNew, nil
}

func (c *ConditionUsecase) UpsertSecondaryCondition(conditionName string, SecondaryType string) (e *ConditionEntity, isNew bool, err error) {

	entity, err := c.GetByCond(And(Eq{"condition_name": conditionName,
		"secondary_type": SecondaryType},
		In("type", Condition_Type_SecondaryCondition)))
	if err != nil {
		return nil, false, err
	}
	if entity == nil {
		entity = &ConditionEntity{
			Type:          Condition_Type_SecondaryCondition,
			SecondaryType: SecondaryType,
			ConditionName: conditionName,
			CreatedAt:     time.Now().Unix(),
			UpdatedAt:     time.Now().Unix(),
		}
		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return nil, false, err
		}
		isNew = true
	}
	return entity, isNew, nil
}

func (c *ConditionUsecase) Upsert(conditionName string) error {
	entity, err := c.GetByCond(Eq{"condition_name": conditionName})
	if err != nil {
		return err
	}
	if entity == nil {
		entity = &ConditionEntity{
			ConditionName: conditionName,
			CreatedAt:     time.Now().Unix(),
			UpdatedAt:     time.Now().Unix(),
		}
	}
	entity.UpdatedAt = time.Now().Unix()
	return c.CommonUsecase.DB().Save(entity).Error
}

func (c *ConditionUsecase) ConditionGet(conditionName string) (*ConditionEntity, error) {
	return c.GetByCond(Eq{"condition_name": conditionName, "type": Condition_Type_Condition})
}

func (c *ConditionUsecase) ConditionUpsert(conditionName string) (*ConditionEntity, error) {
	if conditionName == "" {
		return nil, errors.New("conditionName is empty")
	}
	return c.UpsertNormal(conditionName, Condition_Type_Condition)
}

func (c *ConditionUsecase) UpsertNormal(conditionName string, Type int) (*ConditionEntity, error) {

	entity, err := c.GetByCond(Eq{"condition_name": conditionName, "type": Type})
	if err != nil {
		return nil, err
	}
	if entity == nil {
		entity = &ConditionEntity{
			ConditionName: conditionName,
			Type:          Type,
			CreatedAt:     time.Now().Unix(),
			UpdatedAt:     time.Now().Unix(),
		}
		err = c.CommonUsecase.DB().Save(entity).Error
		if err != nil {
			return nil, err
		}
		return entity, nil
	} else {
		if entity.BizDeletedAt > 0 || entity.DeletedAt > 0 {
			entity.BizDeletedAt = 0
			entity.DeletedAt = 0
			entity.UpdatedAt = time.Now().Unix()
			err = c.CommonUsecase.DB().Save(entity).Error
			if err != nil {
				return nil, err
			}
		}
		return entity, nil
	}
}
