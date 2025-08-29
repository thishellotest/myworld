package biz

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/shopspring/decimal"
	"net/url"
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

type ZohobuzUsecase struct {
	log                 *log.Helper
	CommonUsecase       *CommonUsecase
	conf                *conf.Data
	DataEntryUsecase    *DataEntryUsecase
	ZohoUsecase         *ZohoUsecase
	TUsecase            *TUsecase
	FeeUsecase          *FeeUsecase
	DataComboUsecase    *DataComboUsecase
	UsageStatsUsecase   *UsageStatsUsecase
	AsanaMigrateUsecase *AsanaMigrateUsecase
	MapUsecase          *MapUsecase
	ActionOnceUsecase   *ActionOnceUsecase
	StageTransUsecase   *StageTransUsecase
}

func NewZohobuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	DataEntryUsecase *DataEntryUsecase,
	ZohoUsecase *ZohoUsecase,
	TUsecase *TUsecase,
	FeeUsecase *FeeUsecase,
	DataComboUsecase *DataComboUsecase,
	UsageStatsUsecase *UsageStatsUsecase,
	AsanaMigrateUsecase *AsanaMigrateUsecase,
	MapUsecase *MapUsecase,
	ActionOnceUsecase *ActionOnceUsecase,
	StageTransUsecase *StageTransUsecase) *ZohobuzUsecase {
	uc := &ZohobuzUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		DataEntryUsecase:    DataEntryUsecase,
		ZohoUsecase:         ZohoUsecase,
		TUsecase:            TUsecase,
		FeeUsecase:          FeeUsecase,
		DataComboUsecase:    DataComboUsecase,
		UsageStatsUsecase:   UsageStatsUsecase,
		AsanaMigrateUsecase: AsanaMigrateUsecase,
		MapUsecase:          MapUsecase,
		ActionOnceUsecase:   ActionOnceUsecase,
		StageTransUsecase:   StageTransUsecase,
	}

	return uc
}

func (c *ZohobuzUsecase) SyncClientsByGids(gids []string) error {
	if len(gids) == 0 {
		return nil
	}

	fields := config_zoho.ContactLayout().ContactApiNames()
	params := make(url.Values)
	for _, v := range gids {
		params.Add("ids", v)
	}
	c.UsageStatsUsecase.Stat(UsageType_GetContactRecords, time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Contacts, fields, params)
	if err != nil {
		return err
	}
	listMaps := records.GetTypeList("data")
	if len(listMaps) == 0 {
		return nil
	}
	err = c.SyncClients(listMaps)
	if err != nil {
		return err
	}
	return nil
}

func (c *ZohobuzUsecase) SyncClientCasesByGids(gids []string) error {
	if len(gids) == 0 {
		return nil
	}
	fields := config_zoho.DealLayout().DealApiNames()
	params := make(url.Values)
	for _, v := range gids {
		params.Add("ids", v)
	}
	c.UsageStatsUsecase.Stat(UsageType_GetDealRecords, time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	if err != nil {
		c.log.Error("SyncClientCasesByGids:", err)
		return err
	}
	listMaps := records.GetTypeList("data")
	if len(listMaps) == 0 {
		return nil
	}
	err = c.SyncClientCases(listMaps, "")
	if err != nil {
		c.log.Error("SyncClientCasesByGids:", err)
		return err
	}
	return c.SyncClientCases2ByGids(gids)
}

func (c *ZohobuzUsecase) SyncClientCases2ByGids(gids []string) error {
	fields := config_zoho.DealLayout().DealApiNames2()
	params := make(url.Values)
	for _, v := range gids {
		params.Add("ids", v)
	}
	c.UsageStatsUsecase.Stat(UsageType_GetDealRecords, time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	if err != nil {
		c.log.Info("SyncClientCases2ByGids:", err)
		return err
	}
	listMaps := records.GetTypeList("data")
	if len(listMaps) == 0 {
		return nil
	}
	err = c.SyncClientCases(listMaps, "job2")
	if err != nil {
		c.log.Info("SyncClientCases2ByGids:", err)
		return err
	}
	return nil
}

func (c *ZohobuzUsecase) SyncClientOne(row lib.TypeMap) error {
	var list lib.TypeList
	list = append(list, row)
	return c.SyncClients(list)
}

func (c *ZohobuzUsecase) SyncClients(list lib.TypeList) error {

	var dataList TypeDataEntryList
	for k, _ := range list {
		dataList = append(dataList, TypeDataEntry(config_zoho.ClientMappings(list[k])))
	}
	_, err := c.DataEntryUsecase.Handle(Kind_clients, dataList, Client_FileName_gid, nil)
	return err
}

func (c *ZohobuzUsecase) SyncClientCaseOne(row lib.TypeMap) error {
	var list lib.TypeList
	list = append(list, row)
	return c.SyncClientCases(list, "")
}

func (c *ZohobuzUsecase) SyncClientCases(list lib.TypeList, from string) error {
	var dataList TypeDataEntryList
	for k, _ := range list {
		var typeMaps lib.TypeMap
		var err error
		if from == "job2" {
			typeMaps, err = c.StageTransUsecase.ClientCasesMappings2(list[k])
		} else {
			typeMaps, err = c.StageTransUsecase.ClientCasesMappings(list[k])
		}
		if err != nil {
			c.log.Error(err)
		}
		dataEntry := TypeDataEntry(typeMaps)
		dataEntry[FieldName_biz_deleted_at] = 0
		dataList = append(dataList, dataEntry)
	}
	//return nil
	_, err := c.DataEntryUsecase.Handle(Kind_client_cases, dataList, FileName_client_cases_gid, nil)
	return err
}

func (c *ZohobuzUsecase) SyncTasks(list lib.TypeList) error {
	var dataList TypeDataEntryList
	for k, _ := range list {
		dataEntry := TypeDataEntry(config_zoho.TasksMappings(list[k]))
		dataEntry[FieldName_biz_deleted_at] = 0
		seModule := InterfaceToString(dataEntry["se_module"])
		if seModule == "Deals" {
			dataEntry["re_kind"] = Kind_client_cases + ""
		} else {
			dataEntry["re_kind"] = seModule
		}
		dataList = append(dataList, dataEntry)
	}
	_, err := c.DataEntryUsecase.Handle(Kind_client_tasks, dataList, FieldName_gid, nil)
	return err
}

func (c *ZohobuzUsecase) SyncUserOne(row lib.TypeMap) error {
	var list lib.TypeList
	list = append(list, row)
	return c.SyncUsers(list)
}

func (c *ZohobuzUsecase) SyncUsers(list lib.TypeList) error {
	var dataList TypeDataEntryList
	for k, _ := range list {
		dataList = append(dataList, TypeDataEntry(config_zoho.UserMappings(list[k])))
	}
	_, err := c.DataEntryUsecase.Handle(Kind_users, dataList, FileName_user_gid, nil)
	return err
}

func (c *ZohobuzUsecase) HttpBizSyncUsers(ctx *gin.Context) {
	reply := CreateReply()
	err := c.BizSyncUsers()
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ZohobuzUsecase) BizSyncUsers() error {

	r, err := c.ZohoUsecase.Users()
	if err != nil {
		return err
	}
	if r == nil {
		return errors.New("users response is nil.")
	}
	users := r.GetTypeList("users")
	return c.SyncUsers(users)
}

type ClientsEntity struct {
	ID           int32 `gorm:"primaryKey"`
	Gid          string
	BizDeletedAt int64
	UpdatedAt    int64
}

type ClientCaseEntity struct {
	ID           int32 `gorm:"primaryKey"`
	Gid          string
	BizDeletedAt int64
	UpdatedAt    int64
}

// SyncClientsDeletes 同步状态
func (c *ZohobuzUsecase) SyncClientsDeletes() error {

	sqlRows, err := c.CommonUsecase.DB().Table("clients").Select("id,gid,biz_deleted_at,updated_at").
		Where("biz_deleted_at=0 and deleted_at=0").Rows()
	if err != nil {
		return err
	} else {
		if sqlRows != nil {
			maxLimit := 200
			var records []*ClientsEntity
			for sqlRows.Next() {
				var task ClientsEntity
				err = c.CommonUsecase.DB().ScanRows(sqlRows, &task)
				if err != nil {
					c.log.Error(err)
					continue
				}
				records = append(records, &task)

				if len(records) >= maxLimit {
					err = c.HandleZohoClientDelete(records)
					if err != nil {
						c.log.Error(err)
					}
					records = records[:0]
				}
			}
			if len(records) > 0 {
				err = c.HandleZohoClientDelete(records)
				if err != nil {
					c.log.Error(err, records)
				}
			}
			err = sqlRows.Close()
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}

func (c *ZohobuzUsecase) HandleZohoClientDelete(clients []*ClientsEntity) error {

	if len(clients) == 0 {
		return nil
	}

	fields := config_zoho.ContactLayout().ContactApiNames()
	params := make(url.Values)
	for _, v := range clients {
		params.Add("ids", v.Gid)
	}

	c.UsageStatsUsecase.Stat("ZohobuzUsecase_HandleZohoClientDelete", time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Contacts, fields, params)
	if err != nil {
		return errors.New(err.Error() + ":" + InterfaceToString(params))
	}
	listMaps := records.GetTypeList("data")
	if len(listMaps) == 0 {
		return nil
	}
	maps := make(map[string]bool)
	for _, v := range listMaps {
		maps[v.GetString("id")] = true
	}
	var needsDeleteGid []string
	for _, v := range clients {
		if _, ok := maps[v.Gid]; !ok {
			needsDeleteGid = append(needsDeleteGid, v.Gid)
		}
	}

	//return nil
	if len(needsDeleteGid) > 0 {
		return c.CommonUsecase.DB().Table("clients").
			Where("gid in ?", needsDeleteGid).
			Updates(map[string]interface{}{
				"biz_deleted_at": time.Now().Unix(),
				"updated_at":     time.Now().Unix(),
			}).Error
	}

	return nil
}

// SyncDealsDeletes 同步状态
func (c *ZohobuzUsecase) SyncDealsDeletes() error {

	sqlRows, err := c.CommonUsecase.DB().Table("client_cases").Select("id,gid,biz_deleted_at,updated_at").
		Where("biz_deleted_at=0 and deleted_at=0").Rows()
	if err != nil {
		return err
	} else {
		if sqlRows != nil {
			maxLimit := 200
			var records []*ClientCaseEntity
			for sqlRows.Next() {
				var task ClientCaseEntity
				err = c.CommonUsecase.DB().ScanRows(sqlRows, &task)
				if err != nil {
					c.log.Error(err)
					continue
				}
				records = append(records, &task)

				if len(records) >= maxLimit {
					err = c.HandleZohoDelete(records)
					if err != nil {
						c.log.Error(err)
					}
					records = records[:0]
				}
			}
			if len(records) > 0 {
				err = c.HandleZohoDelete(records)
				if err != nil {
					c.log.Error(err)
				}
			}
			err = sqlRows.Close()
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}

func (c *ZohobuzUsecase) HandleZohoDelete(clientCases []*ClientCaseEntity) error {

	if len(clientCases) == 0 {
		return nil
	}

	fields := config_zoho.DealLayout().DealApiNames()
	params := make(url.Values)
	for _, v := range clientCases {
		params.Add("ids", v.Gid)
	}

	c.UsageStatsUsecase.Stat("ZohobuzUsecase_HandleZohoDelete", time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	if err != nil {
		return err
	}
	listMaps := records.GetTypeList("data")
	if len(listMaps) == 0 {
		return nil
	}
	maps := make(map[string]bool)
	for _, v := range listMaps {
		maps[v.GetString("id")] = true
	}
	var needsDeleteGid []string
	for _, v := range clientCases {
		if _, ok := maps[v.Gid]; !ok {
			needsDeleteGid = append(needsDeleteGid, v.Gid)
		}
	}
	if len(needsDeleteGid) > 0 {
		return c.CommonUsecase.DB().Table("client_cases").
			Where("gid in ?", needsDeleteGid).
			Updates(map[string]interface{}{
				"biz_deleted_at": time.Now().Unix(),
				"updated_at":     time.Now().Unix(),
			}).Error
	}

	return nil
}

func (c *ZohobuzUsecase) HandleAllMan() error {
	cases, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{"biz_deleted_at": 0, "deleted_at": 0})
	if err != nil {
		return err
	}
	for _, v := range cases {
		aa := v.Id()
		err := c.HandleAmount(int32(aa))
		if err != nil {
			c.log.Error(aa)
		} else {
			c.log.Info("HandleAllMan caseId: ", aa)
		}
	}

	return nil
	caseIdStr := "5611,5613,5614,5617,5616,5615,5618,5619,5603,5612,5621,5591,5620,5609,5623,5625,5148,5628,5629,5626,5601,5630,5631,5632,5140,5633,55,5635,5636,5511,5201,5127,274,5170,5574,5171,5475,5179,5311,5129,150,5561,5349,5128,5154,5491,5267,5405,5187,5425,5429,5373,5265,5213,135,5545,280,5548,5572,58,5025,5638,5640,5469,129,5415,5332,5641,5346,5416,5353,355"
	caseIds := strings.Split(caseIdStr, ",")
	for _, v := range caseIds {
		aa, _ := strconv.ParseInt(v, 10, 32)
		err := c.HandleAmount(int32(aa))
		if err != nil {
			c.log.Error(aa)
		} else {
			c.log.Info("HandleAllMan caseId: ", aa)
		}
	}
	return nil
}

func GetFilingDateByItfExpiration(itfExpiration string) (filingDate string) {
	if itfExpiration == "" {
		return ""
	}
	arr := strings.Split(itfExpiration, "-")
	if len(arr) != 3 {
		return ""
	}
	if arr[1] == "02" && arr[2] == "29" {
		arr[2] = "01"
		arr[1] = "03"
	}
	year, err := strconv.ParseInt(arr[0], 10, 32)
	if err != nil {
		return ""
	}
	year = year - 1
	return fmt.Sprintf("%d-%s-%s", year, arr[1], arr[2])
}

// GetDiffDaysFilingDateByItfExpiration 2025-07-07
func GetDiffDaysFilingDateByItfExpiration(itfExpiration string, dateOfApproval string) (diffDays int) {
	fillingDate := GetFilingDateByItfExpiration(itfExpiration)
	if fillingDate == "" {
		return 0
	}
	dateOfApprovalTime, err := time.Parse(time.DateOnly, dateOfApproval)
	if err != nil {
		return 0
	}
	fillingDateTime, err := time.Parse(time.DateOnly, fillingDate)
	if err != nil {
		return 0
	}
	duration := dateOfApprovalTime.Sub(fillingDateTime)
	days := int(duration.Hours() / 24)
	return days
}

func (c *ZohobuzUsecase) HandleAmountForAm(tCase TData) error {
	tClientCaseFields := tCase.CustomFields
	amount := tClientCaseFields.TextValueByNameBasic("amount")
	lib.DPrintln(amount)
	currentRating := tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating)
	newRating := tClientCaseFields.NumberValueByNameBasic(FieldName_new_rating)
	if newRating == 0 {
		newRating = 100
	}

	increaseAmount, err := c.FeeUsecase.GetIncreaseAmount(int(currentRating), int(newRating))
	if err != nil {
		return err
	}

	aa := ((300 / 30) * float32(increaseAmount)) * 0.3333
	newAmount := CentToDollar(aa)
	isOk := false
	if amount != "" {
		amountDecimal, _ := decimal.NewFromString(amount)
		newAmountDeciamal := decimal.NewFromInt(int64(newAmount))
		if !amountDecimal.Equal(newAmountDeciamal) {
			isOk = true
		}

	} else {
		isOk = true
	}
	if isOk {
		destData := make(TypeDataEntry)
		destData[DataEntry_gid] = tClientCaseFields.TextValueByNameBasic("gid")
		destData[FieldName_amount] = newAmount
		c.DataEntryUsecase.HandleOne(Kind_client_cases, destData, DataEntry_gid, nil)
	}

	return nil
}

func (c *ZohobuzUsecase) HandleAmount(clientCaseId int32) error {
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return err
	}
	if tClientCase == nil {
		return errors.New("HandleAmount: tClientCase is nil.")
	}

	tClientCaseFields := tClientCase.CustomFields
	if tClientCaseFields.TextValueByNameBasic(FieldName_stages) == config_vbc.Stages_Completed ||
		tClientCaseFields.TextValueByNameBasic(FieldName_stages) == config_vbc.Stages_AmCompleted {
		return nil
	}
	amount := tClientCaseFields.TextValueByNameBasic("amount")

	if IsAmContract(*tClientCase) {
		return c.HandleAmountForAm(*tClientCase)
	}

	amountDecimal, _ := decimal.NewFromString(amount)
	newAmount, err, noCaseContractBasicDataVo := c.FeeUsecase.ClientCaseAmount(tClientCase)
	c.log.Info("HandleAmount: ", tClientCase.Id(), " newAmount: ", newAmount, " amountDecimal: ", amountDecimal)
	if err != nil {
		return err
	}
	if noCaseContractBasicDataVo {
		c.log.Info("HandleAmount: ", tClientCase.Id(), " noCaseContractBasicDataVo")
		return nil
	}

	newAmountDecimal := decimal.NewFromInt(int64(newAmount))

	if amountDecimal.Equal(newAmountDecimal) {
		c.log.Info("HandleAmount: amountDecimal.Equal(newAmountDecimal) ", clientCaseId, " ", amount, " ", newAmount)
		return nil
	} else {
		c.log.Info("HandleAmount: No Equal ", clientCaseId, " ", amount, " ", newAmount)
	}

	if configs.StoppedZoho {
		destData := make(TypeDataEntry)
		destData[DataEntry_gid] = tClientCaseFields.TextValueByNameBasic("gid")
		destData[FieldName_amount] = newAmountDecimal.String()
		c.DataEntryUsecase.HandleOne(Kind_client_cases, destData, DataEntry_gid, nil)
	} else {
		err = c.ZohoUsecase.ChangeDealAmount(tClientCaseFields.TextValueByNameBasic("gid"), newAmountDecimal.String())
		c.UsageStatsUsecase.Stat(UsageType_HandleAmount, time.Now(), 1)
	}
	return err
}

func ClientCaseNameByCase(firstName string, lastName string, tCase TData) string {
	return ClientCaseName(firstName, lastName, tCase.CustomFields.NumberValueByNameBasic(FieldName_current_rating), tCase.Id())
}
func ClientCaseName(firstName string, lastName string, rating int32, clientCaseId int32) string {
	return fmt.Sprintf("%s %s-%d#%d", firstName, lastName, rating, clientCaseId)
}

func (c *ZohobuzUsecase) HandleClientCaseName(clientCaseId int32) error {
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return err
	}
	if tClientCase == nil {
		return errors.New("HandleClientCaseName: tClientCase is nil.")
	}
	tClientCaseFields := tClientCase.CustomFields
	_, tClientFields, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return err
	}
	if tClientFields == nil {
		return errors.New("HandleClientCaseName: tClientFields is nil.")
	}
	clientCaseName := ClientCaseNameByCase(tClientFields.TextValueByNameBasic("first_name"),
		tClientFields.TextValueByNameBasic("last_name"),
		*tClientCase,
	)
	if clientCaseName != tClientCaseFields.TextValueByNameBasic("deal_name") {
		if configs.StoppedZoho {
			destData := make(lib.TypeMap)
			destData.Set(DataEntry_gid, tClientCaseFields.TextValueByNameBasic("gid"))
			destData.Set(FieldName_deal_name, clientCaseName)
			_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry(destData), DataEntry_gid, nil)
		} else {
			c.log.Info("HandleClientCaseName: do ", clientCaseId, " ", clientCaseName)
			params := make(lib.TypeMap)
			params.Set("Deal_Name", clientCaseName)
			err = c.ZohoUsecase.ChangeDealV1(tClientCaseFields.TextValueByNameBasic("gid"), params)
			c.UsageStatsUsecase.Stat(UsageType_HandleClientCaseName, time.Now(), 1)
		}
		return err
	} else {
		c.log.Info("HandleClientCaseName: info ", clientCaseId, " ", clientCaseName)
		return nil
	}
}

func (c *ZohobuzUsecase) HandleSyncZohoPricingVersion(tCase *TData) error {
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	if configs.StoppedZoho {
		return nil
	}
	tCaseFields := tCase.CustomFields

	if tCaseFields.TextValueByNameBasic(FieldName_pricing_version) != "" {
		return nil
	}

	sPricingVersion := tCaseFields.TextValueByNameBasic(FieldName_s_pricing_version)
	if sPricingVersion == "" {
		return nil
	}

	params := make(lib.TypeMap)
	params.Set("Pricing_Version", sPricingVersion)
	c.log.Debug(params)
	return c.ZohoUsecase.ChangeDealV1(tCaseFields.TextValueByNameBasic("gid"), params)
}

func (c *ZohobuzUsecase) HttpHandleClientCaseName(ctx *gin.Context) {
	reply := CreateReply()
	err := c.BizHttpHandleClientCaseName(ctx.Query("clientCaseIds"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ZohobuzUsecase) BizHttpHandleClientCaseName(clientCaseIds string) error {
	clientCaseIdsArr := strings.Split(clientCaseIds, ",")

	var data lib.TypeList

	for _, v := range clientCaseIdsArr {
		clientCaseId := lib.InterfaceToInt32(v)
		a, err := c.DoClientCaseName(clientCaseId)
		if err != nil {
			c.log.Error("BizHttpHandleClientCaseName: err:", err)
		} else if a != nil {
			data = append(data, a)
		}
	}
	if len(data) > 0 {
		records := make(lib.TypeMap)
		records.Set("data", data)
		_, err := c.ZohoUsecase.PutRecordsV1(config_zoho.Deals, records)
		return err
	}
	return nil
}

// DoClientCaseName 不为nil，且没有err，说明需要同步
func (c *ZohobuzUsecase) DoClientCaseName(clientCaseId int32) (lib.TypeMap, error) {
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return nil, err
	}
	if tClientCase == nil {
		return nil, errors.New("DoClientCaseName: tClientCase is nil.")
	}
	tClientCaseFields := tClientCase.CustomFields
	_, tClientFields, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return nil, err
	}
	if tClientFields == nil {
		return nil, errors.New("DoClientCaseName: tClientFields is nil.")
	}
	clientCaseName := ClientCaseNameByCase(tClientFields.TextValueByNameBasic("first_name"),
		tClientFields.TextValueByNameBasic("last_name"),
		*tClientCase,
	)
	if clientCaseName != tClientCaseFields.TextValueByNameBasic("deal_name") {
		params := make(lib.TypeMap)
		params.Set("id", tClientCaseFields.TextValueByNameBasic("gid"))
		params.Set("Deal_Name", clientCaseName)
		return params, nil
	} else {
		fmt.Println("====", clientCaseName)
		return nil, nil
	}
}

//func (c *ZohobuzUsecase) Http() error {
//	sql := `select asana_migrate.client_case_gid from asana_migrate
//inner join vbcdb.clients c on c.asana_task_gid= asana_migrate.from_asana_gid
//where c.notes !=""`

//}

func (c *ZohobuzUsecase) HttpHandleNotes(ctx *gin.Context) {
	reply := CreateReply()
	err := c.BizHttpHandleNotes(ctx.Query("clientCaseIds"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ZohobuzUsecase) BizHttpHandleNotes(clientCaseIds string) error {
	clientCaseIdsArr := strings.Split(clientCaseIds, ",")

	var data lib.TypeList

	for _, v := range clientCaseIdsArr {
		clientCaseId := lib.InterfaceToInt32(v)
		a, err := c.DoNotes(clientCaseId)
		if err != nil {
			c.log.Error("BizHttpHandleNotes: err:", err)
		} else if a != nil {
			data = append(data, a)
		}
	}
	if len(data) > 0 {
		records := make(lib.TypeMap)
		records.Set("data", data)
		_, err := c.ZohoUsecase.PutRecordsV1(config_zoho.Deals, records)
		return err
	}
	return nil
}

// HttpHandleDataCollection 同步client cases信息（当前同步Case files folder）
func (c *ZohobuzUsecase) HttpHandleDataCollection(ctx *gin.Context) {
	reply := CreateReply()
	err := c.BizHttpHandleDataCollection(ctx.Query("clientCaseIds"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ZohobuzUsecase) BizHttpHandleDataCollection(clientCaseIds string) error {
	clientCaseIdsArr := strings.Split(clientCaseIds, ",")
	for _, v := range clientCaseIdsArr {
		clientCaseId := lib.InterfaceToInt32(v)
		err := c.ActionOnceUsecase.HandleDataCollectionFolder(clientCaseId)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
// Box Folder Link syncs to the Zoho
	params := make(lib.TypeMap)
	params.Set("id", tClientCase.CustomFields.TextValueByNameBasic("gid"))
	params.Set("Data_Collection_Folder", "https://veteranbenefitscenter.app.box.com/folder/"+boxFolderId)
	_, _, err = c.ZohoUsecase.PutRecord(config_zoho.Deals, params)
	if err != nil {
		return err
	}
*/

func (c *ZohobuzUsecase) HttpDCFolderIdToZoho(ctx *gin.Context) {
	reply := CreateReply()
	err := c.BizHttpDCFolderIdToZoho()
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ZohobuzUsecase) BizHttpDCFolderIdToZoho() error {

	list, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{"biz_deleted_at": 0})
	if err != nil {
		return err
	}
	for _, v := range list {
		dataCollectionFolder := v.CustomFields.TextValueByNameBasic("data_collection_folder")
		if dataCollectionFolder != "" {
			continue
		}
		key := MapKeyDataCollectionFolderId(v.CustomFields.NumberValueByNameBasic("id"))
		dcFolderId, err := c.MapUsecase.GetForString(key)
		if err != nil {
			return err
		}
		if dcFolderId == "" {
			continue
		}
		gid := v.CustomFields.TextValueByNameBasic("gid")

		params := make(lib.TypeMap)
		params.Set("id", gid)
		params.Set("Data_Collection_Folder", "https://veteranbenefitscenter.app.box.com/folder/"+dcFolderId)
		_, _, err = c.ZohoUsecase.PutRecordV1(config_zoho.Deals, params)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}

// HttpHandleCaseInfo 同步client cases信息（当前同步Case files folder）
func (c *ZohobuzUsecase) HttpHandleCaseInfo(ctx *gin.Context) {
	reply := CreateReply()
	err := c.BizHttpHandleCaseInfo(ctx.Query("clientCaseIds"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ZohobuzUsecase) BizHttpHandleCaseInfo(clientCaseIds string) error {
	clientCaseIdsArr := strings.Split(clientCaseIds, ",")

	var data lib.TypeList

	for _, v := range clientCaseIdsArr {
		clientCaseId := lib.InterfaceToInt32(v)
		a, err := c.DoCaseInfo(clientCaseId)
		if err != nil {
			c.log.Error("BizHttpHandleCaseInfo: err:", err)
		} else if a != nil {
			data = append(data, a)
		}
	}
	if len(data) > 0 {
		records := make(lib.TypeMap)
		records.Set("data", data)
		_, err := c.ZohoUsecase.PutRecordsV1(config_zoho.Deals, records)
		return err
	}
	return nil
}

// HttpHandleClientInfo 同步地址信息
func (c *ZohobuzUsecase) HttpHandleClientInfo(ctx *gin.Context) {
	reply := CreateReply()
	err := c.BizHttpHandleClientInfo(ctx.Query("clientCaseIds"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ZohobuzUsecase) BizHttpHandleClientInfo(clientCaseIds string) error {
	clientCaseIdsArr := strings.Split(clientCaseIds, ",")

	var data lib.TypeList

	for _, v := range clientCaseIdsArr {
		clientCaseId := lib.InterfaceToInt32(v)
		a, err := c.DoClientInfo(clientCaseId)
		if err != nil {
			c.log.Error("BizHttpHandleNotes: err:", err)
		} else if a != nil {
			data = append(data, a)
		}
	}
	if len(data) > 0 {
		records := make(lib.TypeMap)
		records.Set("data", data)
		_, err := c.ZohoUsecase.PutRecordsV1(config_zoho.Deals, records)
		return err
	}
	return nil
}

func (c *ZohobuzUsecase) DoCaseInfo(clientCaseId int32) (lib.TypeMap, error) {
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return nil, err
	}
	if tClientCase == nil {
		return nil, errors.New("DoCaseInfo: tClientCase is nil.")
	}
	tClientCaseFields := tClientCase.CustomFields
	row := make(lib.TypeMap)

	deal, err := c.ZohoUsecase.GetDeal(tClientCaseFields.TextValueByNameBasic("gid"))
	if err != nil {
		return nil, err
	}
	if deal.GetString("Case_Files_Folder") == "" {
		key := fmt.Sprintf("%s%d", Map_ClientBoxFolderId, clientCaseId)
		boxFolderId, _ := c.MapUsecase.GetForString(key)
		if boxFolderId != "" {
			row.Set("Case_Files_Folder", "https://veteranbenefitscenter.app.box.com/folder/"+boxFolderId)
		}
	}
	//if tClientFields.TextValueByNameBasic("phone") != "" {
	//	if tClientCaseFields.TextValueByNameBasic("phone") == "" {
	//		row.Set("Phone", tClientFields.TextValueByNameBasic("phone"))
	//	}
	//}
	//if tClientFields.TextValueByNameBasic("ssn") != "" {
	//	if tClientCaseFields.TextValueByNameBasic("ssn") == "" {
	//		row.Set("SSN", tClientFields.TextValueByNameBasic("ssn"))
	//	}
	//}
	//if tClientFields.TextValueByNameBasic("dob") != "" {
	//	if tClientCaseFields.TextValueByNameBasic("dob") == "" {
	//		row.Set("Date_of_Birth", tClientFields.TextValueByNameBasic("dob"))
	//	}
	//}
	//if tClientFields.TextValueByNameBasic("state") != "" {
	//	if tClientCaseFields.TextValueByNameBasic("state") == "" {
	//		row.Set("State", tClientFields.TextValueByNameBasic("state"))
	//	}
	//}
	//if tClientFields.TextValueByNameBasic("city") != "" {
	//	if tClientCaseFields.TextValueByNameBasic("city") == "" {
	//		row.Set("City", tClientFields.TextValueByNameBasic("city"))
	//	}
	//}
	//if tClientFields.TextValueByNameBasic("address") != "" {
	//	if tClientCaseFields.TextValueByNameBasic("address") == "" {
	//		row.Set("Street_Address", tClientFields.TextValueByNameBasic("address"))
	//	}
	//}
	//if tClientFields.TextValueByNameBasic("zip_code") != "" {
	//	if tClientCaseFields.TextValueByNameBasic("zip_code") == "" {
	//		row.Set("Zip_Code", tClientFields.TextValueByNameBasic("zip_code"))
	//	}
	//}
	//
	if len(row) == 0 {
		return nil, nil
	}
	row.Set("id", tClientCaseFields.TextValueByNameBasic("gid"))

	return row, nil
}

func (c *ZohobuzUsecase) DoClientInfo(clientCaseId int32) (lib.TypeMap, error) {
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return nil, err
	}
	if tClientCase == nil {
		return nil, errors.New("DoNotes: tClientCase is nil.")
	}
	tClientCaseFields := tClientCase.CustomFields
	clientGid := tClientCase.CustomFields.TextValueByNameBasic("client_gid")

	_, tClientFields, err := c.DataComboUsecase.Client(clientGid)
	if err != nil {
		return nil, err
	}
	row := make(lib.TypeMap)

	if tClientFields.TextValueByNameBasic("email") != "" {
		row.Set("Email", tClientFields.TextValueByNameBasic("email"))
	}
	if tClientFields.TextValueByNameBasic("phone") != "" {
		if tClientCaseFields.TextValueByNameBasic("phone") == "" {
			row.Set("Phone", tClientFields.TextValueByNameBasic("phone"))
		}
	}
	if tClientFields.TextValueByNameBasic("ssn") != "" {
		if tClientCaseFields.TextValueByNameBasic("ssn") == "" {
			row.Set("SSN", tClientFields.TextValueByNameBasic("ssn"))
		}
	}
	if tClientFields.TextValueByNameBasic("dob") != "" {
		if tClientCaseFields.TextValueByNameBasic("dob") == "" {
			row.Set("Date_of_Birth", tClientFields.TextValueByNameBasic("dob"))
		}
	}
	if tClientFields.TextValueByNameBasic("state") != "" {
		if tClientCaseFields.TextValueByNameBasic("state") == "" {
			row.Set("State", tClientFields.TextValueByNameBasic("state"))
		}
	}
	if tClientFields.TextValueByNameBasic("city") != "" {
		if tClientCaseFields.TextValueByNameBasic("city") == "" {
			row.Set("City", tClientFields.TextValueByNameBasic("city"))
		}
	}
	if tClientFields.TextValueByNameBasic("address") != "" {
		if tClientCaseFields.TextValueByNameBasic("address") == "" {
			row.Set("Street_Address", tClientFields.TextValueByNameBasic("address"))
		}
	}
	if tClientFields.TextValueByNameBasic("zip_code") != "" {
		if tClientCaseFields.TextValueByNameBasic("zip_code") == "" {
			row.Set("Zip_Code", tClientFields.TextValueByNameBasic("zip_code"))
		}
	}

	row.Set("id", tClientCaseFields.TextValueByNameBasic("gid"))
	return row, nil
}

func (c *ZohobuzUsecase) DoNotes(clientCaseId int32) (lib.TypeMap, error) {
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return nil, err
	}
	if tClientCase == nil {
		return nil, errors.New("DoNotes: tClientCase is nil.")
	}
	tClientCaseFields := tClientCase.CustomFields
	fmt.Println(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))

	tClientCaseFields.TextValueByNameBasic("notes")

	migrateEntity, err := c.AsanaMigrateUsecase.GetByCond(Eq{"client_case_id": tClientCaseFields.NumberValueByNameBasic("id")})
	if err != nil {
		return nil, err
	}
	if migrateEntity == nil {
		return nil, errors.New("migrateEntity is nil")
	}

	_, list, err := c.AsanaMigrateUsecase.Data(migrateEntity.FromAsanaGid)
	if err != nil {
		return nil, err
	}
	if len(list) <= 0 {
		return nil, errors.New("asana row is nil.")
	}
	for k, _ := range list {
		row := lib.TypeMap(list[k])
		notes := row.GetString("notes")
		if notes != "" {
			params := make(lib.TypeMap)
			params.Set("id", tClientCaseFields.TextValueByNameBasic("gid"))
			params.Set("Description", notes)
			return params, nil
		}
		//row.
	}
	return nil, nil
}
