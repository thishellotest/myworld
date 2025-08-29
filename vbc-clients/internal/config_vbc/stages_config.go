package config_vbc

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

/*
1. Fee Schedule and Contract
2. Getting Started Email
3. Awaiting Client Records
4. STR Request Pending
5. Record Review
6. Schedule Call
7. Statement Notes
8. Statement Drafts
9. Statement Review
10. Statement Finalized
11. Current Treatment
12. Mini-DBQs Draft
13. MedTeam Forms
14. Mini-DBQs Finalized
15. Nexus Letters
16. Medical Team - Private Exams Submitted
17. Medical Team - Exams Scheduled
18. Medical Team - Call Vet
19. DBQ Completed
20. File Claims Draft
21. File Claims with Client
22. Verify Evidence Received
23. Awaiting Decision
24. Awaiting Payment
25. Completed
26. Terminated
27. Dormant
*/

/*
1. Fee Schedule and Contract | Fee Schedule and Contract  |  6159272000000463430
2. Getting Started Email | Getting Started Email  |  6159272000000463433
3. Awaiting Client Records | Awaiting Client Records  |  6159272000000463445
4. STR Request Pending | 4.STR Request Pending  |  6159272000005402001
5. Record Review | Record Review  |  6159272000000463451
6. Schedule Call | Schedule Call  |  6159272000000463454
7. Statement Notes | Statement Notes  |  6159272000000902001
8. Statement Drafts | Statement Drafts  |  6159272000000463459
9. Statement Review | 8. Statement Review  |  6159272000001719014
10. Statement Finalized | Statement Finalized  |  6159272000000463462
11. Current Treatment | Current Treatment  |  6159272000000463465
12. Mini-DBQs Draft | Mini-DBQs  |  6159272000000463468
13. MedTeam Forms | 12. Mini-DBQ Forms  |  6159272000002764001
14. Mini-DBQs Finalized | Mini-DBQs Finalized  |  6159272000000927194
15. Nexus Letters | Nexus Letter  |  6159272000000463471
16. Medical Team - Private Exams Submitted | Medical Team  |  6159272000000463474
17. Medical Team - Exams Scheduled | 15. Medical Team - Exams Scheduled  |  6159272000002167001
18. Medical Team - Call Vet | 14. Medical Team - Call Vet  |  6159272000001381009
19. DBQ Completed | 15. DBQ Completed  |  6159272000001612001
20. File Claims Draft | File Claims (New or Supplement)  |  6159272000000463477
21. File Claims with Client | File Claims with Client  |  6159272000000927197
22. Verify Evidence Received | Verify Evidence Received  |  6159272000000463480
23. Awaiting Decision | Awaiting Decision  |  6159272000000463483
24. Awaiting Payment | Awaiting Payment  |  6159272000000463486
25. Completed | Completed  |  6159272000000496093
26. Terminated | 25. Terminated  |  6159272000004762050
27. Dormant | 26. Suspended  |  6159272000004762053
*/
const (

	// 此处都是DB的配置，与field_options配合使用

	Stages_AmIncomingRequest   = "Am__Incoming Request"
	Stages_AmInformationIntake = "Am__Information Intake"
	Stages_AmContractPending   = "Am__Contract Pending"
	//Stages_AmSendVA2122a           = "Am__Send VA 21-22a"
	Stages_AmAwaitingClientRecords = "Am__Awaiting Client Records"

	Stages_IncomingRequest        = "Incoming Request"
	Stages_FeeScheduleandContract = "Fee Schedule and Contract"

	Stages_GettingStartedEmail   = "Getting Started Email"
	Stages_AwaitingClientRecords = "Awaiting Client Records"
	Stages_STRRequestPending     = "STR Request Pending"
	Stages_AmSTRRequestPending   = "Am__STR Request Pending"
	Stages_RecordReview          = "Record Review"
	Stages_AmRecordReview        = "Am__Record Review"

	Stages_ClaimAnalysis         = "Claim Analysis" // 6
	Stages_AmClaimAnalysis       = "Am__Claim Analysis"
	Stages_ClaimAnalysisReview   = "Claim Analysis Review"     // 7
	Stages_AmClaimAnalysisReview = "Am__Claim Analysis Review" // 7
	Stages_ScheduleCall          = "Schedule Call"             // 8
	Stages_AmScheduleCall        = "Am__Schedule Call"

	Stages_StatementNotes                     = "Statement Notes"
	Stages_AmStatementNotes                   = "Am__Statement Notes"
	Stages_StatementDrafts                    = "Statement Drafts"
	Stages_AmStatementDrafts                  = "Am__Statement Drafts"
	Stages_StatementDrafts_Number             = 10 //需要设置正确，否则影响statement的限制修改
	Stages_StatementReview                    = "Statement Review"
	Stages_AmStatementReview                  = "Am__Statement Review"
	Stages_StatementsFinalized                = "Statement Finalized"
	Stages_AmStatementsFinalized              = "Am__Statement Finalized"
	Stages_CurrentTreatment                   = "Current Treatment"
	Stages_AmCurrentTreatment                 = "Am__Current Treatment"
	Stages_CurrentTreatmentReview             = "Current Treatment Review"
	Stages_StatementUpdates                   = "Statement Updates"
	Stages_StatementUpdatesDraft              = "Statement Updates Draft"
	Stages_StatementUpdatesComplete           = "Statement Updates Complete"
	Stages_PreparingDocumentsTinnitusLetter   = "11. Preparing documents for Tinnitus letter"
	Stages_AmPreparingDocumentsTinnitusLetter = "Am__11. Preparing documents for Tinnitus letter"
	Stages_AwaitingNexusLetter                = "Awaiting Nexus Letter"
	Stages_AmAwaitingNexusLetter              = "Am__Awaiting Nexus Letter"

	Stages_MiniDBQs_Draft           = "Mini-DBQs Draft"
	Stages_AmMiniDBQs_Draft         = "Am__Mini-DBQs Draft"
	Stages_MiniDBQs                 = "Mini-DBQs Finalized"
	Stages_AmMiniDBQs               = "Am__Mini-DBQs Finalized"
	Stages_MiniDBQ_Forms            = "Medical Team - Forms Sent"     // 12. Mini-DBQ Forms
	Stages_AmMiniDBQ_Forms          = "Am__Medical Team - Forms Sent" // 12. Mini-DBQ Forms
	Stages_MedicalTeamFormsSigned   = "Medical Team - Forms Signed"
	Stages_AmMedicalTeamFormsSigned = "Am__Medical Team - Forms Signed"

	//Stages_NexusLetters                = "Nexus Letters"
	Stages_MedicalTeam                   = "Medical Team - Private Exams Submitted"     // 15. Medical Team - Filing
	Stages_AmMedicalTeam                 = "Am__Medical Team - Private Exams Submitted" // 15. Medical Team - Filing
	Stages_MedicalTeamPaymentCollected   = "Medical Team - Payment Collected"
	Stages_AmMedicalTeamPaymentCollected = "Am__Medical Team - Payment Collected"
	Stages_MedicalTeamExamsScheduled     = "Medical Team - Exams Scheduled"
	Stages_AmMedicalTeamExamsScheduled   = "Am__Medical Team - Exams Scheduled"

	Stages_MedicalTeamCallVet   = "Medical Team - Call Vet"
	Stages_AmMedicalTeamCallVet = "Am__Medical Team - Call Vet"

	Stages_MedicalTeamPrefilledFormsReview   = "Medical Team - Prefilled Forms Review"
	Stages_AmMedicalTeamPrefilledFormsReview = "Am__Medical Team - Prefilled Forms Review"

	Stages_DBQ_Completed   = "DBQ Completed"
	Stages_AmDBQ_Completed = "Am__DBQ Completed"

	Stages_StatementFinalChanges   = "Statement Final Changes"
	Stages_AmStatementFinalChanges = "Am__Statement Final Changes"

	Stages_FileClaims_Draft   = "File Claims Draft"
	Stages_AmFileClaims_Draft = "Am__File Claims Draft"
	Stages_FileHLRDraft       = "File HLR Draft"
	Stages_AmFileHLRDraft     = "Am__File HLR Draft"
	Stages_FileClaims         = "File Claims with Client"
	Stages_AmFileClaims       = "Am__File Claims with Client"
	Stages_FileHLRWithClient  = "File HLR with Client"

	Stages_VerifyEvidenceReceived   = "Verify Evidence Received"
	Stages_AmVerifyEvidenceReceived = "Am__Verify Evidence Received"

	Stages_AwaitingDecision   = "Awaiting Decision"
	Stages_AmAwaitingDecision = "Am__Awaiting Decision"

	Stages_AwaitingPayment   = "Awaiting Payment"
	Stages_AmAwaitingPayment = "Am__Awaiting Payment"

	Stages_27_AwaitingBankReconciliation   = "27. Awaiting Bank Reconciliation"
	Stages_Am27_AwaitingBankReconciliation = "Am__27. Awaiting Bank Reconciliation"

	Stages_Completed   = "Completed"
	Stages_AmCompleted = "Am__Completed"

	Stages_Terminated   = "Terminated"
	Stages_Dormant      = "Dormant"
	Stages_AmTerminated = "Am__Terminated"
	Stages_AmDormant    = "Am__Dormant"
)

func ShouldSkipContractUpdate(stages string) bool {
	if stages == Stages_AmIncomingRequest ||
		stages == Stages_IncomingRequest ||
		stages == Stages_Terminated ||
		stages == Stages_Dormant ||
		stages == Stages_AmTerminated ||
		stages == Stages_AmDormant {
		return true
	}
	return false
}

func GetStageNumber(stageOptionLabel string) (int, error) {
	aa := strings.Split(stageOptionLabel, ".")
	if len(aa) <= 1 {
		return 0, errors.New("stageOptionLabel is wrong")
	}
	num := strings.TrimSpace(aa[0])
	c, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(c), nil
}

func StagesToNumber_Deprecated(stages string) (r int, err error) {
	a := strings.Split(stages, ".")
	if len(a) <= 1 {
		return 0, errors.New("StagesToNumber error.")
	}
	res, err := strconv.ParseInt(a[0], 10, 32)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

var StagesSorts = []string{
	Stages_FeeScheduleandContract,
	Stages_GettingStartedEmail,
	Stages_AwaitingClientRecords,
	Stages_RecordReview,
	Stages_ScheduleCall,

	Stages_StatementNotes,
	Stages_StatementDrafts,
	Stages_StatementReview,
	Stages_StatementsFinalized,
	Stages_CurrentTreatment,
	Stages_AwaitingNexusLetter,

	Stages_MiniDBQs_Draft,
	Stages_MiniDBQs,
	Stages_MiniDBQ_Forms,
	Stages_MedicalTeamFormsSigned,
	Stages_MedicalTeam,
	Stages_MedicalTeamPaymentCollected,
	Stages_MedicalTeamExamsScheduled,
	Stages_MedicalTeamCallVet,

	Stages_DBQ_Completed,
	Stages_FileClaims_Draft,
	Stages_FileClaims,
	Stages_VerifyEvidenceReceived,
	Stages_AwaitingDecision,

	Stages_AwaitingPayment,
	Stages_Completed,
}

// GetUnderStages 获取比传参stage更低的stages
func GetUnderStages(stage string) (r []string) {
	for _, v := range StagesSorts {
		if v == stage {
			return
		}
		r = append(r, v)
	}
	return
}

const (
	TaskSubject_ContractFollowUp                          = "Contract Follow-up"
	TaskSubject_WelcomeEmailFollowUp                      = "Welcome Email Follow-up"
	TaskSubject_ClientRecordsFollowUp                     = "Client Records Follow-up"
	TaskSubject_FollowUpOnSTRRecordsStatus                = "Follow-up on STR Records Status"
	TaskSubject_RecordReview                              = "Record Review"
	TaskSubject_ScheduleClientCall                        = "Schedule Client Call"
	TaskSubject_PrepareStatementNotes                     = "Prepare Statement Notes"
	TaskSubject_PrepareStatementDrafts                    = "Prepare Statement Drafts"
	TaskSubject_ReviewStatements                          = "Review Statements"
	TaskSubject_FinalizeStatements                        = "Finalize Statements"
	TaskSubject_CurrentTreatmentFollowUp                  = "Current Treatment Follow-up"
	TaskSubject_PrepareDocumentsForTinnitusLetter         = "Prepare documents for Tinnitus letter"
	TaskSubject_PreparingNexusLetter                      = "Preparing Nexus Letter"
	TaskSubject_CheckNexusLetterStatus                    = "Check Nexus Letter Status"
	TaskSubject_MiniDBQDraftsPreparation                  = "Mini-DBQ Drafts Preparation"
	TaskSubject_FinalizeMiniDQBs                          = "Finalize Mini-DQBs"
	TaskSubject_CheckVetSignedbothMedTeamForms            = "Check Vet signed both Med Team Forms"
	TaskSubject_ReviewSignedDocuments                     = "Review Signed Documents"
	TaskSubject_VerifyFilesUploadedCorrectly              = "Med Team Status - Verify Files Uploaded Correctly"
	TaskSubject_FollowupwithVetforInvoicePayment          = "Follow-up with Vet for Invoice Payment"
	TaskSubject_VerifythatVethasbeenscheduledamedicalexam = "Verify that Vet has been scheduled a medical exam"
	TaskSubject_CallVetForMedicalTeam                     = "Call Vet For Medical Team"
	TaskSubject_DBQCompleted                              = "DBQ Completed"
	TaskSubject_DraftingDocumentsforClaimSubmission       = "Drafting Documents for Claim Submission"
	TaskSubject_FileClaimswithClient                      = "File Claims with Client"
	TaskSubject_VerifyEvidenceReceived                    = "Verify Evidence Received"
	TaskSubject_AwaitingDecisionfromVA                    = "Awaiting Decision from VA"
	TaskSubject_CheckInvoicePaymentStatus                 = "Check Invoice Payment Status"
	TaskSubject_CheckPaymentTransaction                   = "Check that payment hits the bank and reconcile the transaction"
)

type AutomaticCreationTaskSubjectItem struct {
	Subject       string
	PlusDays      int
	AssignUserGid string // 有值时，强制使用
}

// AutomaticCreationTaskSubjectRelationStages 虽然不使用tasks，但依赖此处配置更新due date
var AutomaticCreationTaskSubjectRelationStages = map[string][]AutomaticCreationTaskSubjectItem{
	Stages_IncomingRequest:                  {{"", 0, ""}},
	Stages_FeeScheduleandContract:           {{TaskSubject_ContractFollowUp, 7, ""}},
	Stages_GettingStartedEmail:              {{TaskSubject_WelcomeEmailFollowUp, 1, ""}},
	Stages_AwaitingClientRecords:            {{TaskSubject_ClientRecordsFollowUp, 1, ""}},
	Stages_STRRequestPending:                {{TaskSubject_FollowUpOnSTRRecordsStatus, 30, ""}},
	Stages_RecordReview:                     {{TaskSubject_RecordReview, 3, ""}},
	Stages_ClaimAnalysis:                    {{"", 1, ""}},
	Stages_ScheduleCall:                     {{TaskSubject_ScheduleClientCall, 1, ""}},
	Stages_StatementNotes:                   {{TaskSubject_PrepareStatementNotes, 3, ""}},
	Stages_StatementDrafts:                  {{TaskSubject_PrepareStatementDrafts, 2, ""}},
	Stages_StatementReview:                  {{TaskSubject_ReviewStatements, 3, ""}},
	Stages_StatementsFinalized:              {{TaskSubject_FinalizeStatements, 7, ""}},
	Stages_CurrentTreatment:                 {{TaskSubject_CurrentTreatmentFollowUp, 30, ""}},
	Stages_PreparingDocumentsTinnitusLetter: {{TaskSubject_PrepareDocumentsForTinnitusLetter, 2, ""}},
	Stages_AwaitingNexusLetter:              {{TaskSubject_PreparingNexusLetter, 7, ""}}, //{TaskSubject_CheckNexusLetterStatus, 28, ""}
	Stages_MiniDBQs_Draft:                   {{TaskSubject_MiniDBQDraftsPreparation, 1, ""}},
	Stages_MiniDBQs:                         {{TaskSubject_FinalizeMiniDQBs, 1, ""}},
	Stages_MiniDBQ_Forms:                    {{TaskSubject_CheckVetSignedbothMedTeamForms, 1, ""}},
	Stages_MedicalTeamFormsSigned:           {{TaskSubject_ReviewSignedDocuments, 1, ""}},
	Stages_MedicalTeam:                      {{TaskSubject_VerifyFilesUploadedCorrectly, 1, ""}},
	Stages_MedicalTeamPaymentCollected:      {{TaskSubject_FollowupwithVetforInvoicePayment, 5, ""}},
	Stages_MedicalTeamExamsScheduled:        {{TaskSubject_VerifythatVethasbeenscheduledamedicalexam, 7, ""}},
	Stages_MedicalTeamCallVet:               {{TaskSubject_CallVetForMedicalTeam, 21, ""}},
	Stages_DBQ_Completed:                    {{TaskSubject_DBQCompleted, 21, ""}},
	Stages_FileClaims_Draft:                 {{TaskSubject_DraftingDocumentsforClaimSubmission, 1, ""}},
	Stages_FileClaims:                       {{TaskSubject_FileClaimswithClient, 1, ""}},
	Stages_VerifyEvidenceReceived:           {{TaskSubject_VerifyEvidenceReceived, 10, ""}}, // "ITF Expiration within 90 days"
	Stages_AwaitingDecision:                 {{TaskSubject_AwaitingDecisionfromVA, 45, ""}},
	Stages_AwaitingPayment:                  {{TaskSubject_CheckInvoicePaymentStatus, 0, User_Edward_gid}},
	Stages_27_AwaitingBankReconciliation:    {{TaskSubject_CheckPaymentTransaction, 2, User_Victoria_gid}},

	Stages_AmIncomingRequest:   {{"", 0, ""}},
	Stages_AmInformationIntake: {{TaskSubject_ContractFollowUp, 1, ""}},
	Stages_AmContractPending:   {{TaskSubject_WelcomeEmailFollowUp, 1, ""}},
	//Stages_AmSendVA2122a:                      {{"", 1, ""}},
	Stages_AmAwaitingClientRecords:            {{TaskSubject_ClientRecordsFollowUp, 1, ""}},
	Stages_AmSTRRequestPending:                {{TaskSubject_FollowUpOnSTRRecordsStatus, 30, ""}},
	Stages_AmRecordReview:                     {{TaskSubject_RecordReview, 3, ""}},
	Stages_AmClaimAnalysis:                    {{"", 1, ""}},
	Stages_AmScheduleCall:                     {{TaskSubject_ScheduleClientCall, 1, ""}},
	Stages_AmStatementNotes:                   {{TaskSubject_PrepareStatementNotes, 3, ""}},
	Stages_AmStatementDrafts:                  {{TaskSubject_PrepareStatementDrafts, 2, ""}},
	Stages_AmStatementReview:                  {{TaskSubject_ReviewStatements, 3, ""}},
	Stages_AmStatementsFinalized:              {{TaskSubject_FinalizeStatements, 7, ""}},
	Stages_AmCurrentTreatment:                 {{TaskSubject_CurrentTreatmentFollowUp, 30, ""}},
	Stages_AmPreparingDocumentsTinnitusLetter: {{TaskSubject_PrepareDocumentsForTinnitusLetter, 2, ""}},
	Stages_AmAwaitingNexusLetter:              {{TaskSubject_PreparingNexusLetter, 7, ""}}, //{TaskSubject_CheckNexusLetterStatus, 28, ""}
	Stages_AmMiniDBQs_Draft:                   {{TaskSubject_MiniDBQDraftsPreparation, 1, ""}},
	Stages_AmMiniDBQs:                         {{TaskSubject_FinalizeMiniDQBs, 1, ""}},
	Stages_AmMiniDBQ_Forms:                    {{TaskSubject_CheckVetSignedbothMedTeamForms, 1, ""}},
	Stages_AmMedicalTeamFormsSigned:           {{TaskSubject_ReviewSignedDocuments, 1, ""}},
	Stages_AmMedicalTeam:                      {{TaskSubject_VerifyFilesUploadedCorrectly, 1, ""}},
	Stages_AmMedicalTeamPaymentCollected:      {{TaskSubject_FollowupwithVetforInvoicePayment, 5, ""}},
	Stages_AmMedicalTeamExamsScheduled:        {{TaskSubject_VerifythatVethasbeenscheduledamedicalexam, 7, ""}},
	Stages_AmMedicalTeamCallVet:               {{TaskSubject_CallVetForMedicalTeam, 21, ""}},
	Stages_AmDBQ_Completed:                    {{TaskSubject_DBQCompleted, 21, ""}},
	Stages_AmFileClaims_Draft:                 {{TaskSubject_DraftingDocumentsforClaimSubmission, 1, ""}},
	Stages_AmFileClaims:                       {{TaskSubject_FileClaimswithClient, 1, ""}},
	Stages_AmVerifyEvidenceReceived:           {{TaskSubject_VerifyEvidenceReceived, 10, ""}}, // "ITF Expiration within 90 days"
	Stages_AmAwaitingDecision:                 {{TaskSubject_AwaitingDecisionfromVA, 45, ""}},
	Stages_AmAwaitingPayment:                  {{TaskSubject_CheckInvoicePaymentStatus, 0, User_Edward_gid}},
	Stages_Am27_AwaitingBankReconciliation:    {{TaskSubject_CheckPaymentTransaction, 2, User_Victoria_gid}},
}

var ZohoTaskSubjectRelationStages = map[string][]string{
	Stages_FeeScheduleandContract:           {TaskSubject_ContractFollowUp},
	Stages_GettingStartedEmail:              {TaskSubject_WelcomeEmailFollowUp},
	Stages_AwaitingClientRecords:            {TaskSubject_ClientRecordsFollowUp},
	Stages_STRRequestPending:                {TaskSubject_FollowUpOnSTRRecordsStatus},
	Stages_RecordReview:                     {TaskSubject_RecordReview},
	Stages_ScheduleCall:                     {TaskSubject_ScheduleClientCall},
	Stages_StatementNotes:                   {TaskSubject_PrepareStatementNotes},
	Stages_StatementDrafts:                  {TaskSubject_PrepareStatementDrafts},
	Stages_StatementReview:                  {TaskSubject_ReviewStatements},
	Stages_StatementsFinalized:              {TaskSubject_FinalizeStatements},
	Stages_CurrentTreatment:                 {TaskSubject_CurrentTreatmentFollowUp},
	Stages_PreparingDocumentsTinnitusLetter: {TaskSubject_PrepareDocumentsForTinnitusLetter},
	Stages_AwaitingNexusLetter:              {TaskSubject_PreparingNexusLetter, TaskSubject_CheckNexusLetterStatus},
	Stages_MiniDBQs_Draft:                   {TaskSubject_MiniDBQDraftsPreparation},
	Stages_MiniDBQs:                         {TaskSubject_FinalizeMiniDQBs},
	Stages_MiniDBQ_Forms:                    {TaskSubject_CheckVetSignedbothMedTeamForms},
	Stages_MedicalTeamFormsSigned:           {TaskSubject_ReviewSignedDocuments},
	Stages_MedicalTeam:                      {TaskSubject_VerifyFilesUploadedCorrectly},
	Stages_MedicalTeamPaymentCollected:      {TaskSubject_FollowupwithVetforInvoicePayment},
	Stages_MedicalTeamExamsScheduled:        {TaskSubject_VerifythatVethasbeenscheduledamedicalexam},
	Stages_MedicalTeamCallVet:               {TaskSubject_CallVetForMedicalTeam},
	Stages_DBQ_Completed:                    {TaskSubject_DBQCompleted},
	Stages_FileClaims_Draft:                 {TaskSubject_DraftingDocumentsforClaimSubmission},
	Stages_FileClaims:                       {TaskSubject_FileClaimswithClient},
	Stages_VerifyEvidenceReceived:           {TaskSubject_VerifyEvidenceReceived}, // "ITF Expiration within 90 days"
	Stages_AwaitingDecision:                 {TaskSubject_AwaitingDecisionfromVA},
	Stages_AwaitingPayment:                  {TaskSubject_CheckInvoicePaymentStatus},
	Stages_27_AwaitingBankReconciliation:    {TaskSubject_CheckPaymentTransaction},
}

// 12. MedTeam Forms : Verify that Vet has been scheduled a medical exam
// 16. Medical Team - Exams Scheduled : Check Vet signed both Med Team Forms

// JudgeTaskNeedCompleteBySubject 通过任务名称和当前stage 判断此任务是否需要关闭 true: 需要关闭；false：不需要关闭
func JudgeTaskNeedCompleteBySubject(taskSubject string, stage string) bool {

	if taskSubject == "Contact Lead and move the Lead forward" {
		return true
	}
	stages := GetUnderStages(stage)
	for _, v := range stages {
		if vals, ok := ZohoTaskSubjectRelationStages[v]; ok {
			for _, v1 := range vals {
				if strings.Index(taskSubject, v1) >= 0 {
					return true
				}
			}
		}
	}
	return false
}

// JudgeTaskNeedCompleteBySubjectSpecify 通过任务名称和当前stage 判断此任务是否需要关闭 true: 需要关闭；false：不需要关闭
func JudgeTaskNeedCompleteBySubjectSpecify(taskSubject string, stage string) bool {

	if taskSubject == "Contact Lead and move the Lead forward" {
		return true
	}
	if vals, ok := ZohoTaskSubjectRelationStages[stage]; ok {
		for _, v1 := range vals {
			if strings.Index(taskSubject, v1) >= 0 {
				return true
			}
		}
	}
	return false
}

// JudgeTaskWhetherBelongsStage 通过任务名称判断是否属于此状态
func JudgeTaskWhetherBelongsStage(taskSubject string, stage string) bool {

	if vals, ok := ZohoTaskSubjectRelationStages[stage]; ok {
		for _, v1 := range vals {
			if strings.Index(taskSubject, v1) >= 0 {
				return true
			}
		}
	}
	return false
}
