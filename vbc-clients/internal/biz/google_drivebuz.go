package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"strings"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/lib"
)

type GoogleDrivebuzUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	GoogleDriveUsecase *GoogleDriveUsecase
	BoxbuzUsecase      *BoxbuzUsecase
	BoxUsecase         *BoxUsecase
	LogUsecase         *LogUsecase
	MapUsecase         *MapUsecase
}

func NewGoogleDrivebuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	GoogleDriveUsecase *GoogleDriveUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	BoxUsecase *BoxUsecase,
	LogUsecase *LogUsecase,
	MapUsecase *MapUsecase,
) *GoogleDrivebuzUsecase {
	uc := &GoogleDrivebuzUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		GoogleDriveUsecase: GoogleDriveUsecase,
		BoxbuzUsecase:      BoxbuzUsecase,
		BoxUsecase:         BoxUsecase,
		LogUsecase:         LogUsecase,
		MapUsecase:         MapUsecase,
	}
	return uc
}

func GoogleDrivePaymentFolderName(tClient *TData) (string, error) {
	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	return fmt.Sprintf("%s, %s", lastName, firstName), nil
}

func (c *GoogleDrivebuzUsecase) TransferPaymentForm(ctx context.Context, tCase *TData, tClient *TData) (drivePaymentFormFile *drive.File, err error) {
	if tClient == nil {
		return nil, errors.New("tClient is nil")
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}

	DCPrivateExamsFolderId, err := c.BoxbuzUsecase.DCPrivateExamsFolderId(tCase)
	if err != nil {
		return nil, err
	}
	if DCPrivateExamsFolderId == "" {
		return nil, errors.New("DCPrivateExamsFolderId is empty")
	}

	PatientPaymentFormFileId, err := c.BoxbuzUsecase.GetBoxResIdByCase(DCPrivateExamsFolderId,
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_PE_PatientPaymentForm_File),
		tCase, tClient)
	if err != nil {
		return nil, err
	}
	if PatientPaymentFormFileId == "" {
		return nil, errors.New("PatientPaymentFormFileId is empty")
	}

	paymentsFolderId := c.conf.GoogleDrive.PaymentsFolderId
	folderName, err := GoogleDrivePaymentFolderName(tClient)
	if err != nil {
		return nil, err
	}

	driveFile, err := c.GoogleDriveUsecase.CreateFolder(ctx, paymentsFolderId, folderName)
	if err != nil {
		return nil, err
	}
	if driveFile == nil {
		return nil, errors.New("driveFile is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")
	er := c.MapUsecase.Set(MapKeyGoogleDrivePaymentFolderId(caseId), driveFile.Id)
	if er != nil {
		c.log.Error(er)
	}

	patientPaymentFormFileName, err := PatientPaymentFormFileName(tClient)
	if err != nil {
		return nil, err
	}
	PatientPaymentFormFileRes, err := c.BoxUsecase.DownloadFile(PatientPaymentFormFileId, "")
	if err != nil {
		return nil, err
	}
	defer PatientPaymentFormFileRes.Close()

	drivePaymentFormFile, err = c.GoogleDriveUsecase.UploadFile(ctx, driveFile.Id, patientPaymentFormFileName, PatientPaymentFormFileRes)
	if err != nil {
		return nil, err
	}
	if drivePaymentFormFile == nil {
		return nil, errors.New("drivePaymentFormFile is nil")
	}
	return drivePaymentFormFile, nil
}

// https://veteranbenefitscenter.app.box.com/folder/264219359597
// https://veteranbenefitscenter.app.box.com/folder/263407768664
// 1vHjMD6PNmnKDDqGHuCl-u_x2bdBlZQrX
// https://drive.google.com/drive/folders/1egqWn_OtAOD3iXBje1KroaR1GNDgAgbj
// MEDICAL FILES FOR DOCTOR

func (c *GoogleDrivebuzUsecase) TransferGeneral(ctx context.Context, tCase *TData, tClient *TData) error {

	GeneralFolderId, err := c.BoxbuzUsecase.FolderIdDC_PE_General(tCase, tClient)
	if err != nil {
		return err
	}
	if GeneralFolderId == "" {
		c.log.Debug("GeneralFolderId is empty")
		return nil
	}

	GeneralItems, err := c.BoxUsecase.ListItemsInFolderFormat(GeneralFolderId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if len(GeneralItems) == 0 {
		c.log.Debug("GeneralItems length is 0")
		return nil
	}

	caseId := tCase.CustomFields.NumberValueByNameBasic("id")
	er := c.LogUsecase.SaveLog(caseId,
		Log_FormType_TransferGeneral, map[string]interface{}{
			"GeneralFolderId": GeneralFolderId,
		},
	)
	if er != nil {
		c.log.Error(er)
	}

	folderName, err := GenFolderNameByClient(tClient)
	if err != nil {
		return err
	}

	driveFolder, err := c.GoogleDriveUsecase.CreateFolder(ctx, c.conf.GoogleDrive.MedicalEvaluationsFolderId, folderName)
	if err != nil {
		return err
	}
	if driveFolder == nil {
		return errors.New("driveFolder is nil")
	}
	er = c.MapUsecase.Set(MapKeyGoogleDriveGeneralFolderId(caseId), driveFolder.Id)
	if er != nil {
		c.log.Error(er)
	}

	//medicalDoctorFolder, err := c.GoogleDriveUsecase.CreateFolder(ctx, driveFolder.Id, config_googledrive.MEDICAL_FILES_FOR_DOCTOR)
	//if err != nil {
	//	return err
	//}
	//if medicalDoctorFolder == nil {
	//	return errors.New("medicalDoctorFolder is nil")
	//}

	return c.TransferFiles(ctx, driveFolder.Id, GeneralItems, tCase)
}

func (c *GoogleDrivebuzUsecase) TransferPsych(ctx context.Context, tCase *TData, tClient *TData) error {

	PsychFolderId, err := c.BoxbuzUsecase.FolderIdDC_PE_Psych(tCase, tClient)
	if err != nil {
		return err
	}
	if PsychFolderId == "" {
		c.log.Debug("PsychFolderId is empty")
		return nil
	}

	psychItems, err := c.BoxUsecase.ListItemsInFolderFormat(PsychFolderId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if len(psychItems) == 0 {
		c.log.Debug("psychItems length is 0")
		return nil
	}

	caseId := tCase.CustomFields.NumberValueByNameBasic("id")
	er := c.LogUsecase.SaveLog(caseId,
		Log_FormType_TransferPsych, map[string]interface{}{
			"PsychFolderId": PsychFolderId,
		},
	)
	if er != nil {
		c.log.Error(er)
	}

	folderName, err := GenFolderNameByClient(tClient)
	if err != nil {
		return err
	}

	driveFolder, err := c.GoogleDriveUsecase.CreateFolder(ctx, c.conf.GoogleDrive.PsychEvaluationsFolderId, folderName)
	if err != nil {
		return err
	}
	if driveFolder == nil {
		return errors.New("driveFolder is nil")
	}
	er = c.MapUsecase.Set(MapKeyGoogleDrivePsychFolderId(caseId), driveFolder.Id)
	if er != nil {
		c.log.Error(er)
	}

	//psychDoctorFolder, err := c.GoogleDriveUsecase.CreateFolder(ctx, driveFolder.Id, config_googledrive.PSYCH_FILES_FOR_DOCTOR)
	//if err != nil {
	//	return err
	//}
	//if psychDoctorFolder == nil {
	//	return errors.New("psychDoctorFolder is nil")
	//}

	return c.TransferFiles(ctx, driveFolder.Id, psychItems, tCase)
}

func GenFolderNameByClient(tClient *TData) (string, error) {
	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	return fmt.Sprintf("%s, %s", tClient.CustomFields.TextValueByNameBasic(FieldName_last_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)), nil
}

func (c *GoogleDrivebuzUsecase) TransferFiles(ctx context.Context, googleFolderId string, boxItems lib.TypeList, tCase *TData) error {

	if len(boxItems) == 0 {
		return errors.New("boxItems length is 0")
	}

	for _, v := range boxItems {
		if v.GetString("type") == string(config_box.BoxResType_file) {
			name := v.GetString("name")
			fileRes, err := c.BoxUsecase.DownloadFile(v.GetString("id"), "")
			if err != nil {
				return err
			}
			defer fileRes.Close()
			_, err = c.GoogleDriveUsecase.UploadFile(ctx, googleFolderId, name, fileRes)
			if err != nil {
				return err
			}
		}
	}

	DCPrivateExamsFolderId, err := c.BoxbuzUsecase.DCPrivateExamsFolderId(tCase)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if DCPrivateExamsFolderId == "" {
		return errors.New("DCPrivateExamsFolderId is empty")
	}
	peItems, err := c.BoxUsecase.ListItemsInFolderFormat(DCPrivateExamsFolderId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	for _, v := range peItems {

		if v.GetString("type") == string(config_box.BoxResType_file) {
			name := v.GetString("name")

			if strings.Index(name, PatientPaymentForm_Postfix) >= 0 {
				continue
			}

			fileRes, err := c.BoxUsecase.DownloadFile(v.GetString("id"), "")
			if err != nil {
				c.log.Error(err)
				return err
			}
			defer fileRes.Close()
			_, err = c.GoogleDriveUsecase.UploadFile(ctx, googleFolderId, name, fileRes)
			if err != nil {
				c.log.Error(err)
				return err
			}
		}
	}
	return nil
}
