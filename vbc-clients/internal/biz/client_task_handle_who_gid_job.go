package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	. "vbc/lib/builder"
	//. "vbc/lib/builder"
)

type ClientTaskHandleWhoGidJobParams struct {
	WhoGid string
}

type ClientTaskHandleWhoGidJobUsecase struct {
	CustomTaskBatch
	conf              *conf.Data
	CommonUsecase     *CommonUsecase
	log               *log.Helper
	DataEntryUsecase  *DataEntryUsecase
	DataComboUsecase  *DataComboUsecase
	LogUsecase        *LogUsecase
	TUsecase          *TUsecase
	ClientTaskUsecase *ClientTaskUsecase
}

func NewClientTaskHandleWhoGidJobUsecase(CommonUsecase *CommonUsecase,
	logger log.Logger,
	DataEntryUsecase *DataEntryUsecase,
	DataComboUsecase *DataComboUsecase,
	TUsecase *TUsecase,
	LogUsecase *LogUsecase,
	conf *conf.Data,
	ClientTaskUsecase *ClientTaskUsecase,
) *ClientTaskHandleWhoGidJobUsecase {

	clientTaskHandleWhoGidJobUsecase := &ClientTaskHandleWhoGidJobUsecase{
		CommonUsecase:     CommonUsecase,
		log:               log.NewHelper(logger),
		DataEntryUsecase:  DataEntryUsecase,
		DataComboUsecase:  DataComboUsecase,
		TUsecase:          TUsecase,
		LogUsecase:        LogUsecase,
		conf:              conf,
		ClientTaskUsecase: ClientTaskUsecase,
	}
	clientTaskHandleWhoGidJobUsecase.CustomTaskBatch.RedisQueue = Redis_client_task_handle_who_gid_queue
	clientTaskHandleWhoGidJobUsecase.CustomTaskBatch.RedisProcessing = Redis_client_task_handle_who_gid_processing
	clientTaskHandleWhoGidJobUsecase.CustomTaskBatch.RedisClient = CommonUsecase.RedisClient()
	clientTaskHandleWhoGidJobUsecase.CustomTaskBatch.Log = log.NewHelper(logger)
	clientTaskHandleWhoGidJobUsecase.CustomTaskBatch.Handle = clientTaskHandleWhoGidJobUsecase.HandleTask
	clientTaskHandleWhoGidJobUsecase.CustomTaskBatch.MaxBatchLimit = 1000
	clientTaskHandleWhoGidJobUsecase.CustomTaskBatch.WindowSeconds = 3

	return clientTaskHandleWhoGidJobUsecase
}

func (c *ClientTaskHandleWhoGidJobUsecase) HandleTask(ctx context.Context, customTaskParams []CustomTaskParams) error {
	err := c.BizHandleTask(ctx, customTaskParams)
	if err != nil {
		c.log.Error(err, ":", InterfaceToString(customTaskParams))
	}
	return err
}

// BizHandleTask tasks
func (c *ClientTaskHandleWhoGidJobUsecase) BizHandleTask(ctx context.Context, customTaskParams []CustomTaskParams) error {

	c.log.Debug("customTaskParams:", customTaskParams)
	var whoGids []string
	for _, v := range customTaskParams {
		var clientTaskHandleWhoGidJobParams ClientTaskHandleWhoGidJobParams
		json.Unmarshal([]byte(v.Params), &clientTaskHandleWhoGidJobParams)
		if clientTaskHandleWhoGidJobParams.WhoGid != "" {
			whoGids = append(whoGids, clientTaskHandleWhoGidJobParams.WhoGid)
		}
	}
	return c.Do(whoGids)
}

func (c *ClientTaskHandleWhoGidJobUsecase) Do(whoGids []string) error {

	c.log.Debug("whoGids:", whoGids)
	dueDatesResult, err := c.ClientTaskUsecase.DueDatesByWhoGids(whoGids)
	if err != nil {
		return err
	}

	res, err := c.TUsecase.ListByCond(Kind_clients, In(DataEntry_gid, whoGids))
	if err != nil {
		return err
	}
	for _, v := range res {
		if _, ok := dueDatesResult[v.Gid()]; ok {
			data := make(TypeDataEntry)
			data[DataEntry_gid] = v.Gid()
			data[DataEntry_sys__due_date] = dueDatesResult[v.Gid()].DueDate
			c.log.Debug(data)
			_, err := c.DataEntryUsecase.HandleOne(Kind_clients, data, DataEntry_gid, nil)
			if err != nil {
				c.log.Error(err, "data:", InterfaceToString(data))
			}
		} else {
			c.log.Error("error")
		}
	}

	return nil
}
