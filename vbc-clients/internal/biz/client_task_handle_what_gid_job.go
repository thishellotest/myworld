package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	. "vbc/lib/builder"
	//. "vbc/lib/builder"
)

type ClientTaskHandleWhatGidJobParams struct {
	WhatGid string
}

type ClientTaskHandleWhatGidJobUsecase struct {
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

func NewClientTaskHandleWhatGidJobUsecase(CommonUsecase *CommonUsecase,
	logger log.Logger,
	DataEntryUsecase *DataEntryUsecase,
	DataComboUsecase *DataComboUsecase,
	TUsecase *TUsecase,
	LogUsecase *LogUsecase,
	conf *conf.Data,
	ClientTaskUsecase *ClientTaskUsecase,
) *ClientTaskHandleWhatGidJobUsecase {

	clientTaskHandleWhatGidJobUsecase := &ClientTaskHandleWhatGidJobUsecase{
		CommonUsecase:     CommonUsecase,
		log:               log.NewHelper(logger),
		DataEntryUsecase:  DataEntryUsecase,
		DataComboUsecase:  DataComboUsecase,
		TUsecase:          TUsecase,
		LogUsecase:        LogUsecase,
		conf:              conf,
		ClientTaskUsecase: ClientTaskUsecase,
	}
	clientTaskHandleWhatGidJobUsecase.CustomTaskBatch.RedisQueue = Redis_client_task_handle_what_gid_queue
	clientTaskHandleWhatGidJobUsecase.CustomTaskBatch.RedisProcessing = Redis_client_task_handle_what_gid_processing
	clientTaskHandleWhatGidJobUsecase.CustomTaskBatch.RedisClient = CommonUsecase.RedisClient()
	clientTaskHandleWhatGidJobUsecase.CustomTaskBatch.Log = log.NewHelper(logger)
	clientTaskHandleWhatGidJobUsecase.CustomTaskBatch.Handle = clientTaskHandleWhatGidJobUsecase.HandleTask
	clientTaskHandleWhatGidJobUsecase.CustomTaskBatch.MaxBatchLimit = 1000
	clientTaskHandleWhatGidJobUsecase.CustomTaskBatch.WindowSeconds = 3

	return clientTaskHandleWhatGidJobUsecase
}

func (c *ClientTaskHandleWhatGidJobUsecase) HandleTask(ctx context.Context, customTaskParams []CustomTaskParams) error {
	err := c.BizHandleTask(ctx, customTaskParams)
	if err != nil {
		c.log.Error(err, ":", InterfaceToString(customTaskParams))
	}
	return err
}

// BizHandleTask tasks
func (c *ClientTaskHandleWhatGidJobUsecase) BizHandleTask(ctx context.Context, customTaskParams []CustomTaskParams) error {

	c.log.Debug("customTaskParams:", customTaskParams)
	var whatGids []string
	for _, v := range customTaskParams {
		var clientTaskHandleWhatGidJobParams ClientTaskHandleWhatGidJobParams
		json.Unmarshal([]byte(v.Params), &clientTaskHandleWhatGidJobParams)
		if clientTaskHandleWhatGidJobParams.WhatGid != "" {
			whatGids = append(whatGids, clientTaskHandleWhatGidJobParams.WhatGid)
		}
	}
	return c.Do(whatGids)
}

func (c *ClientTaskHandleWhatGidJobUsecase) Do(whatGids []string) error {

	c.log.Debug("whatGids:", whatGids)
	dueDatesResult, err := c.ClientTaskUsecase.DueDatesByWhatGids(whatGids)
	if err != nil {
		return err
	}

	res, err := c.TUsecase.ListByCond(Kind_client_cases, In(DataEntry_gid, whatGids))
	if err != nil {
		return err
	}
	for _, v := range res {
		if _, ok := dueDatesResult[v.Gid()]; ok {
			data := make(TypeDataEntry)
			data[DataEntry_gid] = v.Gid()
			data[DataEntry_sys__due_date] = dueDatesResult[v.Gid()].DueDate
			c.log.Debug(data)
			_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
			if err != nil {
				c.log.Error(err, "data:", InterfaceToString(data))
			}
		} else {
			c.log.Error("error")
		}
	}

	return nil
}
