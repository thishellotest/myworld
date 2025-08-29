package biz

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ZohoinfoSyncUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	TUsecase         *TUsecase
	DataComboUsecase *DataComboUsecase
	ZohoUsecase      *ZohoUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewZohoinfoSyncUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	DataComboUsecase *DataComboUsecase,
	ZohoUsecase *ZohoUsecase,
	DataEntryUsecase *DataEntryUsecase) *ZohoinfoSyncUsecase {
	uc := &ZohoinfoSyncUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		DataComboUsecase: DataComboUsecase,
		ZohoUsecase:      ZohoUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

func (c *ZohoinfoSyncUsecase) HttpSync(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	err := c.Sync(ctx.Query("kind"), lib.InterfaceToInt32(ctx.Query("id")), ctx.Query("field_name"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

// Sync 同步信息 changeFieldName: 发生改变的字段值
func (c *ZohoinfoSyncUsecase) Sync(kind string, incrId int32, changeFieldName string) error {
	if changeFieldName == "" {
		return errors.New("changeFieldName is empty.")
	}
	if Kind_clients == kind {
		tClient, err := c.TUsecase.DataById(kind, incrId)
		if err != nil {
			return err
		}
		if tClient == nil {
			return errors.New("tClient is nil")
		}
		tClientCases, err := c.TUsecase.ListByCond(Kind_client_cases,
			Eq{"client_gid": tClient.CustomFields.TextValueByNameBasic("gid"),
				"biz_deleted_at": 0})
		if err != nil {
			return err
		}
		for _, tClientCase := range tClientCases {
			syncBaseInfo, err := ZohoBaseInfoSync(tClient, tClientCase, changeFieldName)
			if err != nil {
				return err
			}
			if syncBaseInfo != nil { // 需要同步

				lib.DPrintln("ZohoinfoSyncUsecase syncBaseInfo:"+kind, syncBaseInfo)

				err = c.SyncRow(Kind_client_cases,
					tClientCase.CustomFields.TextValueByNameBasic("gid"),
					syncBaseInfo)
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else if Kind_client_cases == kind {

		tClientCaseBase, err := c.TUsecase.DataById(kind, incrId)
		if err != nil {
			return err
		}
		if tClientCaseBase == nil {
			return errors.New("tClientCase is nil.")
		}
		tClient, err := c.TUsecase.Data(Kind_clients, Eq{"gid": tClientCaseBase.CustomFields.TextValueByNameBasic(FieldName_client_gid)})
		if err != nil {
			return err
		}
		if tClient == nil {
			return errors.New("tClient is nil.")
		}
		syncBaseInfo, err := ZohoBaseInfoSync(tClientCaseBase, tClient, changeFieldName)
		if err != nil {
			return err
		}
		if syncBaseInfo != nil { // 需要同步

			lib.DPrintln("ZohoinfoSyncUsecase syncBaseInfo Kind_clients:"+kind, syncBaseInfo)

			err = c.SyncRow(Kind_clients,
				tClient.CustomFields.TextValueByNameBasic("gid"),
				syncBaseInfo)
			if err != nil {
				return err
			}
		}

		// 同步此人的其它client cases
		tClientCases, err := c.TUsecase.ListByCond(Kind_client_cases,
			And(Eq{"client_gid": tClient.CustomFields.TextValueByNameBasic("gid"),
				"biz_deleted_at": 0}, Neq{"gid": tClientCaseBase.CustomFields.TextValueByNameBasic("gid")}))
		if err != nil {
			return err
		}
		for _, tClientCase := range tClientCases {
			syncBaseInfo, err := ZohoBaseInfoSync(tClientCaseBase, tClientCase, changeFieldName)
			if err != nil {
				return err
			}
			if syncBaseInfo != nil { // 需要同步
				lib.DPrintln("ZohoinfoSyncUsecase syncBaseInfo Kind_client_cases:"+kind, syncBaseInfo)
				err = c.SyncRow(Kind_client_cases,
					tClientCase.CustomFields.TextValueByNameBasic("gid"),
					syncBaseInfo)
				if err != nil {
					return err
				}
			}
		}

		return nil
	} else {
		return errors.New("Sync: kind: " + kind + " does not support.")
	}
}

// SyncRow 此处不能处理stage字段，需要注意
func (c *ZohoinfoSyncUsecase) SyncRow(kind string, gid string, syncBaseInfo lib.TypeMap) error {
	if syncBaseInfo == nil || len(syncBaseInfo) == 0 {
		return errors.New("SyncRow: syncBaseInfo length is 0 or nil")
	}
	if gid == "" {
		return errors.New("SyncRow: gid is empty.")
	}
	if kind == Kind_client_cases {
		zohoData := make(lib.TypeMap)
		for k, v := range syncBaseInfo {
			str := config_zoho.ZohoDealFieldNameByVbcFieldName(k)
			if str == "" {
				return errors.New("SyncRow: ZohoDealFieldName is wrong")
			}
			zohoData.Set(str, v)
		}
		zohoData.Set("id", gid)
		lib.DPrintln("ZohoinfoSyncUsecase Kind_client_cases:", zohoData)
		_, _, err := c.ZohoUsecase.PutRecordV1(config_zoho.Deals, zohoData)
		if err != nil {
			return err
		}
		// 当同步zoho成功，本地直接修改，防止多次重启触发zoho api , 小概率自动同步冲突
		syncBaseInfo.Set(FileName_client_cases_gid, gid)
		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry(syncBaseInfo), FileName_client_cases_gid, nil)
		return err
	} else if kind == Kind_clients {

		zohoData := make(lib.TypeMap)
		for k, v := range syncBaseInfo {
			str := config_zoho.ZohoContactFieldNameByVbcFieldName(k)
			if str == "" {
				return errors.New("SyncRow: ZohoContactFieldName is wrong")
			}
			zohoData.Set(str, v)
		}
		zohoData.Set("id", gid)
		_, _, err := c.ZohoUsecase.PutRecordV1(config_zoho.Contacts, zohoData)
		if err != nil {
			return err
		}
		// 当同步zoho成功，本地直接修改，防止多次重启触发zoho api , 小概率自动同步冲突
		syncBaseInfo.Set(Client_FileName_gid, gid)
		_, err = c.DataEntryUsecase.HandleOne(Kind_clients, TypeDataEntry(syncBaseInfo), Client_FileName_gid, nil)
		return err

	} else {
		return errors.New("kind: " + kind + " does not support.")
	}
}

func ZohoBaseInfoSync(master *TData, slave *TData, changeFieldName string) (lib.TypeMap, error) {
	if master == nil || slave == nil {
		return nil, errors.New("master or slave is nil")
	}
	row := make(lib.TypeMap)

	if changeFieldName == FieldName_email && master.CustomFields.TextValueByNameBasic(FieldName_email) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_email) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_email) {
			row.Set(FieldName_email, master.CustomFields.TextValueByNameBasic(FieldName_email))
		}
	}
	if changeFieldName == FieldName_phone && master.CustomFields.TextValueByNameBasic(FieldName_phone) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_phone) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_phone) {
			row.Set(FieldName_phone, master.CustomFields.TextValueByNameBasic(FieldName_phone))
		}
	}
	if changeFieldName == FieldName_ssn && master.CustomFields.TextValueByNameBasic(FieldName_ssn) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_ssn) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_ssn) {
			row.Set(FieldName_ssn, master.CustomFields.TextValueByNameBasic(FieldName_ssn))
		}
	}
	if changeFieldName == FieldName_dob && master.CustomFields.TextValueByNameBasic(FieldName_dob) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_dob) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_dob) {
			row.Set(FieldName_dob, master.CustomFields.TextValueByNameBasic(FieldName_dob))
		}
	}
	if changeFieldName == FieldName_state && master.CustomFields.TextValueByNameBasic(FieldName_state) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_state) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_state) {
			row.Set(FieldName_state, master.CustomFields.TextValueByNameBasic(FieldName_state))
		}
	}
	if changeFieldName == FieldName_city && master.CustomFields.TextValueByNameBasic(FieldName_city) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_city) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_city) {
			row.Set(FieldName_city, master.CustomFields.TextValueByNameBasic(FieldName_city))
		}
	}
	if changeFieldName == FieldName_address && master.CustomFields.TextValueByNameBasic(FieldName_address) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_address) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_address) {
			row.Set(FieldName_address, master.CustomFields.TextValueByNameBasic(FieldName_address))
		}
	}
	if changeFieldName == FieldName_zip_code && master.CustomFields.TextValueByNameBasic(FieldName_zip_code) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_zip_code) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_zip_code) {
			row.Set(FieldName_zip_code, master.CustomFields.TextValueByNameBasic(FieldName_zip_code))
		}
	}

	if changeFieldName == FieldName_place_of_birth_city && master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_city) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_city) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_city) {
			row.Set(FieldName_place_of_birth_city, master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_city))
		}
	}

	if changeFieldName == FieldName_place_of_birth_country && master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_country) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_country) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_country) {
			row.Set(FieldName_place_of_birth_country, master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_country))
		}
	}

	if changeFieldName == FieldName_place_of_birth_state_province && master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_state_province) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_state_province) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_state_province) {
			row.Set(FieldName_place_of_birth_state_province, master.CustomFields.TextValueByNameBasic(FieldName_place_of_birth_state_province))
		}
	}
	if changeFieldName == FieldName_current_occupation && master.CustomFields.TextValueByNameBasic(FieldName_current_occupation) != "" {
		if master.CustomFields.TextValueByNameBasic(FieldName_current_occupation) !=
			slave.CustomFields.TextValueByNameBasic(FieldName_current_occupation) {
			row.Set(FieldName_current_occupation, master.CustomFields.TextValueByNameBasic(FieldName_current_occupation))
		}
	}

	if len(row) > 0 {
		return row, nil
	}
	return nil, nil
}
