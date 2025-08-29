package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"net/url"
	"vbc/internal/conf"
	"vbc/lib"
)

const Zoom_API_HOST = "https://api.zoom.us"

/*
Zoom汇总：
1. 2024-08-21T03:00:09Z response status code 400: 此错误有可能是Zoom的Recording出错了，该uuid meeting只返回3个记录：mrZDW/30SEme7qPSrf0ZtQ==， 都不能下载
2. 2024-08-20T17:01:00Z response status code 404
- 情况一：重试就可以
- 情况二：uuid meeting已经删除了
*/

type ZoomUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	Oauth2TokenUsecase *Oauth2TokenUsecase
}

func NewZoomUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	Oauth2TokenUsecase *Oauth2TokenUsecase) *ZoomUsecase {
	uc := &ZoomUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		Oauth2TokenUsecase: Oauth2TokenUsecase,
	}
	return uc
}

func (c *ZoomUsecase) Headers() (map[string]string, error) {

	token, err := c.Oauth2TokenUsecase.GetAccessToken(Oauth2_AppId_zoom)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Authorization": "Bearer " + token}, nil
}

// GetAMeeting {"agenda":"","assistant_id":"","created_at":"2024-09-07T02:41:27Z","duration":30,"encrypted_password":"4Dz9wpAuubvNKdq1nYSyw8l4uf1cSe.1","h323_password":"184942","host_email":"engineering@vetbenefitscenter.com","host_id":"pJor6BJUThaZQSODx8BMZg","id":83530941612,"join_url":"https://us06web.zoom.us/j/5188796123?pwd=4Dz9wpAuubvNKdq1nYSyw8l4uf1cSe.1\u0026omn=83530941612","password":"184942","pmi":"5188796123","pre_schedule":false,"pstn_password":"184942","settings":{"allow_multiple_devices":false,"alternative_host_update_polls":false,"alternative_hosts":"","alternative_hosts_email_notification":true,"approval_type":2,"approved_or_denied_countries_or_regions":{"enable":false},"audio":"both","auto_recording":"none","auto_start_ai_companion_questions":false,"auto_start_meeting_summary":true,"breakout_room":{"enable":false},"close_registration":false,"cn_meeting":false,"continuous_meeting_chat":{"auto_add_invited_external_users":false,"enable":false},"device_testing":false,"email_in_attendee_report":false,"email_notification":true,"enable_dedicated_group_chat":false,"encryption_type":"enhanced_encryption","enforce_login":false,"enforce_login_domains":"","focus_mode":false,"global_dial_in_countries":["US"],"global_dial_in_numbers":[{"city":"New York","country":"US","country_name":"US","number":"+1 646 558 8656","type":"toll"},{"country":"US","country_name":"US","number":"+1 646 931 3860","type":"toll"},{"country":"US","country_name":"US","number":"+1 669 444 9171","type":"toll"},{"country":"US","country_name":"US","number":"+1 689 278 1000","type":"toll"},{"country":"US","country_name":"US","number":"+1 719 359 4580","type":"toll"},{"city":"Denver","country":"US","country_name":"US","number":"+1 720 707 2699","type":"toll"},{"country":"US","country_name":"US","number":"+1 253 205 0468","type":"toll"},{"city":"Tacoma","country":"US","country_name":"US","number":"+1 253 215 8782","type":"toll"},{"city":"Washington DC","country":"US","country_name":"US","number":"+1 301 715 8592","type":"toll"},{"country":"US","country_name":"US","number":"+1 305 224 1968","type":"toll"},{"country":"US","country_name":"US","number":"+1 309 205 3325","type":"toll"},{"city":"Chicago","country":"US","country_name":"US","number":"+1 312 626 6799","type":"toll"},{"city":"Houston","country":"US","country_name":"US","number":"+1 346 248 7799","type":"toll"},{"country":"US","country_name":"US","number":"+1 360 209 5623","type":"toll"},{"country":"US","country_name":"US","number":"+1 386 347 5053","type":"toll"},{"country":"US","country_name":"US","number":"+1 507 473 4847","type":"toll"},{"country":"US","country_name":"US","number":"+1 564 217 2000","type":"toll"}],"host_save_video_order":false,"host_video":true,"in_meeting":false,"internal_meeting":false,"jbh_time":0,"join_before_host":false,"meeting_authentication":false,"meeting_invitees":[{"email":"gengling.liao@hotmail.com"},{"email":"lialing@foxmail.com"},{"email":"liaogling@gmail.com"}],"mute_upon_entry":false,"participant_focused_meeting":false,"participant_video":false,"private_meeting":false,"push_change_to_calendar":false,"registrants_confirmation_email":true,"registrants_email_notification":true,"request_permission_to_unmute_participants":false,"resources":[],"show_join_info":false,"show_share_button":false,"sign_language_interpretation":{"enable":false},"use_pmi":true,"waiting_room":false,"watermark":false},"start_time":"2024-09-10T03:00:00Z","start_url":"https://us06web.zoom.us/s/5188796123?zak=eyJ0eXAiOiJKV1QiLCJzdiI6IjAwMDAwMSIsInptX3NrbSI6InptX28ybSIsImFsZyI6IkhTMjU2In0.eyJpc3MiOiJ3ZWIiLCJjbHQiOjAsIm1udW0iOiI1MTg4Nzk2MTIzIiwiYXVkIjoiY2xpZW50c20iLCJ1aWQiOiJwSm9yNkJKVVRoYVpRU09EeDhCTVpnIiwiemlkIjoiMDJjZTEwZTkxN2M4NDIxZDlhODFhZjFlYTJjMDFiNzYiLCJzayI6IjQ3NTU5MzIwOTk4MTI1MDc4MTMiLCJzdHkiOjEwMCwid2NkIjoidXMwNiIsImV4cCI6MTcyNTk1MDAzMywiaWF0IjoxNzI1OTQyODMzLCJhaWQiOiJwemYzcDlyRFRtZXZMb2JvLUxhNEhRIiwiY2lkIjoiIn0.7qdlvJc4DPsDOaSsFPm5psW3w_QcoItiP09CgScPdVs","status":"waiting","timezone":"Asia/Shanghai","topic":"Zoom meeting invitation - Yannan Wang的Zoom会议-个人会议号02","type":2,"uuid":"ewFXquTBRG+/ic60Mpq00g=="}
func (c *ZoomUsecase) GetAMeeting(meetingId int64) (body lib.TypeMap, rawBody string, httpCode int, err error) {

	api := fmt.Sprintf("%s/v2/meetings/%d", Zoom_API_HOST, meetingId)

	//fmt.Println("api:", api)
	headers, err := c.Headers()
	if err != nil {
		return nil, "", 0, err
	}
	res, httpCode, err := lib.Request("GET", api, nil, headers)
	if res != nil {
		rawBody = *res
	}
	if err != nil {
		return nil, rawBody, httpCode, err
	}
	if res == nil {
		return nil, rawBody, httpCode, errors.New("res is nil")
	}
	return lib.ToTypeMapByString(rawBody), rawBody, httpCode, nil
}

func (c *ZoomUsecase) MeetingRecordings(meetingId string) (body lib.TypeMap, rawBody string, httpCode int, err error) {

	api := fmt.Sprintf("%s/v2/meetings/%s/recordings", Zoom_API_HOST, meetingId)
	headers, err := c.Headers()
	if err != nil {
		return nil, "", 0, err
	}
	res, httpCode, err := lib.Request("GET", api, nil, headers)
	if res != nil {
		rawBody = *res
	}
	if err != nil {
		return nil, rawBody, httpCode, err
	}
	if res == nil {
		return nil, rawBody, httpCode, errors.New("res is nil")
	}
	return lib.ToTypeMapByString(rawBody), rawBody, httpCode, nil
}

func (c *ZoomUsecase) DeleteMeetingRecordings(meetingId string) (body lib.TypeMap, rawBody string, httpCode int, err error) {

	api := fmt.Sprintf("%s/v2/meetings/%s/recordings", Zoom_API_HOST, meetingId)
	headers, err := c.Headers()
	if err != nil {
		return nil, "", 0, err
	}
	res, httpCode, err := lib.Request("DELETE", api, nil, headers)
	if res != nil {
		rawBody = *res
	}
	if err != nil {
		return nil, rawBody, httpCode, err
	}
	if res == nil {
		return nil, rawBody, httpCode, errors.New("res is nil")
	}
	return lib.ToTypeMapByString(rawBody), rawBody, httpCode, nil
}

func (c *ZoomUsecase) UsersMeetings(zoomUserId string, fromDate string, toDate string) (body lib.TypeMap, rawBody string, httpCode int, err error) {

	// pJor6BJUThaZQSODx8BMZg YW
	// 4P41WITpTYedmIqdtNFuhw ED
	api := fmt.Sprintf("%s/v2/users/%s/meetings", Zoom_API_HOST, zoomUserId)
	headers, err := c.Headers()
	if err != nil {
		return nil, "", 0, err
	}
	query := make(url.Values)
	//query.Add("from", "2024-07-14")
	query.Add("from", fromDate) // 2024-01-01
	query.Add("page_size", "300")
	if toDate != "" {
		query.Add("to", toDate)
	}
	//t := time.Now()
	//query.Add("to", t.Format(time.DateOnly))
	res, httpCode, err := lib.RequestGet(api, query, headers)
	if err != nil {
		return nil, "", httpCode, err
	}
	if res == nil {
		return nil, "", httpCode, errors.New("res is nil")
	}
	return lib.ToTypeMapByString(*res), *res, httpCode, nil
}

func (c *ZoomUsecase) UsersRecordings(zoomUserId string, fromDate string, toDate string) (*string, error) {

	// 4P41WITpTYedmIqdtNFuhw ED
	// pJor6BJUThaZQSODx8BMZg

	api := fmt.Sprintf("%s/v2/users/%s/recordings", Zoom_API_HOST, zoomUserId)
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}
	query := make(url.Values)
	query.Add("from", fromDate) // 2024-01-01
	query.Add("page_size", "300")
	if toDate != "" {
		query.Add("to", toDate)
	}
	res, _, err := lib.RequestGet(api, query, headers)
	return res, err
}

func (c *ZoomUsecase) User() (res *string, err error) {
	api := fmt.Sprintf("%s/v2/users", Zoom_API_HOST)
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}
	query := make(url.Values)
	query.Add("page_size", "300")

	res, _, err = lib.RequestWithQuery("GET", api, query, nil, headers)
	return res, err
	/*


		{
			"page_count": 1,
			"page_number": 1,
			"page_size": 30,
			"total_records": 5,
			"next_page_token": "",
			"users": [{
				"id": "4P41WITpTYedmIqdtNFuhw",
				"first_name": "Edward",
				"last_name": "Bunting",
				"display_name": "Edward Bunting",
				"email": "ebunting@vetbenefitscenter.com",
				"type": 2,
				"pmi": 9143783503,
				"timezone": "America/Los_Angeles",
				"verified": 0,
				"created_at": "2023-12-08T18:52:10Z",
				"last_login_time": "2024-08-16T01:37:16Z",
				"last_client_version": "6.1.6.23395(android)",
				"pic_url": "https://us06web.zoom.us/p/v2/2c1ef36ccef9f0dde85408eb311f9f4556c560934a2fa42997b02bbb61e4e72d/9aca98ee-3364-44a1-9e0d-4f1608c63236-658",
				"language": "en-US",
				"phone_number": "+1 6198005543",
				"status": "active",
				"role_id": "0",
				"user_created_at": "2023-01-26T04:27:29Z"
			}, {
				"id": "anohPS26QqevqpEkxPX3Bw",
				"first_name": "Victoria",
				"last_name": "E.",
				"display_name": "Victoria E.",
				"email": "venriquez@vetbenefitscenter.com",
				"type": 2,
				"pmi": 4538768334,
				"timezone": "America/Los_Angeles",
				"verified": 1,
				"created_at": "2024-01-11T23:52:29Z",
				"last_login_time": "2024-08-15T23:24:07Z",
				"last_client_version": "6.0.11.39959(win)",
				"pic_url": "https://us06web.zoom.us/p/v2/d0d72c93e79deac0ec3ca81547093c34cb9998cb4dd8772be4d81ac62b540626/0cdea746-b0fd-46af-a09e-af4991d0c0fb-1703",
				"language": "en-US",
				"phone_number": "",
				"status": "active",
				"role_id": "2",
				"user_created_at": "2024-01-11T23:52:29Z"
			}, {
				"id": "iJG16487QZyALnTDrt7oyg",
				"first_name": "Donald",
				"last_name": "Pratko",
				"display_name": "Donald Pratko",
				"email": "dpratko@vetbenefitscenter.com",
				"type": 2,
				"pmi": 5025436184,
				"timezone": "America/New_York",
				"verified": 0,
				"dept": "Veteran Services",
				"created_at": "2024-08-12T22:41:54Z",
				"last_login_time": "2024-08-15T21:55:47Z",
				"last_client_version": "6.1.7.17235(iphone)",
				"pic_url": "https://us06web.zoom.us/p/v2/757fa068903d86a573d0d582c39b876bce6f79cb6fbc8fa71114cb80a11d5493/dd8175f7-adf2-4ca6-bf8f-70a4aa45360a-2191",
				"language": "en-US",
				"phone_number": "+1 8652631500",
				"status": "active",
				"role_id": "2",
				"user_created_at": "2024-03-18T17:12:32Z"
			}, {
				"id": "pJor6BJUThaZQSODx8BMZg",
				"first_name": "Yannan",
				"last_name": "Wang",
				"display_name": "Yannan Wang",
				"email": "engineering@vetbenefitscenter.com",
				"type": 2,
				"pmi": 5188796123,
				"timezone": "America/Los_Angeles",
				"verified": 1,
				"dept": "",
				"created_at": "2024-08-15T15:04:28Z",
				"last_login_time": "2024-08-16T02:43:05Z",
				"pic_url": "https://us06web.zoom.us/p/v2/e6800f683a0d6e482c88e9248a4f3a6fd0a90af7cde90d40ceb7a0472419e985/7a02f234-17da-44dd-a883-c76edce437be-6794",
				"language": "en-US",
				"phone_number": "",
				"status": "active",
				"role_id": "1",
				"user_created_at": "2024-06-25T02:44:28Z"
			}, {
				"id": "RhMUW5cxTICN9r6Sumn5lQ",
				"first_name": "Andrea",
				"last_name": "Ladd",
				"display_name": "Andrea Ladd",
				"email": "aladd@vetbenefitscenter.com",
				"type": 2,
				"pmi": 7045919919,
				"timezone": "America/New_York",
				"verified": 1,
				"created_at": "2024-07-23T00:29:42Z",
				"last_login_time": "2024-08-15T18:12:19Z",
				"last_client_version": "6.1.7.17235(iphone)",
				"pic_url": "https://us06web.zoom.us/p/v2/89b6354b011808d8e179c86d3114824f333af45bacbcbfa4784dfdf00910fb8f/7ef9a19d-4aa5-48c1-950a-315bb396a389-5379",
				"language": "en-US",
				"status": "active",
				"role_id": "2",
				"user_created_at": "2024-07-23T00:27:28Z"
			}]
		}

	*/
}
