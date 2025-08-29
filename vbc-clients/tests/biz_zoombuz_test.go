package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_ZoombuzUsecaseTest(t *testing.T) {

	userRecords, err := UT.ZoomUsecase.UsersRecordings("RhMUW5cxTICN9r6Sumn5lQ", "2024-07-01", "")
	if err != nil {
		lib.DPrintln(err)
		return
	}
	if userRecords == nil {
		return
	}
	records := lib.ToTypeMapByString(*userRecords)
	meetings := records.GetTypeList("meetings")
	err = UT.ZoombuzUsecase.SyncRecords(meetings, false)
	if err != nil {
		lib.DPrintln(err)
		return
	}
}

func Test_ZoombuzUsecase_InitSyncRecords(t *testing.T) {
	err := UT.ZoombuzUsecase.InitSyncRecords()
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_ExecuteSyncRecords(t *testing.T) {
	err := UT.ZoombuzUsecase.ExecuteSyncRecords()
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_BizExecuteSyncRecords(t *testing.T) {
	err := UT.ZoombuzUsecase.BizExecuteSyncRecords("RhMUW5cxTICN9r6Sumn5lQ")
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_SyncRecords(t *testing.T) {
	str := `[{"uuid":"UaOnHojhRZueVNCVcK8lEA==","id":87287728097,"account_id":"pzf3p9rDTmevLobo-La4HQ","host_id":"4P41WITpTYedmIqdtNFuhw","topic":"Conference with Timothy Fortson (I)","type":2,"start_time":"2024-08-16T23:27:36Z","timezone":"America/Los_Angeles","duration":44,"total_size":196200715,"recording_count":4,"share_url":"https://us06web.zoom.us/rec/share/lNw00vo7xyC4wsIWNz3cwhvZ7BSVPU2fJNrCj7FitAQ0_nKjvqR9dOPFBCa6_LrZ.4txU15rydT86hB4Z","recording_files":[{"id":"65f52c07-55da-419d-948b-ab75f55b2b99","meeting_id":"UaOnHojhRZueVNCVcK8lEA==","recording_start":"2024-08-16T23:27:36Z","recording_end":"2024-08-17T00:12:00Z","file_type":"TRANSCRIPT","file_extension":"VTT","file_size":45865,"play_url":"https://us06web.zoom.us/rec/play/wCYiob-8WZGGc8hH0d1SBPRe1ds1GUIkAhAtz_LQVn1TnGUwNtMFuXp08xYWJiD-vnpR1RnF2F3MyHNE.WxMEJVFSFOO8Tb0k","download_url":"https://us06web.zoom.us/rec/download/wCYiob-8WZGGc8hH0d1SBPRe1ds1GUIkAhAtz_LQVn1TnGUwNtMFuXp08xYWJiD-vnpR1RnF2F3MyHNE.WxMEJVFSFOO8Tb0k","status":"completed","recording_type":"audio_transcript"},{"id":"663f4f9c-761f-49a9-bf4a-fabc42785908","meeting_id":"UaOnHojhRZueVNCVcK8lEA==","recording_start":"2024-08-16T23:27:36Z","recording_end":"2024-08-17T00:12:00Z","file_type":"SUMMARY","file_extension":"JSON","file_size":5850,"play_url":"https://us06web.zoom.us/rec/play/BAWqAFhog5Va6sjIT9cKiZyQdEkb804eGRcHxkxoMHWk2CSgyhi2dWL51pNGHwXpD9D6JjzJi4pcD_mZ.vTp-23AOrcitHgW5","download_url":"https://us06web.zoom.us/rec/download/BAWqAFhog5Va6sjIT9cKiZyQdEkb804eGRcHxkxoMHWk2CSgyhi2dWL51pNGHwXpD9D6JjzJi4pcD_mZ.vTp-23AOrcitHgW5","status":"completed","recording_type":"summary"},{"id":"6868bff4-09cf-4f8c-8f1a-b0232cd3884c","meeting_id":"UaOnHojhRZueVNCVcK8lEA==","recording_start":"2024-08-16T23:27:36Z","recording_end":"2024-08-17T00:12:00Z","file_type":"M4A","file_extension":"M4A","file_size":42392752,"play_url":"https://us06web.zoom.us/rec/play/r5ZrfmKcofwOHPZBVCAuYbsnkDvTct9RspowYBUyOvgn8-hiQ4q0ITaVkG0TUCdOJA0j1HWDXblEc25u.14w9rYTtM5JyoVZI","download_url":"https://us06web.zoom.us/rec/download/r5ZrfmKcofwOHPZBVCAuYbsnkDvTct9RspowYBUyOvgn8-hiQ4q0ITaVkG0TUCdOJA0j1HWDXblEc25u.14w9rYTtM5JyoVZI","status":"completed","recording_type":"audio_only"},{"id":"1df78c0d-8039-4800-8a6a-818e1e60271c","meeting_id":"UaOnHojhRZueVNCVcK8lEA==","recording_start":"2024-08-16T23:27:36Z","recording_end":"2024-08-17T00:12:00Z","file_type":"TIMELINE","file_extension":"JSON","file_size":995956,"download_url":"https://us06web.zoom.us/rec/download/AA-47-OwOvWr0Gfo1_FY6v7tx5g0WNDTTSNof5uJATKMPdMufMIGUjSHx7E5zTccjzOf2VOYU2eQhVKX.CukzJC53CoTOYuZ6","status":"completed","recording_type":"timeline"},{"id":"766f6577-0e86-4b1d-92e7-ed62924d9cdf","meeting_id":"UaOnHojhRZueVNCVcK8lEA==","recording_start":"2024-08-16T23:27:36Z","recording_end":"2024-08-17T00:12:00Z","file_type":"MP4","file_extension":"MP4","file_size":152766142,"play_url":"https://us06web.zoom.us/rec/play/nU3BxKGSLL9i7lxIGbhk6aiNImtgus_90QAkM8U3fTBIMRce9IhavBEdlMw2Wx4EavPkBzreooswnEtV.eI9M3Pv66Rm3Odri","download_url":"https://us06web.zoom.us/rec/download/nU3BxKGSLL9i7lxIGbhk6aiNImtgus_90QAkM8U3fTBIMRce9IhavBEdlMw2Wx4EavPkBzreooswnEtV.eI9M3Pv66Rm3Odri","status":"completed","recording_type":"shared_screen_with_speaker_view"},{"id":"d67844c9-168e-4925-97b4-1598f66cd0ab","meeting_id":"UaOnHojhRZueVNCVcK8lEA==","recording_start":"2024-08-16T23:27:36Z","recording_end":"2024-08-17T00:12:00Z","file_type":"SUMMARY","file_extension":"JSON","file_size":2953,"play_url":"https://us06web.zoom.us/rec/play/JaNcYzci4V2rDt7LL_c41z0_lt2oUBTLI95mL4I-dGlw3WK4z1RZnxAqV3jVrSqH7t_y1My4uAZZqO8z.aLC0Iar0-95aJvIe","download_url":"https://us06web.zoom.us/rec/download/JaNcYzci4V2rDt7LL_c41z0_lt2oUBTLI95mL4I-dGlw3WK4z1RZnxAqV3jVrSqH7t_y1My4uAZZqO8z.aLC0Iar0-95aJvIe","status":"completed","recording_type":"summary_next_steps"}],"recording_play_passcode":"YVeo6PyI40urDPhsjhmXcoslqmwhFz-O"}]`
	list := lib.ToTypeListByString(str)
	err := UT.ZoombuzUsecase.SyncRecords(list, false)
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_ExecuteSyncZoomUsers(t *testing.T) {
	err := UT.ZoombuzUsecase.ExecuteSyncZoomUsers()
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_SyncZoomUsers(t *testing.T) {
	str := `[{"id":"4P41WITpTYedmIqdtNFuhw","first_name":"Edward","last_name":"Bunting","display_name":"Edward Bunting","email":"ebunting@vetbenefitscenter.com","type":2,"pmi":9143783503,"timezone":"America/Los_Angeles","verified":0,"created_at":"2023-12-08T18:52:10Z","last_login_time":"2024-08-16T01:37:16Z","last_client_version":"6.1.6.23395(android)","pic_url":"https://us06web.zoom.us/p/v2/2c1ef36ccef9f0dde85408eb311f9f4556c560934a2fa42997b02bbb61e4e72d/9aca98ee-3364-44a1-9e0d-4f1608c63236-658","language":"en-US","phone_number":"+1 6198005543","status":"active","role_id":"0","user_created_at":"2023-01-26T04:27:29Z"},{"id":"anohPS26QqevqpEkxPX3Bw","first_name":"Victoria","last_name":"E.","display_name":"Victoria E.","email":"venriquez@vetbenefitscenter.com","type":2,"pmi":4538768334,"timezone":"America/Los_Angeles","verified":1,"created_at":"2024-01-11T23:52:29Z","last_login_time":"2024-08-15T23:24:07Z","last_client_version":"6.0.11.39959(win)","pic_url":"https://us06web.zoom.us/p/v2/d0d72c93e79deac0ec3ca81547093c34cb9998cb4dd8772be4d81ac62b540626/0cdea746-b0fd-46af-a09e-af4991d0c0fb-1703","language":"en-US","phone_number":"","status":"active","role_id":"2","user_created_at":"2024-01-11T23:52:29Z"},{"id":"iJG16487QZyALnTDrt7oyg","first_name":"Donald","last_name":"Pratko","display_name":"Donald Pratko","email":"dpratko@vetbenefitscenter.com","type":2,"pmi":5025436184,"timezone":"America/New_York","verified":0,"dept":"Veteran Services","created_at":"2024-08-12T22:41:54Z","last_login_time":"2024-08-15T21:55:47Z","last_client_version":"6.1.7.17235(iphone)","pic_url":"https://us06web.zoom.us/p/v2/757fa068903d86a573d0d582c39b876bce6f79cb6fbc8fa71114cb80a11d5493/dd8175f7-adf2-4ca6-bf8f-70a4aa45360a-2191","language":"en-US","phone_number":"+1 8652631500","status":"active","role_id":"2","user_created_at":"2024-03-18T17:12:32Z"},{"id":"pJor6BJUThaZQSODx8BMZg","first_name":"Yannan","last_name":"Wang","display_name":"Yannan Wang","email":"engineering@vetbenefitscenter.com","type":2,"pmi":5188796123,"timezone":"America/Los_Angeles","verified":1,"dept":"","created_at":"2024-08-15T15:04:28Z","last_login_time":"2024-08-16T02:43:05Z","pic_url":"https://us06web.zoom.us/p/v2/e6800f683a0d6e482c88e9248a4f3a6fd0a90af7cde90d40ceb7a0472419e985/7a02f234-17da-44dd-a883-c76edce437be-6794","language":"en-US","phone_number":"","status":"active","role_id":"1","user_created_at":"2024-06-25T02:44:28Z"},{"id":"RhMUW5cxTICN9r6Sumn5lQ","first_name":"Andrea","last_name":"Ladd","display_name":"Andrea Ladd","email":"aladd@vetbenefitscenter.com","type":2,"pmi":7045919919,"timezone":"America/New_York","verified":1,"created_at":"2024-07-23T00:29:42Z","last_login_time":"2024-08-15T18:12:19Z","last_client_version":"6.1.7.17235(iphone)","pic_url":"https://us06web.zoom.us/p/v2/89b6354b011808d8e179c86d3114824f333af45bacbcbfa4784dfdf00910fb8f/7ef9a19d-4aa5-48c1-950a-315bb396a389-5379","language":"en-US","status":"active","role_id":"2","user_created_at":"2024-07-23T00:27:28Z"}]`
	list := lib.ToTypeListByString(str)
	err := UT.ZoombuzUsecase.SyncZoomUsers(list)
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_UpdateZoomRecordingFile(t *testing.T) {
	entity, _ := UT.ZoomRecordingFileUsecase.GetByCond(builder.Eq{"id": 1})
	err := UT.ZoombuzUsecase.UpdateZoomRecordingFile(entity)
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_Meeting(t *testing.T) {
	meetings := make(map[string]biz.HttpResponseBody)
	ccc, err := UT.ZoombuzUsecase.Meeting(meetings, "m7jrIhUQQJ6XckqNHbypog==")
	lib.DPrintln(err, ccc)
	ccc, err = UT.ZoombuzUsecase.Meeting(meetings, "m7jrIhUQQJ6XckqNHbypog==")
	lib.DPrintln(err, ccc)

	ccc, err = UT.ZoombuzUsecase.Meeting(meetings, "m7jrIhUQQJ6XckqNHbypog==1")
	lib.DPrintln(err, ccc)
	ccc, err = UT.ZoombuzUsecase.Meeting(meetings, "m7jrIhUQQJ6XckqNHbypog==1")
	lib.DPrintln(err, ccc)
}

func Test_ZoombuzUsecase_ExecuteDeleteMeetingRecording(t *testing.T) {
	err := UT.ZoombuzUsecase.ExecuteDeleteMeetingRecording()
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_DeleteMeetingRecording(t *testing.T) {
	err := UT.ZoombuzUsecase.DeleteMeetingRecording("QZskvwPQQwCLXF3SnY+sTw==")
	lib.DPrintln(err)
}

func Test_ZoombuzUsecase_ListMeetingForSmsNotice(t *testing.T) {
	c, err := UT.ZoombuzUsecase.ListMeetingForSmsNotice("pJor6BJUThaZQSODx8BMZg")
	lib.DPrintln(c, err)
}
