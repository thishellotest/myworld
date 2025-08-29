package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/lib"
	. "vbc/lib/builder"
)

const (
	EsignVendor_docusign = ""
	EsignVendor_box      = "box"
)

const (
	Type_FeeContract      = ""
	Type_MedicalTeamForms = "MedicalTeamForms"
	Type_AmContract       = "AmContract"
	//Type_ReleaseOfInformation = "ReleaseOfInformation"
	//Type_PatientPaymentForm   = "PatientPaymentForm"
)

const (
	ClientEnvelope_IsSigned_Cancelled = 2
	ClientEnvelope_IsSigned_Yes       = 1
	ClientEnvelope_IsSigned_No        = 0
)

type ClientEnvelopeEntity struct {
	ID             int32 `gorm:"primaryKey"`
	EsignVendor    string
	Type           string
	ClientId       int32
	EnvelopeId     string
	Uri            string
	Status         string
	StatusDatetime string
	ResponseText   string
	IsSigned       int
	AttorneyId     int32
	SignStatus     string
	CreatedAt      int64
	UpdatedAt      int64
	DeletedAt      int64
}

func (ClientEnvelopeEntity) TableName() string {
	return "client_envelopes"
}

func (c *ClientEnvelopeEntity) BoxContactFileId() string {
	ResponseTextMap := lib.ToTypeMapByString(c.ResponseText)
	list := ResponseTextMap.GetTypeList("sign_files.files")
	for _, v := range list {
		fileId := v.GetString("id")
		if fileId != "" {
			return fileId
		}
	}
	return ""
}

func (c *ClientEnvelopeEntity) ContractDateOn() (time.Time, error) {
	if c.CreatedAt <= 0 {
		return time.Unix(c.CreatedAt, 0), errors.New("CreateAt is wrong")
	}
	return time.Unix(c.CreatedAt, 0), nil
}

type ClientEnvelopeUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	DBUsecase[ClientEnvelopeEntity]
}

func NewClientEnvelopeUsecase(logger log.Logger, CommonUsecase *CommonUsecase) *ClientEnvelopeUsecase {
	uc := &ClientEnvelopeUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *ClientEnvelopeUsecase) Add(clientCaseId int32, esignVendor string, EnvelopeId string, responseText string, typ string, attorneyId int32) error {
	entity := &ClientEnvelopeEntity{
		EsignVendor:  esignVendor,
		ClientId:     clientCaseId,
		EnvelopeId:   EnvelopeId,
		ResponseText: responseText,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		AttorneyId:   attorneyId,
		Type:         typ,
	}
	return c.CommonUsecase.DB().Save(entity).Error
}

func (c *ClientEnvelopeUsecase) GetByEnvelopeId(esignVendor string, envelopeId string) (*ClientEnvelopeEntity, error) {
	return c.GetByCond(Eq{"esign_vendor": esignVendor, "envelope_id": envelopeId})
}

func (c *ClientEnvelopeUsecase) GetBoxSignByCaseId(caseId int32, typ string) (*ClientEnvelopeEntity, error) {
	return c.GetByCondWithOrderBy(Eq{"client_id": caseId, "type": typ, "deleted_at": 0, "esign_vendor": "box"}, "id desc")
	//return c.GetByCond(Eq{"client_id": caseId, "type": typ, "deleted_at": 0, "esign_vendor": "box"})
}

// ContractDateOn 获取合同信息 Jan. 2, 2006
func (c *ClientEnvelopeUsecase) ContractDateOn(tCaseId int32, isAmContract bool) (string, error) {

	typ := Type_FeeContract
	if isAmContract {
		typ = Type_AmContract
	}
	boxSign, err := c.GetBoxSignByCaseId(tCaseId, typ)
	if err != nil {
		c.log.Error(err)
		return "", err
	}
	if boxSign == nil {
		return "", errors.New("boxSign is nil")
	}
	contractDateOn, err := boxSign.ContractDateOn()
	if err != nil {
		c.log.Error(err)
		return "", err
	}
	contractDate := contractDateOn.In(configs.GetVBCDefaultLocation()).Format("January 2, 2006")
	return contractDate, nil
}

func (c *ClientEnvelopeUsecase) AmContractBoxFileId(caseId int32) (boxFileId string, err error) {

	clientEnvelope, err := c.GetByCond(Eq{"deleted_at": 0, "esign_vendor": EsignVendor_box, "type": Type_AmContract, "client_id": caseId, "is_signed": ClientEnvelope_IsSigned_Yes})
	if err != nil {
		return "", err
	}
	if clientEnvelope == nil {
		return "", errors.New("clientEnvelope is nil")
	}
	typeMap := lib.ToTypeMapByString(clientEnvelope.ResponseText)
	typeList := typeMap.GetTypeList("sign_files.files")
	if len(typeList) > 0 {
		return typeList[0].GetString("id"), nil
	}
	return "", nil
}
