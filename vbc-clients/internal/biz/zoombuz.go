package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

/*
[{
		"uuid": "UaOnHojhRZueVNCVcK8lEA==",
		"id": 87287728097,
		"account_id": "pzf3p9rDTmevLobo-La4HQ",
		"host_id": "4P41WITpTYedmIqdtNFuhw",
		"topic": "Conference with Timothy Fortson (I)",
		"type": 2,
		"start_time": "2024-08-16T23:27:36Z",
		"timezone": "America/Los_Angeles",
		"duration": 44,
		"total_size": 196200715,
		"recording_count": 4,
		"share_url": "https://us06web.zoom.us/rec/share/lNw00vo7xyC4wsIWNz3cwhvZ7BSVPU2fJNrCj7FitAQ0_nKjvqR9dOPFBCa6_LrZ.4txU15rydT86hB4Z",
		"recording_files": [{
			"id": "65f52c07-55da-419d-948b-ab75f55b2b99",
			"meeting_id": "UaOnHojhRZueVNCVcK8lEA==",
			"recording_start": "2024-08-16T23:27:36Z",
			"recording_end": "2024-08-17T00:12:00Z",
			"file_type": "TRANSCRIPT",
			"file_extension": "VTT",
			"file_size": 45865,
			"play_url": "https://us06web.zoom.us/rec/play/wCYiob-8WZGGc8hH0d1SBPRe1ds1GUIkAhAtz_LQVn1TnGUwNtMFuXp08xYWJiD-vnpR1RnF2F3MyHNE.WxMEJVFSFOO8Tb0k",
			"download_url": "https://us06web.zoom.us/rec/download/wCYiob-8WZGGc8hH0d1SBPRe1ds1GUIkAhAtz_LQVn1TnGUwNtMFuXp08xYWJiD-vnpR1RnF2F3MyHNE.WxMEJVFSFOO8Tb0k",
			"status": "completed",
			"recording_type": "audio_transcript"
		}, {
			"id": "663f4f9c-761f-49a9-bf4a-fabc42785908",
			"meeting_id": "UaOnHojhRZueVNCVcK8lEA==",
			"recording_start": "2024-08-16T23:27:36Z",
			"recording_end": "2024-08-17T00:12:00Z",
			"file_type": "SUMMARY",
			"file_extension": "JSON",
			"file_size": 5850,
			"play_url": "https://us06web.zoom.us/rec/play/BAWqAFhog5Va6sjIT9cKiZyQdEkb804eGRcHxkxoMHWk2CSgyhi2dWL51pNGHwXpD9D6JjzJi4pcD_mZ.vTp-23AOrcitHgW5",
			"download_url": "https://us06web.zoom.us/rec/download/BAWqAFhog5Va6sjIT9cKiZyQdEkb804eGRcHxkxoMHWk2CSgyhi2dWL51pNGHwXpD9D6JjzJi4pcD_mZ.vTp-23AOrcitHgW5",
			"status": "completed",
			"recording_type": "summary"
		}, {
			"id": "6868bff4-09cf-4f8c-8f1a-b0232cd3884c",
			"meeting_id": "UaOnHojhRZueVNCVcK8lEA==",
			"recording_start": "2024-08-16T23:27:36Z",
			"recording_end": "2024-08-17T00:12:00Z",
			"file_type": "M4A",
			"file_extension": "M4A",
			"file_size": 42392752,
			"play_url": "https://us06web.zoom.us/rec/play/r5ZrfmKcofwOHPZBVCAuYbsnkDvTct9RspowYBUyOvgn8-hiQ4q0ITaVkG0TUCdOJA0j1HWDXblEc25u.14w9rYTtM5JyoVZI",
			"download_url": "https://us06web.zoom.us/rec/download/r5ZrfmKcofwOHPZBVCAuYbsnkDvTct9RspowYBUyOvgn8-hiQ4q0ITaVkG0TUCdOJA0j1HWDXblEc25u.14w9rYTtM5JyoVZI",
			"status": "completed",
			"recording_type": "audio_only"
		}, {
			"id": "1df78c0d-8039-4800-8a6a-818e1e60271c",
			"meeting_id": "UaOnHojhRZueVNCVcK8lEA==",
			"recording_start": "2024-08-16T23:27:36Z",
			"recording_end": "2024-08-17T00:12:00Z",
			"file_type": "TIMELINE",
			"file_extension": "JSON",
			"file_size": 995956,
			"download_url": "https://us06web.zoom.us/rec/download/AA-47-OwOvWr0Gfo1_FY6v7tx5g0WNDTTSNof5uJATKMPdMufMIGUjSHx7E5zTccjzOf2VOYU2eQhVKX.CukzJC53CoTOYuZ6",
			"status": "completed",
			"recording_type": "timeline"
		}, {
			"id": "766f6577-0e86-4b1d-92e7-ed62924d9cdf",
			"meeting_id": "UaOnHojhRZueVNCVcK8lEA==",
			"recording_start": "2024-08-16T23:27:36Z",
			"recording_end": "2024-08-17T00:12:00Z",
			"file_type": "MP4",
			"file_extension": "MP4",
			"file_size": 152766142,
			"play_url": "https://us06web.zoom.us/rec/play/nU3BxKGSLL9i7lxIGbhk6aiNImtgus_90QAkM8U3fTBIMRce9IhavBEdlMw2Wx4EavPkBzreooswnEtV.eI9M3Pv66Rm3Odri",
			"download_url": "https://us06web.zoom.us/rec/download/nU3BxKGSLL9i7lxIGbhk6aiNImtgus_90QAkM8U3fTBIMRce9IhavBEdlMw2Wx4EavPkBzreooswnEtV.eI9M3Pv66Rm3Odri",
			"status": "completed",
			"recording_type": "shared_screen_with_speaker_view"
		}, {
			"id": "d67844c9-168e-4925-97b4-1598f66cd0ab",
			"meeting_id": "UaOnHojhRZueVNCVcK8lEA==",
			"recording_start": "2024-08-16T23:27:36Z",
			"recording_end": "2024-08-17T00:12:00Z",
			"file_type": "SUMMARY",
			"file_extension": "JSON",
			"file_size": 2953,
			"play_url": "https://us06web.zoom.us/rec/play/JaNcYzci4V2rDt7LL_c41z0_lt2oUBTLI95mL4I-dGlw3WK4z1RZnxAqV3jVrSqH7t_y1My4uAZZqO8z.aLC0Iar0-95aJvIe",
			"download_url": "https://us06web.zoom.us/rec/download/JaNcYzci4V2rDt7LL_c41z0_lt2oUBTLI95mL4I-dGlw3WK4z1RZnxAqV3jVrSqH7t_y1My4uAZZqO8z.aLC0Iar0-95aJvIe",
			"status": "completed",
			"recording_type": "summary_next_steps"
		}],
		"recording_play_passcode": "YVeo6PyI40urDPhsjhmXcoslqmwhFz-O"
	}
]
*/

type ZoombuzUsecase struct {
	log                      *log.Helper
	CommonUsecase            *CommonUsecase
	conf                     *conf.Data
	ZoomMeetingUsecase       *ZoomMeetingUsecase
	ZoomRecordingFileUsecase *ZoomRecordingFileUsecase
	ZoomUserUsecase          *ZoomUserUsecase
	ZoomUsecase              *ZoomUsecase
}

func NewZoombuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ZoomMeetingUsecase *ZoomMeetingUsecase,
	ZoomRecordingFileUsecase *ZoomRecordingFileUsecase,
	ZoomUserUsecase *ZoomUserUsecase,
	ZoomUsecase *ZoomUsecase) *ZoombuzUsecase {
	uc := &ZoombuzUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		ZoomMeetingUsecase:       ZoomMeetingUsecase,
		ZoomRecordingFileUsecase: ZoomRecordingFileUsecase,
		ZoomUserUsecase:          ZoomUserUsecase,
		ZoomUsecase:              ZoomUsecase,
	}

	return uc
}

func (c *ZoombuzUsecase) InitSyncRecords() error {
	users, err := c.ZoomUserUsecase.AllByCond(Eq{"deleted_at": 0})
	if err != nil {
		return err
	}
	for _, v := range users {
		userRecords, err := c.ZoomUsecase.UsersRecordings(v.UserId, "2024-07-01", "")
		if err != nil {
			return err
		}
		if userRecords == nil {
			return errors.New("userRecords is nil")
		}
		records := lib.ToTypeMapByString(*userRecords)
		meetings := records.GetTypeList("meetings")
		err = c.SyncRecords(meetings, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ZoombuzUsecase) ExecuteSyncRecords() error {
	users, err := c.ZoomUserUsecase.AllByCond(Eq{"deleted_at": 0})
	if err != nil {
		return err
	}
	for _, v := range users {
		err := c.BizExecuteSyncRecords(v.UserId)
		if err != nil {
			//c.log.Error(err)
		}
	}
	return nil
}

func (c *ZoombuzUsecase) BizExecuteSyncRecords(zoomUserId string) error {

	ti := time.Now()
	ti = ti.AddDate(0, 0, -5)

	userRecords, err := c.ZoomUsecase.UsersRecordings(zoomUserId, ti.Format(time.DateOnly), "")
	if err != nil {
		//c.log.Error(err, "zoomUserId:", zoomUserId)
		return err
	}
	if userRecords == nil {
		return errors.New("userRecords is nil")
	}
	records := lib.ToTypeMapByString(*userRecords)
	meetings := records.GetTypeList("meetings")
	err = c.SyncRecords(meetings, false)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}

func (c *ZoombuzUsecase) UpdateZoomRecordingFile(zoomRecordingFileEntity *ZoomRecordingFileEntity) error {
	if zoomRecordingFileEntity == nil {
		return nil
	}
	meeting, err := c.ZoomMeetingUsecase.GetByCond(Eq{"deleted_at": 0, "meeting_uuid": zoomRecordingFileEntity.MeetingUuid})
	if err != nil {
		return err
	}
	if meeting == nil {
		return errors.New("meeting is nil")
	}
	//startTime, err := time.Parse(time.RFC3339, meeting.StartTime)
	//if err != nil {
	//	c.log.Error(err)
	//	return err
	//}
	// 可以使用：MeetingRecordings 替代，此处有问题 todo:lgl
	meetingMap, _, _, err := c.ZoomUsecase.MeetingRecordings(meeting.MeetingUuid)
	if err != nil {
		return err
	}
	/*
		recordings, err := c.ZoomUsecase.UsersRecordings(meeting.HostId, startTime.Format(time.DateOnly), startTime.Format(time.DateOnly))
		if err != nil {
			return err
		}
		if recordings == nil {
			return errors.New("recordings is nil")
		}
		records := lib.ToTypeMapByString(*recordings)
		meetings := records.GetTypeList("meetings")
		meetingMap := ZoomGetMeeting(meetings, meeting.MeetingId)*/

	if meetingMap == nil {
		return errors.New("meetingMap is nil")
	}
	recordingFileMap := ZoomGetRecordingFile(meetingMap, zoomRecordingFileEntity.RecordingFileId)
	if recordingFileMap == nil {
		return errors.New("recordingFileMap is nil")
	}
	err = c.CommonUsecase.DB().Model(zoomRecordingFileEntity).Updates(&ZoomRecordingFileEntity{
		PlayUrl:       recordingFileMap.GetString("play_url"),
		DownloadUrl:   recordingFileMap.GetString("download_url"),
		Status:        recordingFileMap.GetString("status"),
		RecordingType: recordingFileMap.GetString("recording_type"),
		FileExtension: recordingFileMap.GetString("file_extension"),
		FileSize:      recordingFileMap.GetString("file_size"),
		FileType:      recordingFileMap.GetString("file_type"),
		UpdatedAt:     time.Now().Unix(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func ZoomGetRecordingFile(meetingMap lib.TypeMap, recordingFileId string) lib.TypeMap {
	recordingFiles := meetingMap.GetTypeList("recording_files")
	for k, v := range recordingFiles {
		if v.GetString("id") == recordingFileId {
			return recordingFiles[k]
		}
	}
	return nil
}

func ZoomGetMeeting(meetings lib.TypeList, meetingId string) lib.TypeMap {
	for k, v := range meetings {
		if v.GetString("id") == meetingId {
			return meetings[k]
		}
	}
	return nil
}

func (c *ZoombuzUsecase) SyncRecords(records lib.TypeList, needUpdateFiles bool) error {

	if len(records) == 0 {
		return nil
	}
	var meetingUuids []string
	var recordingFiles lib.TypeList
	for _, v := range records {
		meetingUuids = append(meetingUuids, v.GetString("uuid"))
		recordingFiles.AppendList(v.GetTypeList("recording_files"))
	}
	meetings, err := c.ZoomMeetingUsecase.AllByCond(In("meeting_uuid", meetingUuids))
	if err != nil {
		return err
	}
	existMeetings := make(map[string]*ZoomMeetingEntity)
	for k, v := range meetings {
		existMeetings[v.MeetingUuid] = meetings[k]
	}

	for _, v := range records {
		meetingUuid := v.GetString("uuid")
		if _, ok := existMeetings[meetingUuid]; !ok {
			meeting := &ZoomMeetingEntity{
				MeetingUuid:           meetingUuid,
				MeetingId:             v.GetString("id"),
				AccountId:             v.GetString("account_id"),
				HostId:                v.GetString("host_id"),
				Topic:                 v.GetString("topic"),
				Type:                  v.GetString("type"),
				StartTime:             v.GetString("start_time"),
				Timezone:              v.GetString("timezone"),
				Duration:              v.GetString("duration"),
				TotalSize:             v.GetString("total_size"),
				RecordingCount:        int(v.GetInt("recording_count")),
				ShareUrl:              v.GetString("share_url"),
				RecordingPlayPasscode: v.GetString("recording_play_passcode"),
				CreatedAt:             time.Now().Unix(),
				UpdatedAt:             time.Now().Unix(),
			}
			err = c.CommonUsecase.DB().Save(&meeting).Error
			if err != nil {
				c.log.Error(err)
				return err
			}
		}
	}
	return c.SyncRecordingFiles(recordingFiles, needUpdateFiles)
}

func (c *ZoombuzUsecase) SyncRecordingFiles(recordingFiles lib.TypeList, needUpdateFiles bool) error {

	var ids []string
	for _, v := range recordingFiles {
		ids = append(ids, v.GetString("id"))
	}
	records, err := c.ZoomRecordingFileUsecase.AllByCond(In("recording_file_id", ids))
	if err != nil {
		return err
	}
	existRecords := make(map[string]*ZoomRecordingFileEntity)
	for k, v := range records {
		existRecords[v.RecordingFileId] = records[k]
	}
	for _, v := range recordingFiles {
		RecordingFileId := v.GetString("id")
		var zoomRecordingFileEntity *ZoomRecordingFileEntity
		ok := false
		if zoomRecordingFileEntity, ok = existRecords[RecordingFileId]; !ok {
			zoomRecordingFileEntity = &ZoomRecordingFileEntity{
				RecordingFileId: RecordingFileId,
				MeetingUuid:     v.GetString("meeting_id"),
				RecordingStart:  v.GetString("recording_start"),
				RecordingEnd:    v.GetString("recording_end"),
				FileType:        v.GetString("file_type"),
				FileExtension:   v.GetString("file_extension"),
				FileSize:        v.GetString("file_size"),
				PlayUrl:         v.GetString("play_url"),
				DownloadUrl:     v.GetString("download_url"),
				Status:          v.GetString("status"),
				RecordingType:   v.GetString("recording_type"),
				CreatedAt:       time.Now().Unix(),
				UpdatedAt:       time.Now().Unix(),
			}
			err = c.ZoomRecordingFileUsecase.CommonUsecase.DB().Save(&zoomRecordingFileEntity).Error
			if err != nil {
				return err
			}
		} else {
			if needUpdateFiles {
				err = c.ZoomRecordingFileUsecase.CommonUsecase.DB().Model(zoomRecordingFileEntity).Updates(&ZoomRecordingFileEntity{}).Error
				if err != nil {
					return err
				}
			}

			if zoomRecordingFileEntity.HandleResult == HandleResult_HandleProcessing_Error &&
				zoomRecordingFileEntity.Status == "processing" && v.GetString("status") == "completed" {

				zoomRecordingFileEntity.HandleResult = 0
				zoomRecordingFileEntity.RecordingEnd = v.GetString("recording_end")
				zoomRecordingFileEntity.FileType = v.GetString("file_type")
				zoomRecordingFileEntity.FileExtension = v.GetString("file_extension")
				zoomRecordingFileEntity.FileSize = v.GetString("file_size")
				zoomRecordingFileEntity.PlayUrl = v.GetString("play_url")
				zoomRecordingFileEntity.DownloadUrl = v.GetString("download_url")
				zoomRecordingFileEntity.Status = v.GetString("status")
				zoomRecordingFileEntity.RecordingType = v.GetString("recording_type")
				err = c.ZoomRecordingFileUsecase.CommonUsecase.DB().Save(&zoomRecordingFileEntity).Error
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *ZoombuzUsecase) ExecuteSyncZoomUsers() error {

	res, err := c.ZoomUsecase.User()
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("res is nil")
	}
	resMap := lib.ToTypeMapByString(*res)
	users := resMap.GetTypeList("users")

	return c.SyncZoomUsers(users)
}

func (c *ZoombuzUsecase) SyncZoomUsers(list lib.TypeList) error {

	for _, v := range list {
		userId := v.GetString("id")
		var user *ZoomUserEntity
		var err error
		user, err = c.ZoomUserUsecase.GetByCond(Eq{"user_id": userId})
		if err != nil {
			return err
		}
		if user == nil {
			user = &ZoomUserEntity{
				CreatedAt: time.Now().Unix(),
			}
		}
		user.UpdatedAt = time.Now().Unix()
		user.UserId = userId
		user.FirstName = v.GetString("first_name")
		user.LastName = v.GetString("last_name")
		user.DisplayName = v.GetString("display_name")
		user.Email = v.GetString("email")
		user.Type = v.GetString("type")
		user.Pmi = v.GetString("pmi")
		user.Timezone = v.GetString("timezone")
		user.Verified = v.GetString("verified")
		user.CcreatedAt = v.GetString("created_at")
		user.LastLoginTime = v.GetString("last_login_time")
		user.LastClientVersion = v.GetString("last_client_version")
		user.PicUrl = v.GetString("pic_url")
		user.Language = v.GetString("language")
		user.Status = v.GetString("status")
		user.RoleId = v.GetString("role_id")
		user.UserCreatedAt = v.GetString("user_created_at")
		err = c.CommonUsecase.DB().Save(&user).Error
		if err != nil {
			return err
		}
	}

	return nil
}

type HttpResponseBody struct {
	HttpCode int
	RawBody  string
	Body     lib.TypeMap
}

func (c *ZoombuzUsecase) Meeting(meetings map[string]HttpResponseBody, meetingId string) (httpResponseBody HttpResponseBody, err error) {
	if meetings == nil {
		return httpResponseBody, errors.New("meetings is nil")
	}
	if _, ok := meetings[meetingId]; ok {
		return meetings[meetingId], nil
	} else {
		meeting, rawBody, httpCode, err := c.ZoomUsecase.MeetingRecordings(meetingId)
		httpResponseBody.RawBody = rawBody
		httpResponseBody.Body = meeting
		httpResponseBody.HttpCode = httpCode
		meetings[meetingId] = httpResponseBody
		return meetings[meetingId], err
	}
}

func (c *ZoombuzUsecase) ExecuteDeleteMeetingRecording() error {
	sql := fmt.Sprintf("select * from zoom_meetings where  zoom_deleted_at=0 and deleted_at=0 and start_time <='%s' and  NOT exists (select * from zoom_recording_files  where handle_status=0 and handle_result=0 and zoom_recording_files.meeting_uuid=zoom_meetings.meeting_uuid)", time.Now().UTC().AddDate(0, 0, -10).Format(time.RFC3339))
	c.log.Info("ExecuteDeleteMeetingRecording sql:", sql)
	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		c.log.Error(err)
		return err
	}
	defer sqlRows.Close()
	meetings, err := lib.SqlRowsToEntities[ZoomMeetingEntity](c.CommonUsecase.DB(), sqlRows)
	if err != nil {
		return err
	}
	for _, v := range meetings {
		time.Sleep(time.Second)
		err := c.DeleteMeetingRecording(v.MeetingUuid)
		if err != nil {
			c.log.Error(err, " : ", v.MeetingUuid, " : ", v.ID)
		} else {
			c.log.Info(" : ", v.MeetingUuid, " : ", v.ID)
		}
	}
	return nil
}

func (c *ZoombuzUsecase) DeleteMeetingRecording(meetingUuid string) error {

	meetingEntity, err := c.ZoomMeetingUsecase.GetByCond(Eq{"meeting_uuid": meetingUuid})
	if err != nil {
		return err
	}
	if meetingEntity == nil {
		return errors.New("meetingEntity is nil")
	}

	meeting, _, httpCode, err := c.ZoomUsecase.MeetingRecordings(meetingUuid)
	if err != nil && httpCode == 404 {
		meetingEntity.ZoomDeletedAt = 404
		return c.CommonUsecase.DB().Save(&meetingEntity).Error
	}
	if err != nil {
		return err
	}
	recordingFiles := meeting.GetTypeList("recording_files")
	dbRecordingFiles, err := c.ZoomRecordingFileUsecase.AllByCond(Eq{"meeting_uuid": meetingUuid})
	if err != nil {
		return err
	}
	needDeleted := true
	for _, v := range recordingFiles {
		id := v.GetString("id")
		if v.GetString("file_type") == "MP4" {
			continue
		}
		isOk := false
		for _, v1 := range dbRecordingFiles {
			if id == v1.RecordingFileId {
				if v1.HandleStatus == HandleStatus_done && v1.HandleResult == 0 {
					isOk = true
				}
				break
			}
		}
		if !isOk {
			needDeleted = false
			break
		}
	}
	if !needDeleted {
		return errors.New("meetingUuid : " + meetingUuid + " The synchronization is incomplete and needs to be manually handled ")
	}
	_, _, httpCode, err = c.ZoomUsecase.DeleteMeetingRecordings(meetingUuid)
	if err != nil && httpCode != 404 {
		c.log.Error(err)
		return err
	}
	meetingEntity.ZoomDeletedAt = time.Now().Unix()
	err = c.CommonUsecase.DB().Save(&meetingEntity).Error

	return err
}

func (c *ZoombuzUsecase) ListMeetingForSmsNotice(zoomUserId string) (lib.TypeMap, error) {
	now := time.Now()
	begin := now.AddDate(0, 0, -2)
	end := now.AddDate(0, 0, +2)
	r, _, _, err := c.ZoomUsecase.UsersMeetings(zoomUserId, begin.Format(time.DateOnly), end.Format(time.DateOnly))
	return r, err
}
