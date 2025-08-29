package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"sync"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ZoomMeetingEntity struct {
	ID                    int32 `gorm:"primaryKey"`
	MeetingUuid           string
	MeetingId             string // 一个MeetingId，可能对应多个MeetingUuid
	AccountId             string
	HostId                string // 这个对应Zoom UserId
	Topic                 string
	Type                  string
	StartTime             string
	Timezone              string
	Duration              string
	TotalSize             string
	RecordingCount        int
	ShareUrl              string
	RecordingPlayPasscode string
	BoxResId              string
	ZoomDeletedAt         int64
	CreatedAt             int64
	UpdatedAt             int64
	DeletedAt             int64
}

func (ZoomMeetingEntity) TableName() string {
	return "zoom_meetings"
}

//
//func (c *ZoomMeetingEntity) FolderName() string {
//	return fmt.Sprintf("%s_%s #%d", c.StartTime, lib.TrimCharacterForFileName(c.Topic), c.ID)
//}

type ZoomMeetingUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[ZoomMeetingEntity]
	BoxFolderIdLock sync.Mutex
	BoxUsecase      *BoxUsecase
	ZoomUserUsecase *ZoomUserUsecase
}

func NewZoomMeetingUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	BoxUsecase *BoxUsecase,
	ZoomUserUsecase *ZoomUserUsecase) *ZoomMeetingUsecase {
	uc := &ZoomMeetingUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		BoxUsecase:      BoxUsecase,
		ZoomUserUsecase: ZoomUserUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *ZoomMeetingUsecase) FolderName(zoomMeetingEntity *ZoomMeetingEntity) (string, error) {

	if zoomMeetingEntity == nil {
		return "", errors.New("zoomMeetingEntity is nil")
	}
	zoomUser, err := c.ZoomUserUsecase.GetByCond(Eq{"user_id": zoomMeetingEntity.HostId})
	if err != nil {
		c.log.Error(err)
		return "", err
	}
	if zoomUser == nil {
		return "", errors.New("zoomUser is nil")
	}
	fullName := GenFullName(zoomUser.FirstName, zoomUser.LastName)
	folderName := fmt.Sprintf("%s_%s_%s #%d", zoomMeetingEntity.StartTime, fullName, lib.TrimCharacterForFileName(zoomMeetingEntity.Topic), zoomMeetingEntity.ID)
	return folderName, nil
}

func (c *ZoomMeetingUsecase) BoxFolderId(meetingUuid string) (string, error) {
	c.BoxFolderIdLock.Lock()
	defer c.BoxFolderIdLock.Unlock()
	meeting, err := c.GetByCond(Eq{"meeting_uuid": meetingUuid, "deleted_at": 0})
	if err != nil {
		return "", err
	}
	if meeting == nil {
		return "", errors.New("meeting is nil")
	}
	if meeting.BoxResId != "" {
		return meeting.BoxResId, nil
	}
	folderName, err := c.FolderName(meeting)
	if err != nil {
		c.log.Error(err)
		return "", err
	}

	folderId, err := c.BoxUsecase.CreateFolder(folderName, ZoomBoxFolderId())
	if err != nil {
		return "", err
	}
	err = c.CommonUsecase.DB().Model(meeting).Updates(&ZoomMeetingEntity{
		BoxResId:  folderId,
		UpdatedAt: time.Now().Unix(),
	}).Error
	if err != nil {
		return "", err
	}
	return folderId, nil
}

func (c *ZoomMeetingUsecase) RenameBoxFolderName(zoomMeetingEntity *ZoomMeetingEntity) error {
	if zoomMeetingEntity == nil {
		return errors.New("zoomMeetingEntity is nil")
	}
	if zoomMeetingEntity.BoxResId == "" {
		return errors.New("zoomMeetingEntity.BoxResId is empty")
	}
	folderName, err := c.FolderName(zoomMeetingEntity)
	if err != nil {
		return err
	}
	_, err = c.BoxUsecase.UpdateFolderName(zoomMeetingEntity.BoxResId, folderName)
	if err != nil {
		return err
	}
	return nil
}
