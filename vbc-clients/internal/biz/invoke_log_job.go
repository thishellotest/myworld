package biz

import (
	"context"
	"database/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"regexp"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

type InvokeLogJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[InvokeLogEntity]
	//BaseHandle[InvokeLogEntity]
	BaseHandleCustom[InvokeLogEntity]
	NotesUsecase *NotesUsecase
	TUsecase     *TUsecase
}

func NewInvokeLogJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	NotesUsecase *NotesUsecase,
	TUsecase *TUsecase) *InvokeLogJobUsecase {
	uc := &InvokeLogJobUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		NotesUsecase:  NotesUsecase,
		TUsecase:      TUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()
	//uc.BaseHandle.Log = log.NewHelper(logger)
	//uc.BaseHandle.TableName = InvokeLogEntity{}.TableName()
	//uc.BaseHandle.DB = CommonUsecase.DB()
	//uc.BaseHandle.Handle = uc.Handle

	return uc
}

func (c *InvokeLogJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	return c.CommonUsecase.DB().
		Table(InvokeLogEntity{}.TableName()).
		Where("handle_status=? and next_at<=? and deleted_at=0",
			HandleStatus_waiting, time.Now().Unix()).Rows()

}

func (c *InvokeLogJobUsecase) Handle(ctx context.Context, task *InvokeLogEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	err := c.HandleExec(ctx, task)
	task.HandleStatus = HandleStatus_done
	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
	if err != nil {
		task.HandleResult = HandleResult_failure
		task.AppendHandleResultDetail(err.Error())
	} else {
		task.HandleResult = HandleResult_ok
	}
	return c.CommonUsecase.DB().Save(task).Error
}

func (c *InvokeLogJobUsecase) FormatZohoNoteContent(content string) string {
	// 正则匹配 `user#` 后的第一串数字
	re := regexp.MustCompile(`crm\[user#(\d+)#\d+\]crm`)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		userID := matches[1]
		tUser, _ := c.TUsecase.DataByGid(Kind_users, userID)
		userFullName := ""
		if tUser != nil {
			userFullName = tUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname)
		}
		re1 := regexp.MustCompile(`crm\[user#(\d+)#\d+\]crm`)
		result := re1.ReplaceAllString(content, "@["+userFullName+"]($1)")
		return result
	} else {
		return content
	}
}

func (c *InvokeLogJobUsecase) HandleExec(ctx context.Context, task *InvokeLogEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	contentMap := task.GetContentMap()
	relateId := contentMap.GetString("Parent_Id.id")
	moduleName := contentMap.GetString("Parent_Id.module.api_name")
	if moduleName != config_zoho.Deals && moduleName != config_zoho.Contacts {
		return errors.New("no handle moduleName: " + moduleName)
	}
	Note_Content := c.FormatZohoNoteContent(contentMap.GetString("Note_Content"))
	Created_By := contentMap.GetString("Created_By.id")
	Modified_By := contentMap.GetString("Modified_By.id")
	Modified_Time := contentMap.GetString("Modified_Time")
	Created_Time := contentMap.GetString("Created_Time")
	Record_Status__s := contentMap.GetString("Record_Status__s")
	if Record_Status__s != "Available" {
		return errors.New("no handle Record_Status__s:" + Record_Status__s)
	}

	db := c.CommonUsecase.DB().Table("notes")

	result := map[string]interface{}{}
	er := db.Where("gid=?", task.RecordUniqId).Take(&result).Error
	isNewData := false
	if er == gorm.ErrRecordNotFound {
		isNewData = true
	}

	destMap := make(lib.TypeMap)
	ModifiedAt, err := time.Parse(time.RFC3339, Modified_Time)
	if err != nil {
		return err
	}
	CreatedAt, err := time.Parse(time.RFC3339, Created_Time)
	if err != nil {
		return err
	}
	if moduleName == config_zoho.Contacts {
		destMap.Set(Notes_FieldName_kind, Kind_clients)
	}
	destMap.Set(Notes_FieldName_kind_gid, relateId)
	destMap.Set(DataEntry_gid, task.RecordUniqId)
	destMap.Set(Notes_FieldName_content, Note_Content)
	destMap.Set("created_by", Created_By)
	destMap.Set("modified_by", Modified_By)
	destMap.Set("created_at", CreatedAt.Unix())
	destMap.Set("updated_at", ModifiedAt.Unix())
	if isNewData {
		db := c.CommonUsecase.DB().Table("notes")
		err = db.Create(map[string]interface{}(destMap)).Error
		if err != nil {
			return err
		}
	} else {

		err = c.CommonUsecase.DB().Table("notes").
			Where("gid=?", task.RecordUniqId).
			Updates(map[string]interface{}(destMap)).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
