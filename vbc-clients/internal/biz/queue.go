package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

type QueueUsecase struct {
	log                          *log.Helper
	CommonUsecase                *CommonUsecase
	conf                         *conf.Data
	TUsecase                     *TUsecase
	EventBus                     *EventBus
	ClientTaskUsecase            *ClientTaskUsecase
	AutomaticTaskCreationUsecase *AutomaticTaskCreationUsecase
}

func NewQueueUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	EventBus *EventBus,
	ClientTaskUsecase *ClientTaskUsecase,
	AutomaticTaskCreationUsecase *AutomaticTaskCreationUsecase) *QueueUsecase {
	uc := &QueueUsecase{
		log:                          log.NewHelper(logger),
		CommonUsecase:                CommonUsecase,
		conf:                         conf,
		TUsecase:                     TUsecase,
		EventBus:                     EventBus,
		ClientTaskUsecase:            ClientTaskUsecase,
		AutomaticTaskCreationUsecase: AutomaticTaskCreationUsecase,
	}
	uc.EventBus.Subscribe(EventBus_AfterInsertData, uc.HandleAfterInsertData)
	uc.EventBus.Subscribe(EventBus_AfterHandleUpdate, uc.HandleAfterHandleUpdate)
	return uc
}

func (c *QueueUsecase) HandleAfterInsertData(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList, modifiedBy string) {
	if !configs.Enable_Client_Task_ForCRM {
		return
	}

	c.AutomaticTaskCreationUsecase.HandleAutomaticTaskCreation(kindEntity, structField, recognizeFieldName, dataEntryOperResult, sourceData)

	if kindEntity.Kind == Kind_client_tasks {
		if recognizeFieldName == DataEntry_gid {
			c.log.Info("AfterInsertData:1")
			var whatGids []string
			var whoGids []string
			var needHandleCloseTimeTaskGids []string
			for taskGid, v := range dataEntryOperResult {
				needHandleCloseTimeTaskGids = append(needHandleCloseTimeTaskGids, taskGid)
				if v.IsNewRecord {
					row := lib.TypeMap(sourceData.Get(DataEntry_gid, taskGid))
					if row != nil {
						whatGid := row.GetString(TaskFieldName_what_id_gid)
						whoGid := row.GetString(TaskFieldName_who_id_gid)
						if whatGid != "" {
							whatGids = append(whatGids, whatGid)
						}
						if whoGid != "" {
							whoGids = append(whoGids, whoGid)
						}
					}
				}
			}
			c.log.Info("AfterInsertData:2 whatGids", whatGids)
			c.log.Info("AfterInsertData:2 whoGids", whoGids)
			err := c.PushClientTaskHandleWhatGidJobTasks(context.TODO(), whatGids)
			if err != nil {
				c.log.Error(err, "whatGids:", whatGids)
			}
			err = c.PushClientTaskHandleWhoGidJobTasks(context.TODO(), whoGids)
			if err != nil {
				c.log.Error(err, "whoGids:", whatGids)
			}
			err = c.ClientTaskUsecase.HandleCloseTime(needHandleCloseTimeTaskGids)
			if err != nil {
				c.log.Error(err, "HandleCloseTime needHandleCloseTimeTaskGids:", needHandleCloseTimeTaskGids)
			}
		}
	}

}

func (c *QueueUsecase) HandleAfterHandleUpdate(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList, modifiedBy string) {

	if !configs.Enable_Client_Task_ForCRM {
		return
	}

	c.AutomaticTaskCreationUsecase.HandleAutomaticTaskCreation(kindEntity, structField, recognizeFieldName, dataEntryOperResult, sourceData)

	//c.log.Info("HandleAfterHandleUpdate Kind:", kindEntity.Kind)
	if structField == nil {
		c.log.Error("structField is nil")
		return
	}
	if kindEntity.Kind == Kind_client_tasks {
		if recognizeFieldName == DataEntry_gid && len(dataEntryOperResult) > 0 {
			c.log.Info("AfterHandleUpdate:1 dataEntryOperResult:", InterfaceToString(dataEntryOperResult))
			var whatGids []string
			var whoGids []string
			var needHandleCloseTimeTaskGids []string
			for gid, v := range dataEntryOperResult {
				if v.IsUpdated {
					for k1, v1 := range v.DataEntryModifyDataMap {
						if k1 == TaskFieldName_what_id_gid {
							newVal := v1.GetNewVal(FieldType_text)
							oldNew := v1.GetOldVal(FieldType_text)
							if newVal != "" {
								whatGids = append(whatGids, newVal)
							}
							if oldNew != "" {
								whatGids = append(whatGids, oldNew)
							}
						} else if k1 == TaskFieldName_who_id_gid {
							newVal := v1.GetNewVal(FieldType_text)
							oldNew := v1.GetOldVal(FieldType_text)
							if newVal != "" {
								whoGids = append(whoGids, newVal)
							}
							if oldNew != "" {
								whoGids = append(whoGids, oldNew)
							}
						} else if k1 == TaskFieldName_due_date {
							err := c.TaskDueDateChange([]string{gid})
							if err != nil {
								c.log.Error(err)
							}
						} else if k1 == TaskFieldName_status {
							newVal := v1.GetNewVal(FieldType_dropdown)
							if newVal == config_zoho.ClientTaskStatus_Completed {
								needHandleCloseTimeTaskGids = append(needHandleCloseTimeTaskGids, gid)
							}
							err := c.TaskDueDateChange([]string{gid})
							if err != nil {
								c.log.Error(err)
							}
						}
					}
				}
			}
			c.log.Info("AfterHandleUpdate:2 whatGids", whatGids)
			c.log.Info("AfterHandleUpdate:2 whoGids", whoGids)
			err := c.PushClientTaskHandleWhatGidJobTasks(context.TODO(), whatGids)
			if err != nil {
				c.log.Error(err, "AfterHandleUpdate whatGids:", whatGids)
			}
			err = c.PushClientTaskHandleWhoGidJobTasks(context.TODO(), whoGids)
			if err != nil {
				c.log.Error(err, "AfterHandleUpdate whoGids:", whoGids)
			}
			err = c.ClientTaskUsecase.HandleCloseTime(needHandleCloseTimeTaskGids)
			if err != nil {
				c.log.Error(err, "AfterHandleUpdate HandleCloseTime needHandleCloseTimeTaskGids:", needHandleCloseTimeTaskGids)
			}
		}
	}
}

func (c *QueueUsecase) TaskDueDateChange(gids []string) error {
	c.log.Info("TaskDueDateChange gids:", gids)
	if len(gids) == 0 {
		return nil
	}
	res, err := c.TUsecase.ListByCond(Kind_client_tasks, In(DataEntry_gid, gids))
	if err != nil {
		return err
	}
	var whatGids []string
	var whoGids []string
	for _, v := range res {
		whatGid := v.CustomFields.TextValueByNameBasic(TaskFieldName_what_id_gid)
		if whatGid != "" {
			whatGids = append(whatGids, whatGid)
		}
		whoGid := v.CustomFields.TextValueByNameBasic(TaskFieldName_who_id_gid)
		if whoGid != "" {
			whoGids = append(whoGids, whoGid)
		}
	}
	err = c.PushClientTaskHandleWhatGidJobTasks(context.TODO(), whatGids)
	if err != nil {
		c.log.Error(err)
	}
	err = c.PushClientTaskHandleWhoGidJobTasks(context.TODO(), whoGids)
	if err != nil {
		c.log.Error(err)
	}
	c.log.Info("TaskDueDateChange whatGids:", whatGids)
	c.log.Info("TaskDueDateChange whatGids:", whoGids)

	return nil
}

func (c *QueueUsecase) PushClientTaskHandleWhatGidJobTasks(ctx context.Context, whatGids []string) error {
	if len(whatGids) == 0 {
		return nil
	}
	var tasks []CustomTaskParams
	for _, v := range whatGids {
		tasks = append(tasks, CustomTaskParams{
			UniqueKey: v,
			Params: InterfaceToString(ClientTaskHandleWhatGidJobParams{
				WhatGid: v,
			}),
		})
	}

	var dest []string
	for _, v := range tasks {
		bytes, _ := v.MarshalBinary()
		dest = append(dest, string(bytes))
	}
	return c.CommonUsecase.RedisClient().LPush(ctx, Redis_client_task_handle_what_gid_queue, dest).Err()
}

func (c *QueueUsecase) PushClientTaskHandleWhoGidJobTasks(ctx context.Context, whoGids []string) error {
	if len(whoGids) == 0 {
		return nil
	}
	var tasks []CustomTaskParams
	for _, v := range whoGids {
		tasks = append(tasks, CustomTaskParams{
			UniqueKey: v,
			Params: InterfaceToString(ClientTaskHandleWhoGidJobParams{
				WhoGid: v,
			}),
		})
	}

	var dest []string
	for _, v := range tasks {
		bytes, _ := v.MarshalBinary()
		dest = append(dest, string(bytes))
	}
	return c.CommonUsecase.RedisClient().LPush(ctx, Redis_client_task_handle_who_gid_queue, dest).Err()
}

func (c *QueueUsecase) PushClientNameChangeJobTasks(ctx context.Context, caseGids []string) error {
	if len(caseGids) == 0 {
		return nil
	}
	var tasks []CustomTaskParams
	for _, v := range caseGids {
		tasks = append(tasks, CustomTaskParams{
			UniqueKey: v,
			Params: InterfaceToString(ClientNameChangeJobParams{
				CaseGid: v,
			}),
		})
	}

	var dest []string
	for _, v := range tasks {
		bytes, _ := v.MarshalBinary()
		dest = append(dest, string(bytes))
	}
	return c.CommonUsecase.RedisClient().LPush(ctx, Redis_client_name_change_job_queue, dest).Err()
}
