package biz

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

type FlowHttpUsecase struct {
	log         *log.Helper
	conf        *conf.Data
	JWTUsecase  *JWTUsecase
	UserUsecase *UserUsecase
}

func NewFlowHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	UserUsecase *UserUsecase) *FlowHttpUsecase {
	return &FlowHttpUsecase{
		log:         log.NewHelper(logger),
		conf:        conf,
		JWTUsecase:  JWTUsecase,
		UserUsecase: UserUsecase,
	}
}

type FlowHttpStagesTriggerRequest struct {
	Stages string `json:"stages"`
}

type FlowHttpStagesTriggers []FlowHttpStagesTriggerItem

type FlowHttpStagesTriggerItem struct {
	Content string `json:"content"`
}

var FlowHttpStagesTriggerConfigs = map[string]FlowHttpStagesTriggers{
	config_vbc.Stages_FeeScheduleandContract: {{
		Content: "Sending an email: \"Your Veteran Benefits Fee Schedule\"",
	}, {
		Content: "Sending contract: \"Your Veteran Benefits Center Contract\"",
	}},
	config_vbc.Stages_GettingStartedEmail: {{
		Content: "Sending an email: \"Welcome to the Veterans Benefits Center: Next Steps\"",
	}, {
		Content: "Send immediately: Text message \"You'll receive a meeting invitation shortly to walk through your Welcome Guide email...\"",
	}, {
		Content: "Delay sending: Text message \"Checking in about the Welcome email. Any questions?...\"",
	}},

	config_vbc.Stages_AwaitingClientRecords: {
		{
			Content: "Delay sending: Text message \"Checking in about uploading of your records. Any questions?\"",
		},
	},
	config_vbc.Stages_AmAwaitingClientRecords: {
		{
			Content: "Delay sending: Text message \"Checking in about uploading of your records. Any questions?\"",
		},
	},
	config_vbc.Stages_STRRequestPending: {
		{
			Content: "Delay sending: Text message \"Please check your records request status...\"",
		},
	},
	config_vbc.Stages_AmSTRRequestPending: {
		{
			Content: "Delay sending: Text message \"Please check your records request status...\"",
		},
	},
	config_vbc.Stages_RecordReview: {
		{
			Content: "Send immediately: Text message \"We have initiated the record review phase of your case...\"",
		},
		{
			Content: "Sending an email: \"Update: Your Records Review Process Has Begun\"",
		},
	},
	config_vbc.Stages_AmRecordReview: {
		{
			Content: "Send immediately: Text message \"We have initiated the record review phase of your case...\"",
		},
		{
			Content: "Sending an email: \"Update: Your Records Review Process Has Begun\"",
		},
	},
	config_vbc.Stages_ScheduleCall: {
		{
			Content: "Sending an email: \"Veteran Services Contact\"",
		},
		{
			Content: "Send immediately: Text message: \"I wanted to let you know that [Name] will be reaching out to you in the coming days\"",
		},
	},
	config_vbc.Stages_AmScheduleCall: {
		{
			Content: "Sending an email: \"Veteran Services Contact\"",
		},
		{
			Content: "Send immediately: Text message: \"I wanted to let you know that [Name] will be reaching out to you in the coming days\"",
		},
	},
	config_vbc.Stages_StatementsFinalized: {
		{
			Content: "[Client PW Inactive]Send immediately: Text message \"Your personal statements are ready for review in the \"Personal Statements\" folder...\"",
		},
		{
			Content: "[Client PW Inactive]Delay sending: Text message \"I hope you're doing well. Just checking in about your personal statements in the shared folder...\"",
		},
		{
			Content: "[Client PW Inactive]Sending an email: \"Personal Statements Ready for Your Review\"",
		},
		{
			Content: "[Client PW Inactive]Delay sending an email: \"Please Review Your Personal Statements in Shared Folder\"",
		},
		{
			Content: "",
		},
		{
			Content: "[Client PW Active]Send immediately: Text message \"Your personal statements are ready for review...\"",
		},
		{
			Content: "[Client PW Active]Delay sending: Text message \"I hope you're doing well. Just checking in about your personal statements, which are now available at the following URL...\"",
		},
		{
			Content: "[Client PW Active]Sending an email: \"Personal Statements Ready for Your Review\"",
		},
		{
			Content: "[Client PW Active]Delay sending an email: \"Please Review Your Personal Statements – Access via URL and Password\"",
		},
	},
	config_vbc.Stages_AmStatementsFinalized: {
		{
			Content: "[Client PW Inactive]Send immediately: Text message \"Your personal statements are ready for review in the \"Personal Statements\" folder...\"",
		},
		{
			Content: "[Client PW Inactive]Delay sending: Text message \"I hope you're doing well. Just checking in about your personal statements in the shared folder...\"",
		},
		{
			Content: "[Client PW Inactive]Sending an email: \"Personal Statements Ready for Your Review\"",
		},
		{
			Content: "[Client PW Inactive]Delay sending an email: \"Please Review Your Personal Statements in Shared Folder\"",
		},
		{
			Content: "",
		},
		{
			Content: "[Client PW Active]Send immediately: Text message \"Your personal statements are ready for review...\"",
		},
		{
			Content: "[Client PW Active]Delay sending: Text message \"I hope you're doing well. Just checking in about your personal statements, which are now available at the following URL...\"",
		},
		{
			Content: "[Client PW Active]Sending an email: \"Personal Statements Ready for Your Review\"",
		},
		{
			Content: "[Client PW Active]Delay sending an email: \"Please Review Your Personal Statements – Access via URL and Password\"",
		},
	},
	config_vbc.Stages_CurrentTreatment: {
		{
			Content: "Send immediately: Text message \"New documents are available in your Box.com Personal Statements folder...\"",
		},
		{
			Content: "Sending an email: \"Action Required: Confirm Document Review Before Proceeding\"",
		},
		//{
		//	Content: "Delay sending: Text message \"Checking in about your current treatment records. Have you seen your doctor yet?...\"",
		//},
		//{
		//	Content: "Delay sending: Text message \"Following up on your current treatment records...\"",
		//},
	},
	config_vbc.Stages_AmCurrentTreatment: {
		{
			Content: "Send immediately: Text message \"New documents are available in your Box.com Personal Statements folder...\"",
		},
		{
			Content: "Sending an email: \"Action Required: Confirm Document Review Before Proceeding\"",
		},
		//{
		//	Content: "Delay sending: Text message \"Checking in about your current treatment records. Have you seen your doctor yet?...\"",
		//},
		//{
		//	Content: "Delay sending: Text message \"Following up on your current treatment records...\"",
		//},
	},
	config_vbc.Stages_MiniDBQs_Draft: {
		{
			Content: "Send immediately: Text message \"We’re currently preparing your case for private medical exams as part of the next steps...\"",
		},
		{
			Content: "Sending an email: \"Preparing for Your Private Medical Exams\"",
		},
	},
	config_vbc.Stages_AmMiniDBQs_Draft: {
		{
			Content: "Send immediately: Text message \"We’re currently preparing your case for private medical exams as part of the next steps...\"",
		},
		{
			Content: "Sending an email: \"Preparing for Your Private Medical Exams\"",
		},
	},
	config_vbc.Stages_AwaitingDecision: {
		{
			Content: "Delay sending: Text message \"Have you heard anything from the VA? Remember to check your claims every two weeks to ensure...\"",
		},
		{
			Content: "Delay sending: Text message \"friendly reminder to check your claims every two weeks to ensure they remain open...\"",
		},
	},
	config_vbc.Stages_AmAwaitingDecision: {
		{
			Content: "Delay sending: Text message \"Have you heard anything from the VA? Remember to check your claims every two weeks to ensure...\"",
		},
		{
			Content: "Delay sending: Text message \"friendly reminder to check your claims every two weeks to ensure they remain open...\"",
		},
	},
	config_vbc.Stages_AwaitingPayment: {
		{
			Content: "Send immediately: Text message \"Congratulations on your new rating, {first_name}...\"",
		},
		{
			Content: "Delay sending: Text message \"Just a gentle reminder about the invoice payment...\"",
		},
		{
			Content: "Delay sending: Text message \"Your invoice is currently overdue. Please make the payment as soon as possible...\"",
		},
		{
			Content: "Sending an email: \"Congratulations on Your New VA Rating\"",
		},
	},
	config_vbc.Stages_AmAwaitingPayment: {
		{
			Content: "Send immediately: Text message \"Congratulations on your new rating, {first_name}...\"",
		},
		{
			Content: "Delay sending: Text message \"Just a gentle reminder about the invoice payment...\"",
		},
		{
			Content: "Delay sending: Text message \"Your invoice is currently overdue. Please make the payment as soon as possible...\"",
		},
		{
			Content: "Sending an email: \"Congratulations on Your New VA Rating\"",
		},
	},
	config_vbc.Stages_MiniDBQ_Forms: {
		{
			Content: "Sending contract: \"Medical Team Forms\"",
		},
		{
			Content: "Send immediately: Text message \"You will receive important medical exam documents by email shortly...\"",
		},
	},
	config_vbc.Stages_AmMiniDBQ_Forms: {
		{
			Content: "Sending contract: \"Medical Team Forms\"",
		},
		{
			Content: "Send immediately: Text message \"You will receive important medical exam documents by email shortly...\"",
		},
	},
}

func (c *FlowHttpUsecase) StagesTrigger(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	var flowHttpStagesTriggerRequest FlowHttpStagesTriggerRequest
	json.Unmarshal(rawData, &flowHttpStagesTriggerRequest)

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	data, err := c.BizStagesTrigger(userFacade, flowHttpStagesTriggerRequest)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *FlowHttpUsecase) BizStagesTrigger(userFacade UserFacade,
	flowHttpStagesTriggerRequest FlowHttpStagesTriggerRequest) (lib.TypeMap, error) {

	tProfile, _ := c.UserUsecase.GetProfile(&userFacade.TData)

	if flowHttpStagesTriggerRequest.Stages == "" {
		return nil, errors.New("Incorrect parameter")
	}
	//var flowHttpStagesTriggers FlowHttpStagesTriggers

	//if _,ok:=FlowHttpStagesTriggerConfigs[flowHttpStagesTriggerRequest.Stages];ok {
	flowHttpStagesTriggers, _ := FlowHttpStagesTriggerConfigs[flowHttpStagesTriggerRequest.Stages]
	//}
	if tProfile.CustomFields.NumberValueByNameBasic(Profile_FieldName_is_admin) != Profile_IsAdmin_Yes {
		flowHttpStagesTriggers = nil
	}
	data := make(lib.TypeMap)
	data.Set("triggers", flowHttpStagesTriggers)
	return data, nil
}
