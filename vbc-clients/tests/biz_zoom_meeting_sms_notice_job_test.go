package tests

import (
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_ZoomMeetingSmsNoticeJobUsecase_Hanlde(t *testing.T) {

	// YN: pJor6BJUThaZQSODx8BMZg
	zoomUserIds := []string{"4P41WITpTYedmIqdtNFuhw", "iJG16487QZyALnTDrt7oyg"}
	for _, v := range zoomUserIds {
		err := UT.ZoomMeetingSmsNoticeJobUsecase.Handle(v)
		lib.DPrintln(err)
	}
}

func Test_ZoomMeetingSmsNoticeJobUsecase_MeetingStartTime(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5217)
	if err != nil {
		panic(err)
	}

	tClient, err := UT.TUsecase.Data(biz.Kind_clients, builder.Eq{"id": 5120})
	if err != nil {
		lib.DPrintln(err)
	} else {
		now := time.Now()
		a, err := UT.ZoomMeetingSmsNoticeJobUsecase.MeetingStartTime(tClient, tCase, now)
		lib.DPrintln(a, err)
	}
}

func Test_ZoomMeetingSmsNoticeJobUsecase_VerifyTime(t *testing.T) {
	now := time.Now()
	ti, _ := time.Parse(time.RFC3339, "2024-09-10T10:30:00Z")
	f := UT.ZoomMeetingSmsNoticeJobUsecase.VerifyTime(now, ti)
	lib.DPrintln(f)
}

func Test_ZoomMeetingSmsNoticeJobUsecase_TriggerLogKey(t *testing.T) {
	key := UT.ZoomMeetingSmsNoticeJobUsecase.TriggerLogKey(11, time.Now())
	lib.DPrintln(key)
}

func Test_ZoomMeetingSmsNoticeJobUsecase_GetCasesFromMeeting(t *testing.T) {
	str := `{"agenda":"","assistant_id":"","created_at":"2024-09-07T02:41:27Z","duration":30,"encrypted_password":"4Dz9wpAuubvNKdq1nYSyw8l4uf1cSe.1","h323_password":"184942","host_email":"engineering@vetbenefitscenter.com","host_id":"pJor6BJUThaZQSODx8BMZg","id":83530941612,"join_url":"https://us06web.zoom.us/j/5188796123?pwd=4Dz9wpAuubvNKdq1nYSyw8l4uf1cSe.1\u0026omn=83530941612","password":"184942","pmi":"5188796123","pre_schedule":false,"pstn_password":"184942","settings":{"allow_multiple_devices":false,"alternative_host_update_polls":false,"alternative_hosts":"","alternative_hosts_email_notification":true,"approval_type":2,"approved_or_denied_countries_or_regions":{"enable":false},"audio":"both","auto_recording":"none","auto_start_ai_companion_questions":false,"auto_start_meeting_summary":true,"breakout_room":{"enable":false},"close_registration":false,"cn_meeting":false,"continuous_meeting_chat":{"auto_add_invited_external_users":false,"enable":false},"device_testing":false,"email_in_attendee_report":false,"email_notification":true,"enable_dedicated_group_chat":false,"encryption_type":"enhanced_encryption","enforce_login":false,"enforce_login_domains":"","focus_mode":false,"global_dial_in_countries":["US"],"global_dial_in_numbers":[{"city":"New York","country":"US","country_name":"US","number":"+1 646 558 8656","type":"toll"},{"country":"US","country_name":"US","number":"+1 646 931 3860","type":"toll"},{"country":"US","country_name":"US","number":"+1 669 444 9171","type":"toll"},{"country":"US","country_name":"US","number":"+1 689 278 1000","type":"toll"},{"country":"US","country_name":"US","number":"+1 719 359 4580","type":"toll"},{"city":"Denver","country":"US","country_name":"US","number":"+1 720 707 2699","type":"toll"},{"country":"US","country_name":"US","number":"+1 253 205 0468","type":"toll"},{"city":"Tacoma","country":"US","country_name":"US","number":"+1 253 215 8782","type":"toll"},{"city":"Washington DC","country":"US","country_name":"US","number":"+1 301 715 8592","type":"toll"},{"country":"US","country_name":"US","number":"+1 305 224 1968","type":"toll"},{"country":"US","country_name":"US","number":"+1 309 205 3325","type":"toll"},{"city":"Chicago","country":"US","country_name":"US","number":"+1 312 626 6799","type":"toll"},{"city":"Houston","country":"US","country_name":"US","number":"+1 346 248 7799","type":"toll"},{"country":"US","country_name":"US","number":"+1 360 209 5623","type":"toll"},{"country":"US","country_name":"US","number":"+1 386 347 5053","type":"toll"},{"country":"US","country_name":"US","number":"+1 507 473 4847","type":"toll"},{"country":"US","country_name":"US","number":"+1 564 217 2000","type":"toll"}],"host_save_video_order":false,"host_video":true,"in_meeting":false,"internal_meeting":false,"jbh_time":0,"join_before_host":false,"meeting_authentication":false,"meeting_invitees":[{"email":"gengling.liao@hotmail.com"},{"email":"lialing@foxmail.com"},{"email":"liaogling@gmail.com"}],"mute_upon_entry":false,"participant_focused_meeting":false,"participant_video":false,"private_meeting":false,"push_change_to_calendar":false,"registrants_confirmation_email":true,"registrants_email_notification":true,"request_permission_to_unmute_participants":false,"resources":[],"show_join_info":false,"show_share_button":false,"sign_language_interpretation":{"enable":false},"use_pmi":true,"waiting_room":false,"watermark":false},"start_time":"2024-09-10T03:00:00Z","start_url":"https://us06web.zoom.us/s/5188796123?zak=eyJ0eXAiOiJKV1QiLCJzdiI6IjAwMDAwMSIsInptX3NrbSI6InptX28ybSIsImFsZyI6IkhTMjU2In0.eyJpc3MiOiJ3ZWIiLCJjbHQiOjAsIm1udW0iOiI1MTg4Nzk2MTIzIiwiYXVkIjoiY2xpZW50c20iLCJ1aWQiOiJwSm9yNkJKVVRoYVpRU09EeDhCTVpnIiwiemlkIjoiNDdjOGM1MjZmNWVkNGI3OTg3NmNhZTc4ZTFiMmEyM2QiLCJzayI6IjQ3NTU5MzIwOTk4MTI1MDc4MTMiLCJzdHkiOjEwMCwid2NkIjoidXMwNiIsImV4cCI6MTcyNTg4NjgyNiwiaWF0IjoxNzI1ODc5NjI2LCJhaWQiOiJwemYzcDlyRFRtZXZMb2JvLUxhNEhRIiwiY2lkIjoiIn0.GrWGm8LrCvD288aoRsYAAgqlnxihOA4qcvhId6UWTmA","status":"waiting","timezone":"Asia/Shanghai","topic":"Zoom meeting invitation - Yannan Wang的Zoom会议-个人会议号02","type":2,"uuid":"ewFXquTBRG+/ic60Mpq00g=="}`
	s := lib.ToTypeMapByString(str)
	cases, err := UT.ZoomMeetingSmsNoticeJobUsecase.GetCasesFromMeeting(s)
	for _, v := range cases {
		lib.DPrintln(v.Id(), v.Gid())
	}
	lib.DPrintln(err)
	lib.DPrintln(cases)
}
