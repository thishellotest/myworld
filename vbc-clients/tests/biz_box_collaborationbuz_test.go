package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_BoxCollaborationBuzUsecase_HandleAddCollaboration(t *testing.T) {
	// "327431320500", "41426608287"
	//  VBC Team Box User ID: 30888625898  报错：{"type":"error","status":400,"code":"Cannot invite self as a collaborator","help_url":"http:\/\/developers.box.com\/docs\/#errors","message":"Cannot invite self as a collaborator","request_id":"y14fv6i4cadr7w75"}
	//  Normal Box User ID: 39217319801
	// Lili Box User ID：32499237039 ： 成功
	/// Ed Box User ID：30690469672 报错：{"type":"error","status":400,"code":"user_already_collaborator","help_url":"http:\/\/developers.box.com\/docs\/#errors","message":"User is already a collaborator","request_id":"qi9dhdi4cahx8t58"}
	//
	err := UT.BoxCollaborationBuzUsecase.HandleAddCollaboration(biz.Box_collaboration_ow, "327431320500", "30690469672", 12, "userGid02")
	lib.DPrintln(err)
}

func Test_BoxCollaborationBuzUsecase_HandleDeleteCollaboration(t *testing.T) {
	// "327431320500", "41426608287"
	err := UT.BoxCollaborationBuzUsecase.HandleDeleteCollaboration(biz.Box_collaboration_ow, "327431320500", "32499237039", "userGid02")
	lib.DPrintln(err)
}

func Test_BoxCollaborationBuzUsecase_GetRelatedBoxUserIds(t *testing.T) {
	res, err := UT.BoxCollaborationBuzUsecase.GetRelatedBoxUserIds([]string{"0073f2144c224635931a890a5c0536ee",
		"05e9febe5f07463c96526745f2541e8c",
		"3d8666ba206a4c7ba63b792ef8eb3699",
		"6159272000011723055",
		"6159272000000453669",
		"6159272000000453640",
		"6159272000000453001"})
	lib.DPrintln(err)
	//lib.DPrintln(res)
	for _, v := range res {
		if v.BoxUser != nil {
			lib.DPrintln(v.BoxUser.BoxUserId, v.BoxUser.Login)
		}
		lib.DPrintln(v.TUser.Gid(), " v: ", v.TUser.CustomFields.TextValueByNameBasic(biz.UserFieldName_email))
	}
}

func Test_BoxCollaborationBuzUsecase_HandleCollaborationFromCase(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5814)
	err := UT.BoxCollaborationBuzUsecase.HandleCollaborationFromCase(*tCase, biz.HandleCollaborationFromCase_BizType_ClientFolder)
	lib.DPrintln(err)
}

func Test_BoxCollaborationBuzUsecase_HandleCollaborationFromCase1(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5814)
	err := UT.BoxCollaborationBuzUsecase.HandleCollaborationFromCase(*tCase, biz.HandleCollaborationFromCase_BizType_DCFolder)
	lib.DPrintln(err)
}

func Test_BoxCollaborationBuzUsecase_HandleUseVBCActiveCases(t *testing.T) {
	err := UT.BoxCollaborationBuzUsecase.HandleUseVBCActiveCases(5814)
	lib.DPrintln(err)
}
