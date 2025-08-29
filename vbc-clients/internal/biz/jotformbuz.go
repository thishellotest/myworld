package biz

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"sort"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

type JotformbuzUsecase struct {
	log                      *log.Helper
	CommonUsecase            *CommonUsecase
	conf                     *conf.Data
	JotformUsecase           *JotformUsecase
	JotformSubmissionUsecase *JotformSubmissionUsecase
	TUsecase                 *TUsecase
	DataComboUsecase         *DataComboUsecase
	ZohoUsecase              *ZohoUsecase
	ActionOnceUsecase        *ActionOnceUsecase
	BehaviorUsecase          *BehaviorUsecase
	BoxbuzUsecase            *BoxbuzUsecase
	LogUsecase               *LogUsecase
	BoxUsecase               *BoxUsecase
	MapUsecase               *MapUsecase
	DataEntryUsecase         *DataEntryUsecase
	FieldOptionUsecase       *FieldOptionUsecase
	FieldUsecase             *FieldUsecase
}

func NewJotformbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	JotformUsecase *JotformUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	TUsecase *TUsecase,
	DataComboUsecase *DataComboUsecase,
	ZohoUsecase *ZohoUsecase,
	ActionOnceUsecase *ActionOnceUsecase,
	BehaviorUsecase *BehaviorUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	LogUsecase *LogUsecase,
	BoxUsecase *BoxUsecase,
	MapUsecase *MapUsecase,
	DataEntryUsecase *DataEntryUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
	FieldUsecase *FieldUsecase) *JotformbuzUsecase {
	uc := &JotformbuzUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		JotformUsecase:           JotformUsecase,
		JotformSubmissionUsecase: JotformSubmissionUsecase,
		TUsecase:                 TUsecase,
		DataComboUsecase:         DataComboUsecase,
		ZohoUsecase:              ZohoUsecase,
		ActionOnceUsecase:        ActionOnceUsecase,
		BehaviorUsecase:          BehaviorUsecase,
		BoxbuzUsecase:            BoxbuzUsecase,
		LogUsecase:               LogUsecase,
		BoxUsecase:               BoxUsecase,
		MapUsecase:               MapUsecase,
		DataEntryUsecase:         DataEntryUsecase,
		FieldOptionUsecase:       FieldOptionUsecase,
		FieldUsecase:             FieldUsecase,
	}

	return uc
}

func GetSubmissionInfo(jotformInfo lib.TypeMap) (formId, VBCCaseID, firstName, lastName string, err error) {

	answers := jotformInfo.GetTypeMap("content.answers")
	formId = jotformInfo.GetString("content.form_id")
	for _, v := range answers {
		fieldInfo := lib.ToTypeMap(v)
		name := fieldInfo.GetString("name")
		if name == "vbcCase" {
			VBCCaseID = fieldInfo.GetString("answer")
		} else if name == "name" {
			firstName = fieldInfo.GetString("answer.first")
			lastName = fieldInfo.GetString("answer.last")
		}
	}
	if VBCCaseID == "" {
		//return "", "", "", errors.New("VBCCaseID is empty")
	}
	if firstName == "" {
		//return "", "", "", errors.New("firstName is empty")
	}
	if lastName == "" {
		//return "", "", "", errors.New("lastName is empty")
	}
	return
}

func (c *JotformbuzUsecase) HandleSubmission(submissionId string, VBCCaseID string, formId string) error {
	jotformInfo, err := c.JotformUsecase.GetSubmission(submissionId)
	if err != nil {
		return err
	}

	newFormId, VBCCaseIDFromInfo, firstName, lastName, err := GetSubmissionInfo(jotformInfo)
	if formId == "" {
		formId = newFormId
	}
	if err != nil {
		c.log.Error(err)
		return err
	}
	//needHandleMapping := false
	isVBCIntakeFormID := false
	isAmIntakeFormID := false
	if VBCCaseID == "" { // 说明是特别的intake form

		if formId == JotformIntakeFormID {
			//needHandleMapping = true
			isVBCIntakeFormID = true
		}
		if formId == JotformAmIntakeFormID {
			isAmIntakeFormID = true
		}
		VBCCaseID = VBCCaseIDFromInfo
	}

	entity := &JotformSubmissionEntity{
		SubmissionId: submissionId,
		FormId:       formId,
		Uniqcode:     VBCCaseID,
		Notes:        InterfaceToString(jotformInfo),
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = c.JotformSubmissionUsecase.CommonUsecase.DB().Save(entity).Error
	if err != nil {
		return err
	}

	answers := jotformInfo.GetTypeMap("content.answers")
	//for k, v := range answers {
	//	fieldInfo := lib.ToTypeMap(v)
	//	name := fieldInfo.GetString("name")
	//	if name == "vbcCase" {
	//		VBCCaseID = fieldInfo.GetString("answer")
	//	} else if name == "name" {
	//		firstName = fieldInfo.GetString("answer.first")
	//		lastName = fieldInfo.GetString("answer.last")
	//	}
	//}

	if isVBCIntakeFormID || isAmIntakeFormID {

		if len(firstName) == 0 {
			return errors.New("firstName is empty")
		}
		if len(lastName) == 0 {
			return errors.New("lastName is empty")
		}

		tClientCase, err := c.TUsecase.Data(Kind_client_cases, And(Eq{FieldName_uniqcode: VBCCaseID}))
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("The tClientCase does not exist via VBC Case ID.")
		}
		clientCaseGid := tClientCase.CustomFields.TextValueByNameBasic("gid")
		if len(clientCaseGid) == 0 {
			return errors.New("clientCaseGid is empty")
		}

		_, tContactFields, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
		if err != nil {
			return err
		}
		if tContactFields == nil {
			return errors.New("tContactFields is nil.")
		}
		clientGid := tContactFields.TextValueByNameBasic("gid")
		if clientGid == "" {
			return errors.New("clientGid is empty.")
		}

		contactMap, dbContactMap := c.HandeMappingContact(answers, isAmIntakeFormID)
		if len(contactMap) == 0 {
			return errors.New("contactMap Nothing need modify.")
		}

		lib.DPrintln("HandleSubmission contactMap:", contactMap, dbContactMap)

		destData := make(TypeDataEntry)
		for k, v := range contactMap {
			vbcFieldName := config_zoho.ZohoContactVbcFieldNameByZohoFieldName(k)
			if vbcFieldName != "" {
				destData[vbcFieldName] = v
			}
		}
		for k, v := range dbContactMap {
			destData[k] = v
		}
		if len(destData) > 0 {
			destData[DataEntry_gid] = clientGid
		}
		_, er := c.DataEntryUsecase.HandleOne(Kind_clients, destData, DataEntry_gid, nil)
		if er != nil {
			c.log.Error(er, " ", InterfaceToString(destData), " ", InterfaceToString(contactMap))
		}

		err = c.HandleExecDeal(clientCaseGid, answers, VBCCaseID, isAmIntakeFormID)
		if err != nil {
			return err
		}

		if isVBCIntakeFormID {
			err = c.BehaviorUsecase.Add(tClientCase.CustomFields.NumberValueByNameBasic("id"),
				BehaviorType_complete_intake_form, time.Now(), "")
			if err != nil {
				return err
			} else {
				// - Move task from "Getting Started Email" to " Awaiting Client Files"
				err = c.ActionOnceUsecase.StageGettingStartedEmailToAwaitingClientFiles(tClientCase.CustomFields.NumberValueByNameBasic("id"))
				if err != nil {
					return err
				}
			}
		} else if isAmIntakeFormID {

			err = c.BehaviorUsecase.Add(tClientCase.CustomFields.NumberValueByNameBasic("id"),
				BehaviorType_complete_am_intake_form, time.Now(), "")
			if err != nil {
				return err
			} else {
				// - Move task from "Getting Started Email" to " Awaiting Client Files"
				err = c.ActionOnceUsecase.StageInformationIntakeToContractPending(tClientCase.CustomFields.NumberValueByNameBasic("id"))
				if err != nil {
					return err
				}
			}

		}

		lib.DPrintln("VBCCaseID:", VBCCaseID)

	}
	return nil
}

func (c *JotformbuzUsecase) HandeMappingContact(answers lib.TypeMap, isAmContract bool) (destMap lib.TypeMap, dbDestMap lib.TypeMap) {
	destMap = make(lib.TypeMap) // 已经zoho使用，慢慢可以转为dbDestMap
	dbDestMap = make(lib.TypeMap)

	for _, v := range answers {
		fieldInfo := lib.ToTypeMap(v)
		name := fieldInfo.GetString("name")
		if name == "address" {
			addrLine1 := fieldInfo.GetString("answer.addr_line1")
			addrLine2 := fieldInfo.GetString("answer.addr_line2")
			//if len(addrLine2) > 0 {
			//	if len(addrLine1) > 0 {
			//		addrLine1 += ", " + addrLine2
			//	} else {
			//		addrLine1 = addrLine2
			//	}
			//}
			city := fieldInfo.GetString("answer.city")
			zipCode := fieldInfo.GetString("answer.postal")
			state := fieldInfo.GetString("answer.state")

			if len(city) > 0 {
				destMap.Set("Mailing_City", city)
			}
			if len(state) > 0 {
				destMap.Set("Mailing_State", state)
			}
			if len(zipCode) > 0 {
				destMap.Set("Mailing_Zip", zipCode)
			}
			if len(addrLine1) > 0 {
				destMap.Set("Mailing_Street", addrLine1)
			}
			if len(addrLine2) > 0 {
				dbDestMap.Set("apt_number", addrLine2)
			}

		} else if name == "phoneNumber" {
			phoneNumber := fieldInfo.GetString("answer.full")
			if phoneNumber != "" {
				destMap.Set("Mobile", phoneNumber)
			}
		} else if name == "ssn49" {
			ssn := fieldInfo.GetString("answer")
			if ssn != "" {
				destMap.Set("SSN", ssn)
			}
		} else if name == "dateOf" {
			year := fieldInfo.GetString("answer.year")
			month := fieldInfo.GetString("answer.month")
			day := fieldInfo.GetString("answer.day")
			if year != "" {
				destMap.Set("Date_of_Birth", fmt.Sprintf("%s-%s-%s", year, month, day))
			}
		} else if name == Jotform_overallRating {
			if len(fieldInfo.GetString("answer")) > 0 {
				if !isAmContract {
					//destMap.Set("Current_Rating", fieldInfo.GetInt("answer"))
				}
			}
		} else if name == "overallRating65" {
			if len(fieldInfo.GetString("answer")) > 0 {
				dbDestMap[FieldName_current_rating] = fieldInfo.GetInt("answer")
				//destMap.Set("Current_Rating", fieldInfo.GetInt("answer"))
			}
		} else if name == Jotform_branchOf {
			answer := fieldInfo.GetString("answer")
			answer = c.FormResponseBranch(answer)
			if answer != "" {
				destMap.Set("Branch", answer)
			}
		} else if name == "placeOf" {
			city := fieldInfo.GetString("answer.city")
			state := fieldInfo.GetString("answer.state")
			country := fieldInfo.GetString("answer.country")
			if len(city) > 0 {
				destMap.Set("Place_of_Birth_City", city)
			}
			if len(state) > 0 {
				destMap.Set("Place_of_Birth_State_Province", state)
			}
			if len(country) > 0 {
				destMap.Set("Place_of_Birth_Country", country)
			}
		} else if name == "currentOccupation" {
			answer := fieldInfo.GetString("answer")
			if answer != "" {
				destMap.Set("Current_Occupation", answer)
			}
		} else if name == "timeZone" {
			answer := fieldInfo.GetString("answer")
			if answer != "" {
				timezoneFieldEntity, _ := c.FieldUsecase.GetByFieldName(Kind_clients, FieldName_timezone_id)
				if timezoneFieldEntity != nil {
					optionEntity, _ := c.FieldOptionUsecase.GetByEntityAndLabel(*timezoneFieldEntity, answer)
					if optionEntity != nil {
						dbDestMap.Set(FieldName_timezone_id, optionEntity.OptionValue)
					}
				}
			}
		}
		if isAmContract {
			if name == "doYou54" {
				answer := fieldInfo.GetString("answer")
				if answer == "Yes" {
					dbDestMap.Set(FieldName_pending_claims, "Yes")
				} else {
					dbDestMap.Set(FieldName_pending_claims, "No")
				}
			}
		}
	}
	return destMap, dbDestMap
}

const (
	Jotform_branchOf      = "branchOf"
	Jotform_overallRating = "overallRating"
)

func (c *JotformbuzUsecase) FormResponseBranch(jotformBranch string) string {
	// todo:lgl jotformBranch
	return jotformBranch
}

func (c *JotformbuzUsecase) HandleExecDeal(clientCaseGid string, answers lib.TypeMap, uniqcode string, isAmIntakeFormID bool) error {

	//DealFieldInfos := config_zoho.DealLayout().DealFieldInfos()

	//military := lib.InterfaceToTDef[[]string](dataMap.Get("Military Toxic Exposures"), nil)
	destMap := make(lib.TypeMap)

	//
	//destMap.Set("Agent_Orange_Exposure", "NO")
	//destMap.Set("Gulf_War_Illness", "NO")
	//destMap.Set("Burn_Pits_and_Other_Airborne_Hazards", "NO")
	//destMap.Set("Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun", "NO")
	//destMap.Set("Atomic_Veterans_and_Radiation_Exposure", "NO")
	//destMap.Set("Amyotrophic_Lateral_Sclerosis_ALS", "NO")

	for _, v := range answers {
		fieldInfo := lib.ToTypeMap(v)
		name := fieldInfo.GetString("name")
		if name == Jotform_branchOf {
			answer := fieldInfo.GetString("answer")
			answer = c.FormResponseBranch(answer)
			if answer != "" {
				destMap.Set("Branch", answer)
			}
		} else if name == Jotform_overallRating {
			if len(fieldInfo.GetString("answer")) > 0 {
				if !isAmIntakeFormID {
					//destMap.Set("Current_Rating", fieldInfo.GetInt("answer"))
				}
			}
		} else if name == "overallRating65" {
			if len(fieldInfo.GetString("answer")) > 0 {
				destMap.Set("Current_Rating", fieldInfo.GetInt("answer"))
			}
		} else if name == "haveYou" {
			haveYou := fieldInfo.GetString("answer")
			if len(haveYou) > 0 {
				destMap.Set("Retired", haveYou)
			}
		} else if name == "areYou" {
			answer := fieldInfo.GetString("answer")
			if len(answer) > 0 {
				destMap.Set("Active_Duty", answer)
			}
		} else if name == "whatYear" {
			answer := fieldInfo.GetString("answer")
			if len(answer) > 0 {
				destMap.Set("Year_Entering_Service", answer)
			}
		} else if name == "whichYear11" {
			answer := fieldInfo.GetString("answer")
			if len(answer) > 0 {
				destMap.Set("Year_Separate_from_Service", answer)
			}
		} else if name == "typeA62" {
			answer := fieldInfo.GetString("answer")
			if answer != "" {
				destMap.Set("Occupation_during_Service", answer)
			}
		}
	}

	//for _, v := range military {
	//	apiName := config_zoho.ApiNameByFieldLabel(DealFieldInfos, v)
	//	if apiName != "" {
	//		destMap.Set(apiName, "YES")
	//	} else {
	//		if strings.Index(v, "Illness") >= 0 {
	//			destMap.Set("Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun", "YES")
	//		} else if strings.Index(v, "Atomic") >= 0 {
	//			destMap.Set("Atomic_Veterans_and_Radiation_Exposure", "YES")
	//		}
	//	}
	//}

	destMap.Set("Contact_Form", "Yes")
	destMap.Set("Answer_to_Presumptive_Questions", "https://base.vetbenefitscenter.com/questions?uniqcode="+uniqcode)

	//lib.DPrintln("HandleSubmission destMap:", destMap)

	if configs.StoppedZoho {

		destData := make(TypeDataEntry)
		for k, v := range destMap {
			vbcFieldName := config_zoho.ZohoDealVbcFieldNameByZohoFieldName(k)
			if vbcFieldName != "" {
				destData[vbcFieldName] = v
			}
		}
		if len(destData) > 0 {
			destData[DataEntry_gid] = clientCaseGid
			_, er := c.DataEntryUsecase.HandleOne(Kind_client_cases, destData, DataEntry_gid, nil)
			if er != nil {
				c.log.Error(er, " ", InterfaceToString(destData), " ", InterfaceToString(destMap))
			}
		}
	} else {
		destMap.Set("id", clientCaseGid)
		_, _, err := c.ZohoUsecase.PutRecordV1(config_zoho.Deals, destMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *JotformbuzUsecase) HttpQuestionList(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpQuestionList(body.GetString("uniqcode"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *JotformbuzUsecase) BizHttpQuestionList(uniqcode string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)

	var entity *JotformSubmissionEntity
	var err error

	tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{FieldName_uniqcode: uniqcode, "biz_deleted_at": 0, "deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	if tCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource) == ContractSource_AM {
		entity, err = c.JotformSubmissionUsecase.GetByCondWithOrderBy(Eq{"uniqcode": uniqcode, "form_id": JotformAmIntakeFormID}, "id desc")
	} else {
		entity, err = c.JotformSubmissionUsecase.GetByCondWithOrderBy(Eq{"uniqcode": uniqcode, "form_id": JotformIntakeFormID}, "id desc")
	}
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("entity is nil")
	}
	notes := lib.ToTypeMapByString(entity.Notes)
	answers := notes.GetTypeMap("content.answers")

	var questions lib.TypeList
	for _, v := range answers {

		fieldInfo := lib.ToTypeMap(v)
		//row := make(lib.TypeMap)
		//row.Set("order", fieldInfo.GetString("order"))
		//row.Set("question", fieldInfo.GetString("text"))
		//row.Set("answer", fieldInfo.GetString("answer"))
		//row.Set("type", fieldInfo.GetString("type"))

		questions = append(questions, fieldInfo)
	}

	sort.SliceStable(questions, func(i, j int) bool {
		return questions[i].GetInt("order") < questions[j].GetInt("order")
	})

	//tCase, err = c.TUsecase.Data(Kind_client_cases, Eq{"uniqcode": uniqcode})
	//if err != nil {
	//	return nil, err
	//}
	dealName := ""
	if tCase != nil {
		dealName = tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	}
	data.Set("data.deal_name", dealName)
	data.Set("questions", questions)
	return data, nil
}

func (c *JotformbuzUsecase) HandleQuestionnairesJotformHistory() error {

	cases, err := c.TUsecase.ListByCond(Kind_client_cases, And(
		Eq{"biz_deleted_at": 0, "deleted_at": 0},
		Neq{"data_collection_folder": ""}))
	if err != nil {
		return err
	}
	for k, v := range cases {

		key := MapHandleQuestionnairesJotformHistory + InterfaceToString(v.Id())
		mval, err := c.MapUsecase.GetForString(key)
		if err != nil {
			return err
		}
		if mval == "1" {
			continue
		}
		clientCaseId := v.Id()
		QuestionnairesFolderId, err := c.BoxbuzUsecase.GetDCSubFolderId(MapKeyBuildAutoBoxDCQuestionnairesFolderId(clientCaseId), cases[k])
		if err != nil {
			c.LogUsecase.SaveLog(clientCaseId, "HandleQuestionnairesJotformHistory", err)
			continue
		}
		if QuestionnairesFolderId == "" {
			c.LogUsecase.SaveLog(clientCaseId, "HandleQuestionnairesJotformHistory", "QuestionnairesFolderId is empty.")
			continue
		}
		folders, err := c.BoxUsecase.ListItemsInFolderFormat(QuestionnairesFolderId)
		if err != nil {
			c.LogUsecase.SaveLog(clientCaseId, "HandleQuestionnairesJotformHistory", err)
		}
		for _, v1 := range folders {
			files, err := c.BoxUsecase.ListItemsInFolderFormat(v1.GetString("id"))
			if err != nil {
				c.LogUsecase.SaveLog(clientCaseId, "HandleQuestionnairesJotformHistory", map[string]interface{}{
					"id":  v1.GetString("id"),
					"err": err.Error(),
				})
				continue
			}
			for _, v2 := range files {
				pdfName := v2.GetString("name")

				submissionId, _ := lib.FileExt(pdfName, false)
				if submissionId == "" {
					c.LogUsecase.SaveLog(clientCaseId, "HandleQuestionnairesJotformHistory", map[string]interface{}{
						"id":  v2.GetString("id"),
						"err": "submissionId is empty",
					})
					break
				}
				c.log.Info("HandleQuestionnairesJotformHistory:", " submissionId:", submissionId, " uniqcode:", v.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
				err = c.HandleSubmission(submissionId, v.CustomFields.TextValueByNameBasic(FieldName_uniqcode), "")
				if err != nil {
					c.LogUsecase.SaveLog(clientCaseId, "HandleQuestionnairesJotformHistory", map[string]interface{}{
						"id":   v2.GetString("id"),
						"err":  err.Error(),
						"type": "HandleSubmission",
					})
					break
				}
			}
		}
		c.MapUsecase.Set(key, "1")
		break
	}
	return nil
}
