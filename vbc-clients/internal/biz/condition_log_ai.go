package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	ConditionLogAi_LogType_ConditionSource     = 0 // condition的来源；使用from_value
	ConditionLogAi_LogType_ConditionCategories = 1 // condition的分类来源；DestId存的是分类ID

	ConditionLogAi_FromType_PromptAi_Condition = 0 // 来源类型 0: 使用prompt_key ai from_value：来源于表conditions type=0
	ConditionLogAi_FromType_Import             = 1 // 导入（excel）
	ConditionLogAi_FromType_Manual             = 2 // 人工操作
	ConditionLogAi_FromType_ManualDelete       = 3 // 人工操作删除
)

type ConditionLogAiEntity struct {
	ID           int32 `gorm:"primaryKey"`
	ConditionId  int32
	Count        int
	PromptKey    string
	DestId       string
	LogType      int
	FromValue    string
	FromType     int
	CreatedBy    string
	BizDeletedAt int64
	CreatedAt    int64
	UpdatedAt    int64
}

func (ConditionLogAiEntity) TableName() string {
	return "condition_log_ai"
}

func (c *ConditionLogAiEntity) FromValueForInt32() int32 {
	r, _ := strconv.ParseInt(c.FromValue, 10, 32)
	return int32(r)
}

type ConditionLogAiVo struct {
	FromId      string `json:"from_id"`
	FromContent string `json:"from_content"` // 来源内容
}

func (c *ConditionLogAiEntity) ToApi(log *log.Helper, ConditionUsecase *ConditionUsecase) (conditionLogAiVo *ConditionLogAiVo) {

	if c.FromType == ConditionLogAi_FromType_PromptAi_Condition {
		fromConditionSourceId := c.FromValueForInt32()
		if fromConditionSourceId > 0 {
			entity, err := ConditionUsecase.GetByCond(Eq{"id": fromConditionSourceId})
			if entity != nil {
				return &ConditionLogAiVo{
					FromId:      InterfaceToString(entity.ID),
					FromContent: entity.ConditionName,
				}
			} else {
				log.Warn(err, c.FromValue, c.ID)
			}
		}
	}
	return nil
}

type ConditionLogAiUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[ConditionLogAiEntity]
}

func NewConditionLogAiUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *ConditionLogAiUsecase {
	uc := &ConditionLogAiUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

// AddConditionSource 添加condition的来源
func (c *ConditionLogAiUsecase) AddConditionSource(conditionId int32, fromType int, fromValue string, createdBy string, promptKey string) (*ConditionLogAiEntity, error) {
	return c.UpsertCenter(ConditionLogAi_LogType_ConditionSource, conditionId, "", fromType, fromValue, createdBy, promptKey)
}

// AddCategoryOfConditionSource 添加condition的分类来源
func (c *ConditionLogAiUsecase) AddCategoryOfConditionSource(conditionId int32, conditionCategoryId int32, fromType int, fromValue string, createdBy string) (*ConditionLogAiEntity, error) {
	return c.UpsertCenter(ConditionLogAi_LogType_ConditionCategories, conditionId, InterfaceToString(conditionCategoryId), fromType, fromValue, createdBy, "")
}

func (c *ConditionLogAiUsecase) UpsertCenter(logType int, conditionId int32, destId string, fromType int, fromValue string, createdBy string, promptKey string) (*ConditionLogAiEntity, error) {

	entity, err := c.GetByCond(Eq{"condition_id": conditionId,
		"from_value": fromValue,
		"from_type":  fromType,
		"prompt_key": promptKey,
		"log_type":   logType,
		"dest_id":    destId,
	})
	if err != nil {
		return nil, err
	}
	if entity == nil {
		entity = &ConditionLogAiEntity{
			ConditionId: conditionId,
			DestId:      destId,
			LogType:     logType,
			FromType:    fromType,
			FromValue:   fromValue,
			PromptKey:   promptKey,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now().Unix(),
		}
	}
	entity.Count += 1
	entity.UpdatedAt = time.Now().Unix()

	err = c.CommonUsecase.DB().Save(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// AddLogConditionSourceFromImport 添加condition来源， 从导入
func (c *ConditionLogAiUsecase) AddLogConditionSourceFromImport(conditionId int32, fileName string, createdBy string) error {
	_, err := c.AddConditionSource(conditionId, ConditionLogAi_FromType_Import, fileName, createdBy, "")
	return err
}

// AddLogPromptAiCondition 添加condition来源， 从ai
func (c *ConditionLogAiUsecase) AddLogPromptAiCondition(conditionId int32, fromSourceConditionId int32, promptKey string) error {

	_, err := c.AddConditionSource(conditionId, ConditionLogAi_FromType_PromptAi_Condition, InterfaceToString(fromSourceConditionId), "", promptKey)
	return err
}
