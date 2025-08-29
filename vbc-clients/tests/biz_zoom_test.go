package tests

import (
	"testing"
	"time"
	"vbc/lib"
)

func Test_ZoomUsecase_GetAMeeting(t *testing.T) {
	a, _, _, err := UT.ZoomUsecase.GetAMeeting(83530941612)
	lib.DPrintln(a)
	lib.DPrintln(err)
}

func Test_ZoomUsecase_MeetingRecordings(t *testing.T) {
	meetingId := "m7jrIhUQQJ6XckqNHbypog=="
	meetingId = "/3Cw3vEaSSWG4v8VFBoRlw=="

	params, _, code, err := UT.ZoomUsecase.MeetingRecordings(meetingId)
	lib.DPrintln(err, code)
	lib.DPrintln(params)
}

func Test_ZoomUsecase_DMeetingRecordings(t *testing.T) {
	meetingId := "xfz/CgDVQ3WpiIyifQFyVQ=="

	params, _, code, err := UT.ZoomUsecase.DeleteMeetingRecordings(meetingId)
	// code 返回204说明删除成功
	lib.DPrintln(err, code)
	lib.DPrintln(params)
}

func Test_ZoomUsecase_UsersRecordings(t *testing.T) {

	ti := time.Now()
	ti = ti.AddDate(0, 0, -5)

	userRecords, err := UT.ZoomUsecase.UsersRecordings("4P41WITpTYedmIqdtNFuhw", ti.Format(time.DateOnly), "")

	lib.DPrintln(userRecords, err)
}

func Test_ZoomUsecase_UsersMeetings(t *testing.T) {
	now := time.Now()
	beginNow := now.AddDate(0, 0, -3)
	endNow := now.AddDate(0, 0, +5)

	a, b, c, err := UT.ZoomUsecase.UsersMeetings("4P41WITpTYedmIqdtNFuhw", beginNow.Format(time.DateOnly), endNow.Format(time.DateOnly))
	lib.DPrintln(err)
	lib.DPrintln(a)
	lib.DPrintln(b)
	lib.DPrintln(c)
}

func Test_ZoomUsecase_User(t *testing.T) {
	res, err := UT.ZoomUsecase.User()
	lib.DPrintln(err)
	lib.DPrintln(res)
}
