package biz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"runtime"
	"strconv"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

func ZoomBoxFolderId() string {
	if configs.IsProd() {
		return "280127932304"
	}
	return "280914505349"
}

type ZoomRecordingFileJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[ZoomRecordingFileEntity]
	BaseHandleCustom[ZoomRecordingFileEntity]
	ZoomMeetingUsecase       *ZoomMeetingUsecase
	BoxUsecase               *BoxUsecase
	ZoomUsecase              *ZoomUsecase
	ZoombuzUsecase           *ZoombuzUsecase
	ZoomRecordingFileUsecase *ZoomRecordingFileUsecase
	ZoomUploadBoxUsecase     *ZoomUploadBoxUsecase
}

func NewZoomRecordingFileJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ZoomMeetingUsecase *ZoomMeetingUsecase,
	BoxUsecase *BoxUsecase,
	ZoomUsecase *ZoomUsecase,
	ZoombuzUsecase *ZoombuzUsecase,
	ZoomRecordingFileUsecase *ZoomRecordingFileUsecase,
	ZoomUploadBoxUsecase *ZoomUploadBoxUsecase) *ZoomRecordingFileJobUsecase {
	uc := &ZoomRecordingFileJobUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		ZoomMeetingUsecase:       ZoomMeetingUsecase,
		BoxUsecase:               BoxUsecase,
		ZoomUsecase:              ZoomUsecase,
		ZoombuzUsecase:           ZoombuzUsecase,
		ZoomRecordingFileUsecase: ZoomRecordingFileUsecase,
		ZoomUploadBoxUsecase:     ZoomUploadBoxUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}

// ExecuteCrontabHandleProcessingRecording 处理zoom还没有完成的任务
func (c *ZoomRecordingFileJobUsecase) ExecuteCrontabHandleProcessingRecording() error {

	//processing的数据，可以是异常的 todo:lgl
	sql := `select zoom_recording_files.* from zoom_recording_files 
inner join zoom_meetings on zoom_meetings.meeting_uuid=zoom_recording_files.meeting_uuid
where zoom_recording_files.handle_status=0 and zoom_recording_files.deleted_at=0 and 
zoom_recording_files.status='processing' and zoom_recording_files.handle_result=0`
	//sql = fmt.Sprintf("%s and zoom_recording_files.", sql)

	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		c.log.Error(err)
		return err
	}
	defer sqlRows.Close()

	entities, err := lib.SqlRowsToEntities[ZoomRecordingFileEntity](c.CommonUsecase.DB(), sqlRows)
	if err != nil {
		c.log.Error(err)
		return err
	}

	//lib.DPrintln(entities)
	for k, v := range entities {
		meetings := make(map[string]HttpResponseBody)
		meeting, err := c.ZoombuzUsecase.Meeting(meetings, v.MeetingUuid)
		if err != nil {
			if (meeting.HttpCode == 404 &&
				((entities[k].CreatedAt + 7200) > time.Now().Unix())) ||
				meeting.HttpCode != 404 {
				entities[k].HandleResult = HandleResult_HandleProcessing_Error
				entities[k].AppendHandleResultDetail(err.Error())
				err = c.CommonUsecase.DB().Save(&entities[k]).Error
				if err != nil {
					c.log.Error(err)
				}
			}
		} else {
			//recordingFiles := meeting.GetTypeList("recording_files")
			recordingFile := ZoomGetRecordingFile(meeting.Body, v.RecordingFileId)
			if recordingFile == nil {
				entities[k].HandleResult = HandleResult_HandleProcessing_Recording_Missing
				entities[k].AppendHandleResultDetail("Missing Recording File")
				err = c.CommonUsecase.DB().Save(&entities[k]).Error
				if err != nil {
					c.log.Error(err)
				}
			} else {
				entities[k].RecordingStart = recordingFile.GetString("recording_start")
				entities[k].RecordingEnd = recordingFile.GetString("recording_end")
				entities[k].FileType = recordingFile.GetString("file_type")
				entities[k].FileExtension = recordingFile.GetString("file_extension")
				entities[k].FileSize = recordingFile.GetString("file_size")
				entities[k].PlayUrl = recordingFile.GetString("play_url")
				entities[k].DownloadUrl = recordingFile.GetString("download_url")
				entities[k].Status = recordingFile.GetString("status")
				entities[k].RecordingType = recordingFile.GetString("recording_type")
				entities[k].UpdatedAt = time.Now().Unix()
				err = c.CommonUsecase.DB().Save(&entities[k]).Error
				if err != nil {
					c.log.Error(err)
				}
			}
		}
	}

	return nil
}

func (c *ZoomRecordingFileJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	var sql string
	if configs.IsJobTypeDefault() {
		sql = `select zoom_recording_files.* from zoom_recording_files 
inner join zoom_meetings on zoom_meetings.meeting_uuid=zoom_recording_files.meeting_uuid
where zoom_recording_files.handle_status=0 and zoom_recording_files.deleted_at=0 and 
      zoom_recording_files.status='completed' and (zoom_recording_files.file_size<=203838190 or zoom_recording_files.file_size>289395552 )`

		sql = `select zoom_recording_files.* from zoom_recording_files 
inner join zoom_meetings on zoom_meetings.meeting_uuid=zoom_recording_files.meeting_uuid
where zoom_recording_files.handle_status=0 and zoom_recording_files.deleted_at=0 and 
      zoom_recording_files.status='completed' `

	} else if configs.IsJobTypeLargeMemory() {
		sql = `select zoom_recording_files.* from zoom_recording_files 
inner join zoom_meetings on zoom_meetings.meeting_uuid=zoom_recording_files.meeting_uuid
where zoom_recording_files.handle_status=0 and zoom_recording_files.deleted_at=0 and 
      zoom_recording_files.status='completed' and zoom_recording_files.file_size>203838190 and zoom_recording_files.file_size<=289395552`

		// 关闭， 统一使用上面
		sql = `select zoom_recording_files.* from zoom_recording_files 
inner join zoom_meetings on zoom_meetings.meeting_uuid=zoom_recording_files.meeting_uuid
where zoom_recording_files.handle_status=0 and zoom_recording_files.deleted_at=0 and 
      zoom_recording_files.status='completed' and zoom_recording_files.file_size>203838190 and zoom_recording_files.file_size<=289395552 and 1!=1`

		// 不能超过：289395552 了，
	} else {
		panic("ZoomRecordingFileJobUsecase:WaitingTasks config")
	}

	return c.CommonUsecase.DB().Raw(sql).Rows()
	//
	//return c.CommonUsecase.DB().
	//	Table(ZoomRecordingFileEntity{}.TableName()).
	//	Where("handle_status=? ",
	//		HandleStatus_waiting).Rows()
}

func (c *ZoomRecordingFileJobUsecase) Handle(ctx context.Context, task *ZoomRecordingFileEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}

	var err error
	if task.GetFileSize() >= 103838190 {
		err = c.HandleExecUseChunks(ctx, task)
	} else {
		err = c.HandleExec(ctx, task)
	}
	runtime.GC()

	task.HandleStatus = HandleStatus_done
	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
	if err != nil {
		task.HandleResult = HandleResult_failure
		task.AppendHandleResultDetail(err.Error())
	} else {
		task.HandleResult = HandleResult_ok
	}
	return c.CommonUsecase.DB().Model(&task).Updates(
		map[string]interface{}{
			"handle_status":        task.HandleStatus,
			"handle_result":        task.HandleResult,
			"handle_result_detail": task.HandleResultDetail,
			"updated_at":           task.UpdatedAt},
	).Error
}

func (c *ZoomRecordingFileJobUsecase) HandleExec(ctx context.Context, task *ZoomRecordingFileEntity) error {

	if task == nil {
		return nil
	}

	headers, err := c.ZoomUsecase.Headers()
	if err != nil {
		return err
	}

	err = c.ZoombuzUsecase.UpdateZoomRecordingFile(task)
	if err != nil {
		c.log.Error(err, " ID: ", task.ID)
		return err
	}

	lastTask, err := c.ZoomRecordingFileUsecase.GetByCond(Eq{"id": task.ID})
	if err != nil {
		return err
	}
	if lastTask == nil {
		return errors.New("lastTask is nil")
	}

	httpResponse, err := lib.RequestDoTimeout("GET", lastTask.DownloadUrl, nil, headers, time.Hour)
	//httpResponse, err := http.Get(task.DownloadUrl)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if httpResponse == nil {
		return errors.New("httpResponse is nil")
	}
	defer httpResponse.Body.Close()

	folderId, err := c.ZoomMeetingUsecase.BoxFolderId(lastTask.MeetingUuid)
	if err != nil {
		c.log.Error(err, " MeetingUuid:", lastTask.MeetingUuid, " ID:", lastTask.ID)
		return err
	}
	fileId, err := c.BoxUsecase.UploadFile(folderId, httpResponse.Body, lastTask.FileName())
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CommonUsecase.DB().Model(lastTask).Updates(&ZoomRecordingFileEntity{
		BoxResId: fileId,
	}).Error
}

func (c *ZoomRecordingFileJobUsecase) HandleExecUseChunks(ctx context.Context, task *ZoomRecordingFileEntity) error {
	if task == nil {
		return nil
	}
	time.Sleep(time.Second * 20)
	c.log.Debug(fmt.Sprintf("HandleExecUseChunks MeetingUuid: %s ID: %d FileSize: %s RecordingType: %s ", task.MeetingUuid, task.ID, task.FileSize, task.RecordingType))
	err := c.ZoombuzUsecase.UpdateZoomRecordingFile(task)
	if err != nil {
		c.log.Error(err, " ID: ", task.ID)
		return err
	}

	lastTask, err := c.ZoomRecordingFileUsecase.GetByCond(Eq{"id": task.ID})
	if err != nil {
		return err
	}
	if lastTask == nil {
		return errors.New("lastTask is nil")
	}

	folderId, err := c.ZoomMeetingUsecase.BoxFolderId(lastTask.MeetingUuid)
	if err != nil {
		c.log.Error(err, " MeetingUuid:", lastTask.MeetingUuid, " ID:", lastTask.ID)
		return err
	}
	boxFileName := lastTask.FileName()

	fileSize, _ := strconv.ParseInt(lastTask.FileSize, 10, 32)
	fileId, err := c.ZoomUploadBoxUsecase.UploadToBoxAndCheck(lastTask.DownloadUrl, folderId, boxFileName, fileSize)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CommonUsecase.DB().Model(lastTask).Updates(&ZoomRecordingFileEntity{
		BoxResId: fileId,
	}).Error
}
