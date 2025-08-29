package biz

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

type TaskHttpUsecase struct {
	log                    *log.Helper
	conf                   *conf.Data
	JWTUsecase             *JWTUsecase
	RecordbuzSearchUsecase *RecordbuzSearchUsecase
	KindUsecase            *KindUsecase
	ClientTaskUsecase      *ClientTaskUsecase
	TimezonesUsecase       *TimezonesUsecase
	FieldUsecase           *FieldUsecase
	TUsecase               *TUsecase
	DataEntryUsecase       *DataEntryUsecase
	QueueUsecase           *QueueUsecase
}

func NewTaskHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	RecordbuzSearchUsecase *RecordbuzSearchUsecase,
	KindUsecase *KindUsecase,
	ClientTaskUsecase *ClientTaskUsecase,
	TimezonesUsecase *TimezonesUsecase,
	FieldUsecase *FieldUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
	QueueUsecase *QueueUsecase) *TaskHttpUsecase {
	return &TaskHttpUsecase{
		log:                    log.NewHelper(logger),
		conf:                   conf,
		JWTUsecase:             JWTUsecase,
		RecordbuzSearchUsecase: RecordbuzSearchUsecase,
		KindUsecase:            KindUsecase,
		ClientTaskUsecase:      ClientTaskUsecase,
		TimezonesUsecase:       TimezonesUsecase,
		FieldUsecase:           FieldUsecase,
		TUsecase:               TUsecase,
		DataEntryUsecase:       DataEntryUsecase,
		QueueUsecase:           QueueUsecase,
	}
}

type TaskHttpListRequest struct {
	Type string `json:"type"`
}

const (
	TaskHttpListRequest_Type_Default = ""      // 首页获取有到期时间正在进行的任务
	TaskHttpListRequest_Type_Open    = "open"  // 首页获取正在进行的任务
	TaskHttpListRequest_Type_Close   = "close" // 首页获取关闭的任务
)

func (c *TaskHttpUsecase) Complete(ctx *gin.Context) {
	reply := CreateReply()
	// 通过路由获取的
	moduleName := ctx.Param("module_name")
	gid := ctx.Param("gid")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.BizComplete(userFacade, ModuleConvertToKind(moduleName), gid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *TaskHttpUsecase) BizComplete(userFacade UserFacade, kind string, gid string) (lib.TypeMap, error) {

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}

	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}
	searchCls := c.RecordbuzSearchUsecase.NewRecordbuzSearchCls(false)
	hasPermission, err := searchCls.HasPermissionRow(userFacade, *kindEntity, gid)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("Records do not exist or have no permission to access")
	}

	tTask, err := c.TUsecase.Data(Kind_client_tasks, Eq{DataEntry_biz_deleted_at: 0, DataEntry_gid: gid})
	if err != nil {
		return nil, err
	}
	if tTask == nil {
		return nil, errors.New("tTask is nil")
	}
	data := make(TypeDataEntry)
	data[DataEntry_gid] = gid
	data[TaskFieldName_status] = config_zoho.ClientTaskStatus_Completed
	data[TaskFieldName_closed_at] = time.Now().Unix()
	data[DataEntry_modified_by] = userFacade.Gid()
	_, err = c.DataEntryUsecase.HandleOne(Kind_client_tasks, data, DataEntry_gid, &userFacade.TData)
	if err != nil {
		return nil, err
	}
	whatGid := tTask.CustomFields.TextValueByNameBasic(TaskFieldName_what_id_gid)
	if whatGid != "" {
		c.QueueUsecase.PushClientTaskHandleWhatGidJobTasks(context.TODO(), []string{whatGid})
	}
	whoGid := tTask.CustomFields.TextValueByNameBasic(TaskFieldName_who_id_gid)
	if whoGid != "" {
		c.QueueUsecase.PushClientTaskHandleWhoGidJobTasks(context.TODO(), []string{whoGid})
	}

	//searchCls := c.RecordbuzSearchUsecase.NewRecordbuzSearchCls()
	//hasPermission, err := searchCls.HasPermissionRow(userFacade, *kindEntity, kindGid)
	//if err != nil {
	//	return nil, err
	//}
	//if !hasPermission {
	//	return nil, errors.New("Records do not exist or have no permission to access")
	//}
	return nil, nil
}

func (c *TaskHttpUsecase) List(ctx *gin.Context) {
	reply := CreateReply()
	// 通过路由获取的
	moduleName := ctx.Param("module_name")
	kindGid := ctx.Param("kind_gid")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	var taskHttpListRequest TaskHttpListRequest
	rawData, _ := ctx.GetRawData()
	json.Unmarshal(rawData, &taskHttpListRequest)

	data, err := c.BizList(userFacade, ModuleConvertToKind(moduleName), kindGid, taskHttpListRequest)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type TaskCard struct {
	Num   int            `json:"num"`
	Items []TaskCardItem `json:"items"`
}
type TaskCardItem struct {
	Gid        string                    `json:"gid"`
	Subject    string                    `json:"subject"`
	SysDueDate TFieldExtendForSysDueDate `json:"sys_due_date"`
	TaskOwner  *FabUser                  `json:"task_owner"`
	List       []TaskCardItemUnit        `json:"list"`
}

type TaskCardItemUnit struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (c *TaskHttpUsecase) BizList(userFacade UserFacade, kind string, kindGid string, taskHttpListRequest TaskHttpListRequest) (lib.TypeMap, error) {

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}
	searchCls := c.RecordbuzSearchUsecase.NewRecordbuzSearchCls(false)
	hasPermission, err := searchCls.HasPermissionRow(userFacade, *kindEntity, kindGid)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("Records do not exist or have no permission to access")
	}

	var tasks TDataList
	if taskHttpListRequest.Type == TaskHttpListRequest_Type_Default {
		tasks, err = c.ClientTaskUsecase.TasksForIndex(*kindEntity, kindGid)
		if err != nil {
			return nil, err
		}
	} else if taskHttpListRequest.Type == TaskHttpListRequest_Type_Open {
		tasks, err = c.ClientTaskUsecase.TasksForOpen(*kindEntity, kindGid)
		if err != nil {
			return nil, err
		}
	} else if taskHttpListRequest.Type == TaskHttpListRequest_Type_Close {
		tasks, err = c.ClientTaskUsecase.TasksForClose(*kindEntity, kindGid)
		if err != nil {
			return nil, err
		}
	}
	taskCard := &TaskCard{
		Num: len(tasks),
	}

	fieldStruct, err := c.FieldUsecase.CacheStructByKind(Kind_client_tasks)
	if err != nil {
		return nil, err
	}
	if fieldStruct == nil {
		return nil, errors.New("fieldStruct is nil")
	}

	// 格式化时间
	for k, _ := range tasks {
		err = tasks[k].HandleTDataTimezone(c.TimezonesUsecase, *fieldStruct, &userFacade)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range tasks {
		dueData, _ := GenTFieldExtendForSysDueDate(v.CustomFields.TextValueByNameBasic(TaskFieldName_due_date))
		item := TaskCardItem{
			Gid:     v.Gid(),
			Subject: v.CustomFields.TextValueByNameBasic(TaskFieldName_subject),
		}
		if dueData != nil {
			item.SysDueDate = *dueData
		}
		userGid := v.CustomFields.TextValueByNameBasic(TaskFieldName_user_gid)
		if userGid != "" {
			item.TaskOwner = &FabUser{
				Gid:      userGid,
				FullName: v.CustomFields.DisplayValueByNameBasic(TaskFieldName_user_gid),
			}
		}
		if taskHttpListRequest.Type == TaskHttpListRequest_Type_Open {
			unit := TaskCardItemUnit{
				Title:   "Status",
				Content: v.CustomFields.DisplayValueByNameBasic(TaskFieldName_status),
			}
			unit1 := TaskCardItemUnit{
				Title:   "Priority",
				Content: v.CustomFields.DisplayValueByNameBasic(TaskFieldName_priority),
			}
			item.List = append(item.List, unit)
			item.List = append(item.List, unit1)
		} else if taskHttpListRequest.Type == TaskHttpListRequest_Type_Close {
			unit := TaskCardItemUnit{
				Title:   "Closed Time",
				Content: v.CustomFields.DisplayValueByNameBasic(TaskFieldName_closed_at),
			}
			item.List = append(item.List, unit)
		}

		taskCard.Items = append(taskCard.Items, item)
	}
	data := make(lib.TypeMap)
	data.Set("tasks", taskCard)

	return data, nil
}
