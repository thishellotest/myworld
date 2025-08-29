package biz

import (
	"bytes"
	"github.com/gen2brain/go-fitz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.uber.org/zap/buffer"
	"image/jpeg"
	"io"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	. "vbc/lib/builder"
)

type SendVa2122aUsecase struct {
	log                   *log.Helper
	conf                  *conf.Data
	CommonUsecase         *CommonUsecase
	BoxcontractUsecase    *BoxcontractUsecase
	BoxUsecase            *BoxUsecase
	BoxbuzUsecase         *BoxbuzUsecase
	GopdfUsecase          *GopdfUsecase
	MapUsecase            *MapUsecase
	ClientEnvelopeUsecase *ClientEnvelopeUsecase
	TUsecase              *TUsecase
	AttorneyUsecase       *AttorneyUsecase
	MiscUsecase           *MiscUsecase
	DataComboUsecase      *DataComboUsecase
}

func NewSendVa2122aUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	BoxcontractUsecase *BoxcontractUsecase,
	BoxUsecase *BoxUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	GopdfUsecase *GopdfUsecase,
	MapUsecase *MapUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	TUsecase *TUsecase,
	AttorneyUsecase *AttorneyUsecase,
	MiscUsecase *MiscUsecase,
	DataComboUsecase *DataComboUsecase,
) *SendVa2122aUsecase {
	uc := &SendVa2122aUsecase{
		log:                   log.NewHelper(logger),
		CommonUsecase:         CommonUsecase,
		conf:                  conf,
		BoxcontractUsecase:    BoxcontractUsecase,
		BoxUsecase:            BoxUsecase,
		BoxbuzUsecase:         BoxbuzUsecase,
		GopdfUsecase:          GopdfUsecase,
		MapUsecase:            MapUsecase,
		ClientEnvelopeUsecase: ClientEnvelopeUsecase,
		TUsecase:              TUsecase,
		AttorneyUsecase:       AttorneyUsecase,
		MiscUsecase:           MiscUsecase,
		DataComboUsecase:      DataComboUsecase,
	}

	return uc
}

func (c *SendVa2122aUsecase) HandleSeparateAmContract(caseId int32, amContractBoxFileId string) error {

	amSignedVA2122aBoxFolderId, err := c.BoxcontractUsecase.AmSignedVA2122aBoxFolderId(caseId)
	if err != nil {
		return err
	}
	if amSignedVA2122aBoxFolderId == "" {
		return errors.New("amSignedVA2122aBoxFolderId is empty")
	}

	amSignedAgreementFolderId, err := c.BoxcontractUsecase.AmSignedAgreementFolderId(caseId)
	if err != nil {
		return err
	}
	if amSignedAgreementFolderId == "" {
		return errors.New("amSignedAgreementFolderId is empty")
	}
	fileReader, err := c.BoxUsecase.DownloadFile(amContractBoxFileId, "")
	if err != nil {
		return err
	}

	doc, err := fitz.NewFromReader(fileReader)
	if err != nil {
		return err
	}
	defer doc.Close()

	var SignedAgreementBytes []byte
	var SignedVA2122aBytes []byte

	var imgs [][]byte
	for n := 0; n < doc.NumPage(); n++ {

		img, err := doc.Image(n)
		if err != nil {
			return err
		}
		var aaac buffer.Buffer
		jpeg.Encode(&aaac, img, &jpeg.Options{jpeg.DefaultQuality})
		imgs = append(imgs, aaac.Bytes())
		if n == 4 {
			SignedAgreementBytes, err = c.GopdfUsecase.CreatePdfFromImg(imgs)
			if err != nil {
				return err
			}
			imgs = nil
		} else if n == 7 {
			SignedVA2122aBytes, err = c.GopdfUsecase.CreatePdfFromImg(imgs)
			if err != nil {
				return err
			}
		}
	}
	// AM - Signed Agreement, AM - Signed VA 21-22a
	fileId, err := c.BoxUsecase.UploadFile(amSignedVA2122aBoxFolderId, bytes.NewReader(SignedVA2122aBytes), "AM - Signed VA 21-22a.pdf")
	if err != nil {
		return err
	}
	err = c.MapUsecase.Set(MapKeyClientCaseAmSignedVA2122aBoxFileId(caseId), fileId)
	if err != nil {
		return err
	}
	fileId1, err := c.BoxUsecase.UploadFile(amSignedAgreementFolderId, bytes.NewReader(SignedAgreementBytes), "AM - Signed Agreement.pdf")
	if err != nil {
		return err
	}
	err = c.MapUsecase.Set(MapKeyClientCaseAmSignedAgreementBoxFileId(caseId), fileId1)
	if err != nil {
		return err
	}
	//lib.DPrintln(fileId)
	//lib.DPrintln(fileId1)
	return nil
}

func (c *SendVa2122aUsecase) RunHandleSeparateAmContract(caseId int32) error {

	amContractBoxFieldId, err := c.ClientEnvelopeUsecase.AmContractBoxFileId(caseId)
	if err != nil {
		return err
	}
	if amContractBoxFieldId == "" {
		return errors.New("HandleSeparateAmContract amContractBoxFieldId is empty")
	}

	key := MapKeyClientCaseAmSignedAgreementBoxFileId(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err = c.HandleSeparateAmContract(caseId, amContractBoxFieldId)
		return err
	}
	return nil
}

func (c *SendVa2122aUsecase) GetAmSignedVA2122aBytes(caseId int32) ([]byte, error) {

	fileId, err := c.MapUsecase.GetForString(MapKeyClientCaseAmSignedVA2122aBoxFileId(caseId))
	if err != nil {
		return nil, err
	}
	if fileId == "" {
		return nil, errors.New("The AM Signed VA21-22a does not exist.")
	}
	reader, err := c.BoxUsecase.DownloadFile(fileId, "")
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (c *SendVa2122aUsecase) Download21Pdf() error {
	cases, err := c.TUsecase.ListByCond(Kind_client_cases, And(Eq{"biz_deleted_at": 0,
		FieldName_ContractSource: ContractSource_AM}, NotIn(FieldName_stages,
		config_vbc.Stages_AmIncomingRequest,
		config_vbc.Stages_AmInformationIntake,
		config_vbc.Stages_AmContractPending,
		config_vbc.Stages_AmTerminated,
		config_vbc.Stages_AmDormant,
	), Expr("deal_name not like '%Test%'")))
	if err != nil {
		return err
	}

	for k, v := range cases {
		err = c.RunHandleSeparateAmContract(v.Id())
		if err != nil {
			c.log.Error(err, v.Id())
			continue
		}
		tClient, _, _ := c.DataComboUsecase.ClientWithCase(*cases[k])
		if tClient == nil {
			c.log.Error("tClient is nil: ", v.Id())
			continue
		}
		aFileName, err := c.MiscUsecase.Gen2122aFileNameForMisc(*cases[k], *tClient)
		if err != nil {
			c.log.Error(err, v.Id())
			continue
		}

		caseId := v.Id()
		key := MapKeyClientCaseAmSignedVA2122aBoxFileId(caseId)
		boxFileId, _ := c.MapUsecase.GetForString(key)
		if boxFileId == "" {
			c.log.Error("boxFileId is empty: ", v.Id())
			continue
		}
		_, err = c.BoxUsecase.CopyFileNewFileNameReturnFileId(boxFileId, aFileName, "334466451846")
		if err != nil {
			c.log.Error(err, v.Id())
			continue
		}

		//attorneyUniqid := v.CustomFields.TextValueByNameBasic(FieldName_attorney_uniqid)
		//
		//attorneyEntity, err := c.AttorneyUsecase.GetByGid(attorneyUniqid)
		//if err != nil {
		//	c.log.Error(err)
		//}
		//if attorneyEntity == nil {
		//	c.log.Error(attorneyEntity.Gid)
		//}

	}
	return nil
}

func (c *SendVa2122aUsecase) HandleHandleMoving2122aFile() error {
	cases, err := c.TUsecase.ListByCond(Kind_client_cases, And(Eq{"biz_deleted_at": 0,
		FieldName_ContractSource: ContractSource_AM}, NotIn(FieldName_stages,
		config_vbc.Stages_AmIncomingRequest,
		config_vbc.Stages_AmInformationIntake,
		config_vbc.Stages_AmContractPending,
		config_vbc.Stages_AmTerminated,
		config_vbc.Stages_AmDormant,
	), Expr("deal_name not like '%Test%'")))
	if err != nil {
		return err
	}

	for _, v := range cases {
		boxFileId, err := c.MiscUsecase.HandleMoving2122aFile(v.Id())
		if err != nil {
			c.log.Error(err, " caseId: ", v.Id())
			continue
		}
		if boxFileId == "" {
			c.log.Error(" caseId: ", v.Id(), " boxFileId is empty")
			continue
		}

		//attorneyUniqid := v.CustomFields.TextValueByNameBasic(FieldName_attorney_uniqid)
		//
		//attorneyEntity, err := c.AttorneyUsecase.GetByGid(attorneyUniqid)
		//if err != nil {
		//	c.log.Error(err)
		//}
		//if attorneyEntity == nil {
		//	c.log.Error(attorneyEntity.Gid)
		//}

	}
	return nil
}
