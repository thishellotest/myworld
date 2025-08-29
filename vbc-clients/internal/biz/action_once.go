package biz

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/internal/utils"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

type ActionOnceUsecase struct {
	log                        *log.Helper
	CommonUsecase              *CommonUsecase
	conf                       *conf.Data
	MapUsecase                 *MapUsecase
	TUsecase                   *TUsecase
	ZohoUsecase                *ZohoUsecase
	FeeUsecase                 *FeeUsecase
	ClientCaseUsecase          *ClientCaseUsecase
	BoxUsecase                 *BoxUsecase
	DataComboUsecase           *DataComboUsecase
	BoxbuzUsecase              *BoxbuzUsecase
	DbqsUsecase                *DbqsUsecase
	BoxcontractUsecase         *BoxcontractUsecase
	ClientEnvelopeUsecase      *ClientEnvelopeUsecase
	RollpoingUsecase           *RollpoingUsecase
	PdfcpuUsecase              *PdfcpuUsecase
	BehaviorUsecase            *BehaviorUsecase
	TaskCreateUsecase          *TaskCreateUsecase
	BoxcUsecase                *BoxcUsecase
	GoogleDrivebuzUsecase      *GoogleDrivebuzUsecase
	StageTransUsecase          *StageTransUsecase
	PricingVersionUsecase      *PricingVersionUsecase
	DataEntryUsecase           *DataEntryUsecase
	GopdfUsecase               *GopdfUsecase
	RemindUsecase              *RemindUsecase
	PersonalWebformUsecase     *PersonalWebformUsecase
	BoxCollaborationBuzUsecase *BoxCollaborationBuzUsecase
}

func NewActionOnceUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	TUsecase *TUsecase,
	ZohoUsecase *ZohoUsecase,
	FeeUsecase *FeeUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	BoxUsecase *BoxUsecase,
	DataComboUsecase *DataComboUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	DbqsUsecase *DbqsUsecase,
	BoxcontractUsecase *BoxcontractUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	RollpoingUsecase *RollpoingUsecase,
	PdfcpuUsecase *PdfcpuUsecase,
	BehaviorUsecase *BehaviorUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	BoxcUsecase *BoxcUsecase,
	GoogleDrivebuzUsecase *GoogleDrivebuzUsecase,
	StageTransUsecase *StageTransUsecase,
	PricingVersionUsecase *PricingVersionUsecase,
	DataEntryUsecase *DataEntryUsecase,
	GopdfUsecase *GopdfUsecase,
	RemindUsecase *RemindUsecase,
	PersonalWebformUsecase *PersonalWebformUsecase,
	BoxCollaborationBuzUsecase *BoxCollaborationBuzUsecase) *ActionOnceUsecase {
	uc := &ActionOnceUsecase{
		log:                        log.NewHelper(logger),
		CommonUsecase:              CommonUsecase,
		conf:                       conf,
		MapUsecase:                 MapUsecase,
		TUsecase:                   TUsecase,
		ZohoUsecase:                ZohoUsecase,
		FeeUsecase:                 FeeUsecase,
		ClientCaseUsecase:          ClientCaseUsecase,
		BoxUsecase:                 BoxUsecase,
		DataComboUsecase:           DataComboUsecase,
		BoxbuzUsecase:              BoxbuzUsecase,
		DbqsUsecase:                DbqsUsecase,
		BoxcontractUsecase:         BoxcontractUsecase,
		ClientEnvelopeUsecase:      ClientEnvelopeUsecase,
		RollpoingUsecase:           RollpoingUsecase,
		PdfcpuUsecase:              PdfcpuUsecase,
		BehaviorUsecase:            BehaviorUsecase,
		TaskCreateUsecase:          TaskCreateUsecase,
		BoxcUsecase:                BoxcUsecase,
		GoogleDrivebuzUsecase:      GoogleDrivebuzUsecase,
		StageTransUsecase:          StageTransUsecase,
		PricingVersionUsecase:      PricingVersionUsecase,
		DataEntryUsecase:           DataEntryUsecase,
		GopdfUsecase:               GopdfUsecase,
		RemindUsecase:              RemindUsecase,
		PersonalWebformUsecase:     PersonalWebformUsecase,
		BoxCollaborationBuzUsecase: BoxCollaborationBuzUsecase,
	}
	return uc
}

// HttpInitClientCase 更新
func (c *ActionOnceUsecase) HttpInitClientCase(ctx *gin.Context) {

	reply := CreateReply()
	err := c.InitClientCase(lib.InterfaceToInt32(ctx.Query("id")))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

// InitClientCase - 初始化clientCase
func (c *ActionOnceUsecase) InitClientCase(clientCaseId int32) error {
	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "InitClientCase", clientCaseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val != "1" {
		tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("InitClientCase: tClientCase is nil")
		}
		clientCaseGid := tClientCase.CustomFields.TextValueByNameBasic("gid")
		if clientCaseGid == "" {
			return errors.New("InitClientCase: clientCaseGid is empty")
		}
		clientGid := tClientCase.CustomFields.TextValueByNameBasic("client_gid")
		if clientGid == "" {
			return errors.New("clientGid is empty")
		}
		tClient, err := c.TUsecase.DataByGid(Kind_clients, clientGid)
		if err != nil {
			return err
		}
		if tClient == nil {
			return errors.New("tClient is nil")
		}

		fields := needClientSyncCaseFields
		destMap := make(TypeDataEntry)
		// 处理跳转Stage:FeeSchedule的case的价格版本
		pricingVersion := tClientCase.CustomFields.TextValueByNameBasic(FieldName_s_pricing_version)
		contractSource := tClientCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource)

		if pricingVersion == "" || contractSource == "" {
			isPrimaryCaseCalc, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
			if err != nil {
				c.log.Error(err, " caseId:", tClientCase.Id())
			} else {
				if pricingVersion == "" {
					if isPrimaryCaseCalc {
						_, versionEntity, err := c.PricingVersionUsecase.CurrentVersionConfig()
						if err != nil {
							c.log.Error(err, " caseId:", tClientCase.Id())
						} else {
							if versionEntity != nil {
								// todo:lgl 此处注释是为了兼容 AM流程
								//pricingVersion = versionEntity.Version
							}
						}
					} else {
						if primaryCase != nil {
							pricingVersion = primaryCase.CustomFields.TextValueByNameBasic(FieldName_s_pricing_version)
						}
					}
				}
				if contractSource == "" {
					if !isPrimaryCaseCalc {
						if primaryCase != nil {
							contractSource = primaryCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource)
						}
					}
				}
			}
		}
		if pricingVersion != "" || contractSource != "" {

			dataEntry := make(TypeDataEntry)
			dataEntry[DataEntry_gid] = tClientCase.Gid()
			if pricingVersion != "" {
				dataEntry[FieldName_s_pricing_version] = pricingVersion
				destMap[FieldName_pricing_version] = pricingVersion
			}
			if contractSource != "" {
				dataEntry[FieldName_ContractSource] = contractSource
			}

			_, er := c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
			if er != nil {
				c.log.Error(er, " caseId:", tClientCase.Id(), " pricingVersion:", pricingVersion)
			}

		}

		for _, v := range fields {
			dealVal := tClientCase.CustomFields.TextValueByNameBasic(v)
			contactVal := tClient.CustomFields.TextValueByNameBasic(v)
			if dealVal == "" && contactVal != "" {
				destMap[v] = contactVal
			}
		}
		dealVal := tClientCase.CustomFields.TextValueByNameBasic("claims_online")
		if dealVal == "" {
			destMap["claims_online"] = "New claims:\n\n\n\nIncrease:"
		}
		destMap[FieldName_personal_statement_type] = Personal_statement_type_Webform

		if len(destMap) > 0 {
			destMap[DataEntry_gid] = clientCaseGid
			_, er := c.DataEntryUsecase.HandleOne(Kind_client_cases, destMap, DataEntry_gid, nil)
			if er != nil {
				c.log.Error(er, " ", InterfaceToString(destMap))
			}
		}

		err = c.NoPrimaryCaseInit(tClientCase)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ActionOnceUsecase) NoPrimaryCaseInit(tClientCase *TData) error {
	isPrimaryCaseCalc, _, err := c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
	if err != nil {
		return err
	}
	if !isPrimaryCaseCalc {
		primaryCase, err := c.ClientCaseUsecase.PrimaryCase(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
		if err != nil {
			return err
		}
		if primaryCase == nil {
			return errors.New("primaryCase is nil")
		}
		primaryCaseId := primaryCase.CustomFields.NumberValueByNameBasic("id")

		folderKey := fmt.Sprintf("%s%d", Map_ClientBoxFolderId, primaryCaseId)
		folderId, err := c.MapUsecase.GetForString(folderKey)
		if err != nil {
			return err
		}
		if folderId == "" {
			return errors.New("Primary Case box folder is wrong.")
		}

		caseId := tClientCase.CustomFields.NumberValueByNameBasic("id")
		caseCurrentRating := tClientCase.CustomFields.NumberValueByNameBasic(FieldName_current_rating)
		NewEvidence := fmt.Sprintf("New Evidence #%d", caseId)
		NewClaims := fmt.Sprintf("New Claims %d #%d", caseCurrentRating, caseId)
		_, err = c.BoxUsecase.CreateFolder(NewEvidence, folderId)
		if err != nil {
			return err
		}
		_, err = c.BoxUsecase.CreateFolder(NewClaims, folderId)
		if err != nil {
			return err
		}

	}
	return nil
}

// StageGettingStartedEmailToAwaitingClientFiles - Move task from "Getting Started Email" to " Awaiting Client Files"
func (c *ActionOnceUsecase) StageGettingStartedEmailToAwaitingClientFiles(clientCaseId int32) error {
	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "StageGettingStartedEmailToAwaitingClientFiles", clientCaseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val != "1" {
		tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("StageGettingStartedEmailToAwaitingClientFiles: tClientCase is nil")
		}
		clientCaseGid := tClientCase.CustomFields.TextValueByNameBasic("gid")
		if clientCaseGid == "" {
			return errors.New("StageGettingStartedEmailToAwaitingClientFiles: clientCaseGid is empty")
		}

		c.MapUsecase.Set(key, "1")
		if configs.StoppedZoho {
			dbStage := tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages)
			if dbStage == config_vbc.Stages_GettingStartedEmail {
				destData := make(TypeDataEntry)
				destData[DataEntry_gid] = clientCaseGid
				destData[FieldName_stages] = config_vbc.Stages_AwaitingClientRecords
				_, er := c.DataEntryUsecase.HandleOne(Kind_client_cases, destData, DataEntry_gid, nil)
				if er != nil {
					c.log.Error(er, " ", InterfaceToString(destData))
				}
				return er
			} else {
				return errors.New(fmt.Sprintf("CaseId:%d Stage is not \"%s\".", clientCaseId, config_vbc.Stages_GettingStartedEmail))
			}
		} else {
			dealMap, err := c.ZohoUsecase.GetDeal(clientCaseGid)
			if err != nil {
				return err
			}
			zohoStage := dealMap.GetString("Stage")
			dbStage, err := c.StageTransUsecase.BizZohoStageToDBStage(zohoStage)
			if err != nil {
				return err
			}
			if dbStage == config_vbc.Stages_GettingStartedEmail {

				zohoStage1, err := c.StageTransUsecase.DBStageToZohoStage(config_vbc.Stages_AwaitingClientRecords)
				if err != nil {
					return err
				}
				destMap := make(lib.TypeMap)
				destMap.Set("id", clientCaseGid)
				destMap.Set("Stage", zohoStage1)
				_, _, err = c.ZohoUsecase.PutRecordV1(config_zoho.Deals, destMap)
				return err
			} else {
				return errors.New(fmt.Sprintf("CaseId:%d Stage is not \"%s\".", clientCaseId, config_vbc.Stages_GettingStartedEmail))
			}
		}
	}
	return nil
}

// StageInformationIntakeToContractPending -
func (c *ActionOnceUsecase) StageInformationIntakeToContractPending(clientCaseId int32) error {

	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "StageInformationIntakeToContractPending", clientCaseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val != "1" {
		tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("StageInformationIntakeToContractPending: tClientCase is nil")
		}
		clientCaseGid := tClientCase.CustomFields.TextValueByNameBasic("gid")
		if clientCaseGid == "" {
			return errors.New("StageInformationIntakeToContractPending: clientCaseGid is empty")
		}

		c.MapUsecase.Set(key, "1")

		dbStage := tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages)
		if dbStage == config_vbc.Stages_AmInformationIntake {
			destData := make(TypeDataEntry)
			destData[DataEntry_gid] = clientCaseGid
			destData[FieldName_stages] = config_vbc.Stages_AmContractPending
			_, er := c.DataEntryUsecase.HandleOne(Kind_client_cases, destData, DataEntry_gid, nil)
			if er != nil {
				c.log.Error(er, " ", InterfaceToString(destData))
			}
			return er
		} else {
			return errors.New(fmt.Sprintf("CaseId:%d Stage is not \"%s\".", clientCaseId, config_vbc.Stages_AmInformationIntake))
		}

	}
	return nil
}

// MultiCasesBaseInfoSync - 同一个client多个client cases基本信息同步
func (c *ActionOnceUsecase) MultiCasesBaseInfoSync(clientCaseId int32) error {
	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "MultiCasesBaseInfoSync", clientCaseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val != "1" {

		tClientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": clientCaseId})
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("tClientCase  is nil")
		}

		clientGid := tClientCase.CustomFields.TextValueByNameBasic("client_gid")
		if HasEnabledPrimaryCase(clientGid) {
			usePrimaryCaseCalc, _, err := c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
			if err != nil {
				return err
			}
			if !usePrimaryCaseCalc {
				primaryCase, err := c.ClientCaseUsecase.PrimaryCase(clientGid)
				if err != nil {
					return err
				}
				if primaryCase == nil {
					return errors.New("primaryCase is nil")
				}
				caseFilesFolder := primaryCase.CustomFields.TextValueByNameBasic("case_files_folder")
				if caseFilesFolder != "" {

					if configs.StoppedZoho {

						dataEntry := make(TypeDataEntry)
						dataEntry[DataEntry_gid] = tClientCase.CustomFields.TextValueByNameBasic("gid")
						dataEntry["case_files_folder"] = caseFilesFolder
						_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
						if err != nil {
							return err
						}
					} else {
						row := make(lib.TypeMap)
						row.Set("Case_Files_Folder", caseFilesFolder)
						row.Set("id", tClientCase.CustomFields.TextValueByNameBasic("gid"))
						fmt.Println("MultiCasesBaseInfoSync:", row)
						_, _, err = c.ZohoUsecase.PutRecordV1(config_zoho.Deals, row)
						if err != nil {
							return err
						}
					}
				}
			}
		}

		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandlePersonalStatementsFile - 生成
func (c *ActionOnceUsecase) HandlePersonalStatementsFile(caseId int32) error {
	key := MapKeyPersonalStatementsFile(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err = c.BoxbuzUsecase.DoPersonalStatementsFile(caseId)
		if err != nil {
			c.log.Error(err)
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandleClaimsAnalysisFile - 生成
func (c *ActionOnceUsecase) HandleClaimsAnalysisFile(caseId int32) error {
	key := MapKeyClaimsAnalysisFile(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err = c.BoxbuzUsecase.DoClaimsAnalysisFile(caseId)
		if err != nil {
			c.log.Error(err)
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandleDoDocEmailFile - 生成
func (c *ActionOnceUsecase) HandleDoDocEmailFile(caseId int32) error {
	key := MapKeyDocEmailFile(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err = c.BoxbuzUsecase.DoDocEmailFile(caseId)
		if err != nil {
			c.log.Error(err)
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandleDoCopyDocEmailFile -
func (c *ActionOnceUsecase) HandleDoCopyDocEmailFile(caseId int32) error {
	key := MapKeyCopyDocEmailFile(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err = c.BoxbuzUsecase.DoCopyDocEmailFile(caseId)
		if err != nil {
			c.log.Error(err, " caseId: ", caseId)
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandleDoCopyReadPriorToYourDoctorVisitFile -
func (c *ActionOnceUsecase) HandleDoCopyReadPriorToYourDoctorVisitFile(caseId int32) error {
	key := MapKeyDoCopyReadPriorToYourDoctorVisitFile(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err = c.BoxbuzUsecase.DoCopyReadPriorToYourDoctorVisitFile(caseId)
		if err != nil {
			c.log.Error(err)
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandleDataCollectionFolder - 生成DataCollectionFolder
func (c *ActionOnceUsecase) HandleDataCollectionFolder(clientCaseId int32) error {
	//key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "DataCollectionFolderId", clientCaseId)
	key := MapKeyDataCollectionFolderId(clientCaseId)
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
			return errors.New("tClientCase  is nil")
		}

		_, tContactFields, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
		if err != nil {
			return err
		}
		if tContactFields == nil {
			return errors.New("tContactFields is nil.")
		}

		newFolderName := ClientCaseDataCollectionFolderNameForBox(tContactFields.TextValueByNameBasic("first_name"),
			tContactFields.TextValueByNameBasic("last_name"), clientCaseId)
		useVBCActiveFolder, parentFolderId := c.BoxbuzUsecase.GetDataCollectionFolderRootId(*tClientCase)
		boxFolderId, _, err := c.BoxUsecase.CopyFolder(c.conf.Box.NewDataCollectionFolderStructure, newFolderName, parentFolderId)

		if err != nil {
			return err
		}
		if useVBCActiveFolder {
			er := c.BoxCollaborationBuzUsecase.HandleUseVBCActiveCases(clientCaseId)
			if er != nil {
				c.log.Error(er, " HandleUseVBCActiveCases clientCaseId: ", clientCaseId)
			}
		}

		if configs.StoppedZoho {
			params := make(lib.TypeMap)
			params.Set(DataEntry_gid, tClientCase.CustomFields.TextValueByNameBasic("gid"))
			params.Set(FieldName_data_collection_folder, "https://veteranbenefitscenter.app.box.com/folder/"+boxFolderId)
			_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry(params), DataEntry_gid, nil)
			if err != nil {
				return err
			}
		} else {
			// Box Folder Link syncs to the Zoho
			params := make(lib.TypeMap)
			params.Set("id", tClientCase.CustomFields.TextValueByNameBasic("gid"))
			params.Set("Data_Collection_Folder", "https://veteranbenefitscenter.app.box.com/folder/"+boxFolderId)
			_, _, err = c.ZohoUsecase.PutRecordV1(config_zoho.Deals, params)
			if err != nil {
				return err
			}
		}

		c.MapUsecase.Set(key, boxFolderId)
		c.BoxCollaborationBuzUsecase.DoAddPermissionForBox(tClientCase.Id())
	}
	return nil
}

// HandleCopyRecordReviewFiles 处理文件
func (c *ActionOnceUsecase) HandleCopyRecordReviewFiles(clientCaseId int32) error {
	//key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "CopyRecordReviewFiles", clientCaseId)
	key := MapKeyCopyRecordReviewFiles(clientCaseId)
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
			return errors.New("tClientCase  is nil")
		}

		dCRecordReviewFolderId, err := c.BoxbuzUsecase.DCRecordReviewFolderId(tClientCase)
		if err != nil {
			return err
		}

		primaryCase, err := c.ClientCaseUsecase.PrimaryCase(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
		if err != nil {
			return err
		}
		var destCaseId int32
		if primaryCase != nil {
			destCaseId = primaryCase.CustomFields.NumberValueByNameBasic("id")

			// 需要处理New Evidence
			if primaryCase.CustomFields.NumberValueByNameBasic("id") !=
				tClientCase.CustomFields.NumberValueByNameBasic("id") {

				newEvidenceFolderId, err := c.BoxcUsecase.GetNewEvidenceFolderId(primaryCase, tClientCase)
				if err != nil {
					return err
				}
				if newEvidenceFolderId == "" {
					return errors.New("newEvidenceFolderId is empty")
				}
				subItems, err := c.BoxUsecase.ListItemsInFolderFormat(newEvidenceFolderId)
				if err != nil {
					return err
				}
				err = c.BoxbuzUsecase.CopyBoxResItemsToFolder(dCRecordReviewFolderId, subItems, "")
				if err != nil {
					return err
				}
			}

		} else {
			destCaseId = tClientCase.CustomFields.NumberValueByNameBasic("id")
		}

		PrivateMedicalRecordsFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
			MapKeyBuildAutoBoxCPrivateMedicalRecordsFolderId(destCaseId),
			tClientCase)
		if err != nil {
			return err
		}
		ServiceTreatmentRecordsFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
			MapKeyBuildAutoBoxCServiceTreatmentRecordsFolderId(destCaseId),
			tClientCase)
		if err != nil {
			return err
		}
		VAMedicalRecordsFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
			MapKeyBuildAutoBoxCVAMedicalRecordsFolderId(destCaseId),
			tClientCase)
		if err != nil {
			return err
		}
		_, _, err = c.BoxUsecase.CopyFolder(PrivateMedicalRecordsFolderId, config_box.FolderName_PrivateMedicalRecords, dCRecordReviewFolderId)
		if err != nil {
			return err
		}
		_, _, err = c.BoxUsecase.CopyFolder(ServiceTreatmentRecordsFolderId, config_box.FolderName_ServiceTreatmentRecords, dCRecordReviewFolderId)
		if err != nil {
			return err
		}
		_, _, err = c.BoxUsecase.CopyFolder(VAMedicalRecordsFolderId, config_box.FolderName_VAMedicalRecords, dCRecordReviewFolderId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandleMedicalTeamFormsTest (本地测试使用)发送Release Of Information 合同
func (c *ActionOnceUsecase) HandleMedicalTeamFormsTest(clientCaseId int32) error {
	//key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "ReleaseOfInformation", clientCaseId)

	tClientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": clientCaseId})
	if err != nil {
		return err
	}
	if tClientCase == nil {
		return errors.New("tClientCase  is nil")
	}

	tClient, _, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	prefillTags, boxSignTplId, err := c.DbqsUsecase.MedicalTeamFormsPrefillTagsV2(tClientCase, tClient)
	if err != nil {
		return err
	}
	boxContractFolderId, err := c.BoxcontractUsecase.ContractFolderId(clientCaseId)
	if err != nil {
		return err
	}
	signerEmail := tClient.CustomFields.TextValueByNameBasic(FieldName_email)
	leadVSEmail, err := c.DbqsUsecase.BizLeadVSEmail(tClientCase)

	// test code
	signerEmail = "liaogling@gmail.com"
	leadVSEmail = "glliao@vetbenefitscenter.com"
	boxContractFolderId = "288896562235"

	if err != nil {
		return err
	}
	//return nil
	res, contractId, err := c.BoxUsecase.MedicalTeamFormsSignRequests(boxContractFolderId, signerEmail, leadVSEmail, prefillTags, MedicalTeamFormsV2, boxSignTplId)
	if err != nil {
		return err
	}
	c.log.Info("res: ", res, " contractId: ", contractId, " err: ", err)

	return nil
}

// HandleMedicalTeamForms 发送Release Of Information 合同
func (c *ActionOnceUsecase) HandleMedicalTeamForms(clientCaseId int32) error {

	return c.HandleMedicalTeamFormsWithoutTemplate(clientCaseId)

	//key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "ReleaseOfInformation", clientCaseId)
	key := MapKeyMedicalTeamForms(clientCaseId)
	// key的值存的是：合同ID，很重要
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
			return errors.New("tClientCase  is nil")
		}

		tClient, _, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
		if err != nil {
			return err
		}
		if tClient == nil {
			return errors.New("tClient is nil")
		}

		var prefillTags lib.TypeList
		var boxSignTplId string

		if true || tClient.Gid() == "6159272000005519042" { // 发布到线上

			prefillTags, boxSignTplId, err = c.DbqsUsecase.MedicalTeamFormsPrefillTagsV2(tClientCase, tClient)
		} else {
			prefillTags, boxSignTplId, err = c.DbqsUsecase.MedicalTeamFormsPrefillTags(tClientCase, tClient)
		}

		if err != nil {
			return err
		}
		boxContractFolderId, err := c.BoxcontractUsecase.ContractFolderId(clientCaseId)
		if err != nil {
			return err
		}
		signerEmail := tClient.CustomFields.TextValueByNameBasic(FieldName_email)
		leadVSEmail, err := c.DbqsUsecase.BizLeadVSEmail(tClientCase)
		if err != nil {
			return err
		}
		res, contractId, err := c.BoxUsecase.MedicalTeamFormsSignRequests(boxContractFolderId, signerEmail, leadVSEmail, prefillTags, "", boxSignTplId)
		if err != nil {
			return err
		}
		str := ""
		if res != nil {
			str = *res
		}
		err = c.ClientEnvelopeUsecase.Add(clientCaseId, EsignVendor_box, contractId, str, Type_MedicalTeamForms, 0)
		if err == nil {
			err = c.RollpoingUsecase.Upsert(Rollpoing_Vendor_boxsign, contractId)
		}

		c.BehaviorUsecase.Add(clientCaseId, BehaviorType_sent_medical_team_forms_contract, time.Now(), "")

		c.TaskCreateUsecase.CreateTask(clientCaseId,
			map[string]interface{}{"CaseId": clientCaseId},
			Task_Dag_ReminderMedicalTeamFormsContractSent, 0, "", "")

		c.MapUsecase.Set(key, contractId)
	}
	return nil
}

func (c *ActionOnceUsecase) HandleMedicalTeamFormsWithoutTemplate(clientCaseId int32) error {

	//key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "ReleaseOfInformation", clientCaseId)
	key := MapKeyMedicalTeamForms(clientCaseId)
	// key的值存的是：合同ID，很重要
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
			return errors.New("tClientCase  is nil")
		}

		tClient, _, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
		if err != nil {
			return err
		}
		if tClient == nil {
			return errors.New("tClient is nil")
		}

		boxContractFolderId, err := c.BoxcontractUsecase.ContractFolderId(clientCaseId)
		if err != nil {
			return err
		}
		signerEmail := tClient.CustomFields.TextValueByNameBasic(FieldName_email)
		leadVSEmail, err := c.DbqsUsecase.BizLeadVSEmail(tClientCase)
		if err != nil {
			return err
		}
		leadCPEmail, err := c.DbqsUsecase.BizLeadCPEmail(tClientCase)
		if err != nil {
			c.log.Error(err)
			err = nil
		}

		createMedicalTeamFormVo, err := c.DbqsUsecase.MedicalTeamFormsPrefillTagsWithoutTemplate(tClientCase, tClient)
		if err != nil {
			c.log.Error(err)
			return err
		}
		signFileBytes, err := c.GopdfUsecase.CreateMedicalTeamForm(createMedicalTeamFormVo)
		if err != nil {
			c.log.Error(err)
			return err
		}
		folderName := uuid.UuidWithoutStrike()
		signFolderId, err := c.BoxUsecase.CreateFolder(folderName, boxContractFolderId)
		if err != nil {
			c.log.Error(err)
			return err
		}
		signFileId, err := c.BoxUsecase.UploadFile(signFolderId, bytes.NewReader(signFileBytes), "Medical Team Forms.pdf")
		if err != nil {
			return err
		}
		res, contractId, err := c.BoxUsecase.MedicalTeamFormsSignRequestsWithoutTemplate(boxContractFolderId, signerEmail, leadVSEmail, leadCPEmail, signFileId)

		//res, contractId, err := c.BoxUsecase.MedicalTeamFormsSignRequests(boxContractFolderId, signerEmail, leadVSEmail, prefillTags, "", boxSignTplId)
		if err != nil {
			return err
		}
		str := ""
		if res != nil {
			str = *res
		}
		err = c.ClientEnvelopeUsecase.Add(clientCaseId, EsignVendor_box, contractId, str, Type_MedicalTeamForms, 0)
		if err == nil {
			err = c.RollpoingUsecase.Upsert(Rollpoing_Vendor_boxsign, contractId)
		}

		c.BehaviorUsecase.Add(clientCaseId, BehaviorType_sent_medical_team_forms_contract, time.Now(), "")

		c.TaskCreateUsecase.CreateTask(clientCaseId,
			map[string]interface{}{"CaseId": clientCaseId},
			Task_Dag_ReminderMedicalTeamFormsContractSent, 0, "", "")

		c.MapUsecase.Set(key, contractId)
	}
	return nil
}

func (c *ActionOnceUsecase) HandleMedicalTeamFormsReminderEmail(clientCaseId int32) error {

	key := MapKeyMedicalTeamFormsReminderEmail(clientCaseId)
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
			return errors.New("tClientCase  is nil")
		}

		err = c.TaskCreateUsecase.CreateTaskMail(clientCaseId, MailGenre_MedicalExamDocumentsReminder, 0, nil, 0, "", "")
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandlePrivateExamsSubmitted Stage 15.触发
func (c *ActionOnceUsecase) HandlePrivateExamsSubmitted(ctx context.Context, clientCaseId int32) error {

	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "HandlePrivateExamsSubmitted", clientCaseId)

	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	needNotifyCPTeam := false
	if val == "" {

		tCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if err != nil {
			c.log.Error(err)
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}

		//if tCase.CustomFields.NumberValueByNameBasic("id") != 5101 {
		//	return nil
		//}

		tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
		if err != nil {
			c.log.Error(err)
			return err
		}
		if tClient == nil {
			return errors.New("tClient is nil")
		}
		err = c.HandlePrivateExamsSubmittedFirstStep(tCase)
		if err != nil {
			needNotifyCPTeam = true
			c.log.Error("HandlePrivateExamsSubmittedFirstStep: ", tCase.CustomFields.NumberValueByNameBasic("id"), " err: ", err)
		}
		_, err = c.GoogleDrivebuzUsecase.TransferPaymentForm(ctx, tCase, tClient)
		if err != nil {
			needNotifyCPTeam = true
			c.log.Error(err)
			//return err
		}

		time.Sleep(2 * time.Second)
		err = c.GoogleDrivebuzUsecase.TransferPsych(ctx, tCase, tClient)
		if err != nil {
			needNotifyCPTeam = true
			c.log.Error(err)
			//return err
		}
		time.Sleep(2 * time.Second)
		err = c.GoogleDrivebuzUsecase.TransferGeneral(ctx, tCase, tClient)
		if err != nil {
			needNotifyCPTeam = true
			c.log.Error(err)
			//return err
		}

		c.MapUsecase.Set(key, "1")

		if needNotifyCPTeam {
			er := c.RemindUsecase.CreateTaskForSubmissionToGoogleDriveFailed(*tCase)
			if er != nil {
				c.log.Error(er, " ", tCase.Id())
			}
		}

	}
	return nil
}

// HandlePrivateExamsSubmittedFirstStep 处理第一步
func (c *ActionOnceUsecase) HandlePrivateExamsSubmittedFirstStep(tCase *TData) error {
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	clientCaseId := tCase.CustomFields.NumberValueByNameBasic("id")
	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "HandlePrivateExamsSubmittedFirstStep", clientCaseId)

	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		DD214FolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(MapKeyBuildAutoBoxCDD214FolderId(clientCaseId), tCase)
		if err != nil {
			return err
		}
		if DD214FolderId == "" {
			return errors.New("DD214FolderId is empty")
		}

		DisabilityRatingListFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(MapKeyBuildAutoBoxCDisabilityRatingListFolderId(clientCaseId), tCase)
		if err != nil {
			return err
		}
		if DisabilityRatingListFolderId == "" {
			return errors.New("DisabilityRatingListFolderId is empty")
		}

		RatingDecisionLettersFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(MapKeyBuildAutoBoxCRatingDecisionLettersFolderId(clientCaseId), tCase)
		if err != nil {
			return err
		}
		if RatingDecisionLettersFolderId == "" {
			return errors.New("RatingDecisionLettersFolderId is empty")
		}

		DCPrivateExamsFolderId, err := c.BoxbuzUsecase.GetDCSubFolderId(MapKeyBuildAutoBoxDCPrivateExamsFolderId(
			clientCaseId), tCase)
		if err != nil {
			return err
		}
		if DCPrivateExamsFolderId == "" {
			return errors.New("DCPrivateExamsFolderId is empty")
		}

		DD214Items, err := c.BoxUsecase.ListItemsInFolderFormat(DD214FolderId)
		if err != nil {
			return err
		}
		err = c.BoxbuzUsecase.CopyBoxResItemsToFolder(DCPrivateExamsFolderId, DD214Items, CopyBoxResItemsToFolder_Type_file_only_and_ignore_409)
		if err != nil {
			return errors.New(err.Error() + ":DD214Items")
		}

		DisabilityRatingListItems, err := c.BoxUsecase.ListItemsInFolderFormat(DisabilityRatingListFolderId)
		if err != nil {
			return err
		}
		err = c.BoxbuzUsecase.CopyBoxResItemsToFolder(DCPrivateExamsFolderId, DisabilityRatingListItems, CopyBoxResItemsToFolder_Type_file_only_and_ignore_409)
		if err != nil {
			return errors.New(err.Error() + ":DisabilityRatingListItems")
		}

		RatingDecisionLettersItems, err := c.BoxUsecase.ListItemsInFolderFormat(RatingDecisionLettersFolderId)
		if err != nil {
			return err
		}
		err = c.BoxbuzUsecase.CopyBoxResItemsToFolder(DCPrivateExamsFolderId, RatingDecisionLettersItems, CopyBoxResItemsToFolder_Type_file_only_and_ignore_409)
		if err != nil {
			return errors.New(err.Error() + ":RatingDecisionLettersItems")
		}

		PrivateMedicalRecordsFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
			MapKeyBuildAutoBoxCPrivateMedicalRecordsFolderId(clientCaseId),
			tCase)
		if err != nil {
			return err
		}
		ServiceTreatmentRecordsFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
			MapKeyBuildAutoBoxCServiceTreatmentRecordsFolderId(clientCaseId),
			tCase)
		if err != nil {
			return err
		}
		VAMedicalRecordsFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
			MapKeyBuildAutoBoxCVAMedicalRecordsFolderId(clientCaseId),
			tCase)

		PrivateMedicalRecordsFolderIdItems, err := c.BoxUsecase.ListItemsInFolderFormat(PrivateMedicalRecordsFolderId)
		if err != nil {
			return err
		}
		err = c.BoxbuzUsecase.CopyBoxResItemsToFolder(DCPrivateExamsFolderId, PrivateMedicalRecordsFolderIdItems, CopyBoxResItemsToFolder_Type_file_only_and_ignore_409)
		if err != nil {
			return errors.New(err.Error() + ":PrivateMedicalRecordsFolderIdItems")
		}

		ServiceTreatmentRecordsFolderIdItems, err := c.BoxUsecase.ListItemsInFolderFormat(ServiceTreatmentRecordsFolderId)
		if err != nil {
			return err
		}
		err = c.BoxbuzUsecase.CopyBoxResItemsToFolder(DCPrivateExamsFolderId, ServiceTreatmentRecordsFolderIdItems, CopyBoxResItemsToFolder_Type_file_only_and_ignore_409)
		if err != nil {
			return errors.New(err.Error() + ":ServiceTreatmentRecordsFolderIdItems")
		}

		VAMedicalRecordsFolderIdItems, err := c.BoxUsecase.ListItemsInFolderFormat(VAMedicalRecordsFolderId)
		if err != nil {
			return err
		}
		err = c.BoxbuzUsecase.CopyBoxResItemsToFolder(DCPrivateExamsFolderId, VAMedicalRecordsFolderIdItems, CopyBoxResItemsToFolder_Type_file_only_and_ignore_409)
		if err != nil {
			return errors.New(err.Error() + ":VAMedicalRecordsFolderIdItems")
		}

		c.MapUsecase.Set(key, "1")

	}
	return nil
}

func (c *ActionOnceUsecase) HandleUpcomingContactInformation(caseId int32) error {

	key := MapKeyUpcomingContactInformation(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		err = c.DoUpcomingContactInformation(caseId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ActionOnceUsecase) DoUpcomingContactInformation(caseId int32) error {

	c.log.Info("DoUpcomingContactInformation caseId:", caseId)
	leadVSChangeLogKey := MapKeyLeadVSChangeLog(caseId)

	changeLog, err := c.MapUsecase.GetForString(leadVSChangeLogKey)
	if err != nil {
		return err
	}
	if changeLog == "" {
		return nil
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	leadVSChangeLogVo := lib.StringToTDef(changeLog, LeadVSChangeLogVo{})
	if leadVSChangeLogVo.PreviousVSUserGid == tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid) {
		return nil
	}

	err = c.TaskCreateUsecase.CreateTaskMail(caseId, MailGenre_UpcomingContactInformation, 0, lib.TypeMap{
		"LeadVSChangeLog": changeLog,
	}, 0, "", "")
	if err != nil {
		return err
	}

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	// 做延时处理
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)
	c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextUpcomingContactInformation,
		Params: lib.TypeMap{
			"LeadVSChangeLog": changeLog,
		},
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))

	return nil
}

func (c *ActionOnceUsecase) HandleEmailMiniDBQsDrafts(caseId int32) error {

	key := MapKeyEmailMiniDBQsDrafts(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		err = c.DoEmailMiniDBQsDrafts(caseId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ActionOnceUsecase) DoEmailMiniDBQsDrafts(caseId int32) error {
	err := c.TaskCreateUsecase.CreateTaskMail(caseId, MailGenre_MiniDBQsDrafts, 0, nil, 0, "", "")
	if err != nil {
		return err
	}
	return nil
}

func (c *ActionOnceUsecase) HandleYourRecordsReviewProcessHasBegun(caseId int32) error {

	key := MapKeyYourRecordsReviewProcessHasBegun(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		err = c.DoYourRecordsReviewProcessHasBegun(caseId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ActionOnceUsecase) DoYourRecordsReviewProcessHasBegun(caseId int32) error {
	err := c.TaskCreateUsecase.CreateTaskMail(caseId, MailGenre_YourRecordsReviewProcessHasBegun, 0, nil, 0, "", "")
	if err != nil {
		return err
	}
	return nil
}

func (c *ActionOnceUsecase) HandlePleaseScheduleYourDoctorAppointments(caseId int32) error {

	key := MapKeyPleaseScheduleYourDoctorAppointments(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		err = c.DoPleaseScheduleYourDoctorAppointments(caseId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ActionOnceUsecase) DoPleaseScheduleYourDoctorAppointments(caseId int32) error {
	err := c.TaskCreateUsecase.CreateTaskMail(caseId, MailGenre_PleaseScheduleYourDoctorAppointments, 0, nil, 0, "", "")
	if err != nil {
		return err
	}
	return nil
}

func (c *ActionOnceUsecase) HandlePersonalStatementsReadyforYourReview(caseId int32) error {

	key := MapKeyPersonalStatementsReadyforYourReview(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		err = c.DoPersonalStatementsReadyforYourReview(caseId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// DoPersonalStatementsReadyforYourReview 立即发送
func (c *ActionOnceUsecase) DoPersonalStatementsReadyforYourReview(caseId int32) error {
	err := c.TaskCreateUsecase.CreateTaskMail(caseId, MailGenre_PersonalStatementsReadyforYourReview, 0, nil, 0, Task_FromType_AutomationCrontabEmail, InterfaceToString(caseId))
	if err != nil {
		return err
	}
	return nil
}

func (c *ActionOnceUsecase) HandleVAForm2122aSubmission(caseId int32) error {
	// 关闭
	return nil
	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "HandleVAForm2122aSubmission", caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		err = c.DoVAForm2122aSubmission(caseId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// DoVAForm2122aSubmission 立即发送
func (c *ActionOnceUsecase) DoVAForm2122aSubmission(caseId int32) error {
	err := c.TaskCreateUsecase.CreateTaskMail(caseId, MailGenre_VAForm2122aSubmission, 0, nil, 0, "", InterfaceToString(caseId))
	if err != nil {
		return err
	}
	return nil
}

// HandlePleaseReviewYourPersonalStatementsinSharedFolder 延时到14天后发送
func (c *ActionOnceUsecase) HandlePleaseReviewYourPersonalStatementsinSharedFolder(caseId int32) error {

	key := MapKeyPleaseReviewYourPersonalStatementsinSharedFolder(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		err = c.DoPleaseReviewYourPersonalStatementsinSharedFolder(caseId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ActionOnceUsecase) DoPleaseReviewYourPersonalStatementsinSharedFolder(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	timeAt, err := GetStatementsFinalizedEvery14DaysTime(tCase, c.log)
	if err != nil {
		c.log.Error(err)
		return err
	}
	err = c.TaskCreateUsecase.CreateTaskMail(caseId, MailGenre_PleaseReviewYourPersonalStatementsinSharedFolder, 0, nil, timeAt.Unix(), Task_FromType_AutomationCrontabEmail, InterfaceToString(caseId))
	if err != nil {
		return err
	}
	return nil
}

// CancelAutomationCrontabEmailTasks 取消任务，阶段修改后，后续任务不需要了
func (c *ActionOnceUsecase) CancelAutomationCrontabEmailTasks(caseId int32) error {

	return c.CommonUsecase.DB().Model(&TaskEntity{}).
		Where("from_id = ? and from_type =? and task_status=?",
			InterfaceToString(caseId),
			Task_FromType_AutomationCrontabEmail,
			Task_TaskStatus_processing).
		Updates(map[string]interface{}{
			"task_status": Task_TaskStatus_cancel,
			"updated_at":  time.Now().Unix()}).Error

}

func (c *ActionOnceUsecase) HandleCopyPersonalStatementsDoc(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	useNewPersonalWebForm, err := c.PersonalWebformUsecase.IsUseNewPersonalWebForm(caseId)
	if err != nil {
		return err
	}
	if useNewPersonalWebForm {
		return nil
	}

	key := MapKeyCopyPersonalStatementsDoc(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {

		err = c.DoCopyPersonalStatementsDoc(caseId)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ActionOnceUsecase) DoCopyPersonalStatementsDoc(caseId int32) error {

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

	clientPSFolderId, err := c.BoxbuzUsecase.CPersonalStatementsFolderIdByAnyCase(*tCase)
	if err != nil {
		return err
	}
	if clientPSFolderId == "" {
		return errors.New("clientPSFolderId is empty")
	}

	_, psDocBoxFileId, err := c.BoxbuzUsecase.PersonalStatementDocFileBoxFileId(tClient, tCase)
	if err != nil {
		return err
	}
	if psDocBoxFileId == "" {
		return errors.New("psDocBoxFileId is empty")
	}
	_, _, err = c.BoxUsecase.CopyFile(psDocBoxFileId, clientPSFolderId)
	if err != nil {
		return err
	}
	return nil
}

func CanHelpUsImproveSurvey(stages string) bool {
	if stages == config_vbc.Stages_AwaitingPayment ||
		stages == config_vbc.Stages_27_AwaitingBankReconciliation ||
		stages == config_vbc.Stages_AmAwaitingPayment ||
		stages == config_vbc.Stages_Am27_AwaitingBankReconciliation {
		return true
	}
	return false
}

// HandleHelpUsImproveSurvey - Survey
func (c *ActionOnceUsecase) HandleHelpUsImproveSurvey(caseId int32) error {
	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "HelpUsImprove", caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val != "" {
		return nil
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	newRating := tCase.CustomFields.NumberValueByNameBasic(FieldName_new_rating)
	lib.DPrintln(stages, newRating)
	if CanHelpUsImproveSurvey(stages) && newRating == 100 {
		err = c.DoHandleHelpUsImproveSurvey(*tCase)
		if err != nil {
			c.log.Error(err)
			return err
		}
	} else {
		return nil
	}
	c.MapUsecase.Set(key, "1")
	return nil
}

func (c *ActionOnceUsecase) DoHandleHelpUsImproveSurvey(tCase TData) error {

	return c.TaskCreateUsecase.CreateTaskMail(tCase.Id(), MailGenre_HelpUsImproveSurvey, 0, nil, 0, "", "")
}

func FormatFullName(firstName, lastName string) string {
	if lastName == "" {
		return firstName
	}
	return firstName + " " + lastName
}
