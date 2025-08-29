package biz

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type ActionOnceHandleCopyMedicalTeamFormsUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	MapUsecase    *MapUsecase
	TUsecase      *TUsecase
	//ZohoUsecase           *ZohoUsecase
	//FeeUsecase            *FeeUsecase
	//ClientCaseUsecase     *ClientCaseUsecase
	BoxUsecase       *BoxUsecase
	DataComboUsecase *DataComboUsecase
	BoxbuzUsecase    *BoxbuzUsecase
	//DbqsUsecase           *DbqsUsecase
	//BoxcontractUsecase    *BoxcontractUsecase
	ClientEnvelopeUsecase *ClientEnvelopeUsecase
	//RollpoingUsecase      *RollpoingUsecase
	PdfcpuUsecase *PdfcpuUsecase
}

func NewActionOnceHandleCopyMedicalTeamFormsUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	TUsecase *TUsecase,
//ZohoUsecase *ZohoUsecase,
//FeeUsecase *FeeUsecase,
//ClientCaseUsecase *ClientCaseUsecase,
	BoxUsecase *BoxUsecase,
	DataComboUsecase *DataComboUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
//DbqsUsecase *DbqsUsecase,
//BoxcontractUsecase *BoxcontractUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
//RollpoingUsecase *RollpoingUsecase,
	PdfcpuUsecase *PdfcpuUsecase) *ActionOnceHandleCopyMedicalTeamFormsUsecase {
	uc := &ActionOnceHandleCopyMedicalTeamFormsUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		MapUsecase:    MapUsecase,
		TUsecase:      TUsecase,
		//ZohoUsecase:           ZohoUsecase,
		//FeeUsecase:            FeeUsecase,
		//ClientCaseUsecase:     ClientCaseUsecase,
		BoxUsecase:       BoxUsecase,
		DataComboUsecase: DataComboUsecase,
		BoxbuzUsecase:    BoxbuzUsecase,
		//DbqsUsecase:           DbqsUsecase,
		//BoxcontractUsecase:    BoxcontractUsecase,
		ClientEnvelopeUsecase: ClientEnvelopeUsecase,
		//RollpoingUsecase:      RollpoingUsecase,
		PdfcpuUsecase: PdfcpuUsecase,
	}
	return uc
}

// HandleCopyMedicalTeamForms 拷贝指定合同文件到指定目录
func (c *ActionOnceHandleCopyMedicalTeamFormsUsecase) HandleCopyMedicalTeamForms(clientCaseId int32) error {
	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "CopyMedicalTeamForms", clientCaseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		tClientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": clientCaseId})
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("HandleCopyMedicalTeamForms: tClientCase  is nil")
		}
		tClient, _, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
		if err != nil {
			return err
		}
		if tClient == nil {
			return errors.New("HandleCopyMedicalTeamForms: tClient is nil")
		}

		DCPrivateExamsFolderId, err := c.BoxbuzUsecase.GetDCSubFolderId(MapKeyBuildAutoBoxDCPrivateExamsFolderId(
			tClientCase.CustomFields.NumberValueByNameBasic("id")), tClientCase)
		if err != nil {
			return err
		}
		if DCPrivateExamsFolderId == "" {
			return errors.New("DCPrivateExamsFolderId is empty")
		}
		roiContractId, err := c.MapUsecase.GetForString(MapKeyMedicalTeamForms(clientCaseId))
		if err != nil {
			return err
		}
		if roiContractId == "" {
			return errors.New("roiContractId is empty")
		}
		envelopeEntity, err := c.ClientEnvelopeUsecase.GetByEnvelopeId(EsignVendor_box, roiContractId)
		if err != nil {
			return err
		}
		if envelopeEntity == nil {
			return errors.New("envelopeEntity is nil")
		}
		boxContactFileId := envelopeEntity.BoxContactFileId()
		if boxContactFileId == "" {
			return errors.New("boxContactFileId is empty.")
		}

		PatientPaymentFormFileId, ReleaseOfInformationFormFileId, err := c.BizCopyMedicalTeamForms(boxContactFileId, tClient, DCPrivateExamsFolderId)
		if err != nil {
			return err
		}

		c.MapUsecase.Set(MapKeyBuildAutoBoxDCPatientPaymentFormFileId(clientCaseId), PatientPaymentFormFileId)
		c.MapUsecase.Set(MapKeyBuildAutoBoxDCReleaseOfInformationFormFileId(clientCaseId), ReleaseOfInformationFormFileId)

		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ActionOnceHandleCopyMedicalTeamFormsUsecase) BizCopyMedicalTeamForms(boxContactFileId string, tClient *TData, DCPrivateExamsFolderId string) (PatientPaymentFormFileId, ReleaseOfInformationFormFileId string, err error) {

	contactFileBodyReader, err := c.BoxUsecase.DownloadFile(boxContactFileId, "")
	if err != nil {
		return "", "", err
	}
	if contactFileBodyReader == nil {
		return "", "", errors.New("contactFileBodyReader is nil")
	}
	defer contactFileBodyReader.Close()
	contactBytes, err := io.ReadAll(contactFileBodyReader)
	if err != nil {
		return "", "", err
	}
	contactBytesReader := bytes.NewReader(contactBytes)
	result, err := c.PdfcpuUsecase.SplitPdfAndCombine(contactBytesReader, []*SplitPdfAndCombineConf{
		{
			PageBegin: 1,
			PageEnd:   3,
		},
		{
			PageBegin: 4,
			PageEnd:   5, // 第5页有可能不存在
		},
	})
	if err != nil {
		return "", "", err
	}
	if len(result) != 2 {
		return "", "", errors.New("BizCopyMedicalTeamForms: result length is wrong")
	}

	//fullName := FormatFullName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
	//	tClient.CustomFields.TextValueByNameBasic(FieldName_last_name))

	//PatientPaymentForm := fullName + " - Patient Payment Form.pdf"
	//ReleaseOfInformationForm := fullName + " - Release of Information Form.pdf"

	PatientPaymentForm, err := PatientPaymentFormFileName(tClient)
	if err != nil {
		return "", "", err
	}
	ReleaseOfInformationForm, err := ReleaseOfInformationFormFileName(tClient)
	if err != nil {
		return "", "", err
	}

	PatientPaymentFormFileId, err = c.BoxUsecase.UploadFile(DCPrivateExamsFolderId, result[0], PatientPaymentForm)
	if err != nil {
		return "", "", err
	}
	ReleaseOfInformationFormFileId, err = c.BoxUsecase.UploadFile(DCPrivateExamsFolderId, result[1], ReleaseOfInformationForm)
	if err != nil {
		return "", "", err
	}
	return PatientPaymentFormFileId, ReleaseOfInformationFormFileId, nil
}

const PatientPaymentForm_Postfix = " - Patient Payment Form.pdf"
const ReleaseOfInformationForm_Postfix = " - Release of Information Form.pdf"

func PatientPaymentFormFileName(tClient *TData) (fileName string, err error) {

	if tClient == nil {
		return "", errors.New("PatientPaymentFormFileName: tClient is nil")
	}
	fullName := FormatFullName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name))

	PatientPaymentForm := fullName + PatientPaymentForm_Postfix
	return PatientPaymentForm, nil
}

func ReleaseOfInformationFormFileName(tClient *TData) (fileName string, err error) {
	if tClient == nil {
		return "", errors.New("ReleaseOfInformationFormFileName: tClient is nil")
	}
	fullName := FormatFullName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name))
	ReleaseOfInformationForm := fullName + ReleaseOfInformationForm_Postfix
	return ReleaseOfInformationForm, nil
}
