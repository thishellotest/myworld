package tests

import (
	"testing"
	"vbc/lib"
)

func Test_UnsubscribesbuzUsecase_HandleFromDialpadWebhookEvent(t *testing.T) {

	str := `{"contact":{"id":6230136767135744,"name":"Test - VBC","phone_number":"+16192000000"},"created_date":1738808741123,"direction":"inbound","event_timestamp":1738808741527,"from_number":"+16192788886","id":5594630430703616,"is_internal":false,"message_delivery_result":null,"message_status":"pending","mms":false,"mms_url":null,"sender_id":null,"target":{"id":4693435745845248,"name":"Edward Bunting Jr.","phone_number":"(619) 800-0000","type":"user"},"text":" sTOp","text_content":"unsTOp","to_number":["+16198005543"]}`
	data := lib.ToTypeMapByString(str)
	err := UT.UnsubscribesbuzUsecase.HandleFromDialpadWebhookEvent(data)
	lib.DPrintln(err)
}

func Test_UnsubscribesbuzUsecase_NotifyAdmin(t *testing.T) {
	err := UT.UnsubscribesbuzUsecase.NotifyAdmin("+18056604465", "stop")
	lib.DPrintln(err)
}

func Test_UnsubscribesbuzUsecase_Upsert(t *testing.T) {
	err := UT.UnsubscribesbuzUsecase.Upsert("", "+18056604465", "unstop")
	lib.DPrintln(err)
}
