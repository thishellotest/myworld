package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/configs"
	"vbc/internal/conf"
)

type CaseWithoutTaskUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	TUsecase      *TUsecase
	MapUsecase    *MapUsecase
	MailUsecase   *MailUsecase
}

func NewCaseWithoutTaskUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase,
	MailUsecase *MailUsecase,
) *CaseWithoutTaskUsecase {
	uc := &CaseWithoutTaskUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
		MapUsecase:    MapUsecase,
		MailUsecase:   MailUsecase,
	}

	return uc
}

type CaseWithoutTaskVo struct {
	StagesName string
	Items      []CaseWithoutTaskItem
}
type CaseWithoutTaskItem struct {
	ClientCaseName string
	CreatedTime    string
	Gid            string
	IsNew          bool
}

func (c *CaseWithoutTaskUsecase) IsNewAndSet(caseId int32, stage string) (bool, error) {
	key := MapKeyCaseWithoutTaskFlag(caseId, stage)
	a, err := c.MapUsecase.GetForString(key)
	if err != nil {
		c.log.Error(err)
		return false, err
	}
	if a == "" {
		c.MapUsecase.Set(key, "1")
		return true, nil
	}
	return false, nil
}

func (c *CaseWithoutTaskUsecase) GetCases() (result []*CaseWithoutTaskVo, err error) {

	list, err := c.GetCasesForNotify()
	if err != nil {
		return nil, err
	}
	res := make(map[string]*CaseWithoutTaskVo)
	for _, v := range list {
		fields := v.CustomFields
		stage := fields.TextValueByNameBasic(FieldName_stages)
		if _, ok := res[stage]; !ok {
			caseWithoutTaskVo := &CaseWithoutTaskVo{
				StagesName: stage,
			}
			res[stage] = caseWithoutTaskVo
			result = append(result, caseWithoutTaskVo)
		}

		createTime, _ := TimeToVBCDisplay(fields.TextValueByNameBasic("created_time"))

		isNewAndSet, err := c.IsNewAndSet(fields.NumberValueByNameBasic("id"), fields.TextValueByNameBasic(FieldName_stages))
		if err != nil {
			return nil, err
		}

		res[stage].Items = append(res[stage].Items, CaseWithoutTaskItem{
			ClientCaseName: fields.TextValueByNameBasic(FieldName_deal_name),
			CreatedTime:    createTime,
			IsNew:          isNewAndSet,
			Gid:            fields.TextValueByNameBasic("gid"),
		})
	}
	return
}

func (c *CaseWithoutTaskUsecase) GetCasesForNotify() (tDataList TDataList, err error) {

	//return
	//	sql := `select * from client_cases where email not in ('liaogling@gmail.com', 'lialing@foxmail.com') and id in (
	//select t.id as id  from (
	//select id, SUBSTRING_INDEX(stages, '.', 1) as stages, stages as source_stages,gid  from client_cases
	//where biz_deleted_at=0 and deleted_at=0 ) as t
	//left join client_tasks on t.gid=client_tasks.what_id_gid and client_tasks.biz_deleted_at=0 and client_tasks.status!='Completed'
	//where t.stages!=25 and  t.stages!=26 and t.stages!=27 and (  client_tasks.status is null)
	//) order by  CONVERT(stages, SIGNED INTEGER) asc ,id  asc`

	sql := `select * from client_cases where email not in ('liaogling@gmail.com', 'lialing@foxmail.com') and id in (
select t.id as id  from (
select id, stages as stages, stages as source_stages,gid  from client_cases 
where biz_deleted_at=0 and deleted_at=0 ) as t 
left join client_tasks on t.gid=client_tasks.what_id_gid and client_tasks.biz_deleted_at=0 and client_tasks.status!='Completed'
where t.stages!='Completed' and  t.stages!='Terminated' and t.stages!='Dormant' and (  client_tasks.status is null)
) order by  CONVERT(stages, SIGNED INTEGER) asc ,id  asc`

	return c.TUsecase.ListByRawSql(Kind_client_cases, sql)
}

func (c *CaseWithoutTaskUsecase) NotifyEmailBody() (vo *MailMessageVo, err error) {

	subject := "VBC: List of Client Cases without a Task"
	cases, err := c.GetCases()
	if err != nil {
		c.log.Error(err)
		return nil, err
	}

	body, err := CaseWithoutTasksEmailBody(subject, cases)
	if err != nil {
		return nil, err
	}
	vo = &MailMessageVo{
		Email:   "ywang@vetbenefitscenter.com;ebunting@vetbenefitscenter.com",
		Subject: subject,
		Body:    body,
	}

	if configs.IsDev() {
		vo.Email = "liaogling@gmail.com;lialing@foxmail.com"
	}

	return vo, nil
}

func (c *CaseWithoutTaskUsecase) ReminderManager() error {
	vo, err := c.NotifyEmailBody()

	if err != nil {
		return err
	}

	mailServiceConfig := InitMailServiceConfig()
	mailMessage := &MailMessage{
		To:      vo.Email,
		Subject: vo.Subject,
		Body:    vo.Body,
	}

	err = c.MailUsecase.SendEmail(mailServiceConfig, mailMessage, "", nil)
	c.CommonUsecase.DB().Save(&EmailLogEntity{
		ClientId:   0,
		Email:      vo.Email,
		TaskId:     0,
		Tpl:        "ListClientCasesWithoutTasks",
		SubId:      0,
		SenderMail: mailServiceConfig.Username,
		SenderName: mailServiceConfig.Name,
		Subject:    mailMessage.Subject,
		Body:       mailMessage.Body,
	})

	return nil
}
