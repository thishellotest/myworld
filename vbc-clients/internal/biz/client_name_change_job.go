package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
	//. "vbc/lib/builder"
)

type ClientNameChangeJobParams struct {
	CaseGid string
}

type ClientNameChangeJobUsecase struct {
	CustomTaskBatch
	conf              *conf.Data
	CommonUsecase     *CommonUsecase
	log               *log.Helper
	DataEntryUsecase  *DataEntryUsecase
	DataComboUsecase  *DataComboUsecase
	LogUsecase        *LogUsecase
	TUsecase          *TUsecase
	ClientTaskUsecase *ClientTaskUsecase
	MapUsecase        *MapUsecase
	BoxUsecase        *BoxUsecase
}

func NewClientNameChangeJobUsecase(CommonUsecase *CommonUsecase,
	logger log.Logger,
	DataEntryUsecase *DataEntryUsecase,
	DataComboUsecase *DataComboUsecase,
	TUsecase *TUsecase,
	LogUsecase *LogUsecase,
	conf *conf.Data,
	ClientTaskUsecase *ClientTaskUsecase,
	MapUsecase *MapUsecase,
	BoxUsecase *BoxUsecase,
) *ClientNameChangeJobUsecase {

	clientNameChangeJobUsecase := &ClientNameChangeJobUsecase{
		CommonUsecase:     CommonUsecase,
		log:               log.NewHelper(logger),
		DataEntryUsecase:  DataEntryUsecase,
		DataComboUsecase:  DataComboUsecase,
		TUsecase:          TUsecase,
		LogUsecase:        LogUsecase,
		conf:              conf,
		ClientTaskUsecase: ClientTaskUsecase,
		MapUsecase:        MapUsecase,
		BoxUsecase:        BoxUsecase,
	}
	clientNameChangeJobUsecase.CustomTaskBatch.RedisQueue = Redis_client_name_change_job_queue
	clientNameChangeJobUsecase.CustomTaskBatch.RedisProcessing = Redis_client_name_change_job_processing
	clientNameChangeJobUsecase.CustomTaskBatch.RedisClient = CommonUsecase.RedisClient()
	clientNameChangeJobUsecase.CustomTaskBatch.Log = log.NewHelper(logger)
	clientNameChangeJobUsecase.CustomTaskBatch.Handle = clientNameChangeJobUsecase.HandleTask
	clientNameChangeJobUsecase.CustomTaskBatch.MaxBatchLimit = 1000
	clientNameChangeJobUsecase.CustomTaskBatch.WindowSeconds = 30

	return clientNameChangeJobUsecase
}

func (c *ClientNameChangeJobUsecase) HandleTask(ctx context.Context, customTaskParams []CustomTaskParams) error {
	err := c.BizHandleTask(ctx, customTaskParams)
	if err != nil {
		c.log.Error(err, ":", InterfaceToString(customTaskParams))
	}
	return err
}

// BizHandleTask tasks
func (c *ClientNameChangeJobUsecase) BizHandleTask(ctx context.Context, customTaskParams []CustomTaskParams) error {

	c.log.Debug("customTaskParams:", customTaskParams)
	var caseGids []string
	for _, v := range customTaskParams {
		var clientNameChangeJobParams ClientNameChangeJobParams
		json.Unmarshal([]byte(v.Params), &clientNameChangeJobParams)
		if clientNameChangeJobParams.CaseGid != "" {
			caseGids = append(caseGids, clientNameChangeJobParams.CaseGid)
		}
	}
	return c.Do(caseGids)
}

func (c *ClientNameChangeJobUsecase) Do(caseGids []string) error {

	c.log.Debug("ClientNameChangeJobUsecase caseGids:", caseGids)
	res, err := c.TUsecase.ListByCond(Kind_client_cases, In(DataEntry_gid, caseGids))
	if err != nil {
		return err
	}
	for k, v := range res {
		tClient, _ := c.TUsecase.DataByGid(Kind_clients, v.CustomFields.TextValueByNameBasic(FieldName_client_gid))
		err = c.HandleBoxClientFolder(tClient, res[k])
		if err != nil {
			c.log.Error(err, " # ", v.Id())
		}
		err = c.HandleDCBoxClientFolder(tClient, res[k])
		if err != nil {
			c.log.Error(err, " # ", v.Id())
		}
	}
	return nil
}

func (c *ClientNameChangeJobUsecase) HandleDCBoxClientFolder(tClient *TData, tCase *TData) error {
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	dcFolderId, err := c.MapUsecase.GetForString(MapKeyDataCollectionFolderId(tCase.Id()))
	if err != nil {
		return err
	}
	if dcFolderId != "" {
		boxFolderRes, _, err := c.BoxUsecase.GetFolderInfo(dcFolderId)
		if err != nil {
			return err
		}
		if boxFolderRes == nil {
			c.log.Error("boxFolderRes is nil")
		} else {
			boxFolderResMap := lib.ToTypeMapByString(*boxFolderRes)
			// VBC - TestL, Test1 #5511
			dcFolderName := boxFolderResMap.GetString("name")
			newDcFolderName := ClientCaseDataCollectionFolderNameForBox(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
				tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

			if newDcFolderName != dcFolderName {
				_, err = c.BoxUsecase.UpdateFolderName(dcFolderId, newDcFolderName)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *ClientNameChangeJobUsecase) HandleBoxClientFolder(tClient *TData, tCase *TData) error {
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	clientBoxFolderId, err := c.MapUsecase.GetForString(MapKeyClientBoxFolderId(tCase.Id()))
	if err != nil {
		return err
	}
	if clientBoxFolderId != "" {
		boxFolderRes, _, err := c.BoxUsecase.GetFolderInfo(clientBoxFolderId)
		if err != nil {
			return err
		}
		if boxFolderRes == nil {
			c.log.Error("boxFolderRes is nil")
		} else {
			boxFolderResMap := lib.ToTypeMapByString(*boxFolderRes)
			// VBC - TestL, Test1 #5511
			clientFolderName := boxFolderResMap.GetString("name")
			newClientFolderName := ClientFolderNameForBox(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
				tClient.CustomFields.TextValueByNameBasic(FieldName_last_name))
			if strings.Index(clientFolderName, "#") > 0 {
				newClientFolderName = fmt.Sprintf("%s #%d", newClientFolderName, tCase.Id())
			}
			if newClientFolderName != clientFolderName {
				_, err = c.BoxUsecase.UpdateFolderName(clientBoxFolderId, newClientFolderName)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
