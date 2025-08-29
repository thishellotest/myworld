package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/internal/utils"
	. "vbc/lib/builder"
)

/*
历史需要发送的任务补发
*/

type ReissueTriggerStrRequestPendingUsecase struct {
	log               *log.Helper
	conf              *conf.Data
	CommonUsecase     *CommonUsecase
	TUsecase          *TUsecase
	TaskUsecase       *TaskUsecase
	TaskCreateUsecase *TaskCreateUsecase
}

func NewReissueTriggerStrRequestPendingUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	TaskUsecase *TaskUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
) *ReissueTriggerStrRequestPendingUsecase {
	uc := &ReissueTriggerStrRequestPendingUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		TUsecase:          TUsecase,
		TaskUsecase:       TaskUsecase,
		TaskCreateUsecase: TaskCreateUsecase,
	}

	return uc
}

/*
select id  from client_cases where stages='STR Request Pending' and biz_deleted_at=0 and deleted_at=0
and id not in (

select incr_id from tasks where   incr_id in (

	select id
	from client_cases
	where stages='STR Request Pending' and biz_deleted_at=0 and deleted_at=0

)
and task_input in ('{"HandleSendSMSType":"TextSTRRequestPendingLongerThan30Days"}', '{"HandleSendSMSType":"TextSTRRequestPending45Days"}')
and task_status=0
group by incr_id)
*/

func (c *ReissueTriggerStrRequestPendingUsecase) NeedHandleCases() ([]int32, error) {

	return []int32{42, 66, 73, 261, 349, 5005, 5078, 5084, 5102, 5106, 5113, 5126, 5130, 5142, 5159, 5172, 5179, 5181, 5210, 5221, 5222, 5225, 5243, 5249, 5252, 5260, 5274, 5275, 5294, 5297, 5347}, nil
}

// ExistTask 判断任务是否存在， 存在就不要重复创建了
func (c *ReissueTriggerStrRequestPendingUsecase) ExistTask(caseId int32) (bool, error) {
	e, err := c.TaskUsecase.GetByCond(And(Eq{"incr_id": caseId, "task_status": 0}, In("task_input",
		`{"HandleSendSMSType":"TextSTRRequestPendingLongerThan30Days"}`,
		`{"HandleSendSMSType":"TextSTRRequestPending45Days"}`)))
	if err != nil {
		return false, err
	}
	if e != nil {
		return true, nil
	}
	return false, nil

}

func (c *ReissueTriggerStrRequestPendingUsecase) LastTask(caseId int32) (*TaskEntity, error) {
	return c.TaskUsecase.GetByCondWithOrderBy(And(Eq{"incr_id": caseId}, In("task_input",
		`{"HandleSendSMSType":"TextSTRRequestPendingLongerThan30Days"}`,
		`{"HandleSendSMSType":"TextSTRRequestPending45Days"}`), In("task_status", 1)),
		"id desc")
}

func (c *ReissueTriggerStrRequestPendingUsecase) Handle() error {
	caseIds, err := c.NeedHandleCases()
	if err != nil {
		return err
	}
	for _, caseId := range caseIds {
		time.Sleep(1)
		err := c.DoOne(caseId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ReissueTriggerStrRequestPendingUsecase) DoOne(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil : " + InterfaceToString(caseId))
	}
	exist, err := c.ExistTask(caseId)
	if err != nil {
		return err
	}
	if !exist {
		task, err := c.LastTask(caseId)
		if err != nil {
			return err
		}
		var calTime time.Time
		if task != nil {
			calTime = time.Unix(task.NextAt, 0)
		} else {
			calTime = time.Now()
		}
		timeLocation := GetCaseTimeLocation(tCase, c.log)
		nextTime, err := utils.CalIntervalDayTime(calTime, 30, "08:00", timeLocation)
		if err != nil {
			return err
		}
		if nextTime.Before(time.Now()) {
			nextTime = time.Now()
		}
		err = c.TaskCreateUsecase.CreateTaskWithFrom(caseId,
			CronTriggerVo{
				HandleSendSMSType: HandleSendSMSTextSTRRequestPending45Days,
			}, Task_Dag_CronTrigger, nextTime.Unix(),
			Task_FromType_DialpadSMS, InterfaceToString(caseId),
		)
		if err != nil {
			c.log.Error(err)
			return err
		}
	}
	return nil
}
