package biz

import (
	"fmt"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/lib"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type BoxbuzUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	BoxUsecase       *BoxUsecase
	MapUsecase       *MapUsecase
	TUsecase         *TUsecase
	DataComboUsecase *DataComboUsecase
	FeeUsecase       *FeeUsecase
	WordUsecase      *WordUsecase
}

func NewBoxbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	BoxUsecase *BoxUsecase,
	MapUsecase *MapUsecase,
	TUsecase *TUsecase,
	DataComboUsecase *DataComboUsecase,
	FeeUsecase *FeeUsecase,
	WordUsecase *WordUsecase) *BoxbuzUsecase {
	uc := &BoxbuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		BoxUsecase:       BoxUsecase,
		MapUsecase:       MapUsecase,
		TUsecase:         TUsecase,
		DataComboUsecase: DataComboUsecase,
		FeeUsecase:       FeeUsecase,
		WordUsecase:      WordUsecase,
	}

	return uc
}

const (
	CopyBoxResItemsToFolder_Type_file_only                = "file_only"
	CopyBoxResItemsToFolder_Type_file_only_and_ignore_409 = "file_only_and_ignore_409"

	BoxPersonalStatementsTemplateFileId = "1682431103455"
	BoxDocEmailTemplateFileId           = "1711991037359"
	BoxClaimsAnalysisTemplateFileId     = "1798477791938"
)

func GenUpdatePersonalStatementsFileName(firstName string, lastName string, caseId int32) string {
	var destFirstName string
	if firstName != "" {
		destFirstName = firstName[:1]
	}
	return fmt.Sprintf("DoNotUseAutoTest_Updated_Personal Statements_%s%s_%d.docx", destFirstName, lastName, caseId)
}

func GenPersonalStatementsFileNameForUpdateStatement(firstName string, lastName string, caseId int32) string {
	name := GenPersonalStatementsFileName(firstName, lastName, caseId)
	return fmt.Sprintf("ClientPS_Source_%s", name)
}

func GenPersonalStatementsFileName(firstName string, lastName string, caseId int32) string {
	return GenPersonalStatementsFileNameWithType(firstName, lastName, caseId, "docx")
}

func GenPersonalStatementsFileNamePdf(firstName string, lastName string, caseId int32) string {
	return GenPersonalStatementsFileNameWithType(firstName, lastName, caseId, "pdf")
}

func GenPersonalStatementsFileNameWithType(firstName string, lastName string, caseId int32, fileType string) string {
	var destFirstName string
	if firstName != "" {
		destFirstName = firstName[:1]
	}
	return fmt.Sprintf("Personal Statements_%s%s_%d.%s", destFirstName, lastName, caseId, fileType)
}

// GenPersonalStatementsFileNameAiAuto togo:lgl 正式上线后文件名改为与GenPersonalStatementsFileName一样
func GenPersonalStatementsFileNameAiAuto(firstName string, lastName string, caseId int32) string {
	var destFirstName string
	if firstName != "" {
		destFirstName = firstName[:1]
	}
	return fmt.Sprintf("DoNotUseAutoTest_Personal Statements_%s%s_%d.docx", destFirstName, lastName, caseId)
}

func GenDocEmailFileName(firstName string, lastName string, caseId int32) string {
	var destFirstName string
	if firstName != "" {
		destFirstName = firstName[:1]
	}
	return fmt.Sprintf("ConditionsSummary_%s%s_%d.docx", destFirstName, lastName, caseId)
}

func GenDocPDFEmailFileName(firstName string, lastName string, caseId int32) string {
	var destFirstName string
	if firstName != "" {
		destFirstName = firstName[:1]
	}
	return fmt.Sprintf("ConditionsSummary_%s%s_%d.pdf", destFirstName, lastName, caseId)
}

// GenDocEmailFileNameAuto togo:lgl 正式上线后文件名改为与GenDocEmailFileName一样
func GenDocEmailFileNameAuto(firstName string, lastName string, caseId int32) string {
	var destFirstName string
	if firstName != "" {
		destFirstName = firstName[:1]
	}
	return fmt.Sprintf("ConditionsSummary_%s%s_%d.docx", destFirstName, lastName, caseId)
}

func GenClaimsAnalysisFileName(firstName string, lastName string, caseId int32) string {
	var destFirstName string
	if firstName != "" {
		destFirstName = firstName[:1]
	}
	return fmt.Sprintf("Claims Analysis_%s%s_%d.docx", destFirstName, lastName, caseId)
}

func (c *BoxbuzUsecase) DoClaimsAnalysisFile(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	fileName := GenClaimsAnalysisFileName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	folderId, err := c.GetOrMakeClaimsAnalysisFolderId(tCase)
	if err != nil {
		return err
	}

	_, err = c.BoxUsecase.CopyFileNewFileName(BoxClaimsAnalysisTemplateFileId, fileName, folderId)
	if err != nil {
		return err
	}
	return nil
}

func (c *BoxbuzUsecase) DoDocEmailFile(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	fileName := GenDocEmailFileName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	dCPersonalStatementsFolderId, err := c.DCPersonalStatementsFolderId(tCase)
	if err != nil {
		return err
	}

	_, err = c.BoxUsecase.CopyFileNewFileName(BoxDocEmailTemplateFileId, fileName, dCPersonalStatementsFolderId)
	if err != nil {
		return err
	}
	return nil
}

func (c *BoxbuzUsecase) DoPersonalStatementsFile(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	fileName := GenPersonalStatementsFileName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	dCPersonalStatementsFolderId, err := c.DCPersonalStatementsFolderId(tCase)
	if err != nil {
		return err
	}

	wordReader, err := c.WordUsecase.DoPersonalStatementsWord(tCase)
	if err != nil {
		return err
	}
	_, err = c.BoxUsecase.UploadFile(dCPersonalStatementsFolderId, wordReader, fileName)

	return err
}

func (c *BoxbuzUsecase) DoPersonalStatementsFileBackup(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	fileName := GenPersonalStatementsFileName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	dCPersonalStatementsFolderId, err := c.DCPersonalStatementsFolderId(tCase)
	if err != nil {
		return err
	}

	_, err = c.BoxUsecase.CopyFileNewFileName(BoxPersonalStatementsTemplateFileId, fileName, dCPersonalStatementsFolderId)
	if err != nil {
		return err
	}
	return nil
}

func (c *BoxbuzUsecase) CopyBoxResItemsToFolder(folderId string, resItems lib.TypeList, Type string) error {
	for _, v := range resItems {
		resId := v.GetString("id")
		resType := v.GetString("type")
		resName := v.GetString("name")
		if Type == CopyBoxResItemsToFolder_Type_file_only ||
			Type == CopyBoxResItemsToFolder_Type_file_only_and_ignore_409 {
			if resType == string(config_box.BoxResType_file) {
				_, httpCode, err := c.BoxUsecase.CopyFile(resId, folderId)
				if Type == CopyBoxResItemsToFolder_Type_file_only_and_ignore_409 && httpCode == config_box.HttpCode_409 {
					c.log.Info(err)
				} else {
					if err != nil {
						return err
					}
				}
			}
		} else {
			if resType == string(config_box.BoxResType_folder) {
				_, _, err := c.BoxUsecase.CopyFolder(resId, resName, folderId)
				if err != nil {
					return err
				}
			} else {
				_, _, err := c.BoxUsecase.CopyFile(resId, folderId)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// MergeFolder 从源文件夹合并到新文件夹, 文件冲突时，删除目标文件夹的文件或文件夹
// 文件冲突就报错
// 此方法较危险，因为会删除文件夹(删除方法)
// 递归方法，安全-当前方式
func (c *BoxbuzUsecase) MergeFolder(sourceFolderId string, destFolderId string) error {
	sourceRes, err := c.BoxUsecase.ListItemsInFolder(sourceFolderId)
	if err != nil {
		return err
	}
	if sourceRes == nil {
		return errors.New("sourceRes is nil")
	}
	destRes, err := c.BoxUsecase.ListItemsInFolder(destFolderId)
	if err != nil {
		return err
	}
	if destRes == nil {
		return errors.New("destRes is nil")
	}
	sourceResMap := lib.ToTypeMapByString(*sourceRes)
	sourceResEntries := sourceResMap.GetTypeList("entries")

	destResMap := lib.ToTypeMapByString(*destRes)
	destResEntries := destResMap.GetTypeList("entries")
	existData := func(typeMaps lib.TypeMap, destResEntries lib.TypeList) (existType string, existId string) { // existType:folder/file
		for _, v := range destResEntries {
			if v.GetString("type") == typeMaps.GetString("type") &&
				v.GetString("name") == typeMaps.GetString("name") {
				return v.GetString("type"), v.GetString("id")
			}
		}
		return "", ""
	}
	for _, v := range sourceResEntries {
		existType, existId := existData(v, destResEntries)
		if existId != "" && existType == "folder" {

			// 递归方法
			err = c.MergeFolder(v.GetString("id"), existId)
			if err != nil {
				return err
			}

			// 不使用删除方法
			//_, err = c.BoxUsecase.DeleteFolder(existId)
			//if err != nil {
			//	return err
			//}
		} else if v.GetString("type") == "folder" {
			_, _, err = c.BoxUsecase.CopyFolder(v.GetString("id"), v.GetString("name"), destFolderId)
			c.log.Info("CopyFolder:", v.GetString("id"), v.GetString("name"), destFolderId)
			if err != nil {
				return err
			}
		} else {
			if existId == "" { // 文件已经存在就不再拷贝了
				_, _, err = c.BoxUsecase.CopyFile(v.GetString("id"), destFolderId)
				c.log.Info("CopyFile:", v.GetString("id"), v.GetString("name"), destFolderId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// RenameDataCollection 按新命名规则重命名
func (c *BoxbuzUsecase) RenameDataCollection() error {
	destRes, err := c.BoxUsecase.ListItemsInFolder("263406803830")
	if err != nil {
		return err
	}
	if destRes == nil {
		return errors.New("destRes is nil")
	}
	destResMap := lib.ToTypeMapByString(*destRes)
	destResEntries := destResMap.GetTypeList("entries")
	for _, v := range destResEntries {

		if v.GetString("type") == "folder" {
			sourceName := v.GetString("name")
			newName, err := DataCollectionTidy(sourceName)
			if err != nil {
				return err
			}
			if newName != sourceName {
				_, err := c.BoxUsecase.UpdateFolderName(v.GetString("id"), newName)
				if err != nil {
					return err
				}
				time.Sleep(1 * time.Second)
				//break
			}
		}

	}
	return nil
}

// RenameClientCasesFolderName 按新命名规则重命名
func (c *BoxbuzUsecase) RenameClientCasesFolderName() error {
	destRes, err := c.BoxUsecase.ListItemsInFolder("241109085470")
	if err != nil {
		return err
	}
	if destRes == nil {
		return errors.New("destRes is nil")
	}
	destResMap := lib.ToTypeMapByString(*destRes)
	destResEntries := destResMap.GetTypeList("entries")
	for _, v := range destResEntries {

		if v.GetString("type") == "folder" {
			sourceName := v.GetString("name")
			newName, err := ClientCaseFolderTidy(sourceName)
			if err != nil {
				return err
			}
			if newName != sourceName {
				_, err := c.BoxUsecase.UpdateFolderName(v.GetString("id"), newName)
				if err != nil {
					return err
				}
				time.Sleep(1 * time.Second)
				//break
			}
		}

	}
	return nil
}

// CPersonalStatementsFolderIdByAnyCase 支持Primary Case和Seconds Case
func (c *BoxbuzUsecase) CPersonalStatementsFolderIdByAnyCase(tCase TData) (subFolderId string, err error) {

	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(&tCase)
	if err != nil {
		return "", err
	}
	if !isPrimaryCase {
		return c.GetClientSubFolderId(
			MapKeyBuildAutoBoxCPersonalStatementsFolderId(primaryCase.Id()),
			primaryCase)
	}
	return c.GetClientSubFolderId(
		MapKeyBuildAutoBoxCPersonalStatementsFolderId(tCase.Id()),
		&tCase)
}

func (c *BoxbuzUsecase) CPersonalStatementsFolderId(tCase *TData) (subFolderId string, err error) {
	return c.GetClientSubFolderId(
		MapKeyBuildAutoBoxCPersonalStatementsFolderId(tCase.Id()),
		tCase)
}

func (c *BoxbuzUsecase) GetClientSubFolderId(key string, tClientCase *TData) (subFolderId string, err error) {

	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}
	subFolderId, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if subFolderId == "" {
		err = c.HandleClientFolder(tClientCase)
		if err != nil {
			return "", err
		}
		subFolderId, err = c.MapUsecase.GetForString(key)
	}
	return
}

func (c *BoxbuzUsecase) GetClientBoxFolderId(tClientCase *TData) (clientBoxFolderId string, err error) {
	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}
	clientCaseId := tClientCase.CustomFields.NumberValueByNameBasic("id")
	clientBoxFolderId, err = c.MapUsecase.GetForString(MapKeyClientBoxFolderId(clientCaseId))
	return
}

// HandleClientFolder 处理客户文件夹下的文件、文件夹与DB关系
func (c *BoxbuzUsecase) HandleClientFolder(tClientCase *TData) error {
	if tClientCase == nil {
		return errors.New("tClientCase is nil")
	}
	clientCaseId := tClientCase.CustomFields.NumberValueByNameBasic("id")
	clientBoxFolderId, err := c.MapUsecase.GetForString(MapKeyClientBoxFolderId(clientCaseId))
	if err != nil {
		return err
	}
	if clientBoxFolderId == "" {
		return errors.New("HandleClientFolder: clientBoxFolderId is nil")
	}
	res, err := c.BoxUsecase.ListItemsInFolderFormat(clientBoxFolderId)
	if err != nil {
		return err
	}
	var PrivateMedicalRecordsFolderId string
	var VAMedicalRecordsFolderId string
	var ServiceTreatmentRecordsFolderId string
	var MiscFolderId string
	var DD214FolderId string
	var DisabilityRatingListFolderId string
	var RatingDecisionLettersFolderId string
	var PersonalStatementsFolderId string

	for _, v := range res {
		if v.GetString("type") == "folder" {
			if v.GetString("name") == config_box.FolderName_PrivateMedicalRecords {
				PrivateMedicalRecordsFolderId = v.GetString("id")
			} else if v.GetString("name") == config_box.FolderName_VAMedicalRecords {
				VAMedicalRecordsFolderId = v.GetString("id")
			} else if v.GetString("name") == config_box.FolderName_ServiceTreatmentRecords {
				ServiceTreatmentRecordsFolderId = v.GetString("id")
			} else if v.GetString("name") == config_box.FolderName_C_Misc {
				MiscFolderId = v.GetString("id")
			} else if v.GetString("name") == config_box.FolderName_C_DD214 {
				DD214FolderId = v.GetString("id")
			} else if v.GetString("name") == config_box.FolderName_C_DisabilityRatingList {
				DisabilityRatingListFolderId = v.GetString("id")
			} else if v.GetString("name") == config_box.FolderName_C_RatingDecisionLetters {
				RatingDecisionLettersFolderId = v.GetString("id")
			} else if v.GetString("name") == config_box.FolderName_C_PersonalStatements {
				PersonalStatementsFolderId = v.GetString("id")
			}
		}
	}
	if PrivateMedicalRecordsFolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxCPrivateMedicalRecordsFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxCPrivateMedicalRecordsFolderId(clientCaseId), PrivateMedicalRecordsFolderId)
		}
	}
	if VAMedicalRecordsFolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxCVAMedicalRecordsFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxCVAMedicalRecordsFolderId(clientCaseId), VAMedicalRecordsFolderId)
		}
	}
	if ServiceTreatmentRecordsFolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxCServiceTreatmentRecordsFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxCServiceTreatmentRecordsFolderId(clientCaseId), ServiceTreatmentRecordsFolderId)
		}
	}
	if MiscFolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxCMiscFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxCMiscFolderId(clientCaseId), MiscFolderId)
		}
	}
	if DD214FolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxCDD214FolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxCDD214FolderId(clientCaseId), DD214FolderId)
		}
	}
	if DisabilityRatingListFolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxCDisabilityRatingListFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxCDisabilityRatingListFolderId(clientCaseId), DisabilityRatingListFolderId)
		}
	}
	if RatingDecisionLettersFolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxCRatingDecisionLettersFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxCRatingDecisionLettersFolderId(clientCaseId), RatingDecisionLettersFolderId)
		}
	}
	if PersonalStatementsFolderId != "" {
		key := MapKeyBuildAutoBoxCPersonalStatementsFolderId(clientCaseId)
		val, err := c.MapUsecase.GetForString(key)
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(key, PersonalStatementsFolderId)
		}
	}
	return nil
}

func (c *BoxbuzUsecase) DCPrivateExamsFolderId(tCase *TData) (DCPrivateExamsFolderId string, err error) {
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	clientCaseId := tCase.CustomFields.NumberValueByNameBasic("id")
	return c.GetDCSubFolderId(MapKeyBuildAutoBoxDCPrivateExamsFolderId(
		clientCaseId), tCase)
}

func (c *BoxbuzUsecase) DCPersonalStatementsFolderId(tCase *TData) (DCFolderId string, err error) {
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	clientCaseId := tCase.CustomFields.NumberValueByNameBasic("id")
	return c.GetDCSubFolderId(MapKeyBuildAutoBoxDCPersonalStatementsFolderId(
		clientCaseId), tCase)
}

func (c *BoxbuzUsecase) PersonalStatementDocFileBoxFileId(tClient *TData, tCase *TData) (dCPersonalStatementsFolderId string, boxFileId string, err error) {

	if tClient == nil {
		return "", "", errors.New("tClient is nil")
	}
	if tCase == nil {
		return "", "", errors.New("tCase is nil")
	}

	psFileName := GenPersonalStatementsFileName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name), tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	dCPersonalStatementsFolderId, err = c.DCPersonalStatementsFolderId(tCase)
	if err != nil {
		return "", "", err
	}
	if dCPersonalStatementsFolderId == "" {
		return "", "", errors.New("dCPersonalStatementsFolderId is empty")
	}
	resItems, err := c.BoxUsecase.ListItemsInFolderFormat(dCPersonalStatementsFolderId)
	if err != nil {
		return "", "", err
	}
	for _, v := range resItems {
		resId := v.GetString("id")
		resType := v.GetString("type")
		resName := v.GetString("name")
		if resType == string(config_box.BoxResType_file) {
			if resName == psFileName {
				boxFileId = resId
				break
			}
		}
	}
	return
}

func (c *BoxbuzUsecase) FolderIdDC_PE_Psych(tCase *TData, tClient *TData) (PsychFolderId string, err error) {
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	DCPrivateExamsFolderId, err := c.DCPrivateExamsFolderId(tCase)
	if err != nil {
		return "", err
	}
	if DCPrivateExamsFolderId == "" {
		return "", errors.New("DCPrivateExamsFolderId is empty")
	}
	return c.GetBoxResIdByCase(DCPrivateExamsFolderId,
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_PE_Psych_Folder),
		tCase,
		tClient)
}

func (c *BoxbuzUsecase) FolderIdDC_PE_General(tCase *TData, tClient *TData) (GeneralFolderId string, err error) {
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	DCPrivateExamsFolderId, err := c.DCPrivateExamsFolderId(tCase)
	if err != nil {
		return "", err
	}
	if DCPrivateExamsFolderId == "" {
		return "", errors.New("DCPrivateExamsFolderId is empty")
	}
	return c.GetBoxResIdByCase(DCPrivateExamsFolderId,
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_PE_General_Folder),
		tCase,
		tClient)
}

func (c *BoxbuzUsecase) GetOrMakeClaimsAnalysisFolderId(tClientCase *TData) (folderId string, err error) {
	if tClientCase == nil {
		return "", errors.New("GetClaimsAnalysisFolderId:tClientCase is nil")
	}
	return c.GetDCSubFolderId(MapKeyBuildAutoBoxDCClaimsAnalysisFolderId(tClientCase.Id()), tClientCase)
}

func (c *BoxbuzUsecase) DCQuestionnairesFolderId(tCase TData) (questionnairesFolderId string, err error) {
	return c.GetDCSubFolderId(MapKeyBuildAutoBoxDCQuestionnairesFolderId(tCase.Id()), &tCase)
}

func (c *BoxbuzUsecase) DCGetOrMakeUpdateQuestionnairesFolderId(tCase TData) (subFolderId string, err error) {
	questionnairesFolderId, err := c.DCQuestionnairesFolderId(tCase)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("UpdateQuestionnairesFolderId:%d", tCase.Id())
	return c.GetOrMakeSubFolderId(key, questionnairesFolderId, "Update Questionnaires")
}

func (c *BoxbuzUsecase) GetOrMakeSubFolderId(key string, folderId string, subFolderName string) (subFolderId string, err error) {

	subFolderId, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if subFolderId == "" {
		records, err := c.BoxUsecase.ListItemsInFolderFormat(folderId)
		if err != nil {
			return "", err
		}
		for _, v := range records {
			if v.GetString("type") == "folder" {
				if v.GetString("name") == subFolderName {
					subFolderId = v.GetString("id")
					break
				}
			}
		}
		if subFolderId == "" {
			subFolderId, err = c.BoxUsecase.CreateFolder(subFolderName, folderId)
			if err != nil {
				return "", err
			}
		}
		c.MapUsecase.Set(key, subFolderId)
		return subFolderId, nil
	}
	return subFolderId, nil
}

func (c *BoxbuzUsecase) GetDCFolderId(caseId int32) (dcBoxFolderId string, err error) {
	return c.MapUsecase.GetForString(MapKeyDataCollectionFolderId(caseId))
}

func (c *BoxbuzUsecase) GetDCSubFolderId(key string, tClientCase *TData) (subFolderId string, err error) {

	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}
	subFolderId, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if subFolderId == "" {
		err = c.HandleDCFolder(tClientCase)
		if err != nil {
			return "", err
		}
		subFolderId, err = c.MapUsecase.GetForString(key)
		if err != nil {
			return "", err
		}
		if subFolderId == "" {
			dcBoxFolderId, err := c.MapUsecase.GetForString(MapKeyDataCollectionFolderId(tClientCase.Id()))
			if err != nil {
				return "", err
			}
			if dcBoxFolderId == "" {
				return "", errors.New("GetDCSubFolderId: dcBoxFolderId is nil")
			}
			if key == MapKeyBuildAutoBoxDCClaimsAnalysisFolderId(tClientCase.Id()) { // 走新建文件夹流程
				subFolderId, err = c.CreateFolder(key, config_box.FolderName_DCClaimsAnalysis, dcBoxFolderId)
				if err != nil {
					return "", err
				}
			}
		}
		return
	}
	return
}

func (c *BoxbuzUsecase) CreateFolder(key string, folderName string, destParentId string) (boxFolderId string, err error) {
	boxFolderId, err = c.BoxUsecase.CreateFolder(folderName, destParentId)
	if err != nil {
		return "", err
	}
	if boxFolderId == "" {
		return "", errors.New("boxFolderId is empty")
	}
	c.MapUsecase.Set(key, boxFolderId)
	return boxFolderId, nil
}

const (
	CaseRelaBox_bizType_DC_RecordReview         = "DC_RecordReview"
	CaseRelaBox_bizName_VAMedicalRecords        = "VAMedicalRecords"
	CaseRelaBox_bizName_ServiceTreatmentRecords = "ServiceTreatmentRecords"
	CaseRelaBox_bizName_PrivateMedicalRecords   = "PrivateMedicalRecords"
)

func (c *BoxbuzUsecase) DoCopyDocEmailFile(caseId int32) (err error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		c.log.Error(err)
		return err
	}
	var CPersonalStatementsFolderId string
	if isPrimaryCase {
		CPersonalStatementsFolderId, err = c.CPersonalStatementsFolderId(tCase)
		if err != nil {
			return err
		}
	} else {
		//primaryCaseId := primaryCase.Id()
		//c.log.Debug("primaryCaseId:", primaryCaseId)
		CPersonalStatementsFolderId, err = c.CPersonalStatementsFolderId(primaryCase)
		if err != nil {
			return err
		}
	}
	if CPersonalStatementsFolderId == "" {
		return errors.New("CPersonalStatementsFolderId is empty")
	}

	docEmailFileId, err := c.DCPSDocEmailFileId(tCase)
	if err != nil {
		return err
	}
	lib.DPrintln(docEmailFileId, CPersonalStatementsFolderId)
	_, _, err = c.BoxUsecase.CopyFile(docEmailFileId, CPersonalStatementsFolderId)
	if err != nil {
		return err
	}
	return nil
}

func (c *BoxbuzUsecase) DoCopyReadPriorToYourDoctorVisitFile(caseId int32) (err error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		c.log.Error(err)
		return err
	}
	var CPersonalStatementsFolderId string
	if isPrimaryCase {
		CPersonalStatementsFolderId, err = c.CPersonalStatementsFolderId(tCase)
		if err != nil {
			return err
		}
	} else {
		CPersonalStatementsFolderId, err = c.CPersonalStatementsFolderId(primaryCase)
		if err != nil {
			return err
		}
	}
	if CPersonalStatementsFolderId == "" {
		return errors.New("CPersonalStatementsFolderId is empty")
	}

	_, _, err = c.BoxUsecase.CopyFile(config_box.ReadPriorToYourDoctorVisitFileId, CPersonalStatementsFolderId)
	if err != nil {
		return err
	}
	return nil
}

func (c *BoxbuzUsecase) DCPSDocEmailFileId(tCase *TData) (docEmailFileId string, err error) {

	//caseId := int32(5369)
	//tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, caseId)

	psFolderId, err := c.DCPersonalStatementsFolderId(tCase)
	docEmailFileId, err = c.GetBoxResId(psFolderId,
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_PS_DocEmail_File),
		tCase.Id(),
	)
	return docEmailFileId, err
}

func (c *BoxbuzUsecase) GetBoxResId(folderId string,
	caseRelaBoxVo *config_box.CaseRelaBoxVo,
	clientCaseId int32) (boxResId string, err error) {
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		c.log.Error(err)
		return "", err
	}
	tClient, _, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		c.log.Error(err)
		return "", err
	}

	return c.GetBoxResIdByCase(folderId, caseRelaBoxVo, tClientCase, tClient)
}

// GetBoxResIdKey 获取存在maps的key
func (c *BoxbuzUsecase) GetBoxResIdKey(caseRelaBoxVo *config_box.CaseRelaBoxVo, tClientCase *TData) (string, error) {

	if caseRelaBoxVo == nil {
		return "", errors.New("caseRelaBoxVo is nil")
	}
	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}
	clientCaseId := tClientCase.CustomFields.NumberValueByNameBasic("id")

	if caseRelaBoxVo.UniqueKey == config_box.CaseRelaBox_DC_PE_PatientPaymentForm_File {
		return MapKeyBuildAutoBoxDCPatientPaymentFormFileId(clientCaseId), nil
	}
	return MapKeyCaseRelaBox(caseRelaBoxVo.UniqueKey, clientCaseId), nil
}

// GetBoxResIdByCase 处理客户folderId文件夹下的文件、文件夹与DB关系, 可以是第二个case文件夹
func (c *BoxbuzUsecase) GetBoxResIdByCase(folderId string,
	caseRelaBoxVo *config_box.CaseRelaBoxVo,
	tClientCase *TData, tClient *TData) (boxResId string, err error) {
	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}
	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	clientCaseId := tClientCase.CustomFields.NumberValueByNameBasic("id")

	//key := MapKeyCaseRelaBox(caseRelaBoxVo.UniqueKey, clientCaseId)
	key, err := c.GetBoxResIdKey(caseRelaBoxVo, tClientCase)
	if err != nil {
		return "", err
	}

	boxResId, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if boxResId != "" {
		return boxResId, nil
	}

	res, err := c.BoxUsecase.ListItemsInFolderFormat(folderId)
	if err != nil {
		return "", err
	}

	for _, v := range res {
		if v.GetString("type") == string(caseRelaBoxVo.BoxResType) {
			if len(caseRelaBoxVo.ResPartialNames) > 0 {
				for _, v1 := range caseRelaBoxVo.ResPartialNames {
					if strings.Index(v.GetString("name"), v1) >= 0 {
						boxResId = v.GetString("id")
						break
					}
				}
			} else {

				if caseRelaBoxVo.UniqueKey == config_box.CaseRelaBox_DC_PE_PatientPaymentForm_File {
					if strings.Index(v.GetString("name"), PatientPaymentForm_Postfix) > 0 {
						boxResId = v.GetString("id")
						break
					}
				} else {
					resName := caseRelaBoxVo.ResName
					if caseRelaBoxVo.UniqueKey == config_box.CaseRelaBox_C_New_Evidence_Folder {
						resName = resName + " #" + InterfaceToString(clientCaseId)
					}
					if v.GetString("name") == resName {
						boxResId = v.GetString("id")
						break
					}
				}
			}
		}
	}
	if boxResId != "" {
		c.MapUsecase.Set(key, boxResId)
	}
	return boxResId, nil
}

// HandleDCFolder 处理客户DC文件夹下的文件、文件夹与DB关系
func (c *BoxbuzUsecase) HandleDCFolder(tClientCase *TData) error {
	if tClientCase == nil {
		return errors.New("tClientCase is nil")
	}
	clientCaseId := tClientCase.CustomFields.NumberValueByNameBasic("id")
	dcBoxFolderId, err := c.MapUsecase.GetForString(MapKeyDataCollectionFolderId(clientCaseId))
	if err != nil {
		return err
	}
	if dcBoxFolderId == "" {
		return errors.New("HandleClientFolder: dcBoxFolderId is nil")
	}
	res, err := c.BoxUsecase.ListItemsInFolderFormat(dcBoxFolderId)
	if err != nil {
		return err
	}
	var DCQuestionnairesFolderId string
	var DCPrivateExamsFolderId string
	var DCPersonalStatementsFolderId string
	var DCClaimsAnalysisFolderId string

	for _, v := range res {
		if v.GetString("type") == "folder" {
			if v.GetString("name") == config_box.FolderName_DCQuestionnaires {
				DCQuestionnairesFolderId = v.GetString("id")
			}
			if v.GetString("name") == config_box.FolderName_DCPrivateExams {
				DCPrivateExamsFolderId = v.GetString("id")
			}
			if v.GetString("name") == config_box.FolderName_DCPersonalStatements {
				DCPersonalStatementsFolderId = v.GetString("id")
			}
			if v.GetString("name") == config_box.FolderName_DCClaimsAnalysis {
				DCClaimsAnalysisFolderId = v.GetString("id")
			}
		}
	}
	if DCQuestionnairesFolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxDCQuestionnairesFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxDCQuestionnairesFolderId(clientCaseId), DCQuestionnairesFolderId)
		}
	}
	if DCPrivateExamsFolderId != "" {
		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxDCPrivateExamsFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxDCPrivateExamsFolderId(clientCaseId), DCPrivateExamsFolderId)
		}
	}
	if DCPersonalStatementsFolderId != "" {

		val, err := c.MapUsecase.GetForString(MapKeyBuildAutoBoxDCPersonalStatementsFolderId(clientCaseId))
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(MapKeyBuildAutoBoxDCPersonalStatementsFolderId(clientCaseId), DCPersonalStatementsFolderId)
		}
	}
	if DCClaimsAnalysisFolderId != "" {
		key := MapKeyBuildAutoBoxDCClaimsAnalysisFolderId(clientCaseId)
		val, err := c.MapUsecase.GetForString(key)
		if err != nil {
			return err
		}
		if val == "" {
			c.MapUsecase.Set(key, DCClaimsAnalysisFolderId)
		}
	}
	return nil
}

func (c *BoxbuzUsecase) DCRecordReviewFolderId(tClientCase *TData) (boxFolderId string, err error) {
	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}

	clientCaseId := tClientCase.CustomFields.NumberValueByNameBasic("id")

	boxFolderId, err = c.MapUsecase.GetForString(MapKeyBuildAutoBoxDCRecordReviewFolderId(clientCaseId))
	if err != nil {
		return "", err
	}
	if boxFolderId != "" {
		return boxFolderId, nil
	}

	dataCollectionFolderId, err := c.MapUsecase.GetForString(MapKeyDataCollectionFolderId(clientCaseId))
	if err != nil {
		return "", err
	}
	if dataCollectionFolderId == "" {
		return "", errors.New("dataCollectionFolderId is empty, clientCaseId: " + InterfaceToString(clientCaseId))
	}
	res, err := c.BoxUsecase.ListItemsInFolderFormat(dataCollectionFolderId)
	if err != nil {
		return "", err
	}
	for _, v := range res {
		if v.GetString("type") == "folder" && v.GetString("name") == config_box.FolderName_RecordReview {
			boxFolderId = v.GetString("id")
			break
		}
	}
	if boxFolderId != "" {
		c.MapUsecase.Set(MapKeyBuildAutoBoxDCRecordReviewFolderId(clientCaseId), boxFolderId)
		return boxFolderId, nil
	} else {
		return "", errors.New("DCRecordReviewFolderId is empty.")
	}
}

// SameNameFolderOrFile typ: folder / file, 判断name是否有重名，有的话返回重名id
func (c *BoxbuzUsecase) SameNameFolderOrFile(typ string, name string, parentFolderId string) (sameId string, err error) {
	if typ != "folder" && typ != "file" {
		return "", errors.New("typ is wrong.")
	}
	entities, err := c.BoxUsecase.ListItemsInFolderFormat(parentFolderId)
	if err != nil {
		return "", err
	}
	for _, v := range entities {
		if v.GetString("type") == typ && v.GetString("name") == name {
			return v.GetString("id"), nil
		}
	}
	return "", nil
}

func ClientCaseFolderTidy(sourceName string) (newName string, err error) {
	if strings.Index(sourceName, "#") >= 0 {
		res := strings.Split(sourceName, "#")
		if len(res) != 2 {
			return sourceName, errors.New("sourceName error: " + sourceName)
		}
		a := strings.TrimSpace(res[0])
		b := strings.TrimSpace(res[1])
		newName = fmt.Sprintf("%s #%s", a, b)
	} else {
		newName = sourceName
	}
	//fmt.Println("sourceName_1: ", sourceName, "newName_1: ", newName)
	return
}

func DataCollectionTidy(sourceName string) (newName string, err error) {
	if strings.Index(sourceName, "#") >= 0 {
		res := strings.Split(sourceName, "#")
		if len(res) != 2 {
			return sourceName, errors.New("sourceName error: " + sourceName)
		}
		a := strings.TrimSpace(res[0])
		b := strings.TrimSpace(res[1])
		if lib.InterfaceToInt32(b) < 5000 && strings.Index(sourceName, "Edward David") < 0 {
			newName = fmt.Sprintf("%s", a)
		} else {
			newName = fmt.Sprintf("%s #%s", a, b)
		}
	} else {
		newName = sourceName
	}
	//fmt.Println("sourceName_1: ", sourceName, "newName_1: ", newName)
	return
}

// CreateFolderByEntries 创建或返回文件夹
func (c *BoxbuzUsecase) CreateFolderByEntries(beginIndex int, entries lib.TypeList, parentFolderId string) (bottommostFolderId string, path string, err error) {

	isCreateFolder := false
	LastCreateFolder := ""
	destFolderId := ""
	for i := beginIndex; i < len(entries); i++ {
		destFolderId = ""
		destFolderName := entries[i].GetString("name")
		path += destFolderName + "/"

		if !isCreateFolder {
			items, err := c.BoxUsecase.ListItemsInFolderFormat(parentFolderId)
			if err != nil {
				c.log.Error(err)
				return "", "", err
			}
			for _, v := range items {
				resId := v.GetString("id")
				resType := v.GetString("type")
				resName := v.GetString("name")
				if resType == string(config_box.BoxResType_folder) && resName == destFolderName {
					destFolderId = resId
					break
				}
			}
		} else {
			parentFolderId = LastCreateFolder
		}
		if destFolderId == "" {
			destFolderId, err = c.BoxUsecase.CreateFolder(destFolderName, parentFolderId)
			if err != nil {
				c.log.Error(err)
				return "", "", err
			}
			LastCreateFolder = destFolderId
			isCreateFolder = true
		} else {
			parentFolderId = destFolderId
			LastCreateFolder = destFolderId
		}
		if destFolderId == "" {
			return "", "", errors.New("destFolderId is empty")
		}
	}
	return destFolderId, path, nil
}

// CopyFileToFolderNoCover 拷贝文件到目录，如果文件名已经存在就不进行覆盖
func (c *BoxbuzUsecase) CopyFileToFolderNoCover(folderId string, fileId string, fileName string) (fileNameExist bool, newFileId string, err error) {
	items, err := c.BoxUsecase.ListItemsInFolderFormat(folderId)
	if err != nil {
		return false, "", err
	}
	for _, v := range items {
		if v.GetString("type") == string(config_box.BoxResType_file) &&
			v.GetString("name") == fileName {
			return true, "", nil
		}
	}
	newFileId, err = c.BoxUsecase.CopyFileNewFileNameReturnFileId(fileId, fileName, folderId)
	fileNameExist = false
	return
}

// RealtimeCPersonalStatementsDocxFileId 获取文件id， 当docxFileId不存在时，说明文件不存在
func (c *BoxbuzUsecase) RealtimeCPersonalStatementsDocxFileId(tClient TData, tCase TData) (cPersonalStatementsFolderId string, docxFileName string, docxFileId string, err error) {

	docxFileName = GenPersonalStatementsFileName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name), tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	cPersonalStatementsFolderId, err = c.CPersonalStatementsFolderIdByAnyCase(tCase)
	if err != nil {
		return "", "", "", err
	}
	if cPersonalStatementsFolderId == "" {
		return "", "", "", errors.New("cPersonalStatementsFolderId is empty")
	}
	resItems, err := c.BoxUsecase.ListItemsInFolderFormat(cPersonalStatementsFolderId)
	if err != nil {
		return "", "", "", err
	}
	for _, v := range resItems {
		resId := v.GetString("id")
		resType := v.GetString("type")
		resName := v.GetString("name")
		if resType == string(config_box.BoxResType_file) {
			if resName == docxFileName {
				docxFileId = resId
				break
			}
		}
	}
	return
}

// CopyFolderSubsToFolder 解决手动创建box文件夹，复制文件夹的问题
func (c *BoxbuzUsecase) CopyFolderSubsToFolder(sourceFolderId string, destFolderId string) error {

	list, err := c.BoxUsecase.ListItemsInFolderFormat(sourceFolderId)
	if err != nil {
		return err
	}

	for _, v := range list {
		typ := v.GetString("type")
		id := v.GetString("id")
		name := v.GetString("name")
		if typ == "folder" {
			_, _, err := c.BoxUsecase.CopyFolder(id, name, destFolderId)
			if err != nil {
				lib.DPrintln("err:", err, id, name, destFolderId)
			}
		} else {
			_, _, err := c.BoxUsecase.CopyFile(id, destFolderId)
			if err != nil {
				lib.DPrintln("err:", err, id, name, destFolderId)
			}
		}

	}
	return nil
}

// UseVBCActiveCasesFolder 所有客户都启用VBCActiveCasesFolder
func UseVBCActiveCasesFolder(tCase TData) bool {
	return true
	dealName := tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	if strings.Index(strings.ToLower(dealName), "test") >= 0 {
		return true
	}
	return false
}

func (c *BoxbuzUsecase) GetClientFolderRootId(tCase TData) (useVBCActiveCasesFolder bool, folderId string) {

	isPrimaryCase, primaryCase, _ := c.FeeUsecase.UsePrimaryCaseCalc(&tCase)
	if isPrimaryCase {
		if UseVBCActiveCasesFolder(tCase) {
			return true, c.conf.Box.ClientFolderStructureParentIdV2
		}
	} else {
		if primaryCase != nil {
			key := MapKeyUseVBCActiveCasesFolder(primaryCase.Id())
			val, _ := c.MapUsecase.GetForString(key)
			if val == "1" {
				return true, c.conf.Box.ClientFolderStructureParentIdV2
			}
		}
	}
	return false, c.conf.Box.ClientFolderStructureParentId
}

func (c *BoxbuzUsecase) GetDataCollectionFolderRootId(tCase TData) (useVBCActiveCasesFolder bool, folderId string) {

	isPrimaryCase, primaryCase, _ := c.FeeUsecase.UsePrimaryCaseCalc(&tCase)
	if isPrimaryCase {
		if UseVBCActiveCasesFolder(tCase) {
			return true, c.conf.Box.DataCollectionFolderIdV2
		}
	} else {
		if primaryCase != nil {
			key := MapKeyUseVBCActiveCasesFolder(primaryCase.Id())
			val, _ := c.MapUsecase.GetForString(key)
			if val == "1" {
				return true, c.conf.Box.DataCollectionFolderIdV2
			}
		}
	}
	return false, c.conf.Box.DataCollectionFolderId
}
