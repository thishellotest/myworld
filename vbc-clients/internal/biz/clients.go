package biz

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ClientEntity struct {
	ID       int32 `gorm:"primaryKey"`
	Uniqcode string
}

func (ClientEntity) TableName() string {
	return TableName_client_cases
}

type ClientUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[ClientEntity]
	AsanaUsecase          *AsanaUsecase
	TaskFailureLogUsecase *TaskFailureLogUsecase
	ZohoUsecase           *ZohoUsecase
	StageTransUsecase     *StageTransUsecase
	TUsecase              *TUsecase
	DataEntryUsecase      *DataEntryUsecase
}

func NewClientUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AsanaUsecase *AsanaUsecase,
	TaskFailureLogUsecase *TaskFailureLogUsecase,
	ZohoUsecase *ZohoUsecase,
	StageTransUsecase *StageTransUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase) *ClientUsecase {
	uc := &ClientUsecase{
		log:                   log.NewHelper(logger),
		CommonUsecase:         CommonUsecase,
		conf:                  conf,
		AsanaUsecase:          AsanaUsecase,
		TaskFailureLogUsecase: TaskFailureLogUsecase,
		ZohoUsecase:           ZohoUsecase,
		StageTransUsecase:     StageTransUsecase,
		TUsecase:              TUsecase,
		DataEntryUsecase:      DataEntryUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *ClientUsecase) BizChangeStagesToGettingStartedEmail(gid string) error {

	if configs.StoppedZoho {

		tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		dbStage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
		if dbStage != config_vbc.Stages_FeeScheduleandContract {
			return errors.New("stage is not Stages_FeeScheduleandContract.")
		}
		destMap := make(TypeDataEntry)
		destMap[DataEntry_gid] = gid
		destMap[FieldName_stages] = config_vbc.Stages_GettingStartedEmail
		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, destMap, DataEntry_gid, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ClientUsecase) BizChangeStagesForAm(gid string) error {

	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	dbStage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if dbStage != config_vbc.Stages_AmContractPending {
		return errors.New("stage is not  Stages_AmContractPending.")
	}
	destMap := make(TypeDataEntry)
	destMap[DataEntry_gid] = gid
	destMap[FieldName_stages] = config_vbc.Stages_AmAwaitingClientRecords
	_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, destMap, DataEntry_gid, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientUsecase) HandleChangeStagesForAm(gid string) error {

	err := c.BizChangeStagesForAm(gid)
	if err != nil {
		c.log.Error("BizChangeStagesToGettingStartedEmail:", InterfaceToString(map[string]interface{}{
			"gid": gid,
			"err": err.Error(),
		}))
		return c.TaskFailureLogUsecase.Add(TaskType_ChangeStagesToGettingStartedEmail, 0,
			map[string]interface{}{
				"gid": gid,
				"err": err.Error(),
			})
	}
	return nil
}

func (c *ClientUsecase) HandleChangeStagesToGettingStartedEmail(gid string) error {

	err := c.BizChangeStagesToGettingStartedEmail(gid)
	if err != nil {
		c.log.Error("BizChangeStagesToGettingStartedEmail:", InterfaceToString(map[string]interface{}{
			"gid": gid,
			"err": err.Error(),
		}))
		return c.TaskFailureLogUsecase.Add(TaskType_ChangeStagesToGettingStartedEmail, 0,
			map[string]interface{}{
				"gid": gid,
				"err": err.Error(),
			})
	}
	return nil
}

func (c *ClientUsecase) BizChangeStagesToMiniDBQsFinalized(gid string) error {

	if configs.StoppedZoho {
		tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil")
		}
		dbStage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
		if dbStage != config_vbc.Stages_MiniDBQ_Forms && dbStage != config_vbc.Stages_AmMiniDBQ_Forms {
			return errors.New("stage is not Stages_MiniDBQ_Forms.")
		}
		dataEntry := make(TypeDataEntry)
		dataEntry[DataEntry_gid] = gid

		if IsAmContract(*tCase) {
			dataEntry[FieldName_stages] = config_vbc.Stages_AmMedicalTeamFormsSigned
		} else {
			dataEntry[FieldName_stages] = config_vbc.Stages_MedicalTeamFormsSigned
		}
		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
		if err != nil {
			return err
		}
		return nil
	} else {

		deal, err := c.ZohoUsecase.GetDeal(gid)
		if err != nil {
			return err
		}
		if deal == nil {
			return errors.New("deal is nil")
		}
		zohoStage := deal.GetString("Stage")

		dbStage, err := c.StageTransUsecase.BizZohoStageToDBStage(zohoStage)
		if err != nil {
			c.log.Error(err)
			return err
		}
		if dbStage != config_vbc.Stages_MiniDBQ_Forms {
			return errors.New("stage is not Stages_MiniDBQ_Forms.")
		}

		zohoStage1, err := c.StageTransUsecase.DBStageToZohoStage(config_vbc.Stages_MedicalTeamFormsSigned)
		if err != nil {
			c.log.Error(err)
			return err
		}

		destMap := make(lib.TypeMap)
		destMap.Set("id", gid)
		destMap.Set("Stage", zohoStage1)
		_, _, err = c.ZohoUsecase.PutRecordV1(config_zoho.Deals, destMap)

		return err
	}
}

func (c *ClientUsecase) HandleChangeStagesToMiniDBQsFinalized(gid string) error {
	err := c.BizChangeStagesToMiniDBQsFinalized(gid)
	if err != nil {
		c.log.Error(err, " gid: ", gid)
		return c.TaskFailureLogUsecase.Add(TaskType_ChangeStagesToMiniDBQsFinalized, 0,
			map[string]interface{}{
				"gid": gid,
				"err": err.Error(),
			})
	}
	return nil
}

// GetByPhone +13109719619
func (c *ClientUsecase) GetByPhone(phone string) (*TData, error) {

	phone1, phone2, phone3, err := FormatPhoneNumber(phone)
	if err != nil {
		return nil, err
	}

	if phone1 != "" || phone2 != "" {
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
		return c.TUsecase.Data(Kind_clients, Or(conds...))

	}
	return nil, nil
}

type ClientPipelines map[string]string

func (c ClientPipelines) GetByClientGid(clientGid string) string {
	if v, ok := c[clientGid]; ok {
		return v
	}
	return Pipelines_default
}

func (c *ClientUsecase) GetOneClientPipeline(clientGid string) (string, error) {
	a, err := c.GetClientsPipelines([]string{clientGid})
	if err != nil {
		return "", err
	}
	return a.GetByClientGid(clientGid), nil
}

func (c *ClientUsecase) GetClientsPipelines(clientGids []string) (ClientPipelines, error) {
	data := make(ClientPipelines)
	if len(clientGids) == 0 {
		return data, nil
	}
	sql := fmt.Sprintf(`SELECT cc.*
FROM client_cases cc
INNER JOIN (
    SELECT client_gid, MIN(id) AS min_id
    FROM client_cases
    WHERE client_gid IN ('%s') and biz_deleted_at=0 and deleted_at=0
    GROUP BY client_gid
) t ON cc.client_gid = t.client_gid AND cc.id = t.min_id`, strings.Join(clientGids, "','"))
	res, err := c.TUsecase.ListByRawSql(Kind_client_cases, sql)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		contractSource := v.CustomFields.TextValueByNameBasic(FieldName_ContractSource)
		data[v.CustomFields.TextValueByNameBasic(FieldName_client_gid)] = ContractSourceToPipelines(contractSource)
	}
	return data, nil
}
