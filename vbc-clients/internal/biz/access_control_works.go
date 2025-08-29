package biz

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	"vbc/lib/uuid"
)

const (
	WorkType_remind_fee_contract_signing = "remind_fee_contract_signing"

	AccessControlWorkPayloadTask_Type_remind_by_email = "remind_by_email"

	AccessControlWorkStatus_waiting = 0
	AccessControlWorkStatus_done    = 1
)

type AccessControlWorkPayload struct {
	Title string
	Tasks []AccessControlWorkPayloadTask
}

func (c AccessControlWorkPayload) GetByIndex(index int32) *AccessControlWorkPayloadTask {
	for k, _ := range c.Tasks {
		if k == int(index) {
			return &c.Tasks[k]
		}
	}
	return nil
}

type AccessControlWorkPayloadTask struct {
	Id     string
	Type   string
	Title  string
	Params string
}

type RemindFeeContractSigningParams struct {
}

func SpawnRemindFeeContractSigningByEmail(params RemindFeeContractSigningParams) AccessControlWorkPayloadTask {
	e := AccessControlWorkPayloadTask{
		Id:     uuid.UuidWithoutStrike(),
		Type:   AccessControlWorkPayloadTask_Type_remind_by_email,
		Params: InterfaceToString(params),
	}
	if e.Type == AccessControlWorkPayloadTask_Type_remind_by_email {
		e.Title = "Notify the client by email"
	}
	return e
}

type AccessControlWorkEntity struct {
	ID        int32 `gorm:"primaryKey"`
	Token     string
	WorkType  string
	RelaId    string
	Status    int
	Payload   string
	ExpiredAt int64
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

func (c *AccessControlWorkEntity) ClientId() int32 {
	r, _ := strconv.ParseInt(c.RelaId, 0, 32)
	return int32(r)
}

func (c *AccessControlWorkEntity) GetPayload() (AccessControlWorkPayload, error) {
	r := lib.StringToT[AccessControlWorkPayload](c.Payload)
	if r.IsOk() {
		return r.Unwrap(), nil
	} else {
		return AccessControlWorkPayload{}, r.Err()
	}
}
func (c AccessControlWorkEntity) TableName() string {
	return "access_control_works"
}

type AccessControlWorkUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[AccessControlWorkEntity]
	TUsecase *TUsecase
}

func NewAccessControlWorkUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase) *AccessControlWorkUsecase {
	uc := &AccessControlWorkUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *AccessControlWorkUsecase) AccessControlWorkPayload(work *AccessControlWorkEntity) (
	accessControlWorkPayload *AccessControlWorkPayload,
	err error) {
	if work == nil {
		return nil, errors.New("work is nil.")
	}
	pl, err := work.GetPayload()
	if err != nil {
		return nil, err
	}
	if work.WorkType == WorkType_remind_fee_contract_signing {
		tClient, err := c.TUsecase.DataById(Kind_client_cases, work.ClientId())
		if err != nil {
			return nil, err
		}
		if tClient == nil {
			return nil, errors.New("tClient is nil")
		}

		pl.Title = fmt.Sprintf("Please remind client (%s, %s) a to sign the contract",
			tClient.CustomFields.TextValueByNameBasic("last_name"), tClient.CustomFields.TextValueByNameBasic("first_name"))

	} else {
		return nil, errors.New("The WorkType does not support.")
	}
	return &pl, nil
}

func (c *AccessControlWorkUsecase) CreateAccessControlWork(workType string, relaId string, payload AccessControlWorkPayload, expiredAt time.Time) (token string, err error) {

	token = uuid.UuidWithoutStrike()
	e := &AccessControlWorkEntity{
		Token:     token,
		WorkType:  workType,
		RelaId:    relaId,
		ExpiredAt: expiredAt.Unix(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	e.Payload = InterfaceToString(payload)

	err = c.CommonUsecase.DB().Create(&e).Error
	if err != nil {
		return "", err
	}
	return
}

func (c *AccessControlWorkUsecase) VerifyAccess(work *AccessControlWorkEntity) error {
	if work == nil {
		return errors.New("work is nil.")
	}
	if work.ExpiredAt <= time.Now().Unix() {
		return errors.New("This assignment has expired.")
	}
	return nil
}
