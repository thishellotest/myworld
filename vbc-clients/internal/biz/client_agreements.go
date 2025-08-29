package biz

import (
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib/adobesign"
	. "vbc/lib/builder"
)

type ClientAgreementEntity struct {
	ID          int32 `gorm:"primaryKey"`
	Name        string
	ClientId    int32
	AgreementId string
	Status      string
	Type        string
	SenderEmail string
	Notes       string
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64
}

func (ClientAgreementEntity) TableName() string {
	return "client_agreements"
}

type ClientAgreementUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[ClientAgreementEntity]
}

func NewClientAgreementUsecase(logger log.Logger, CommonUsecase *CommonUsecase, conf *conf.Data) *ClientAgreementUsecase {

	uc := &ClientAgreementUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *ClientAgreementUsecase) Update(agreement *adobesign.Agreement) error {
	if agreement == nil {
		return errors.New("agreement is nil.")
	}
	entity, err := c.GetByCond(And(Eq{"deleted_at": 0}, Eq{"agreement_id": agreement.Id}))
	if err != nil {
		return err
	}
	if entity == nil {
		entity = &ClientAgreementEntity{
			CreatedAt: time.Now().Unix(),
		}
	}
	entity.Name = agreement.Name
	entity.AgreementId = agreement.Id
	entity.Status = agreement.Status
	entity.Type = agreement.Type
	entity.SenderEmail = agreement.SenderEmail
	entity.Notes = InterfaceToString(agreement)
	entity.UpdatedAt = time.Now().Unix()
	return c.CommonUsecase.DB().Save(&entity).Error
}
