package biz

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"math"
	"strconv"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ClientCaseUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	TUsecase         *TUsecase
	MapUsecase       *MapUsecase
	DataEntryUsecase *DataEntryUsecase
	FieldUsecase     *FieldUsecase
	ZohoUsecase      *ZohoUsecase
	UserUsecase      *UserUsecase
}

func NewClientCaseUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase,
	DataEntryUsecase *DataEntryUsecase,
	FieldUsecase *FieldUsecase,
	ZohoUsecase *ZohoUsecase,
	UserUsecase *UserUsecase) *ClientCaseUsecase {
	uc := &ClientCaseUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		MapUsecase:       MapUsecase,
		DataEntryUsecase: DataEntryUsecase,
		FieldUsecase:     FieldUsecase,
		ZohoUsecase:      ZohoUsecase,
		UserUsecase:      UserUsecase,
	}

	return uc
}

func GetCaseTimeLocation(tCase *TData, log *log.Helper) (timeLocation time.Location) {

	if tCase != nil {
		timezoneId := tCase.CustomFields.TextValueByNameBasic(FieldName_timezone_id)
		if timezoneId != "" {

			la, err := time.LoadLocation(timezoneId)
			if err != nil {
				if log != nil {
					log.Error(err)
				}
			} else {
				if la != nil {
					return *la
				}
			}
		}
	}
	a := configs.GetVBCDefaultLocation()
	return *a
}

func CaseToRelaApi(tCase *TData, log *log.Helper) lib.TypeMap {

	if tCase == nil {
		log.Error("tCase is nil")
		return nil
	}
	data := make(lib.TypeMap)
	data.Set("gid", tCase.CustomFields.TextValueByNameBasic("gid"))
	data.Set("deal_name", tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name))
	return data
}

// IsPrimaryCase 判断是否为主要primary case
func IsPrimaryCase(clientCase *TData) bool {
	if clientCase.CustomFields.NumberValueByNameBasic(FieldName_is_primary_case) == Is_primary_case_YES {
		return true
	}
	return false
}

// IsDeletedCase 是否为删除的case
func IsDeletedCase(tCase *TData) (bool, error) {
	if tCase == nil {
		return false, errors.New("tCase, isn il")
	}
	if tCase.CustomFields.NumberValueByNameBasic(FieldName_biz_deleted_at) > 0 {
		return true, nil
	}
	return false, nil
}

// PrimaryCase 获取First case，影响Second case ...
func (c *ClientCaseUsecase) PrimaryCase(clientGid string) (primaryCase *TData, err error) {
	clientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{
		"client_gid":              clientGid,
		FieldName_is_primary_case: Is_primary_case_YES})
	if err != nil {
		return nil, err
	}
	if clientCase == nil {
		return nil, nil
	}
	if clientCase.CustomFields.NumberValueByNameBasic(FieldName_biz_deleted_at) != 0 {
		return nil, errors.New("the primary case was deleted")
	}
	return clientCase, nil
}

// NotPrimaryCases 获取非主要primary cases
func (c *ClientCaseUsecase) NotPrimaryCases(clientGid string, clientCaseGid string) (cases []*TData, err error) {
	return c.TUsecase.ListByCond(Kind_client_cases, And(
		Eq{FieldName_is_primary_case: Is_primary_case_NO, "client_gid": clientGid, DataEntry_biz_deleted_at: 0},
		In(FieldName_stages, config_vbc.Stages_AwaitingPayment,
			config_vbc.Stages_27_AwaitingBankReconciliation,
			config_vbc.Stages_Completed,
			config_vbc.Stages_AmAwaitingPayment,
			config_vbc.Stages_Am27_AwaitingBankReconciliation,
			config_vbc.Stages_AmCompleted),
		Neq{"gid": clientCaseGid}))
}

// CurrentCaseInProgress 获取当前正在进行的case
func (c *ClientCaseUsecase) CurrentCaseInProgress(clientGid string) (currentCase *TData, err error) {
	return c.TUsecase.Data(Kind_client_cases, And(
		Eq{"client_gid": clientGid, FieldName_biz_deleted_at: 0},
		NotIn(FieldName_stages, config_vbc.Stages_AwaitingPayment,
			config_vbc.Stages_27_AwaitingBankReconciliation,
			config_vbc.Stages_Completed,
			config_vbc.Stages_AmAwaitingPayment,
			config_vbc.Stages_Am27_AwaitingBankReconciliation,
			config_vbc.Stages_AmCompleted),
	))
}

// ClientCaseContractBasicDataVoById 获取配置合同信息
func (c *ClientCaseUsecase) ClientCaseContractBasicDataVoById(clientCaseId int32) (*ClientCaseContractBasicDataVo, error) {
	key := fmt.Sprintf("%s%d", Map_ClientCaseContractBasicData, clientCaseId)
	contractBaseDataStr, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return nil, err
	}
	if contractBaseDataStr == "" {
		return nil, errors.New("contractBaseDataStr is empty clientCaseId:" + InterfaceToString(clientCaseId))
	}
	clientCaseContractBasicDataVo := lib.StringToTDef[*ClientCaseContractBasicDataVo](contractBaseDataStr, nil)
	if clientCaseContractBasicDataVo == nil {
		return nil, errors.New("clientCaseContractBasicDataVo is nil clientCaseId:" + InterfaceToString(clientCaseId))
	}
	return clientCaseContractBasicDataVo, nil
}

func (c *ClientCaseUsecase) GetPricingVersion(tCase *TData) (pricingVersion string, err error) {
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	pricingVersion = tCase.CustomFields.TextValueByNameBasic(FieldName_s_pricing_version)
	if pricingVersion == "" {
		pricingVersion = DefaultPricingVersion
		//c.log.Debug("GetPricingVersion: from DefaultPricingVersion ", tCase.CustomFields.NumberValueByNameBasic("id"))
	} else {
		//c.log.Debug("GetPricingVersion: from DB ", tCase.CustomFields.NumberValueByNameBasic("id"))
	}
	return pricingVersion, nil
}

func (c *ClientCaseUsecase) SaveContractSource(gid string, contractSource string) error {

	if contractSource == "" {
		return errors.New("contractSource is empty")
	}
	_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry{
		"gid":                    gid,
		FieldName_ContractSource: contractSource,
	}, DataEntry_gid, nil)
	return err
}

func (c *ClientCaseUsecase) SavePricingVersion(tCase *TData, pricingVersion string) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	if pricingVersion == "" {
		return errors.New("pricingVersion is empty")
	}

	_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry{
		"id":                        tCase.CustomFields.NumberValueByNameBasic("id"),
		FieldName_s_pricing_version: pricingVersion,
		FieldName_pricing_version:   pricingVersion,
	}, "id", nil)
	return err
}

func (c *ClientCaseUsecase) EnabledTwoBySMS(caseId int32) (bool, error) {
	if configs.EnabledTwoBySMS {
		// 上线
		return true, nil

		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			c.log.Error(err)
			return false, err
		}
		if tCase == nil {
			return false, errors.New("tCase is nil")
		}
		// todo:lgl 只有这邮箱开启
		if tCase.CustomFields.TextValueByNameBasic(FieldName_email) == "liaogling@gmail.com" {
			return true, nil
		}

	}
	return false, nil
}

func (c *ClientCaseUsecase) GetCaseWithCache(caches lib.Cache[*TData], caseId int32) (*TData, error) {
	key := InterfaceToString(caseId)
	entity, exists := caches.Get(key)
	if exists {
		return entity, nil
	}

	entity, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return nil, err
	}
	caches.Set(key, entity)
	return entity, nil
}

func (c *ClientCaseUsecase) HttpClaimsInfo(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpClaimsInfo(body.GetString("uniqcode"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ClientCaseUsecase) HttpClaimsInfoTest(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpClaimsInfoTest(body.GetString("uniqcode"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ClientCaseUsecase) HttpDetailById(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpDetailById(body.GetInt("id"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ClientCaseUsecase) BizHttpDetailById(caseId int32) (lib.TypeMap, error) {

	data := make(lib.TypeMap)

	tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": caseId, "deleted_at": 0, "biz_deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	dealName := tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	data.Set("data.deal_name", dealName)
	data.Set("data.gid", tCase.Gid())
	data.Set("data.id", tCase.Id())
	data.Set("data."+FieldName_stages, tCase.CustomFields.ValueByName(FieldName_stages))
	data.Set("data."+FieldName_case_files_folder, tCase.CustomFields.ValueByName(FieldName_case_files_folder))
	data.Set("data."+FieldName_data_collection_folder, tCase.CustomFields.ValueByName(FieldName_data_collection_folder))
	return data, nil
}

func (c *ClientCaseUsecase) HttpDetail(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpDetail(body.GetString("uniqcode"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type DetailWhole struct {
	ShowRecordReview bool           `json:"show_record_review"`
	PrimaryName      string         `json:"primary_name"`
	Amount           string         `json:"amount"` // 暂时client cases专用
	UpdatedAt        int32          `json:"updated_at"`
	CreatedTime      string         `json:"created_time"`
	Extras           interface{}    `json:"extras"`
	Sections         DetailSections `json:"sections"`
	ClientPipelines  string         `json:"client_pipelines"`
}

func GenItfExpirationExtend(tData *TData) *TFieldExtendForSysDueDate {
	itfDate := tData.CustomFields.TextValueByNameBasic(FieldName_itf_expiration)
	val := itfDate
	if val != "" {
		currentTime := time.Now().In(configs.GetVBCDefaultLocation())
		itfTime, _ := time.ParseInLocation(time.DateOnly, val, configs.GetVBCDefaultLocation())
		daysRemaining := int64(math.Ceil(itfTime.Sub(currentTime).Hours() / 24))
		daysRemainingStr := strconv.FormatInt(daysRemaining, 10)

		entity, _ := GenTFieldExtendForSysItfFormula(daysRemainingStr)
		entity.Value = itfDate
		entity.Label = tData.CustomFields.DisplayValueByNameBasic(FieldName_itf_expiration)
		return entity
	}
	return nil
}

type DetailWholeCaseExtras struct {
	StageStartDate TFieldExtendForSysDueDate `json:"stage_start_date"`
	StageDueDate   TFieldExtendForSysDueDate `json:"stage_due_date"`
	ItfFormula     TFieldExtendForSysDueDate `json:"itf_formula"`
	ItfExpiration  TFieldExtendForSysDueDate `json:"itf_expiration"`
	MedicalDbqCost MedicalDbqCost            `json:"medical_dbq_cost"`
	Pipelines      string                    `json:"pipelines"`
}

type DetailSections []DetailSection

type DetailSection struct {
	SectionLabel string     `json:"section_label"`
	SectionName  string     `json:"section_name"`
	Left         SectionTab `json:"left"`
	Right        SectionTab `json:"right"`
}

type SectionTab struct {
	Fields []DetailField `json:"fields"`
}

type DetailField struct {
	FieldLabel string      `json:"field_label"`
	FieldName  string      `json:"field_name"`
	FieldType  string      `json:"field_type"`
	Value      interface{} `json:"value"`
	Extend     interface{} `json:"extend"`
}

func (c *ClientCaseUsecase) BizHttpDetail(uniqcode string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)

	tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"uniqcode": uniqcode, "deleted_at": 0, "biz_deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}

	var detailSections DetailSections
	var detailSection DetailSection
	detailSection.SectionLabel = "Claims Information"
	detailSection.SectionName = "claims_information"
	var left SectionTab
	var right SectionTab

	var fieldNames = []string{
		"service_connections", "previous_denials",
		"claims_online", "claims_next_round",
		"claims_supplemental",
	}

	structField, err := c.FieldUsecase.CacheStructByKind(Kind_client_cases)
	if err != nil {
		return nil, err
	}

	for _, v := range tCase.CustomFields {
		if lib.InArray(v.Name, fieldNames) {
			var detailField DetailField
			field := structField.GetByFieldName(v.Name)
			if field == nil {
				return nil, errors.New("field is nil")
			}
			detailField.FieldLabel = field.FieldLabel
			detailField.FieldName = field.FieldName
			detailField.FieldType = field.FieldType
			detailField.Value = CaseClaimsDivide(tCase.CustomFields.TextValueByNameBasic(field.FieldName))

			if len(left.Fields) == 3 {
				right.Fields = append(right.Fields, detailField)
			} else {
				left.Fields = append(left.Fields, detailField)
			}
		}
	}
	detailSection.Left = left
	detailSection.Right = right
	detailSections = append(detailSections, detailSection)
	dealName := ""
	if tCase != nil {
		dealName = tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	}
	var detailWhole DetailWhole
	detailWhole.Sections = detailSections
	data.Set("detail", detailWhole)
	data.Set("data.deal_name", dealName)
	return data, nil
}

var CaseClaimsTypes = []string{
	"str",
	"new",
	"opinion",
	"increase",
	"previous denial",
	"denial",
}

type CaseClaimsRow struct {
	Rating       string `json:"rating"`
	Condition    string `json:"condition"`
	ClaimsType   string `json:"claims_type"`
	DisplayValue string `json:"display_value"`
}

type CaseClaimsRowV2 struct {
	Rating       string   `json:"rating"`
	Condition    string   `json:"condition"`
	ClaimsTypes  []string `json:"claims_types"`
	DisplayValue string   `json:"display_value"`
}

// CaseClaimsTypeV2 Left knee strain status post medial meniscal tear with limitation of flexion and extension (str, increase)
func CaseClaimsTypeV2(val string) (claims string, caseClaimsTypes []string) {

	i := strings.LastIndex(val, "(")
	if i > 0 {
		typeStr := val[i:]
		typeStr = strings.ReplaceAll(typeStr, "(", "")
		typeStr = strings.ReplaceAll(typeStr, ")", "")
		temps := strings.Split(typeStr, ",")

		isOk := true
		for _, v := range temps {
			v := strings.TrimSpace(v)
			if v != "" {

				if !lib.InArray(v, CaseClaimsTypes) {
					isOk = false
					break
				}

				caseClaimsTypes = append(caseClaimsTypes, v)
			}
		}
		if !isOk {
			return val, nil
		}
		return strings.TrimSpace(val[:i-1]), caseClaimsTypes
	}

	return val, nil
}

func CaseClaimsType(val string) (claims string, caseClaimsType string) {
	for _, v := range CaseClaimsTypes {
		tV := fmt.Sprintf("(%s)", v)
		i := strings.Index(val, tV)
		if i > 0 {
			return strings.TrimSpace(val[:i]), v
		}
	}
	return val, ""
}

func CaseClaimsRowDisplayValue(val *CaseClaimsRow) {
	var r string
	if val.Rating != "" {
		r += val.Rating + " - "
	}
	r += val.Condition
	if val.ClaimsType != "" {
		r += " (" + val.ClaimsType + ")"
	}
	val.DisplayValue = r
}

func CaseClaimsRowDisplayValueV2(val *CaseClaimsRowV2) {
	var r string
	if val.Rating != "" {
		r += val.Rating + " - "
	}
	r += val.Condition

	tV := strings.Join(val.ClaimsTypes, ", ")

	if tV != "" {
		r += " (" + tV + ")"
	}
	val.DisplayValue = r
}

func CaseClaimsDivideV2(val string) (r []CaseClaimsRowV2) {
	res := strings.Split(val, "\n")
	for _, v := range res {
		v := strings.TrimSpace(v)
		if v == "" {
			continue
		}
		rows := strings.Split(v, "-")
		if len(rows) == 1 {
			temp := strings.TrimSpace(rows[0])
			if temp == "" {
				continue
			}
			a, b := CaseClaimsTypeV2(temp)
			caseClaimsRow := CaseClaimsRowV2{
				Condition:    a,
				ClaimsTypes:  b,
				DisplayValue: fmt.Sprintf("%s (%s)", a, b),
			}
			CaseClaimsRowDisplayValueV2(&caseClaimsRow)
			r = append(r, caseClaimsRow)
		} else if len(rows) == 2 {
			temp := strings.TrimSpace(rows[0])
			temp1 := strings.TrimSpace(rows[1])
			if temp == "" && temp1 == "" {
				continue
			}
			a, b := CaseClaimsTypeV2(temp1)
			caseClaimsRow := CaseClaimsRowV2{
				Rating:      temp,
				Condition:   a,
				ClaimsTypes: b,
			}
			CaseClaimsRowDisplayValueV2(&caseClaimsRow)

			r = append(r, caseClaimsRow)

		}
	}
	return
}

func CaseClaimsDivide(val string) (r []CaseClaimsRow) {
	res := strings.Split(val, "\n")
	for _, v := range res {
		v := strings.TrimSpace(v)
		if v == "" {
			continue
		}
		rows := strings.Split(v, "-")
		if len(rows) == 1 {
			temp := strings.TrimSpace(rows[0])
			if temp == "" {
				continue
			}
			a, b := CaseClaimsType(temp)
			caseClaimsRow := CaseClaimsRow{
				Condition:    a,
				ClaimsType:   b,
				DisplayValue: fmt.Sprintf("%s (%s)", a, b),
			}
			CaseClaimsRowDisplayValue(&caseClaimsRow)
			r = append(r, caseClaimsRow)
		} else if len(rows) == 2 {
			temp := strings.TrimSpace(rows[0])
			temp1 := strings.TrimSpace(rows[1])
			if temp == "" && temp1 == "" {
				continue
			}
			a, b := CaseClaimsType(temp1)
			caseClaimsRow := CaseClaimsRow{
				Rating:     temp,
				Condition:  a,
				ClaimsType: b,
			}
			CaseClaimsRowDisplayValue(&caseClaimsRow)

			r = append(r, caseClaimsRow)

		}
	}
	return
}

func (c *ClientCaseUsecase) BizHttpClaimsInfoTest(uniqcode string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)

	tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"uniqcode": uniqcode, "deleted_at": 0, "biz_deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	serviceConnections := tCase.CustomFields.TextValueByNameBasic("service_connections")
	previousDenials := tCase.CustomFields.TextValueByNameBasic("previous_denials")
	claimsOnline := tCase.CustomFields.TextValueByNameBasic("claims_online")
	claimsNextRound := tCase.CustomFields.TextValueByNameBasic("claims_next_round")
	claimsSupplemental := tCase.CustomFields.TextValueByNameBasic("claims_supplemental")

	dealName := ""
	if tCase != nil {
		dealName = tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	}
	var list lib.TypeList
	list = append(list, lib.TypeMap{
		"section": "Service Connections",
		"value":   serviceConnections,
	})

	data.Set("data.deal_name", dealName)
	data.Set("data.service_connections", CaseClaimsDivide(serviceConnections))
	data.Set("data.previous_denials", CaseClaimsDivide(previousDenials))
	data.Set("data.claims_online", CaseClaimsDivide(claimsOnline))
	data.Set("data.claims_next_round", CaseClaimsDivide(claimsNextRound))
	data.Set("data.claims_supplemental", CaseClaimsDivide(claimsSupplemental))

	return data, nil
}

func (c *ClientCaseUsecase) BizHttpClaimsInfo(uniqcode string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)

	tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"uniqcode": uniqcode, "deleted_at": 0, "biz_deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	serviceConnections := tCase.CustomFields.TextValueByNameBasic("service_connections")
	previousDenials := tCase.CustomFields.TextValueByNameBasic("previous_denials")
	claimsOnline := tCase.CustomFields.TextValueByNameBasic("claims_online")
	claimsNextRound := tCase.CustomFields.TextValueByNameBasic("claims_next_round")
	claimsSupplemental := tCase.CustomFields.TextValueByNameBasic("claims_supplemental")

	dealName := ""
	if tCase != nil {
		dealName = tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	}

	data.Set("data.deal_name", dealName)
	data.Set("data.service_connections", CaseClaimsDivide(serviceConnections))
	data.Set("data.previous_denials", CaseClaimsDivide(previousDenials))
	data.Set("data.claims_online", CaseClaimsDivide(claimsOnline))
	data.Set("data.claims_next_round", CaseClaimsDivide(claimsNextRound))
	data.Set("data.claims_supplemental", CaseClaimsDivide(claimsSupplemental))

	return data, nil
}

func (c *ClientCaseUsecase) HttpSave(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpSave(body)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ClientCaseUsecase) BizHttpSave(body lib.TypeMap) (lib.TypeMap, error) {

	lib.DPrintln("body:", body)
	uniqcode := body.GetString("uniqcode")
	if uniqcode == "" {
		return nil, errors.New("uniqcode is empty")
	}

	tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"uniqcode": uniqcode})
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("The client does not exist")
	}

	var keys = []string{
		"claims_next_round",
		"claims_online",
		"claims_supplemental",
		"previous_denials",
		"service_connections",
	}
	destData := make(lib.TypeMap)
	for k, _ := range body {
		if lib.InArray(k, keys) {
			value := body.GetTypeList(k)
			var destVal []string
			for _, v1 := range value {
				a := lib.InterfaceToTDef[*CaseClaimsRow](v1, nil)
				if a != nil {
					CaseClaimsRowDisplayValue(a)
					destVal = append(destVal, a.DisplayValue)
				}
			}

			destData.Set(k, strings.Join(destVal, "\n"))
		}
	}
	if len(destData) > 0 {
		if configs.IsProd() {
			if uniqcode != "2487303324" {
				return nil, nil
			}
			destMap := make(lib.TypeMap)
			destMap.Set("id", tCase.CustomFields.TextValueByNameBasic(FieldName_gid))
			for k, v := range destData {
				zohoFieldName := config_zoho.ZohoDealFieldNameByVbcFieldName(k)
				if zohoFieldName != "" {
					destMap[zohoFieldName] = v
				}
			}
			c.log.Info("BizHttpSave:", destMap)
			_, _, err := c.ZohoUsecase.PutRecordV1(config_zoho.Deals, destMap)
			if err != nil {
				return nil, err
			}
		}
		destData.Set("uniqcode", uniqcode)
		_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry(destData), "uniqcode", nil)
		if err != nil {
			return nil, err
		}
	}
	//data := make(lib.TypeMap)
	//data.Set("data.val", "aaa")
	return nil, nil
}

func (c *ClientCaseUsecase) GetLeadVSByPhone(phone string) (vsTUser *TData, err error) {
	tCase, err := c.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	if tCase != nil {
		leadVsFullName := tCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs)
		if leadVsFullName != "" {
			return c.UserUsecase.GetByFullName(leadVsFullName)
		}
	}
	return
}

// GetByPhone +13109719619 找一个
func (c *ClientCaseUsecase) GetByPhone(phone string) (*TData, error) {

	phone1, phone2, phone3, err := FormatPhoneNumber(phone)
	if err != nil {
		return nil, err
	}

	if phone1 != "" || phone2 != "" || phone3 != "" {
		var conds []Cond
		if phone1 != "" {
			conds = append(conds, Eq{"phone": phone1})
		}
		if phone2 != "" {
			conds = append(conds, Eq{"phone": phone2})
		}
		if phone3 != "" {
			conds = append(conds, Eq{"phone": phone3})
		}
		return c.TUsecase.DataWithOrderBy(Kind_client_cases, And(Eq{FieldName_biz_deleted_at: 0}, Or(conds...)), "id desc")
	}
	return nil, nil
}

func (c *ClientCaseUsecase) ItfCasesByUserGid(userGid string) (TDataList, error) {

	builder := Dialect(MYSQL).Select("*").From("client_cases").Where(Eq{"deleted_at": 0})
	builder.And(Eq{"biz_deleted_at": 0})
	builder.And(NotIn("stages",
		config_vbc.Stages_AwaitingDecision,
		config_vbc.Stages_AwaitingPayment,
		config_vbc.Stages_27_AwaitingBankReconciliation,
		config_vbc.Stages_Completed,
		config_vbc.Stages_Terminated,
		config_vbc.Stages_Dormant,
		config_vbc.Stages_AmAwaitingDecision,
		config_vbc.Stages_AmAwaitingPayment,
		config_vbc.Stages_Am27_AwaitingBankReconciliation,
		config_vbc.Stages_AmCompleted,
		config_vbc.Stages_AmTerminated,
		config_vbc.Stages_AmDormant,
	))
	builder.And(Eq{DataEntry_user_gid: userGid})
	builder.And(Like{FieldName_collaborators, fmt.Sprintf(",%s,", userGid)})

	now := time.Now().In(configs.GetVBCDefaultLocation())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, configs.GetVBCDefaultLocation())
	//begin := now.Format(time.DateOnly)
	end := now.AddDate(0, 0, 91).Format(time.DateOnly)
	builder.And(And(Neq{"itf_expiration": ""}, NotNull{"itf_expiration"}, Lt{"itf_expiration": end}))
	builder.OrderBy("itf_expiration")
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	c.log.Info("sql:", sql)
	return c.TUsecase.ListByRawSql(Kind_client_cases, sql)
}

func (c *ClientCaseUsecase) ItfCases() (TDataList, error) {
	builder := Dialect(MYSQL).Select("*").From("client_cases").Where(Eq{"deleted_at": 0})
	builder.And(Eq{"biz_deleted_at": 0})
	builder.And(NotIn("stages",
		config_vbc.Stages_AwaitingDecision,
		config_vbc.Stages_AwaitingPayment,
		config_vbc.Stages_27_AwaitingBankReconciliation,
		config_vbc.Stages_Completed,
		config_vbc.Stages_Terminated,
		config_vbc.Stages_Dormant,
		config_vbc.Stages_AmAwaitingDecision,
		config_vbc.Stages_AmAwaitingPayment,
		config_vbc.Stages_Am27_AwaitingBankReconciliation,
		config_vbc.Stages_AmCompleted,
		config_vbc.Stages_AmTerminated,
		config_vbc.Stages_AmDormant,
	))
	now := time.Now().In(configs.GetVBCDefaultLocation())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, configs.GetVBCDefaultLocation())
	//begin := now.Format(time.DateOnly)
	end := now.AddDate(0, 0, 91).Format(time.DateOnly)
	builder.And(And(Neq{"itf_expiration": ""}, NotNull{"itf_expiration"}, Lt{"itf_expiration": end}))
	builder.OrderBy("itf_expiration")
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	c.log.Info("sql:", sql)
	return c.TUsecase.ListByRawSql(Kind_client_cases, sql)
}
func (c *ClientCaseUsecase) GetOldestCaseByClientGid(clientGid string) (tClient *TData, err error) {

	return c.TUsecase.DataWithOrderBy(Kind_client_cases, Eq{
		"biz_deleted_at": 0,
		"client_gid":     clientGid,
	}, "id asc")
}
