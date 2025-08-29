package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	RecordLog_RecordType_CRMCases = "crm_cases"
	RecordLog_BizType_Stages      = "crm_stages"
)

type RecordLogEntity struct {
	ID         int32 `gorm:"primaryKey"`
	RecordType string
	RecordId   string
	BizType    string
	BizValue   string
	StartTime  int64
	EndTime    int64
	CloseTime  int64
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  int64
	CreatedBy  string
	ModifiedBy string
}

func (RecordLogEntity) TableName() string {
	return "record_log"
}

type RecordLogUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[RecordLogEntity]
}

func NewRecordLogUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *RecordLogUsecase {
	uc := &RecordLogUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *RecordLogUsecase) UpdateBizCrmStages(caseGid string, userGid string, endTime int64) error {
	entity, err := c.BizCrmStagesLatest(caseGid)
	if err != nil {
		return err
	}
	if entity != nil {
		entity.EndTime = endTime
		entity.UpdatedAt = time.Now().Unix()
		entity.ModifiedBy = userGid
		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RecordLogUsecase) CloseBizCrmStages(caseGid string, userGid string) error {
	entity, err := c.BizCrmStagesLatest(caseGid)
	if err != nil {
		return err
	}
	if entity != nil {
		entity.CloseTime = time.Now().Unix()
		entity.UpdatedAt = time.Now().Unix()
		entity.ModifiedBy = userGid
		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RecordLogUsecase) AddBizCrmStages(caseGid string, stageValue string, startTime int64, endTime int64, userGid string) (*RecordLogEntity, error) {

	return c.Add(RecordLog_RecordType_CRMCases, caseGid, RecordLog_BizType_Stages, stageValue, startTime, endTime, userGid)
}

func (c *RecordLogUsecase) BizCrmStagesLatest(caseGid string) (*RecordLogEntity, error) {
	return c.GetByCondWithOrderBy(Eq{
		"record_type": RecordLog_RecordType_CRMCases,
		"record_id":   caseGid,
		"biz_type":    RecordLog_BizType_Stages,
		//"biz_value":   stageValue,
		"deleted_at": 0,
	}, "id desc")
}

func (c *RecordLogUsecase) Add(RecordType string, RecordId string, BizType string, BizValue string, StartTime int64, EndTime int64, userGid string) (*RecordLogEntity, error) {
	entity := &RecordLogEntity{
		RecordType: RecordType,
		RecordId:   RecordId,
		BizType:    BizType,
		BizValue:   BizValue,
		StartTime:  StartTime,
		EndTime:    EndTime,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
		CreatedBy:  userGid,
		ModifiedBy: userGid,
	}
	err := c.CommonUsecase.DB().Save(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
