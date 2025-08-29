package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type EnvelopeStatusChangeEntity struct {
	ID                          int32 `gorm:"primaryKey"`
	EnvelopeId                  string
	CreatedDatetime             string
	AttachmentsUri              string
	CertificateUri              string
	CustomFieldsUri             string
	DocumentsCombinedUri        string
	EnvelopeLocation            string
	DocumentsUri                string
	IsSignatureProviderEnvelope string
	LastModifiedDatetime        string
	NotificationUri             string
	PurgeState                  string
	RecipientsUri               string
	SenderAccountId             string
	SenderEmail                 string
	SenderUserId                string
	SenderUsername              string
	SentDatetime                string
	SigningLocation             string
	Status                      string
	StatusChangedDatetime       string
	TemplatesUri                string
	EmailBlurb                  string
	EmailSubject                string
	CreatedAt                   int64
	UpdatedAt                   int64
}

func (EnvelopeStatusChangeEntity) TableName() string {
	return "envelope_status_changes"
}

type EnvelopeStatusChangeUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[EnvelopeStatusChangeEntity]
	MapUsecase            *MapUsecase
	ClientEnvelopeUsecase *ClientEnvelopeUsecase
	TUsecase              *TUsecase
	BoxUsecase            *BoxUsecase
}

func NewEnvelopeStatusChangeUsecase(logger log.Logger, CommonUsecase *CommonUsecase, MapUsecase *MapUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	TUsecase *TUsecase,
	BoxUsecase *BoxUsecase,
	conf *conf.Data) *EnvelopeStatusChangeUsecase {

	uc := &EnvelopeStatusChangeUsecase{
		log:                   log.NewHelper(logger),
		CommonUsecase:         CommonUsecase,
		MapUsecase:            MapUsecase,
		ClientEnvelopeUsecase: ClientEnvelopeUsecase,
		TUsecase:              TUsecase,
		BoxUsecase:            BoxUsecase,
		conf:                  conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

// RunEnvelopeStatusChangeJob Job
func (c *EnvelopeStatusChangeUsecase) RunEnvelopeStatusChangeJob(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("EnvelopeStatusChangeUsecase:RunEnvelopeStatusChangeJob:Done")
				return
			default:
				divideId, err := c.MapUsecase.GetForInt(Map_EnvelopeStatusChange_divide)
				if err != nil {
					c.log.Error(err)
				} else {
					sqlRows, err := c.CommonUsecase.DB().Table(EnvelopeStatusChangeEntity{}.TableName()).
						Where("id>?",
							divideId).Rows()
					if err != nil {
						c.log.Error(err)
					} else {
						if sqlRows != nil {
							newDivideId := int32(0)
							for sqlRows.Next() {
								var entity EnvelopeStatusChangeEntity
								err = c.CommonUsecase.DB().ScanRows(sqlRows, &entity)
								if err != nil {
									c.log.Error(err)
								} else {
									newDivideId = entity.ID
									err = c.Handle(&entity)
									if err != nil {
										c.log.Error(err)
									}
									c.MapUsecase.Set(Map_EnvelopeStatusChange_divide, lib.InterfaceToString(newDivideId))
								}
							}
							err = sqlRows.Close()
							if err != nil {
								c.log.Error(err)
							}
						}
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()
	return nil
}

// Handle 处理事件
func (c *EnvelopeStatusChangeUsecase) Handle(entity *EnvelopeStatusChangeEntity) error {
	if entity == nil {
		return errors.New("EnvelopeStatusChangeEntity is nil.")
	}
	if entity.Status == "completed" {
		//err := c.HandleCompleted(entity)
		typeMap := make(lib.TypeMap)
		typeMap.Set("EnvelopeStatusChangeId", entity.ID)

		clientEnvelope, _ := c.ClientEnvelopeUsecase.GetByCond(And(Eq{"envelope_id": entity.EnvelopeId},
			Eq{"esign_vendor": EsignVendor_docusign}))
		incrId := int32(0)
		if clientEnvelope != nil {
			incrId = clientEnvelope.ClientId
		}
		err := c.CommonUsecase.DB().Save(&TaskEntity{
			IncrId:    incrId,
			TaskInput: typeMap.ToString(),
			Event:     Task_Dag_BoxCreateClientContracts,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}).Error

		if err != nil {
			a := GenLog(entity.ID, Log_FromType_Envelope_completed, err.Error())
			err = c.CommonUsecase.DB().Save(a).Error
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}
