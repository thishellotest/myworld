package biz

import (
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

type DialpadbuzUsecase struct {
	log                    *log.Helper
	CommonUsecase          *CommonUsecase
	conf                   *conf.Data
	DialpadUsecase         *DialpadUsecase
	TUsecase               *TUsecase
	DataEntryUsecase       *DataEntryUsecase
	DataComboUsecase       *DataComboUsecase
	ClientEnvelopeUsecase  *ClientEnvelopeUsecase
	BehaviorUsecase        *BehaviorUsecase
	MapUsecase             *MapUsecase
	UserUsecase            *UserUsecase
	ClientCaseUsecase      *ClientCaseUsecase
	LogUsecase             *LogUsecase
	PersonalWebformUsecase *PersonalWebformUsecase
}

func NewDialpadbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	DialpadUsecase *DialpadUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
	DataComboUsecase *DataComboUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	BehaviorUsecase *BehaviorUsecase,
	MapUsecase *MapUsecase,
	UserUsecase *UserUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	LogUsecase *LogUsecase,
	PersonalWebformUsecase *PersonalWebformUsecase) *DialpadbuzUsecase {
	uc := &DialpadbuzUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		DialpadUsecase:         DialpadUsecase,
		TUsecase:               TUsecase,
		DataEntryUsecase:       DataEntryUsecase,
		DataComboUsecase:       DataComboUsecase,
		ClientEnvelopeUsecase:  ClientEnvelopeUsecase,
		BehaviorUsecase:        BehaviorUsecase,
		MapUsecase:             MapUsecase,
		UserUsecase:            UserUsecase,
		ClientCaseUsecase:      ClientCaseUsecase,
		LogUsecase:             LogUsecase,
		PersonalWebformUsecase: PersonalWebformUsecase,
	}
	return uc
}

func AfterTextBodyHandle(text string) (newText string) {
	newText = text + "\n\nMsg&data rates may apply. Reply HELP for help or STOP to opt-out"
	return newText
}

func (c *DialpadbuzUsecase) HandleContractReminder(contractReminderType config_vbc.ContractReminderType, caseId int32, isAmContract bool) error {

	if !configs.EnabledContractReminderBySMS {
		return nil
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}

	tUser, err := c.TUsecase.Data(Kind_users, And(Eq{"gid": tCase.CustomFields.TextValueByNameBasic("user_gid")}))
	if err != nil {
		return err
	}
	if tUser == nil {
		return errors.New("tUser is nil")
	}
	dialpadUserid := tUser.CustomFields.TextValueByNameBasic("dialpad_userid")
	if dialpadUserid == "" {
		return errors.New("dialpadUserid is empty")
	}

	text, err := c.ReminderText(tClient, tCase, contractReminderType, isAmContract)
	if err != nil {
		return err
	}
	phone := tClient.CustomFields.TextValueByNameBasic(FieldName_phone)
	newPhone, err := FormatUSAPhoneHandle(phone)
	if err != nil {
		return err
	}

	email := tClient.CustomFields.TextValueByNameBasic(FieldName_email)
	if false && configs.IsWorkflowDebug(email) { // 测试帐号不发送短信

	} else {
		err = c.DialpadUsecase.SendSms(newPhone, text, dialpadUserid, caseId, "Dialpad:ContractReminder")
		if err != nil {
			return err
		}
	}

	if contractReminderType == config_vbc.ContractReminderFirst {
		behaviorType := BehaviorType_contract_reminder_first_sms
		if isAmContract {
			behaviorType = BehaviorType_am_contract_reminder_first_sms
		}
		err := c.BehaviorUsecase.Add(caseId, behaviorType, time.Now(), "")
		if err != nil {
			return err
		}
	} else if contractReminderType == config_vbc.ContractReminderSecond {
		behaviorType := BehaviorType_contract_reminder_second_sms
		if isAmContract {
			behaviorType = BehaviorType_am_contract_reminder_second_sms
		}
		err := c.BehaviorUsecase.Add(caseId, behaviorType, time.Now(), "")
		if err != nil {
			return err
		}
	} else if contractReminderType == config_vbc.ContractReminderThird {
		behaviorType := BehaviorType_contract_reminder_third_sms
		if isAmContract {
			behaviorType = BehaviorType_am_contract_reminder_third_sms
		}
		err := c.BehaviorUsecase.Add(caseId, behaviorType, time.Now(), "")
		if err != nil {
			return err
		}
	} else if contractReminderType == config_vbc.ContractReminderFourth {
		behaviorType := BehaviorType_contract_reminder_fourth_sms
		if isAmContract {
			behaviorType = BehaviorType_am_contract_reminder_fourth_sms
		}
		err := c.BehaviorUsecase.Add(caseId, behaviorType, time.Now(), "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *DialpadbuzUsecase) ReminderText(tClient *TData, tCase *TData, contractReminderType config_vbc.ContractReminderType, isAmContract bool) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}

	company := "VBC"
	if isAmContract {
		company = "Augustus Miles"
	}

	if isAmContract {
		if contractReminderType == config_vbc.ContractReminderFirst {
			text = "Dear {first_name}, I hope this message finds you well. This is a follow-up from " + company + " regarding the agreement sent on {contract_date}. Have you had an opportunity to review it yet? If you need any clarification or another copy of the agreement, please don't hesitate to reach out. Thanks!"
		} else if contractReminderType == config_vbc.ContractReminderSecond {
			text = "Dear {first_name}, I hope you're well. This is my second follow-up from " + company + " regarding the agreement sent on {contract_date}. Have you had a chance to review it? If you need another copy or have any questions, please let me know. Thanks!"
		} else if contractReminderType == config_vbc.ContractReminderThird {
			text = "Dear {first_name}, I hope you're well. This is our third follow-up from " + company + " regarding the agreement sent on {contract_date}. Have you reviewed it? If you have any questions or need another copy, please let me know. Thanks!"
		} else if contractReminderType == config_vbc.ContractReminderFourth {
			text = "Dear {first_name}, I hope you're well. This is our final follow-up from " + company + " regarding the agreement sent on {contract_date}. We've tried to reach you a few times. If you have any questions or are ready to proceed, please let us know. If we don't hear back, we'll assume you're not interested. Thanks for your time."
		} else {
			return "", errors.New("contractReminderType is wrong")
		}
	} else {
		if contractReminderType == config_vbc.ContractReminderFirst {
			text = "Dear {first_name}, I hope this message finds you well. This is a follow-up from " + company + " regarding the contract sent on {contract_date}. Have you had an opportunity to review it yet? If you need any clarification or another copy of the contract, please don't hesitate to reach out. Thanks!"
		} else if contractReminderType == config_vbc.ContractReminderSecond {
			text = "Dear {first_name}, I hope you're well. This is my second follow-up from " + company + " regarding the contract sent on {contract_date}. Have you had a chance to review it? If you need another copy or have any questions, please let me know. Thanks!"
		} else if contractReminderType == config_vbc.ContractReminderThird {
			text = "Dear {first_name}, I hope you're well. This is our third follow-up from " + company + " regarding the contract sent on {contract_date}. Have you reviewed it? If you have any questions or need another copy, please let me know. Thanks!"
		} else if contractReminderType == config_vbc.ContractReminderFourth {
			text = "Dear {first_name}, I hope you're well. This is our final follow-up from " + company + " regarding the contract sent on {contract_date}. We've tried to reach you a few times. If you have any questions or are ready to proceed, please let us know. If we don't hear back, we'll assume you're not interested. Thanks for your time."
		} else {
			return "", errors.New("contractReminderType is wrong")
		}
	}

	contractDateOn, err := c.ClientEnvelopeUsecase.ContractDateOn(tCase.CustomFields.NumberValueByNameBasic("id"), isAmContract)
	if err != nil {
		c.log.Error(err)
		return "", err
	}

	text = strings.ReplaceAll(text, "{first_name}", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	text = strings.ReplaceAll(text, "{contract_date}", contractDateOn)

	text = AfterTextBodyHandle(text)
	return text, nil

}

func TextReplaceAll(text string, params lib.TypeMap) string {
	for k, v := range params {
		text = strings.ReplaceAll(text, "{"+k+"}", InterfaceToString(v))
	}
	return text
}

func (c *DialpadbuzUsecase) TextAfterSignedContract() (text string, err error) {

	text = "VBC: Thank you for opting in to receiving sms messages from us. Message frequency may vary. Message and Data Rates may apply. Reply STOP to stop receiving messages from us. Reply HELP for more information."

	return text, nil
}

func (c *DialpadbuzUsecase) TextGettingStartedEmail(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	//text = "Hi {first_name},\n\nWe've emailed you our How-to Guide to help you get started. Please complete the first 3 steps ASAP to initiate your back pay date and request your official military records. If you haven't received the email, please check your spam folder.\n\nCan't find the Box.com account invite? Check the instructions at the bottom of step 4. Once you've completed all 5 steps, let us know so we can begin your record review.\n\nIf you have any questions or need help, don't hesitate to reach out. We can schedule a meeting to go through these steps together if needed.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	text = "Hi {first_name},\n\nYou'll receive a meeting invitation shortly to walk through your Welcome Guide email together. However, if you'd like to begin immediately, our How-To Guide attached in the email has all the instructions you need.\n\nQuick tips:\n • Can't find the email? Please check your spam folder\n • Can't find the Box.com account invite? Check the instructions at the bottom of step 4.\n • Completed all five steps in the email before the meeting? Just notify us and we'll cancel the scheduled session\n • Need help with any step? We're happy to assist\n\nEither way, we’re here to support you – don’t hesitate to reach out with any questions! \n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextGettingStartedEmailTaskLongerThan30Days -Stage 2: Task open longer than 30 days
func (c *DialpadbuzUsecase) TextGettingStartedEmailTaskLongerThan30Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nChecking in about the Welcome email. Any questions? We can schedule a call to go over it together if you'd like.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextAwaitingClientRecordsLongerThan30Days - Stage  3: Task open longer than 30 days
func (c *DialpadbuzUsecase) TextAwaitingClientRecordsLongerThan30Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nChecking in about uploading of your records. Any questions? We can schedule a call to go over it together if you'd like.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextSTRRequestPendingLongerThan30Days - Stage 4 – STR Request Pending: Waiting on STRs task open for 30 days (Not a stage yet, may want to add back)
func (c *DialpadbuzUsecase) TextSTRRequestPendingLongerThan30Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nPlease check your records request status. Remember to do this every 2 week until you receive your records. If your FOIA request closes unexpectedly, let me know right away.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextSTRRequestPending45Days -Stage 4 – STR Request Pending: Waiting on STRs every 45 days after ：只在45天发一次。
func (c *DialpadbuzUsecase) TextSTRRequestPending45Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nAny updates on your records? Please remember to check the status every 2 weeks until you receive your records.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextScheduleCall -Stage 4 –
func (c *DialpadbuzUsecase) TextScheduleCall(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nWe have initiated the record review phase of your case. This process typically takes up to 4 weeks to complete thoroughly. Once we've finished the review, we will contact you to discuss the findings and outline the next steps. If you have any questions or concerns during this time, please don't hesitate to reach out.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextStatementFinalized - Stage 10 – Statements Finalized: Immediately after moving to this stage and open task created
func (c *DialpadbuzUsecase) TextStatementFinalized(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nYour personal statements are ready for review in the \"Personal Statements\" folder of your shared folder. Please review carefully for accuracy and ensure they reflect your experiences. Let me know once you've finished reviewing and making any necessary changes.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextStatementFinalizedForWebForm - Stage 10 – Statements Finalized: Immediately after moving to this stage and open task created
func (c *DialpadbuzUsecase) TextStatementFinalizedForWebForm(tClient TData, tCase TData, tUser *TData) (text string, err error) {

	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nYour personal statements are ready for review:\n\nURL: {statement_url}\nPassword: {statement_password}\n\nPlease review carefully for accuracy and ensure they reflect your experiences. Let me know once you’ve finished reviewing and making any necessary comments.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("statement_url", tCase.CustomFields.TextValueByNameBasic(FieldName_personal_statement_manager))
	params.Set("statement_password", tCase.CustomFields.TextValueByNameBasic(FieldName_personal_statement_password))
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextStatementFinalizedEvery14Days -  Statements Finalized: Immediately after moving to this stage and open task created
func (c *DialpadbuzUsecase) TextStatementFinalizedEvery14Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nI hope you're doing well. Just checking in about your personal statements in the shared folder. Please review them for accuracy when you have time, and let me know if you need any help or once you've finished reviewing.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextStatementFinalizedEvery14DaysForWebForm -  Statements Finalized: Immediately after moving to this stage and open task created
func (c *DialpadbuzUsecase) TextStatementFinalizedEvery14DaysForWebForm(tClient TData, tCase TData, tUser *TData) (text string, err error) {

	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nI hope you're doing well. Just checking in about your personal statements, which are now available at the following URL:\n\nURL: {statement_url}\nPassword: {statement_password}\n\nPlease review them for accuracy when you have time, and let me know if you need any help or once you've finished reviewing.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("statement_url", tCase.CustomFields.TextValueByNameBasic(FieldName_personal_statement_manager))
	params.Set("statement_password", tCase.CustomFields.TextValueByNameBasic(FieldName_personal_statement_password))
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextMiniDBQsDrafts - Immediately
func (c *DialpadbuzUsecase) TextMiniDBQsDrafts(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nWe’re currently preparing your case for private medical exams as part of the next steps in your claim process. I’ll keep you updated as things progress and will let you know if we need anything further from you.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextCurrentTreatment - Current Treatment: Immediately after moving to this stage and open task created
func (c *DialpadbuzUsecase) TextCurrentTreatment(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	// First Version
	//text = "Hi {first_name},\n\nNow that your statements are finalized, please schedule appointments with your doctor for all your claimed conditions. Once scheduled, let me know the date so we can meet the day before to go over some important points.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	// Second Version
	//text = "Hi {first_name},\n\nNow that your statements are finalized, please schedule appointments with your doctor for all your claimed conditions. Once scheduled, let me know the date so we can meet the day before to go over some important points.\nAdditionally, if you prefer to email your healthcare provider, we've included an email draft at the end of your statements to help communicate your issues. Feel free to edit this draft as needed.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	//text = "Hi {first_name},\n\nNow that your statements are finalized, please schedule appointments with your doctor for all your claimed conditions. Once scheduled, let me know the date so we can meet the day before to go over some important points.\nAdditionally, if you prefer to email your healthcare provider, we've included an email draft to help communicate your issues. Feel free to edit this draft as needed.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	//text = "Hi {first_name},\n\nNow that your statements are finalized, please schedule appointments with your healthcare provider for any conditions that you have not been seen for within the past 2 years. Once scheduled, let me know the date so we can meet the day before to go over some important points.\n\nAdditionally, if you prefer to email your healthcare provider, we've included an email draft to help communicate your issues. Feel free to edit this draft as needed.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	//text = "Hi {first_name},\n\nNow that your statements are finalized, please schedule doctor appointments for the conditions listed in the email draft. Once scheduled, please share the date so that we can meet prior.\n\nWe've provided:\n- An email template if you prefer contacting your doctor electronically\n- A guide to help you communicate effectively with your doctor\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	//text = "Hi {first_name},\n\nNow that your statements are finalized, please schedule doctor appointments for the conditions listed in the email draft. Once scheduled, please share the date so that we can meet prior.\n\nWe've provided:\n- An email template if you prefer contacting your doctor electronically (located in your Personal Statements folder)\n- A guide to help you communicate effectively with your doctor (located in your Personal Statements folder)\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	//text = "Hi {first_name},\n\nNow that your statements are finalized, please complete these steps:\n 1. Review both documents in your Personal Statements folder:\n     a. Guide for effective doctor communication\n     b. Email template for contacting your doctor electronically\n 2. Reply to me via text or email that you've reviewed both documents\n\nNext Steps:\n • After your confirmation, I will provide specific scheduling authorization\n\nCritical Reminder:\n • Do not schedule any appointments until we speak to avoid unnecessary delays\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	//text = "Hi {first_name},\n\nPlease review these NEW documents in your Box.com (Personal Statements folder):\n • Guide for Effective Doctor Communication - Tips for discussing your conditions\n • Conditions Summary - Comprehensive list of all conditions we're pursuing (in ready-to-use email format)\n\nRequired Next Steps:\n 1. Review both documents carefully\n 2. Reply to confirm you've completed your review\n 3. Schedule appointments for any conditions that:\n     a. Are listed in the Conditions Summary\n     b. Haven't  been treated in the last 12 months\n\nImportant Reminder:\n • Unless it's a medical emergency, please contact me before scheduling any appointments to prevent processing delays\n\nIf you have any questions about the documents or scheduling process, please reply to this message.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	text = "Hi {first_name},\n\nNew documents are available in your Box.com Personal Statements folder:\n• Guide for Effective Doctor Communication - Tips for discussing your conditions with healthcare providers\n• Conditions Summary - Comprehensive list of all conditions we're pursuing (formatted for easy reference)\n\nRequired Next Steps:\n 1. Review both documents thoroughly\n 2. Reply to confirm you've completed your review\n 3. Contact me to discuss scheduling before making any appointments (unless medically urgent)\n 4. Schedule appointments for conditions that:\n     a. Are listed in the Conditions Summary, AND\n     b. Haven't been treated in the last 12 months\n 5. Notify me immediately with:\n     a. All appointment dates and times\n     b. Specific conditions being addressed at each appointment\n\nCritical Reminder: Do not schedule any appointments until we've spoken, unless medically urgent. This prevents delays in your case.\n\nIf you have questions or issues accessing documents, please reply.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextCurrentTreatment30Day （已经停用） -Stage 11 – Current Treatment: 30 days
func (c *DialpadbuzUsecase) TextCurrentTreatment30Day(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nChecking in about your current treatment records. Have you seen your doctor yet? Any updates?\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextCurrentTreatmentFollowingEvery30Day （已经停用）-Stage 11 – Current Treatment: Every following 30 days
func (c *DialpadbuzUsecase) TextCurrentTreatmentFollowingEvery30Day(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nFollowing up on your current treatment records. Could you give me a call when you have a moment?\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextAwaitingDecision30Days -Stage 23 – Awaiting Decision: 30 days
func (c *DialpadbuzUsecase) TextAwaitingDecision30Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nHave you heard anything from the VA? Remember to check your claims every two weeks to ensure they stay open and don't close unexpectedly.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextAwaitingDecisionEveryFollowing30Days -Stage 23 – Awaiting Decision: Every following 45 days （每隔30天）
func (c *DialpadbuzUsecase) TextAwaitingDecisionEveryFollowing30Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nFriendly reminder to check your claims every two weeks to ensure they remain open.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextAwaitingPayment -Stage 24 – Awaiting Payment: Immediately after moving to this stage and open task created
func (c *DialpadbuzUsecase) TextAwaitingPayment(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	//text = "Congratulations on your new rating, {first_name}! I've emailed you the invoice and information about additional benefits you now qualify for. You can pay the invoice via credit card, debit card, or wire transfer (see wiring instructions in the bottom of attached invoice). Please check your spam folder if the email doesn't appear in your inbox.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	text = "Congratulations on your new rating, {first_name}! I've emailed you the invoice and information about additional benefits you now qualify for. You can pay the invoice via credit card, debit card, bank transfer, or wire transfer (instructions included).\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextAwaitingPaymentAfter14Days Stage 24 – Awaiting Payment: Payment Reminder after 14 days
func (c *DialpadbuzUsecase) TextAwaitingPaymentAfter14Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nJust a gentle reminder about the invoice payment. Could you please pay when you get a chance? Thank you!\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextAwaitingPaymentTaskOpen30Days Stage 24 – Awaiting Payment: Task open 30 days
func (c *DialpadbuzUsecase) TextAwaitingPaymentTaskOpen30Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	text = "Hi {first_name},\n\nYour invoice is currently overdue. Please make the payment as soon as possible. If you're having any issues, let me know and I'll be happy to help.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextMedTeamForms Text
func (c *DialpadbuzUsecase) TextMedTeamForms(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	//text = "Hello {first_name},\n\nYou will soon receive an email with crucial documents for your private medical exams. These documents require your review and signature to authorize payment and allow the Medical Team to review your records before scheduling exams. Please check your email and respond promptly.\n\nIf you have any questions, don't hesitate to reach out to us for assistance.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	text = "Hello {first_name},\n\nYou will receive important medical exam documents by email shortly. These require your immediate signature to:\n • Authorize payment\n • Allow our Medical Team to review your records\n • Begin scheduling your exams\nPlease check your spam folder if not received and respond promptly.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextITFDeadlineIn90Days Text
func (c *DialpadbuzUsecase) TextITFDeadlineIn90Days(tClient *TData, tUser *TData) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	//text = "Hello {first_name},\n\nYou will soon receive an email with crucial documents for your private medical exams. These documents require your review and signature to authorize payment and allow the Medical Team to review your records before scheduling exams. Please check your email and respond promptly.\n\nIf you have any questions, don't hesitate to reach out to us for assistance.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	text = "Hi {first_name},\n\nYour Intent to File deadline is approaching in 90 days. If your claim requires independent medical examinations, please be aware of the following timeline requirements:\n\nTimeline Requirements:\n • Independent medical exams require a minimum of 60 days to complete\n • This includes scheduling appointments, gathering documentation, and conducting thorough reviews\n\nImportant: With only 90 days remaining and our 60-day processing requirement, time is of the essence. This is critical to secure your back pay.\n\nNext Steps: Contact us immediately if independent medical exams are needed for your claim and to avoid missing this critical deadline.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

// TextUpcomingContactInformation
func (c *DialpadbuzUsecase) TextUpcomingContactInformation(tClient TData, tUser TData, newUser TData) (text string, err error) {

	//text = "Hello {first_name},\n\nYou will soon receive an email with crucial documents for your private medical exams. These documents require your review and signature to authorize payment and allow the Medical Team to review your records before scheduling exams. Please check your email and respond promptly.\n\nIf you have any questions, don't hesitate to reach out to us for assistance.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	text = "Hi {first_name},\n\nI wanted to let you know that {upcoming_contact_vs:full_name} will be reaching out to you in the coming days. See contact information below for your reference:\n\n{upcoming_contact_vs:full_name}\n{upcoming_contact_vs:title}\nP: {upcoming_contact_vs:mobile}\nE: {upcoming_contact_vs:email}\n\nPlease feel free to connect sooner if you have any questions.\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"
	params := make(lib.TypeMap)

	params.Set("upcoming_contact_vs:full_name", newUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname))
	params.Set("upcoming_contact_vs:title", newUser.CustomFields.TextValueByNameBasic(UserFieldName_title))
	params.Set("upcoming_contact_vs:mobile", newUser.CustomFields.TextValueByNameBasic(UserFieldName_mobile))
	params.Set("upcoming_contact_vs:email", newUser.CustomFields.TextValueByNameBasic(UserFieldName_email))

	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

func (c *DialpadbuzUsecase) TextZoomMeetingNotice(tClient *TData, tUser *TData, meetingTopic string, meetingLink string, meetingStartTime string) (text string, err error) {

	if tClient == nil {
		return "", errors.New("tClient is nil")
	}
	if tUser == nil {
		return "", errors.New("tUser is nil")
	}

	// July 25, 2024, 9:02 PM (PT)
	text = "Dear {first_name},\n\nThis is a reminder for the upcoming Zoom conference.\n\nTopic: {meeting_topic}\nMeeting Start Time: {meeting_start_time}\nMeeting Link: {meeting_link}\n\nPlease join the meeting on time. Looking forward to your participation!\n\n{users:first_name} {users:last_name}\nVeteran Benefits Center"

	params := make(lib.TypeMap)
	params.Set("first_name", tClient.CustomFields.TextValueByNameBasic(FieldName_first_name))
	params.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
	params.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))

	params.Set("meeting_start_time", meetingStartTime)
	params.Set("meeting_topic", meetingTopic)
	params.Set("meeting_link", meetingLink)

	text = TextReplaceAll(text, params)
	text = AfterTextBodyHandle(text)
	return text, nil
}

type HandleSendSMSType string

const (
	HandleSendSMSTextAfterSignedContract = HandleSendSMSType("TextAfterSignedContract")

	HandleSendSMSTextGettingStartedEmail                     = HandleSendSMSType("TextGettingStartedEmail")
	HandleSendSMSTextGettingStartedEmailTaskLongerThan30Days = HandleSendSMSType("TextGettingStartedEmailTaskLongerThan30Days")

	HandleSendSMSTextAwaitingClientRecordsLongerThan30Days = HandleSendSMSType("TextAwaitingClientRecordsLongerThan30Days")

	HandleSendSMSTextSTRRequestPendingLongerThan30Days = HandleSendSMSType("TextSTRRequestPendingLongerThan30Days")
	// HandleSendSMSTextSTRRequestPending45Days 改为60天
	HandleSendSMSTextSTRRequestPending45Days = HandleSendSMSType("TextSTRRequestPending45Days")

	HandleSendSMSTextScheduleCall = HandleSendSMSType("TextScheduleCall")

	HandleSendSMSTextStatementFinalized            = HandleSendSMSType("TextStatementFinalized")
	HandleSendSMSTextStatementFinalizedEvery14Days = HandleSendSMSType("TextStatementFinalizedEvery14Days")

	HandleSendSMSTextMiniDBQsDrafts = HandleSendSMSType("TextMiniDBQsDrafts")

	HandleSendSMSTextCurrentTreatment                    = HandleSendSMSType("TextCurrentTreatment")
	HandleSendSMSTextCurrentTreatment30Day               = HandleSendSMSType("TextCurrentTreatment30Day")
	HandleSendSMSTextCurrentTreatmentFollowingEvery30Day = HandleSendSMSType("TextCurrentTreatmentFollowingEvery30Day")

	HandleSendSMSTextAwaitingDecision30Days               = HandleSendSMSType("TextAwaitingDecision30Days")
	HandleSendSMSTextAwaitingDecisionEveryFollowing30Days = HandleSendSMSType("TextAwaitingDecisionEveryFollowing30Days")

	HandleSendSMSTextAwaitingPayment               = HandleSendSMSType("TextAwaitingPayment")
	HandleSendSMSTextAwaitingPaymentAfter14Days    = HandleSendSMSType("TextAwaitingPaymentAfter14Days")
	HandleSendSMSTextAwaitingPaymentTaskOpen30Days = HandleSendSMSType("TextAwaitingPaymentTaskOpen30Days")
	HandleSendSMSTextMedTeamForms                  = HandleSendSMSType("TextMedTeamForms")

	HandleSendSMSTextUpcomingContactInformation = HandleSendSMSType("TextUpcomingContactInformation")

	HandleSendSMSTextZoomMeetingNotice = HandleSendSMSType("ZoomMeetingNotice")

	HandleSendSMSTextITFDeadlineIn90Days = HandleSendSMSType("TextITFDeadlineIn90Days")
)

func (c *DialpadbuzUsecase) NeedLimit(handleSendSMSType HandleSendSMSType) bool {
	if handleSendSMSType == HandleSendSMSTextGettingStartedEmail ||
		handleSendSMSType == HandleSendSMSTextScheduleCall ||
		handleSendSMSType == HandleSendSMSTextStatementFinalized ||
		handleSendSMSType == HandleSendSMSTextMiniDBQsDrafts ||
		handleSendSMSType == HandleSendSMSTextCurrentTreatment ||
		handleSendSMSType == HandleSendSMSTextAwaitingPayment ||
		handleSendSMSType == HandleSendSMSTextMedTeamForms ||
		handleSendSMSType == HandleSendSMSTextAfterSignedContract ||
		handleSendSMSType == HandleSendSMSTextUpcomingContactInformation {
		return true
	}
	return false
}

func (c *DialpadbuzUsecase) IsLimit(handleSendSMSType HandleSendSMSType, caseId int32) (bool, error) {
	key := fmt.Sprintf("%s%s:%d", "Dialpadbuz:", handleSendSMSType, caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return false, err
	}
	if val == "1" {
		return true, nil
	}
	return false, nil
}

func (c *DialpadbuzUsecase) NeedHandleSendSMS(handleSendSMSType HandleSendSMSType, caseId int32) (bool, error) {
	if c.NeedLimit(handleSendSMSType) {
		isLimit, err := c.IsLimit(handleSendSMSType, caseId)
		if err != nil {
			return false, err
		}
		if isLimit {
			return false, nil
		}
	}
	return true, nil
}

func (c *DialpadbuzUsecase) HandleSendSMS(handleSendSMSType HandleSendSMSType, caseId int32, cronTriggerVo CronTriggerVo) error {
	needHandleSendSMS, err := c.NeedHandleSendSMS(handleSendSMSType, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if !needHandleSendSMS {
		return nil
	}
	err = c.HandleSendSMSBiz(handleSendSMSType, caseId, cronTriggerVo)
	if err != nil {
		return err
	}
	if c.NeedLimit(handleSendSMSType) {
		key := fmt.Sprintf("%s%s:%d", "Dialpadbuz:", handleSendSMSType, caseId)
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

// HandleSendSMSBiz 有些每隔多少天的任务，直接调用此方法
func (c *DialpadbuzUsecase) HandleSendSMSBiz(handleSendSMSType HandleSendSMSType, caseId int32, cronTriggerVo CronTriggerVo) error {
	var err error
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	var tUser *TData
	var text string

	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}

	if handleSendSMSType == HandleSendSMSTextAwaitingPayment ||
		handleSendSMSType == HandleSendSMSTextAwaitingPaymentAfter14Days ||
		handleSendSMSType == HandleSendSMSTextAwaitingPaymentTaskOpen30Days {

		// todo:lgl 以下为测试帐号
		if tCase.CustomFields.TextValueByNameBasic(FieldName_email) == "liaogling@gmail.com" ||
			tCase.CustomFields.TextValueByNameBasic(FieldName_email) == "lialing@foxmail.com" {
			primaryVs := "Engineering Team"
			tUser, err = c.UserUsecase.GetByFullName(primaryVs)
			if err != nil {
				return err
			}
			if tUser == nil {
				return errors.New("tUser is nil")
			}
		} else {
			tUser, err = c.UserUsecase.GetByGid(config_vbc.User_Edward_gid)
			if err != nil {
				return err
			}
			if tUser == nil {
				return errors.New("tUser is nil")
			}
		}
	} else if handleSendSMSType == HandleSendSMSTextUpcomingContactInformation {

		leadVSChangeLogValue := cronTriggerVo.Params.GetString("LeadVSChangeLog")
		leadVSChangeLogVo := lib.StringToTDef(leadVSChangeLogValue, LeadVSChangeLogVo{})

		var newUser *TData
		newUser, _ = c.UserUsecase.GetUserByLeadVS(tCase)
		if newUser == nil {
			return errors.New("LeadVS User is nil: " + InterfaceToString(tCase.Id()))
		}
		if leadVSChangeLogVo.PreviousVSUserGid == "" {
			tUser = newUser
		} else {
			tUser, _ = c.UserUsecase.GetByGid(leadVSChangeLogVo.PreviousVSUserGid)
		}
		if tUser == nil {
			return errors.New("tUser is nil: " + InterfaceToString(tCase.Id()))
		}
		text, err = c.TextUpcomingContactInformation(*tClient, *tUser, *newUser)
	} else {

		if configs.UseOwnerSendingSMS {
			tUser, err = c.UserUsecase.GetByGid(tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid))
			if err != nil {
				c.log.Error(err)
				return err
			}
		} else {
			primaryVs := tCase.CustomFields.TextValueByNameBasic("primary_vs")
			// todo:lgl 以下为测试帐号
			if tCase.CustomFields.TextValueByNameBasic(FieldName_email) == "liaogling@gmail.com" ||
				tCase.CustomFields.TextValueByNameBasic(FieldName_email) == "lialing@foxmail.com" {
				primaryVs = "Engineering Team"
			}
			tUser, err = c.UserUsecase.GetByFullName(primaryVs)
			if err != nil {
				return err
			}
			if tUser == nil {
				return errors.New("tUser is nil")
			}
		}
	}

	if handleSendSMSType == HandleSendSMSTextAfterSignedContract {
		text, err = c.TextAfterSignedContract()
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextGettingStartedEmail {
		text, err = c.TextGettingStartedEmail(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextGettingStartedEmailTaskLongerThan30Days {

		text, err = c.TextGettingStartedEmailTaskLongerThan30Days(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextAwaitingClientRecordsLongerThan30Days {
		text, err = c.TextAwaitingClientRecordsLongerThan30Days(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextSTRRequestPendingLongerThan30Days {
		text, err = c.TextSTRRequestPendingLongerThan30Days(tClient, tUser)
		if err != nil {
			return err
		}

	} else if handleSendSMSType == HandleSendSMSTextSTRRequestPending45Days {
		text, err = c.TextSTRRequestPending45Days(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextScheduleCall {
		text, err = c.TextScheduleCall(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextStatementFinalized {
		useNewPersonalWebForm, err := c.PersonalWebformUsecase.IsUseNewPersonalWebForm(tCase.Id())
		if err != nil {
			return err
		}
		if useNewPersonalWebForm {
			text, err = c.TextStatementFinalizedForWebForm(*tClient, *tCase, tUser)
			if err != nil {
				return err
			}
		} else {
			text, err = c.TextStatementFinalized(tClient, tUser)
			if err != nil {
				return err
			}
		}

	} else if handleSendSMSType == HandleSendSMSTextStatementFinalizedEvery14Days {

		useNewPersonalWebForm, err := c.PersonalWebformUsecase.IsUseNewPersonalWebForm(tCase.Id())
		if err != nil {
			return err
		}
		if useNewPersonalWebForm {
			text, err = c.TextStatementFinalizedEvery14DaysForWebForm(*tClient, *tCase, tUser)
			if err != nil {
				return err
			}
		} else {
			text, err = c.TextStatementFinalizedEvery14Days(tClient, tUser)
			if err != nil {
				return err
			}
		}
	} else if handleSendSMSType == HandleSendSMSTextMiniDBQsDrafts {
		text, err = c.TextMiniDBQsDrafts(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextCurrentTreatment {
		text, err = c.TextCurrentTreatment(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextCurrentTreatment30Day {
		text, err = c.TextCurrentTreatment30Day(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextCurrentTreatmentFollowingEvery30Day {
		text, err = c.TextCurrentTreatmentFollowingEvery30Day(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextAwaitingDecision30Days {
		text, err = c.TextAwaitingDecision30Days(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextAwaitingDecisionEveryFollowing30Days {
		text, err = c.TextAwaitingDecisionEveryFollowing30Days(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextAwaitingPayment {
		text, err = c.TextAwaitingPayment(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextAwaitingPaymentAfter14Days {
		text, err = c.TextAwaitingPaymentAfter14Days(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextAwaitingPaymentTaskOpen30Days {
		text, err = c.TextAwaitingPaymentTaskOpen30Days(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextMedTeamForms {
		text, err = c.TextMedTeamForms(tClient, tUser)
		if err != nil {
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextUpcomingContactInformation {
		// 在前面已处理
	} else if handleSendSMSType == HandleSendSMSTextITFDeadlineIn90Days {
		text, err = c.TextITFDeadlineIn90Days(tClient, tUser)
		if err != nil {
			return err
		}
	} else {
		return errors.New("handleSendSMSType is error")
	}

	return c.DoHandleSendSMSBiz(handleSendSMSType, text, *tUser, *tClient, *tCase)
}

func (c *DialpadbuzUsecase) DoHandleSendSMSBiz(handleSendSMSType HandleSendSMSType, text string, senderTUser TData, receiveClient TData, receiveCase TData) error {

	if text == "" {
		return errors.New("text is empty")
	}
	caseId := receiveCase.Id()
	phone := receiveClient.CustomFields.TextValueByNameBasic(FieldName_phone)
	newPhone, err := FormatUSAPhoneHandle(phone)
	if err != nil {
		return err
	}

	dialpadUserid := senderTUser.CustomFields.TextValueByNameBasic(UserFieldName_dialpad_userid)
	if dialpadUserid == "" {
		return errors.New("dialpadUserid is empty")
	}

	err = c.DialpadUsecase.SendSms(newPhone, text, dialpadUserid, caseId, "Dialpad:"+string(handleSendSMSType))
	if err != nil {
		return err
	}

	// 加入行为
	behTyoe := fmt.Sprintf("%s%s", BehaviorType_prefix_sms, string(handleSendSMSType))
	er := c.BehaviorUsecase.Add(caseId, behTyoe, time.Now(), "")
	if er != nil {
		c.log.Error(er)
	}
	return nil
}

func (c *DialpadbuzUsecase) BizSendSms(typ HandleSendSMSType, tClient *TData, tCase *TData, tUser *TData, smsContent string) error {

	if tClient == nil {
		return errors.New("tClient is nil")
	}
	if tUser == nil {
		return errors.New("tUser is nil")
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.Id()

	phone := tClient.CustomFields.TextValueByNameBasic(FieldName_phone)
	newPhone, err := FormatUSAPhoneHandle(phone)
	if err != nil {
		return err
	}

	dialpadUserid := tUser.CustomFields.TextValueByNameBasic(UserFieldName_dialpad_userid)
	if dialpadUserid == "" {
		return errors.New("dialpadUserid is empty")
	}

	err = c.DialpadUsecase.SendSms(newPhone, smsContent, dialpadUserid, caseId, "Dialpad:"+string(typ))
	if err != nil {
		return err
	}
	return nil
}

// SyncDialpadUserToVBCUser 使用邮箱地址同步匹配
func (c *DialpadbuzUsecase) SyncDialpadUserToVBCUser() error {
	items, err := c.DialpadUsecase.UserList()
	if err != nil {
		return err
	}
	for _, v := range items {
		lib.DPrintln(v)

		emails := lib.InterfaceToTDef[[]string](v.Get("emails"), nil)
		if len(emails) != 1 {
			return errors.New("emails is wrong: " + InterfaceToString(emails))
		}
		email := emails[0]
		dialpadUserid := v.GetString("id")
		phoneNumbers := lib.InterfaceToTDef[[]string](v.Get("phone_numbers"), nil)
		if len(phoneNumbers) != 1 {
			return errors.New("phoneNumbers is wrong: " + InterfaceToString(phoneNumbers))
		}
		phoneNumber := phoneNumbers[0]
		user, err := c.TUsecase.Data(Kind_users, And(Eq{"email": email, "deleted_at": 0}))
		if err != nil {
			c.log.Error(err)
			return err
		}

		lib.DPrintln("email: ", email, "userId: ", dialpadUserid, "phoneNumber: ", phoneNumber)

		if user == nil {
			c.log.Info("email: ", email, " does not exist from users")
			continue
		}
		_, err = c.DataEntryUsecase.HandleOne(Kind_users, map[string]interface{}{
			"id":                  user.CustomFields.NumberValueByNameBasic("id"),
			"dialpad_userid":      dialpadUserid,
			"dialpad_phonenumber": phoneNumber,
		}, "id", nil)
		if err != nil {
			c.log.Error(err)
			return err
		}
	}
	return nil
}

func (c *DialpadbuzUsecase) HandleAfterActionStop(phone string) error {
	text := "VBC: You will receive no further messages from us. If this was in error reply UNSTOP to continue receiving messages."
	leadVsTUser, err := c.ClientCaseUsecase.GetLeadVSByPhone(phone)
	if err != nil {
		return err
	}
	if leadVsTUser == nil {
		return errors.New("leadVsTUser is nil")
	}
	dialpadUserid := leadVsTUser.CustomFields.TextValueByNameBasic(UserFieldName_dialpad_userid)
	if dialpadUserid == "" {
		return errors.New("dialpadUserid is empty")
	}
	c.log.Info("HandleAfterActionStop:phone:", phone)
	c.log.Info("HandleAfterActionStop:text:", text)
	c.log.Info("HandleAfterActionStop:dialpadUserid:", dialpadUserid)

	c.LogUsecase.SaveLog(0, "Dialpadbuz:HandleAfterActionStop", map[string]interface{}{
		"ReceivePhone":        phone,
		"ReceiveText":         text,
		"SenderDialpadUserid": dialpadUserid,
		"SenderUserFullName":  leadVsTUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname),
	})

	if configs.Enable_SMS_New_Version_Debug {
		return nil
	}
	return c.DialpadUsecase.SendSms(phone, text, dialpadUserid, 0, "HandleAfterActionStop")
}

// HandleAfterSignedContract phone
func (c *DialpadbuzUsecase) HandleAfterSignedContract(phone string) error {
	text := "VBC: Thank you for opting in to receiving sms messages from us. Message frequency may vary. Message and Data Rates may apply. Reply STOP to stop receiving messages from us. Reply HELP for more information."
	leadVsTUser, err := c.ClientCaseUsecase.GetLeadVSByPhone(phone)
	if err != nil {
		return err
	}
	if leadVsTUser == nil {
		return errors.New("leadVsTUser is nil")
	}
	dialpadUserid := leadVsTUser.CustomFields.TextValueByNameBasic(UserFieldName_dialpad_userid)
	if dialpadUserid == "" {
		return errors.New("dialpadUserid is empty")
	}
	c.log.Info("HandleAfterSignedContract:phone:", phone)
	c.log.Info("HandleAfterSignedContract:text:", text)
	c.log.Info("HandleAfterSignedContract:dialpadUserid:", dialpadUserid)

	c.LogUsecase.SaveLog(0, "Dialpadbuz:HandleAfterSignedContract", map[string]interface{}{
		"ReceivePhone":        phone,
		"ReceiveText":         text,
		"SenderDialpadUserid": dialpadUserid,
		"SenderUserFullName":  leadVsTUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname),
	})
	if configs.Enable_SMS_New_Version_Debug {
		return nil
	}
	return c.DialpadUsecase.SendSmsNoFilter(phone, text, dialpadUserid, 0, "HandleAfterActionStop")
}
