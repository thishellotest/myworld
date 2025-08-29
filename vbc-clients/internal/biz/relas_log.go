package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	RelasLog_Type_condition_2_jotform     = "condition_2_jotform"     // 老版本不使用 SourceId: conditions.id; TargetId: jotform.form_id
	RelasLog_Type_condition_2_jotform_new = "condition_2_jotform_new" // 新版本
)

type RelasLogEntity struct {
	ID        int32 `gorm:"primaryKey"`
	Type      string
	SourceId  string
	TargetId  string
	Sort      int
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

func (RelasLogEntity) TableName() string {
	return "relas_log"
}

type RelasLogUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[RelasLogEntity]
}

func NewRelasLogUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *RelasLogUsecase {
	uc := &RelasLogUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *RelasLogUsecase) AllBySourceId(Type string, sourceId string) ([]*RelasLogEntity, error) {
	return c.AllByCond(Eq{"type": Type, "source_id": sourceId, "deleted_at": 0})
}

func (c *RelasLogUsecase) ConditionRelas(sourceId string) ([]*RelasLogEntity, error) {
	if configs.NewPSGen {
		return c.AllByCond(Eq{"type": RelasLog_Type_condition_2_jotform_new,
			"source_id": sourceId, "deleted_at": 0})
	}
	return c.AllByCond(Eq{"type": RelasLog_Type_condition_2_jotform,
		"source_id": sourceId, "deleted_at": 0})
}

func (c *RelasLogUsecase) ConditionUpsert(sourceId string, targetId string, sort *int) (*RelasLogEntity, error) {
	if configs.NewPSGen {
		return c.Upsert(RelasLog_Type_condition_2_jotform_new, sourceId, targetId, sort)
	}
	return c.Upsert(RelasLog_Type_condition_2_jotform, sourceId, targetId, sort)
}

func (c *RelasLogUsecase) ConditionRemoveOtherTargetIds(sourceId string, retainTargetIds []string) error {
	if configs.NewPSGen {
		c.RemoveOtherTargetIds(RelasLog_Type_condition_2_jotform_new, sourceId, retainTargetIds)
	}
	return c.RemoveOtherTargetIds(RelasLog_Type_condition_2_jotform, sourceId, retainTargetIds)
}
func (c *RelasLogUsecase) RemoveOtherTargetIds(Type string, sourceId string, retainTargetIds []string) error {
	return c.UpdatesByCond(map[string]interface{}{"deleted_at": time.Now().Unix()}, And(Eq{"deleted_at": 0, "type": Type, "source_id": sourceId}, NotIn("target_id", retainTargetIds)))
}

func (c *RelasLogUsecase) ConditionExists(sourceId string) (bool, error) {
	if configs.NewPSGen {
		return c.Exists(RelasLog_Type_condition_2_jotform_new, sourceId)
	}
	return c.Exists(RelasLog_Type_condition_2_jotform, sourceId)
}

func (c *RelasLogUsecase) Exists(Type string, sourceId string) (bool, error) {
	a, err := c.GetByCond(Eq{"deleted_at": 0, "type": Type, "source_id": sourceId})
	if err != nil {
		return false, err
	}
	if a == nil {
		return false, nil
	}
	return true, nil
}

func (c *RelasLogUsecase) Upsert(Type string, sourceId string, targetId string, sort *int) (*RelasLogEntity, error) {
	entity, err := c.GetByCond(Eq{"type": Type, "source_id": sourceId, "target_id": targetId})
	if err != nil {
		return nil, err
	}
	if entity != nil {
		//if entity.DeletedAt > 0 {
		entity.DeletedAt = 0
		if sort != nil {
			entity.Sort = *sort
		}
		entity.UpdatedAt = time.Now().Unix()
		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return nil, err
		}
		//}
		return entity, nil
	} else {
		entity = &RelasLogEntity{
			Type:     Type,
			SourceId: sourceId,
			TargetId: targetId,

			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		if sort != nil {
			entity.Sort = *sort
		}
		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return nil, err
		}
		return entity, nil
	}
}
