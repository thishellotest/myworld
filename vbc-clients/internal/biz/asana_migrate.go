package biz

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

/*
迁移注册事项：
- 确认asana的邮箱不重复
*/

func AsanaMigrateDbName() string {
	dbName := "vbcdev"
	if configs.IsProd() {
		dbName = "vbcdb"
	}
	return dbName
}

func AsanaMigrateSql(asanaClientGid string) string {
	dbName := AsanaMigrateDbName()
	var sql string
	if asanaClientGid != "" {
		sql = fmt.Sprintf("select * from %s.clients where deleted_at=0 and biz_deleted_at=0 and email!='' and asana_task_gid='%s'", dbName, asanaClientGid)
	} else {
		sql = fmt.Sprintf("select * from %s.clients where deleted_at=0 and biz_deleted_at=0 and stages!='' and email!='' and asana_projects like '%1206472580135542%' ", dbName)
	}
	return sql
}

/*
CREATE TABLE `asana_migrate` (

	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`client_gid` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'zoho',
	`cleint_id` int(11) NOT NULL DEFAULT '0' COMMENT 'zoho',
	`client_case_gid` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'zoho',
	`client_case_id` int(11) NOT NULL DEFAULT '0' COMMENT 'zoho',
	`from_client_id` int(11) NOT NULL DEFAULT '0' COMMENT 'asana',
	`from_asana_gid` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'asana',
	`created_at` int(11) NOT NULL DEFAULT '0',
	`updated_at` int(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`),
	UNIQUE KEY `uniq_asana_gid` (`from_asana_gid`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
*/
type AsanaMigrateEntity struct {
	ID            int32 `gorm:"primaryKey"`
	ClientGid     string
	ClientId      int32
	ClientCaseGid string
	ClientCaseId  int32
	FromClientId  int32
	FromAsanaGid  string
	CreatedAt     int64
	UpdatedAt     int64
}

func (AsanaMigrateEntity) TableName() string {
	return "asana_migrate"
}

type AsanaMigrateUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[AsanaMigrateEntity]
	TUsecase         *TUsecase
	ZohoUsecase      *ZohoUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewAsanaMigrateUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	ZohoUsecase *ZohoUsecase,
	DataEntryUsecase *DataEntryUsecase) *AsanaMigrateUsecase {
	uc := &AsanaMigrateUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		ZohoUsecase:      ZohoUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *AsanaMigrateUsecase) Data(asanaClientGid string) (columns []string, list []map[string]interface{}, err error) {
	sqlRows, err := c.CommonUsecase.DB().Raw(AsanaMigrateSql(asanaClientGid)).Rows()
	if err != nil {
		return nil, nil, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	columns, list, err = lib.SqlRowsTrans(sqlRows)
	return
}

// HttpSyncData 同步未完成迁移的字段
func (c *AsanaMigrateUsecase) HttpSyncData(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	err := c.BizHttpSyncData(ctx.Query("asana_client_gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AsanaMigrateUsecase) BizHttpSyncData(asanaClientGid string) error {

	_, list, err := c.Data(asanaClientGid)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return errors.New("list is length 0")
	}
	for _, v := range list {
		err = c.SyncDataOne(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *AsanaMigrateUsecase) SyncDataOne(asanaData lib.TypeMap) error {

	fromClientId := asanaData.GetInt("id")
	fromClientGid := asanaData.GetString("asana_task_gid")
	//fromEmail := asanaData.GetString("email")
	if fromClientGid == "" {
		return errors.New("fromClientGid is empty.")
	}
	if fromClientId <= 0 {
		return errors.New("fromClientId is wrong")
	}
	migrateEntity, err := c.GetByCond(Eq{"from_asana_gid": fromClientGid})
	if err != nil {
		return err
	}
	if migrateEntity == nil {
		return errors.New("没有迁移不同步")
	}
	zohoContact := make(lib.TypeMap)
	zohoContact.Set("id", migrateEntity.ClientGid)
	//zohoContact.Set("Mailing_Street", asanaData.GetString("address"))
	zohoContact.Set("Branch", asanaData.GetString("branch"))
	_, _, err = c.ZohoUsecase.PutRecordV1(config_zoho.Contacts, zohoContact)
	if err != nil {
		return err
	}

	//deal := make(lib.TypeMap)
	//deal.Set("id", migrateEntity.ClientCaseGid)
	//deal.Set("Stage", AsanaStagesToZohoStages(asanaData.GetString("stages")))
	//_, _, err = c.ZohoUsecase.PutRecord(config_zoho.Deals, deal)

	return err
}

func (c *AsanaMigrateUsecase) HttpMigrateOne(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	err := c.BizHttpMigrateOne(ctx.Query("asana_client_gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AsanaMigrateUsecase) BizHttpMigrateOne(asanaClientGid string) error {
	if asanaClientGid == "" {
		return errors.New("asanaClientGid is empty.")
	}
	_, list, err := c.Data(asanaClientGid)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return errors.New("list is length 0")
	}
	for _, v := range list {
		err = c.MigrateOne(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *AsanaMigrateUsecase) MigrateOne(data lib.TypeMap) error {

	fromClientId := data.GetInt("id")
	fromClientGid := data.GetString("asana_task_gid")
	fromEmail := data.GetString("email")
	if fromClientGid == "" {
		return errors.New("fromClientGid is empty.")
	}
	if fromClientId <= 0 {
		return errors.New("fromClientId is wrong")
	}
	migrateEntity, err := c.GetByCond(Eq{"from_asana_gid": fromClientGid})
	if err != nil {
		return err
	}

	// 处理maps
	err = c.HandleMaps(fromClientId)
	if err != nil {
	}

	err = c.HandleBehaviors(fromClientId)
	if err != nil {
	}

	if migrateEntity != nil {
		return nil
	}

	clientId, clientGid, err := c.HandleZohoContact(fromEmail, data)
	if err != nil {
		return err
	}
	clientCaseId, clientCaseGid, err := c.HandleZohoDeal(clientGid, data)
	if err != nil {
		return err
	}

	asanaMigrateEntity := &AsanaMigrateEntity{
		ClientGid:     clientGid,
		ClientId:      clientId,
		ClientCaseGid: clientCaseGid,
		ClientCaseId:  clientCaseId,
		FromAsanaGid:  fromClientGid,
		FromClientId:  fromClientId,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}
	return c.CommonUsecase.DB().Save(&asanaMigrateEntity).Error
}

func (c *AsanaMigrateUsecase) HandleZohoDeal(clientGid string, asanaData lib.TypeMap) (clientCaseId int32, clientCaseGid string, err error) {

	// 判断记录是否存在
	uniqcode := asanaData.GetString("uniqcode")
	tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"uniqcode": uniqcode})
	if err != nil {
		return 0, "", err
	}
	if tCase != nil {
		return tCase.CustomFields.NumberValueByNameBasic("id"), tCase.CustomFields.TextValueByNameBasic("gid"), nil
	}

	clientCaseGid, err = c.CreateZohoDeal(clientGid, asanaData)
	if err != nil {
		return 0, "", err
	}

	data := make(lib.TypeMap)
	data.Set("id", asanaData.GetInt("id"))
	data.Set("gid", clientCaseGid)
	data.Set("created_at", time.Now().Unix())
	data.Set("updated_at", time.Now().Unix())
	data.Set("uniqcode", uniqcode)
	err = c.CommonUsecase.DB().Table("client_cases").Create(data.ToOrigin()).Error
	if err != nil {
		return 0, "", err
	}
	var uniqueCodeGeneratorEntity UniqueCodeGeneratorEntity
	c.CommonUsecase.DB().Table(AsanaMigrateDbName()+".unique_code_generator").
		Where("uuid=?", uniqcode).
		Find(&uniqueCodeGeneratorEntity)
	er := c.CommonUsecase.DB().Save(&uniqueCodeGeneratorEntity).Error
	if er != nil {
		c.log.Error(er)
	}

	tClientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"gid": clientCaseGid})
	if err != nil {
		return 0, "", err
	}
	if tClientCase == nil {
		return 0, "", errors.New("tClientCase is nil.")
	}
	return tClientCase.CustomFields.NumberValueByNameBasic("id"), clientCaseGid, nil
}

func HandleOwnerId(source lib.TypeMap, asanaData lib.TypeMap) {
	if asanaData.GetString("assignee_gid") == "1205444097333494" { // Edward Bunting
		source.Set("Owner.id", 6159272000000453640)
	} else if asanaData.GetString("assignee_gid") == "1206686474861326" { //  Gary
		source.Set("Owner.id", 6159272000000453669)
	} else {
		source.Set("Owner.id", 6159272000000453001) // Yannan
	}
}

func (c *AsanaMigrateUsecase) CreateZohoDeal(clientGid string, asanaData lib.TypeMap) (clientCaseGid string, err error) {
	deal := make(lib.TypeMap)

	//  {
	//"id" : "425248000000104001"
	//}

	//deal.Set("Contact_Name", map[string]int64{
	//	"id": lib.InterfaceToInt64(clientGid),
	//})

	HandleOwnerId(deal, asanaData)

	deal.Set("Contact_Name", lib.InterfaceToInt64(clientGid))
	deal.Set("Pipeline", "VBC Clients")
	deal.Set("Current_Rating", asanaData.GetInt("current_rating"))
	deal.Set("Effective_Current_Rating", asanaData.GetInt("effective_current_rating"))
	deal.Set("Atomic_Veterans_and_Radiation_Exposure", asanaData.GetString("atomic_veterans"))
	deal.Set("Branch", asanaData.GetString("branch"))
	deal.Set("Burn_Pits_and_Other_Airborne_Hazards", asanaData.GetString("burn_pits"))
	deal.Set("C_File_Submitted", asanaData.GetString("c_file_submitted"))
	deal.Set("Contact_Form", asanaData.GetString("contact_form"))
	deal.Set("DD214", asanaData.GetString("dd214"))
	deal.Set("Deal_Name", asanaData.GetString("last_name")+" "+asanaData.GetString("first_name"))
	deal.Set("Disability_Rating_List_Screenshot", asanaData.GetString("disability_rating"))
	deal.Set("Gulf_War_Illness", asanaData.GetString("gulf_war"))
	deal.Set("Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun", asanaData.GetString("illness_due"))
	deal.Set("Lead_Source", asanaData.GetString("source"))
	deal.Set("New_Rating", asanaData.GetString("new_rating"))
	deal.Set("Rating_Decision_Letters", asanaData.GetString("rating_decision"))
	deal.Set("Referring_Person", asanaData.GetString("referrer"))
	deal.Set("Retired", asanaData.GetString("retired"))
	deal.Set("STRs", asanaData.GetString("strs"))
	deal.Set("Stage", AsanaStagesToZohoStages(asanaData.GetString("stages")))
	deal.Set("TDIU", asanaData.GetString("tdiu"))
	deal.Set("Agent_Orange_Exposure", asanaData.GetString("agent_orange"))
	deal.Set("Amyotrophic_Lateral_Sclerosis_ALS", asanaData.GetString("amyotrophic"))

	if AsanaMigrateVerifyFieldDateTypeValue(asanaData.GetString("itf_expiration")) {
		deal.Set("ITF_Expiration", asanaData.GetString("itf_expiration"))
	}

	clientCaseGid, _, err = c.ZohoUsecase.CreateRecord(config_zoho.Deals, deal)
	return
}

func (c *AsanaMigrateUsecase) HandleMaps(fromClientId int32) error {
	if fromClientId <= 0 {
		return errors.New("fromClientId is wrong.")
	}
	var maps []MapEntity
	err := c.CommonUsecase.DB().Raw(fmt.Sprintf("select * from %s.maps where mkey like '%%:%d'", AsanaMigrateDbName(), fromClientId)).Find(&maps).Error
	if err != nil {
		return err
	}
	er := c.CommonUsecase.DB().Save(maps).Error
	if er != nil {
		c.log.Error(er)
	}
	return nil
}

func (c *AsanaMigrateUsecase) HandleBehaviors(fromClientId int32) error {
	if fromClientId <= 0 {
		return errors.New("fromClientId is wrong.")
	}
	var entities []BehaviorEntity
	err := c.CommonUsecase.DB().
		Raw(fmt.Sprintf("select * from %s.behaviors where incr_id=%d", AsanaMigrateDbName(), fromClientId)).
		Find(&entities).Error
	if err != nil {
		return err
	}
	er := c.CommonUsecase.DB().Save(entities).Error
	if er != nil {
		c.log.Error(er)
	}
	return nil
}

func (c *AsanaMigrateUsecase) HandleZohoContact(email string, asanaData lib.TypeMap) (clientId int32, clientGid string, err error) {

	// 通过email查看contact是否存在
	tContact, err := c.TUsecase.Data(Kind_clients, Eq{"email": email})
	if err != nil {
		return 0, "", err
	}
	if tContact == nil {
		// 开始创建contact
		clientGid, err = c.CreateZohoContact(email, asanaData)
		if err != nil {
			return 0, "", err
		}
		data := make(lib.TypeMap)
		data.Set("id", asanaData.GetInt("id"))
		data.Set("gid", clientGid)
		data.Set("email", email)
		err = c.CommonUsecase.DB().Table("clients").Create(data.ToOrigin()).Error
		if err != nil {
			return 0, "", err
		}

		tContact, err = c.TUsecase.Data(Kind_clients, Eq{"gid": clientGid})
		if err != nil {
			return 0, "", err
		}
		if tContact == nil {
			return 0, "", errors.New("tContact is nil.")
		}
	}
	return tContact.CustomFields.NumberValueByNameBasic("id"), tContact.CustomFields.TextValueByNameBasic("gid"), nil
}

func AsanaMigrateVerifyFieldDateTypeValue(val string) bool {
	if val == "" || val == "0" {
		return false
	}
	return true
}

func (c *AsanaMigrateUsecase) CreateZohoContact(email string, asanaData lib.TypeMap) (clientGid string, err error) {
	zohoContact := make(lib.TypeMap)

	HandleOwnerId(zohoContact, asanaData)

	zohoContact.Set("Email", email)
	zohoContact.Set("Current_Rating", asanaData.GetInt("current_rating"))
	zohoContact.Set("Effective_Current_Rating", asanaData.GetInt("effective_current_rating"))
	zohoContact.Set("First_Name", asanaData.GetString("first_name"))
	zohoContact.Set("Last_Name", asanaData.GetString("last_name"))
	zohoContact.Set("Lead_Source", asanaData.GetString("source"))
	zohoContact.Set("Mailing_City", asanaData.GetString("city"))
	zohoContact.Set("Mailing_State", asanaData.GetString("address_state"))
	zohoContact.Set("Mailing_Street", asanaData.GetString("address"))
	zohoContact.Set("Mailing_Zip", asanaData.GetString("zip_code"))
	zohoContact.Set("Mobile", asanaData.GetString("phone"))
	zohoContact.Set("Referring_Person", asanaData.GetString("referrer"))
	zohoContact.Set("Retired", asanaData.GetString("retired"))
	zohoContact.Set("SSN", asanaData.GetString("ssn"))

	if AsanaMigrateVerifyFieldDateTypeValue(asanaData.GetString("dob")) {
		zohoContact.Set("Date_of_Birth", asanaData.GetString("dob"))
	}

	clientGid, _, err = c.ZohoUsecase.CreateRecord(config_zoho.Contacts, zohoContact)
	return
}

func AsanaStagesToZohoStages(stages string) string {
	return stages

	//if stages == vbc_config.AsanaStages_IncomingRequest {
	//	return vbc_config.Stages_IncomingRequest
	//} else if stages == vbc_config.AsanaStages_FeeScheduleandContract {
	//	return vbc_config.Stages_FeeScheduleandContract
	//} else if stages == vbc_config.AsanaStages_GettingStartedEmail {
	//	return vbc_config.Stages_GettingStartedEmail
	//} else if stages == vbc_config.AsanaStages_UpdateClientIntakeInfo {
	//	return vbc_config.Stages_AwaitingClientRecords
	//} else if stages == vbc_config.AsanaStages_AwaitingCFile {
	//	return vbc_config.Stages_AwaitingClientRecords
	//} else if stages == vbc_config.AsanaStages_RecordReview {
	//	return vbc_config.Stages_RecordReview
	//} else if stages == vbc_config.AsanaStages_ScheduleCall {
	//	return vbc_config.Stages_ScheduleCall
	//} else if stages == vbc_config.AsanaStages_Statement_Notes {
	//	return vbc_config.Stages_StatementNotes
	//} else if stages == vbc_config.AsanaStages_StatementDrafts {
	//	return vbc_config.Stages_StatementDrafts
	//} else if stages == vbc_config.AsanaStages_StatementsFinalized {
	//	return vbc_config.Stages_StatementsFinalized
	//} else if stages == vbc_config.AsanaStages_CurrentTreatment {
	//	return vbc_config.Stages_CurrentTreatment
	//} else if stages == vbc_config.AsanaStages_MiniDBQs {
	//	return vbc_config.Stages_MiniDBQs_Draft
	//} else if stages == vbc_config.AsanaStages_NexusLetters {
	//	return vbc_config.Stages_NexusLetters
	//} else if stages == vbc_config.AsanaStages_MedicalTeam {
	//	return vbc_config.Stages_MedicalTeam
	//} else if stages == vbc_config.AsanaStages_FileClaims {
	//	return vbc_config.Stages_FileClaims_Draft
	//} else if stages == vbc_config.AsanaStages_VerifyEvidenceReceived {
	//	return vbc_config.Stages_VerifyEvidenceReceived
	//} else if stages == vbc_config.AsanaStages_AwaitingDecision {
	//	return vbc_config.Stages_AwaitingDecision
	//} else if stages == vbc_config.AsanaStages_AwaitingPayment {
	//	return vbc_config.Stages_AwaitingPayment
	//} else if stages == vbc_config.AsanaStages_Completed {
	//	return vbc_config.Stages_Completed
	//}
	//return stages
}
