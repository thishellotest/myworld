package biz

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"reflect"
	"strings"
	"sync"
	"time"
	"vbc/configs"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

const Task_TaskStatus_processing = 0
const Task_TaskStatus_finish = 1
const Task_TaskStatus_failure = 2
const Task_TaskStatus_cancel = 3 // 主动取消

const Task_max_retry = 3

const MailGenre_StartYourVADisabilityClaimRepresentation = "StartYourVADisabilityClaimRepresentation"

const MailGenre_InitialContact = "InitialContact"
const MailGenre_FeeScheduleCommunication = "FeeScheduleCommunication"
const MailGenre_GettingStartedEmail = "GettingStartedEmail"
const MailGenre_AmGettingStartedEmail = "AmGettingStartedEmail"

const MailGenre_CongratulationsNewRating = "CongratulationsNewRating"
const MailGenre_AmCongratulationsNewRating = "AmCongratulationsNewRating"

const MailGenre_Custom = "Custom"                             // 使用系统发件箱：提醒用户
const MailGenre_NotifyVsRemindClient = "NotifyVsRemindClient" // 通过VS处理提醒业务
const MailGenre_SignFeeContractFirstRemind = "SignFeeContractFirstRemind"
const MailGenre_MedicalExamDocumentsReminder = "MedicalExamDocumentsReminder"
const MailGenre_AmIntakeFormReminder = "AmIntakeFormReminder"
const MailGenre_ContractReminder = "ContractReminder"
const MailGenre_AmContractReminder = "AmContractReminder"
const MailGenre_UpcomingContactInformation = "UpcomingContactInformation"
const MailGenre_MiniDBQsDrafts = "MiniDBQsDrafts"
const MailGenre_YourRecordsReviewProcessHasBegun = "YourRecordsReviewProcessHasBegun"
const MailGenre_PleaseScheduleYourDoctorAppointments = "PleaseScheduleYourDoctorAppointments"
const MailGenre_PersonalStatementsReadyforYourReview = "PersonalStatementsReadyforYourReview"
const MailGenre_PleaseReviewYourPersonalStatementsinSharedFolder = "PleaseReviewYourPersonalStatementsinSharedFolder"
const MailGenre_HelpUsImproveSurvey = "HelpUsImproveSurvey"
const MailGenre_VAForm2122aSubmission = "VAForm2122aSubmission"
const MailGenre_ITFDeadlineIn90Days = "ITFDeadlineIn90Days"

// NeedUseSystemEmailConfig 是否要使用系统发件箱：Yannan， 其它的都应该使用Lead VS或Client Case Owner, 没有时，通知all collaborators
func NeedUseSystemEmailConfig(mailGenre string) bool {
	if mailGenre == MailGenre_Custom || mailGenre == MailGenre_NotifyVsRemindClient {
		return true
	}
	return false
}

const Task_FromType_AutomationCrontabEmail = "AutomationCrontabEmail"
const Task_FromType_DialpadSMS = "DialpadSMS"

type TaskInputMailVo struct {
	Genre string
}

type TaskEntity struct {
	ID          int32 `gorm:"primaryKey"`
	IncrId      int32
	TaskStatus  int
	TaskInput   string
	NextAt      int64
	NextRetryAt int64
	RetryCount  int32
	Event       string
	FromType    string
	FromId      string
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64
}

func (TaskEntity) TableName() string {
	return "tasks"
}

type TaskUsecase struct {
	log                     *log.Helper
	CommonUsecase           *CommonUsecase
	TUsecase                *TUsecase
	MailUsecase             *MailUsecase
	DocuSignUsecase         *DocuSignUsecase
	ClientEnvelopeUsecase   *ClientEnvelopeUsecase
	Dag                     *Dag
	RegisterDags            map[string]interface{}
	TaskCreateUsecase       *TaskCreateUsecase
	BehaviorUsecase         *BehaviorUsecase
	RemindUsecase           *RemindUsecase
	ContractReminderUsecase *ContractReminderUsecase
	DialpadUsecase          *DialpadUsecase
	DialpadbuzUsecase       *DialpadbuzUsecase
	DBUsecase[TaskEntity]
	MiscUsecase        *MiscUsecase
	SendVa2122aUsecase *SendVa2122aUsecase
}

//
//func registerType(elem interface{}) {
//	t := reflect.TypeOf(elem).Elem()
//	fmt.Println(t.Name(), "====")
//	typeRegistry[t.Name()] = t
//}

func RegisterDagName(dag interface{}) string {
	return reflect.TypeOf(dag).Elem().Name()
}

func NewTaskUsecase(logger log.Logger, CommonUsecase *CommonUsecase, TUsecase *TUsecase, MailUsecase *MailUsecase,
	DocuSignUsecase *DocuSignUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	Dag *Dag,
	TaskCreateUsecase *TaskCreateUsecase,
	BehaviorUsecase *BehaviorUsecase,
	RemindUsecase *RemindUsecase,
	ContractReminderUsecase *ContractReminderUsecase,
	DialpadUsecase *DialpadUsecase,
	DialpadbuzUsecase *DialpadbuzUsecase,
	MiscUsecase *MiscUsecase,
	SendVa2122aUsecase *SendVa2122aUsecase) *TaskUsecase {

	// 注入
	RegisterDags := make(map[string]interface{})
	RegisterDags[RegisterDagName(Dag)] = Dag

	uc := &TaskUsecase{
		log:                     log.NewHelper(logger),
		CommonUsecase:           CommonUsecase,
		TUsecase:                TUsecase,
		MailUsecase:             MailUsecase,
		DocuSignUsecase:         DocuSignUsecase,
		ClientEnvelopeUsecase:   ClientEnvelopeUsecase,
		RegisterDags:            RegisterDags,
		TaskCreateUsecase:       TaskCreateUsecase,
		BehaviorUsecase:         BehaviorUsecase,
		RemindUsecase:           RemindUsecase,
		ContractReminderUsecase: ContractReminderUsecase,
		DialpadUsecase:          DialpadUsecase,
		DialpadbuzUsecase:       DialpadbuzUsecase,
		MiscUsecase:             MiscUsecase,
		SendVa2122aUsecase:      SendVa2122aUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *TaskUsecase) Invoke(task *TaskEntity) error {

	str := task.Event
	dagInfo := strings.Split(str, ".")
	if len(dagInfo) != 2 {
		return errors.New("Dag's name is incorrect.")
	}
	if _, ok := c.RegisterDags[dagInfo[0]]; !ok {
		return errors.New(dagInfo[0] + " doesn't exists.")
	}
	method := reflect.ValueOf(c.RegisterDags[dagInfo[0]]).MethodByName(dagInfo[1])
	args := []reflect.Value{reflect.ValueOf(task)}
	ret := method.Call(args)

	if ret[0].Interface() == nil {
		err := c.HandleDependOn(task)
		if err != nil {
			return err
		}
		return nil
	}
	return ret[0].Interface().(error)
}

func (c *TaskUsecase) HandleDependOn(task *TaskEntity) error {
	if task == nil {
		return errors.New("HandleDependOn task is nil.")
	}
	if task.Event == Task_Dag_CreateEnvelopeAndSent {
		err := c.BehaviorUsecase.Add(task.IncrId, BehaviorType_fee_schedule_contract, time.Now(), "")
		if err != nil {
			return err
		}
		if configs.EnabledContractReminder {

			// 取消此客户，因为修改邮箱地址的历史提提醒
			er := c.CommonUsecase.DB().Model(&TaskEntity{}).
				Where("incr_id = ? and event =? and task_status=?",
					task.IncrId,
					"Dag.HandleContractReminder",
					Task_TaskStatus_processing).
				Updates(map[string]interface{}{
					"task_status": Task_TaskStatus_cancel,
					"updated_at":  time.Now().Unix()}).Error
			if er != nil {
				c.log.Error(er, " IncrId: ", task.IncrId, " ID: ", task.ID)
				return er
			}
			return c.ContractReminderUsecase.FirstReminder(task.IncrId)
		} else {
			return nil
		}

	} else if task.Event == Task_Dag_CreateEnvelopeAndSentFromBoxAm {

		err := c.BehaviorUsecase.Add(task.IncrId, BehaviorType_sent_am_contract, time.Now(), "")
		if err != nil {
			return err
		}
		if configs.EnabledContractReminder {
			// 取消此客户，因为修改邮箱地址的历史提提醒
			er := c.CommonUsecase.DB().Model(&TaskEntity{}).
				Where("incr_id = ? and event =? and task_status=?",
					task.IncrId,
					"Dag.HandleContractReminder",
					Task_TaskStatus_processing).
				Updates(map[string]interface{}{
					"task_status": Task_TaskStatus_cancel,
					"updated_at":  time.Now().Unix()}).Error
			if er != nil {
				c.log.Error(er, " IncrId: ", task.IncrId, " ID: ", task.ID)
				return er
			}
			return c.ContractReminderUsecase.FirstReminder(task.IncrId)
		} else {
			return nil
		}

	} else if task.Event == Task_Dag_BuzEmail {
		taskInput := lib.ToTypeMapByString(task.TaskInput)
		tpl := InterfaceToString(taskInput.Get("Genre"))
		if tpl == MailGenre_GettingStartedEmail {
			err := c.BehaviorUsecase.Add(task.IncrId, BehaviorType_sent_intake_form, time.Now(), "")
			if err != nil {
				return err
			}
			return c.RemindUsecase.CreateUnfinishedIntakeForm(task.IncrId)
		} else if tpl == MailGenre_AmGettingStartedEmail {
			return c.BehaviorUsecase.Add(task.IncrId, BehaviorType_sent_am_getting_started_email, time.Now(), "")
		} else if tpl == MailGenre_ContractReminder || tpl == MailGenre_AmContractReminder {
			if !configs.EnabledContractReminder {
				return nil
			}
			isAmContract := false
			if tpl == MailGenre_AmContractReminder {
				isAmContract = true
			}
			SubId := taskInput.GetInt("SubId")
			if SubId == 1 {
				behaviorType := BehaviorType_contract_reminder_first
				if isAmContract {
					behaviorType = BehaviorType_am_contract_reminder_first
				}
				err := c.BehaviorUsecase.Add(task.IncrId, behaviorType, time.Now(), "")
				if err != nil {
					return err
				}
				err = c.DialpadbuzUsecase.HandleContractReminder(config_vbc.ContractReminderFirst, task.IncrId, isAmContract)
				if err != nil {
					c.log.Error(err)
					return err
				}

			}
			if SubId == 2 {
				behaviorType := BehaviorType_contract_reminder_second
				if isAmContract {
					behaviorType = BehaviorType_am_contract_reminder_second
				}
				err := c.BehaviorUsecase.Add(task.IncrId, behaviorType, time.Now(), "")
				if err != nil {
					return err
				}
				err = c.DialpadbuzUsecase.HandleContractReminder(config_vbc.ContractReminderSecond, task.IncrId, isAmContract)
				if err != nil {
					c.log.Error(err)
					return err
				}
			}
			if SubId == 3 {
				behaviorType := BehaviorType_contract_reminder_third
				if isAmContract {
					behaviorType = BehaviorType_am_contract_reminder_third
				}
				err := c.BehaviorUsecase.Add(task.IncrId, behaviorType, time.Now(), "")
				if err != nil {
					return err
				}
				err = c.DialpadbuzUsecase.HandleContractReminder(config_vbc.ContractReminderThird, task.IncrId, isAmContract)
				if err != nil {
					c.log.Error(err)
					return err
				}
			}
			if SubId == 4 {
				behaviorType := BehaviorType_contract_reminder_fourth
				if isAmContract {
					behaviorType = BehaviorType_am_contract_reminder_fourth
				}
				err := c.BehaviorUsecase.Add(task.IncrId, behaviorType, time.Now(), "")
				if err != nil {
					return err
				}
				err = c.DialpadbuzUsecase.HandleContractReminder(config_vbc.ContractReminderFourth, task.IncrId, isAmContract)
				if err != nil {
					c.log.Error(err)
					return err
				}
			}
		} else if tpl == MailGenre_StartYourVADisabilityClaimRepresentation {
			if configs.EnabledAmIntakeFormReminder {
				// 取消此客户，因为修改邮箱地址的历史提提醒
				er := c.CommonUsecase.DB().Model(&TaskEntity{}).
					Where("incr_id = ? and event =? and task_status=? and from_type=?",
						task.IncrId,
						Task_Dag_HandleReminder,
						Task_TaskStatus_processing,
						config_vbc.GroupReminderAmIntakeForm,
					).
					Updates(map[string]interface{}{
						"task_status": Task_TaskStatus_cancel,
						"updated_at":  time.Now().Unix()}).Error
				if er != nil {
					c.log.Error(er, " IncrId: ", task.IncrId, " ID: ", task.ID)
					return er
				}
				return c.ContractReminderUsecase.AmIntakeFormReminderFirstReminder(task.IncrId)
			}
		} else if tpl == MailGenre_AmIntakeFormReminder {
			SubId := taskInput.GetInt("SubId")
			if SubId == 1 {
				behaviorType := BehaviorType_amintakeform_reminder_first
				err := c.BehaviorUsecase.Add(task.IncrId, behaviorType, time.Now(), "")
				if err != nil {
					return err
				}

			}
			if SubId == 2 {
				behaviorType := BehaviorType_amintakeform_reminder_second
				err := c.BehaviorUsecase.Add(task.IncrId, behaviorType, time.Now(), "")
				if err != nil {
					return err
				}

			}
			if SubId == 3 {
				behaviorType := BehaviorType_amintakeform_reminder_third
				err := c.BehaviorUsecase.Add(task.IncrId, behaviorType, time.Now(), "")
				if err != nil {
					return err
				}

			}

		}
		return nil
	} else if task.Event == Task_Dag_BoxCreateClientContracts { // Client sign contract

		err := c.BehaviorUsecase.Add(task.IncrId, BehaviorType_complete_fee_schedule_contract, time.Now(), "")
		if err != nil {
			c.log.Error(err)
			return err
		}
	} else if task.Event == Task_Dag_BoxCreateFolderForNewClient {

		taskInputParams := lib.ToTypeMapByString(task.TaskInput)
		caseId := taskInputParams.GetInt("ClientId")
		if caseId <= 0 {
			c.log.Error("HandleMoving2122aFile error: caseId: ", caseId)
		} else {

			err := c.SendVa2122aUsecase.RunHandleSeparateAmContract(caseId)
			if err != nil {
				c.log.Error(err, " RunHandleSeparateAmContract: ", caseId)
			}

			_, err = c.MiscUsecase.HandleMoving2122aFile(caseId)
			if err != nil {
				c.log.Error("HandleMoving2122aFile error: ", err, " caseId: ", caseId)
			}
		}
	}
	return nil
}

func (c *TaskUsecase) WaitingTasks() (*sql.Rows, error) {
	sqlRows, err := c.CommonUsecase.DB().Table(TaskEntity{}.TableName()).
		Where("task_status=? and next_at<=?  and next_retry_at<=? and deleted_at=0",
			Task_TaskStatus_processing, time.Now().Unix(), time.Now().Unix()).Rows()
	return sqlRows, err
}

func (c *TaskUsecase) RunTaskJob(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("RunTaskJob Done")
				return
			default:

				//sqlRows, err := c.CommonUsecase.DB().Table(TaskEntity{}.TableName()).
				//	Where("task_status=? and next_at<=?  and next_retry_at<=? and deleted_at=0",
				//		Task_TaskStatus_processing, time.Now().Unix(), time.Now().Unix()).Rows()
				sqlRows, err := c.WaitingTasks()
				if err != nil {
					c.log.Error(err)
				} else {
					if sqlRows != nil {
						gLimit := lib.NewGLimit(3)
						var waitGroup sync.WaitGroup
						for sqlRows.Next() {
							var task TaskEntity
							err = c.CommonUsecase.DB().ScanRows(sqlRows, &task)
							if err != nil {
								c.log.Error(err)
								continue
							}

							waitGroup.Add(1)
							gLimit.Run(func() {
								err = c.ExecTask(&task)
								if err != nil {
									c.log.Error(err)
								}
								waitGroup.Done()
							})
							time.Sleep(1)
						}
						err = sqlRows.Close()
						waitGroup.Wait()
						if err != nil {
							c.log.Error(err)
						}
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	return nil
}

/*
func (c *TaskUsecase) CreateEnvelope(task *TaskEntity) error {
	var taskInput lib.TypeMap
	taskInput = lib.ToTypeMapByString(task.TaskInput)
	docusignTemplateId := InterfaceToString(taskInput.Get("docusignTemplateId"))
	clientName := InterfaceToString(taskInput.Get("clientName"))
	clientEmail := InterfaceToString(taskInput.Get("clientEmail"))
	agentName := InterfaceToString(taskInput.Get("agentName"))
	agentEmail := InterfaceToString(taskInput.Get("agentEmail"))

	if len(docusignTemplateId) == 0 || len(clientEmail) == 0 || len(agentEmail) == 0 {
		return errors.New("Parameters is wrong.")
	}
	envelopeSummary, err := c.DocuSignUsecase.CreateEnvelopeAndSent(docusignTemplateId, clientName, clientEmail, agentName, agentEmail)
	if err != nil {
		return err
	}
	if envelopeSummary == nil {
		return errors.New("envelopeSummary is nil.")
	}

	entity := ClientEnvelopeEntity{
		ClientId:       task.IncrId,
		EnvelopeId:     envelopeSummary.EnvelopeID,
		Uri:            envelopeSummary.URI,
		Status:         envelopeSummary.Status,
		StatusDatetime: envelopeSummary.StatusDateTime.Format(time.RFC3339),
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}
	return c.ClientEnvelopeUsecase.CommonUsecase.DB().Save(&entity).Error
}
*/

/*
func (c *TaskUsecase) BuzEmail(task *TaskEntity) error {

	tCustomer, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": task.IncrId})
	if err != nil {
		return err
	}
	taskInput := lib.ToTypeMapByString(task.TaskInput)
	tpl := InterfaceToString(taskInput.Get("Genre"))
	subId := taskInput.GetInt("SubId")
	tTpl, err := c.TUsecase.Data(Kind_email_tpls, Eq{"tpl": tpl, "sub_id": subId})
	if err != nil {
		return err
	}
	err, mailSubject, mailBody, email, senderEmail, senderName := c.MailUsecase.SendEmailWithData(tCustomer, tTpl)
	if err != nil {
		return err
	}
	task.TaskStatus = Task_TaskStatus_finish
	return c.CommonUsecase.DB().Save(&EmailLogEntity{
		ClientId:   task.IncrId,
		Email:      email,
		TaskId:     task.ID,
		Tpl:        tpl,
		SubId:      subId,
		SenderMail: senderEmail,
		SenderName: senderName,
		Subject:    mailSubject,
		Body:       mailBody,
	}).Error
}
*/

func (c *TaskUsecase) UpdateTaskStatus(err error, task *TaskEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	if err == nil {
		task.TaskStatus = Task_TaskStatus_finish
	} else {
		task.TaskStatus = Task_TaskStatus_failure
	}
	task.UpdatedAt = time.Now().Unix()
	err1 := c.CommonUsecase.DB().Save(&task).Error
	if task.TaskStatus == Task_TaskStatus_failure {
		note := ""
		if err != nil {
			note = err.Error()
		}
		c.AddLog(task.ID, note)
	}
	return err1
}

func (c *TaskUsecase) ExecTask(task *TaskEntity) error {
	err := c.Invoke(task)
	return c.UpdateTaskStatus(err, task)
	return nil
}

// AddLog 添加日志
func (c *TaskUsecase) AddLog(taskId int32, notes string) error {

	return c.CommonUsecase.DB().Create(&TaskLogEntity{
		TaskId:    taskId,
		Notes:     notes,
		CreatedAt: time.Now().Unix(),
	}).Error
}
