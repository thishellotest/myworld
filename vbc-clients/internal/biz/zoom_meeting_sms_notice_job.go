package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ZoomMeetingSmsNoticeJobUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	ZoombuzUsecase    *ZoombuzUsecase
	ZoomUsecase       *ZoomUsecase
	TUsecase          *TUsecase
	MapUsecase        *MapUsecase
	DialpadbuzUsecase *DialpadbuzUsecase
	UserUsecase       *UserUsecase
	DialpadUsecase    *DialpadUsecase
	DataComboUsecase  *DataComboUsecase
	LogUsecase        *LogUsecase
	BUsaStateUsecase  *BUsaStateUsecase
	LogInfoUsecase    *LogInfoUsecase
}

func NewZoomMeetingSmsNoticeJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ZoombuzUsecase *ZoombuzUsecase,
	ZoomUsecase *ZoomUsecase,
	TUsecase *TUsecase,
	DialpadbuzUsecase *DialpadbuzUsecase,
	UserUsecase *UserUsecase,
	MapUsecase *MapUsecase,
	DialpadUsecase *DialpadUsecase,
	DataComboUsecase *DataComboUsecase,
	LogUsecase *LogUsecase,
	BUsaStateUsecase *BUsaStateUsecase,
	LogInfoUsecase *LogInfoUsecase) *ZoomMeetingSmsNoticeJobUsecase {
	uc := &ZoomMeetingSmsNoticeJobUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		ZoombuzUsecase:    ZoombuzUsecase,
		ZoomUsecase:       ZoomUsecase,
		TUsecase:          TUsecase,
		MapUsecase:        MapUsecase,
		DialpadbuzUsecase: DialpadbuzUsecase,
		UserUsecase:       UserUsecase,
		DialpadUsecase:    DialpadUsecase,
		DataComboUsecase:  DataComboUsecase,
		LogUsecase:        LogUsecase,
		BUsaStateUsecase:  BUsaStateUsecase,
		LogInfoUsecase:    LogInfoUsecase,
	}

	return uc
}

func (c *ZoomMeetingSmsNoticeJobUsecase) Run(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ZoomMeetingSmsNoticeJobUsecase:Run:Done")
				return
			default:

				// prod
				zoomUserIds := []string{"4P41WITpTYedmIqdtNFuhw", "iJG16487QZyALnTDrt7oyg"}

				//zoomUserIds := []string{"pJor6BJUThaZQSODx8BMZg"}

				for _, v := range zoomUserIds {
					err := c.Handle(v)
					if err != nil {
						c.log.Error(err)
					}
				}

				// 5分钟执行一次
				time.Sleep(3 * time.Minute)
			}
		}
	}()
	return nil
}

func (c *ZoomMeetingSmsNoticeJobUsecase) Handle(zoomUserId string) error {

	c.log.Info("ZoomMeetingSmsNoticeJobUsecase:Handle", zoomUserId)
	res, err := c.ZoombuzUsecase.ListMeetingForSmsNotice(zoomUserId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if res.GetInt("total_records") > 300 {
		c.log.Error("It's over the maximum length："+res.GetString("total_records"), " zoomUserId:", zoomUserId)
	}
	meetings := res.GetTypeList("meetings")
	now := time.Now()
	//userCaches := lib.CacheInit[*TData]()
	for _, v := range meetings {
		startTimeStr := v.GetString("start_time")
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.log.Error(err)
		} else {
			ok := c.VerifyTime(now, startTime)
			if ok {
				startTime.In(configs.GetVBCDefaultLocation())
				meeting, _, _, err := c.ZoomUsecase.GetAMeeting(v.GetInt64("id"))
				c.LogInfoUsecase.SaveLogInfo(0, "ZoomUsecase:GetAMeeting", map[string]interface{}{"meetingId": v.GetInt64("id")})

				if err != nil {
					c.log.Error(err)
				} else {
					cases, err := c.GetCasesFromMeeting(meeting)
					if err != nil {
						c.log.Error(err)
					} else {
						for k, tCase := range cases {

							triggerLogKey := c.TriggerLogKey(tCase.Id(), startTime)
							triggerLogVal, err := c.MapUsecase.GetForString(triggerLogKey)
							if err != nil {
								c.log.Error(err)
								return err
							}
							if triggerLogVal == "" {
								//tUser, err := c.UserUsecase.GetUserWithCache(userCaches, tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid))
								primaryVs := tCase.CustomFields.TextValueByNameBasic("primary_vs")
								var tUser *TData
								if primaryVs == "" {
									tUser, err = c.UserUsecase.GetByGid(tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid))
									if err != nil {
										return err
									}
								} else {
									if tCase.CustomFields.TextValueByNameBasic(FieldName_email) == "liaogling@gmail.com" {
										primaryVs = "Engineering Team"
									}
									tUser, err = c.UserUsecase.GetByFullName(primaryVs)
									if err != nil {
										return err
									}
								}

								if tUser == nil {
									c.log.Error("tUser is nil", " caseId:", tCase.Id(), " triggerLogKey:", triggerLogKey)
									continue
								}

								if err != nil {
									c.log.Error(err)
								} else {

									tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
									if err != nil {
										c.log.Error(err)
										return err
									}
									// todo:时区和发送短信
									meetingTopic := meeting.GetString("topic")
									meetingLink := meeting.GetString("join_url")
									meetingStartTime, err := c.MeetingStartTime(tClient, cases[k], startTime)
									if err != nil {
										return err
									}
									noticeText, err := c.DialpadbuzUsecase.TextZoomMeetingNotice(tClient, tUser, meetingTopic, meetingLink, meetingStartTime)
									if err != nil {
										c.log.Error(err)
									} else {

										c.log.Info("ZoomMeetingSmsNoticeJobUsecase:noticeText", strings.ReplaceAll(noticeText, "\n", "__"))
										// todo:lgl 暂时注解，先验证
										debugNotSendSms := false
										if !debugNotSendSms {
											err = c.DialpadbuzUsecase.BizSendSms(HandleSendSMSTextZoomMeetingNotice, tClient, tCase, tUser, noticeText)
										}
										if err != nil {
											//c.log.Error(err)
											// 临时关闭
											c.MapUsecase.Set(triggerLogKey, "1")
										} else {
											c.LogUsecase.SaveLog(tCase.Id(), "TextZoomMeetingNotice", map[string]interface{}{
												"meeting":                 InterfaceToString(meeting),
												"meetingTopic":            meetingTopic,
												"meetingLink":             meetingLink,
												"meetingStartTime":        meetingStartTime,
												"meetingStartTime_origin": v.GetString("start_time"),
												"meetingId":               v.GetString("id"),
												"meetingHost_email":       v.GetString("host_email"),
												"meetingHost_id":          v.GetString("host_id"),
												"meetingJoin_url":         v.GetString("join_url"),
												"meetingUuid":             v.GetString("uuid"),
												"noticeText":              noticeText,
												"caseId":                  tCase.Id(),
												"userId":                  tUser.Id(),
												"userFullName":            tUser.CustomFields.TextValueByNameBasic("full_name"),
											})
											c.MapUsecase.Set(triggerLogKey, "1")
										}
									}
								}
							}
						}
					}
				}

			}

		}
	}
	return nil
}

func (c *ZoomMeetingSmsNoticeJobUsecase) MeetingStartTime(tClient *TData, tCase *TData, startTime time.Time) (string, error) {
	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}

	state := tClient.CustomFields.TextValueByNameBasic("state")
	ut, loc, err := c.BUsaStateUsecase.GetTimeLocationByUsaState(state)
	if err != nil {
		c.log.Error("GetTimeLocationByUsaState:", tCase.Id(), " | ", tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name), " | ", err.Error(), " | ", state)
		ut = USA_TIMEZONE_PT
		loc, err = time.LoadLocation("America/Los_Angeles")
		if err != nil {
			c.log.Error(err)
			return "", err
		}
	}

	a := startTime.In(loc).Format("January 2, 2006, 03:04 PM")
	a = fmt.Sprintf("%s (%s)", a, ut)
	return a, nil
}

func (c *ZoomMeetingSmsNoticeJobUsecase) TriggerLogKey(caseId int32, startTime time.Time) string {

	return fmt.Sprintf("%s%d:%s", MapMeetingSmsNotice, caseId, startTime.UTC().Format("2006-01-02_15"))
}

func (c *ZoomMeetingSmsNoticeJobUsecase) GetCasesFromMeeting(meeting lib.TypeMap) ([]*TData, error) {
	invites := meeting.GetTypeList("settings.meeting_invitees")
	var destEmails []string

	vbcTeamEmails := []string{"imayra642@gmail.com"}

	if len(invites) >= 0 {

		for _, v := range invites {
			email := v.GetString("email")
			if strings.Index(email, "vetbenefitscenter.com") < 0 && !lib.InArray(email, vbcTeamEmails) {
				destEmails = append(destEmails, email)
			}
		}
	}
	if len(destEmails) > 0 {
		return c.TUsecase.ListByCond(Kind_client_cases,
			And(
				In("email", destEmails),
				Eq{"deleted_at": 0, "biz_deleted_at": 0},
				NotIn("stages", []string{config_vbc.Stages_Completed,
					config_vbc.Stages_Terminated,
					config_vbc.Stages_Dormant,
					config_vbc.Stages_AmCompleted,
					config_vbc.Stages_AmTerminated,
					config_vbc.Stages_AmDormant})))
	}
	return nil, nil
}

func (c *ZoomMeetingSmsNoticeJobUsecase) VerifyTime(now time.Time, startTime time.Time) bool {

	t1 := now.Add(24 * time.Hour)
	if startTime.After(now) && startTime.Before(t1) {
		return true
	} else {
		return false
	}
}
