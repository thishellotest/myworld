package biz

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/internal/utils"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/esign/v2.1/envelopes"
	"vbc/lib/uuid"
)

const (
	Task_Dag_BuzEmail                             = "Dag.BuzEmail"
	Task_Dag_CreateEnvelopeAndSentFromBoxAm       = "Dag.CreateEnvelopeAndSentFromBoxAm"
	Task_Dag_CreateEnvelopeAndSent                = "Dag.CreateEnvelopeAndSent"
	Task_Dag_BoxCreateClientContracts             = "Dag.BoxCreateClientContracts"
	Task_Dag_BoxCreateFolderForNewClient          = "Dag.BoxCreateFolderForNewClient"
	Task_Dag_GetEnvelopeDocuments                 = "Dag.GetEnvelopeDocuments"
	Task_Dag_SaveSignedContractInBox              = "Dag.SaveSignedContractInBox"
	Task_Dag_ReminderMedicalTeamFormsContractSent = "Dag.ReminderMedicalTeamFormsContractSent"
	Task_Dag_HandleReminder                       = "Dag.HandleReminder"
	Task_Dag_HandleContractReminder               = "Dag.HandleContractReminder"
	Task_Dag_ContractNonResponsive                = "Dag.ContractNonResponsive"
	Task_Dag_AmContractNonResponsive              = "Dag.AmContractNonResponsive"
	Task_Dag_CronTrigger                          = "Dag.CronTrigger" // 定时触发器
)

type Dag struct {
	log                         *log.Helper
	TUsecase                    *TUsecase
	MailUsecase                 *MailUsecase
	CommonUsecase               *CommonUsecase
	DocuSignUsecase             *DocuSignUsecase
	ClientEnvelopeUsecase       *ClientEnvelopeUsecase
	BoxUsecase                  *BoxUsecase
	MapUsecase                  *MapUsecase
	conf                        *conf.Data
	EnvelopeStatusChangeUsecase *EnvelopeStatusChangeUsecase
	TaskCreateUsecase           *TaskCreateUsecase
	AdobeSignUsecase            *AdobeSignUsecase
	AdobeWebhookEventUsecase    *AdobeWebhookEventUsecase
	ClientAgreementUsecase      *ClientAgreementUsecase
	BoxcontractUsecase          *BoxcontractUsecase
	RollpoingUsecase            *RollpoingUsecase
	DataComboUsecase            *DataComboUsecase
	ZohoUsecase                 *ZohoUsecase
	UserUsecase                 *UserUsecase
	RemindUsecase               *RemindUsecase
	BehaviorUsecase             *BehaviorUsecase
	DialpadbuzUsecase           *DialpadbuzUsecase
	CronTriggerUsecase          *CronTriggerUsecase
	StageTransUsecase           *StageTransUsecase
	GopdfUsecase                *GopdfUsecase
	DataEntryUsecase            *DataEntryUsecase
	AttorneybuzUsecase          *AttorneybuzUsecase
	AttorneyUsecase             *AttorneyUsecase
	AmUsecase                   *AmUsecase
	PersonalWebformUsecase      *PersonalWebformUsecase
	BoxbuzUsecase               *BoxbuzUsecase
	BoxCollaborationBuzUsecase  *BoxCollaborationBuzUsecase
}

func NewDag(CommonUsecase *CommonUsecase, TUsecase *TUsecase,
	logger log.Logger,
	MailUsecase *MailUsecase,
	DocuSignUsecase *DocuSignUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	MapUsecase *MapUsecase,
	BoxUsecase *BoxUsecase,
	conf *conf.Data,
	EnvelopeStatusChangeUsecase *EnvelopeStatusChangeUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	AdobeSignUsecase *AdobeSignUsecase,
	AdobeWebhookEventUsecase *AdobeWebhookEventUsecase,
	ClientAgreementUsecase *ClientAgreementUsecase,
	BoxcontractUsecase *BoxcontractUsecase,
	RollpoingUsecase *RollpoingUsecase,
	DataComboUsecase *DataComboUsecase,
	ZohoUsecase *ZohoUsecase,
	UserUsecase *UserUsecase,
	RemindUsecase *RemindUsecase,
	BehaviorUsecase *BehaviorUsecase,
	DialpadbuzUsecase *DialpadbuzUsecase,
	CronTriggerUsecase *CronTriggerUsecase,
	StageTransUsecase *StageTransUsecase,
	GopdfUsecase *GopdfUsecase,
	DataEntryUsecase *DataEntryUsecase,
	AttorneybuzUsecase *AttorneybuzUsecase,
	AttorneyUsecase *AttorneyUsecase,
	AmUsecase *AmUsecase,
	PersonalWebformUsecase *PersonalWebformUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	BoxCollaborationBuzUsecase *BoxCollaborationBuzUsecase) *Dag {
	return &Dag{
		log:                         log.NewHelper(logger),
		CommonUsecase:               CommonUsecase,
		TUsecase:                    TUsecase,
		MailUsecase:                 MailUsecase,
		DocuSignUsecase:             DocuSignUsecase,
		ClientEnvelopeUsecase:       ClientEnvelopeUsecase,
		BoxUsecase:                  BoxUsecase,
		MapUsecase:                  MapUsecase,
		conf:                        conf,
		EnvelopeStatusChangeUsecase: EnvelopeStatusChangeUsecase,
		TaskCreateUsecase:           TaskCreateUsecase,
		AdobeSignUsecase:            AdobeSignUsecase,
		AdobeWebhookEventUsecase:    AdobeWebhookEventUsecase,
		ClientAgreementUsecase:      ClientAgreementUsecase,
		BoxcontractUsecase:          BoxcontractUsecase,
		RollpoingUsecase:            RollpoingUsecase,
		DataComboUsecase:            DataComboUsecase,
		ZohoUsecase:                 ZohoUsecase,
		UserUsecase:                 UserUsecase,
		RemindUsecase:               RemindUsecase,
		BehaviorUsecase:             BehaviorUsecase,
		DialpadbuzUsecase:           DialpadbuzUsecase,
		CronTriggerUsecase:          CronTriggerUsecase,
		StageTransUsecase:           StageTransUsecase,
		GopdfUsecase:                GopdfUsecase,
		DataEntryUsecase:            DataEntryUsecase,
		AttorneybuzUsecase:          AttorneybuzUsecase,
		AttorneyUsecase:             AttorneyUsecase,
		AmUsecase:                   AmUsecase,
		PersonalWebformUsecase:      PersonalWebformUsecase,
		BoxbuzUsecase:               BoxbuzUsecase,
		BoxCollaborationBuzUsecase:  BoxCollaborationBuzUsecase,
	}
}

const Sign_type_box = "box"
const Sign_type_adobe = "adobe"

func (c *Dag) CreateEnvelopeAndSentFromAdobeSign(task *TaskEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	var taskInput lib.TypeMap
	taskInput = lib.ToTypeMapByString(task.TaskInput)

	templateId := InterfaceToString(taskInput.Get("templateId"))
	clientFirstName := InterfaceToString(taskInput.Get("clientFirstName"))
	clientLastName := InterfaceToString(taskInput.Get("clientLastName"))
	clientEmail := InterfaceToString(taskInput.Get("clientEmail"))
	agentFirstName := InterfaceToString(taskInput.Get("agentFirstName"))
	agentLastName := InterfaceToString(taskInput.Get("agentLastName"))
	agentEmail := InterfaceToString(taskInput.Get("agentEmail"))

	if len(templateId) == 0 || len(clientEmail) == 0 || len(agentEmail) == 0 {
		return errors.New("Parameters is wrong.")
	}
	agreementId, err := c.AdobeSignUsecase.CreateAgreement("Your Veteran Benefits Center Contract",
		templateId,
		CreateAgreementMember{
			Email:     clientEmail,
			FirstName: clientFirstName,
			LastName:  clientLastName,
		},
		CreateAgreementMember{
			Email:     agentEmail,
			FirstName: agentFirstName,
			LastName:  agentLastName,
		})
	if err != nil {
		return err
	}

	entity := ClientAgreementEntity{
		ClientId:    task.IncrId,
		AgreementId: agreementId,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	return c.CommonUsecase.DB().Save(&entity).Error
}

type CreateEnvelopeTaskInput struct {
	ContractIndex   int    `json:"contractIndex"`
	SignType        string `json:"signType"`
	TemplateId      string `json:"templateId"`
	ClientFirstName string `json:"clientFirstName"`
	ClientLastName  string `json:"clientLastName"`
	ClientEmail     string `json:"clientEmail"`
	AgentFirstName  string `json:"agentFirstName"`
	AgentLastName   string `json:"agentLastName"`
	AgentEmail      string `json:"agentEmail"`
}

func (c *CreateEnvelopeTaskInput) Verify() error {
	if c.TemplateId == "" || c.ClientEmail == "" || c.AgentEmail == "" {
		return errors.New("CreateEnvelopeTaskInput Verify error.")
	}
	return nil
}

func (c *Dag) CreateEnvelopeAndSentFromBox(task *TaskEntity) error {

	return c.CreateEnvelopeAndSentFromBoxWithTemplate(task)
	if task == nil {
		return errors.New("task is nil.")
	}

	createEnvelopeTaskInput := lib.StringToTDef[*CreateEnvelopeTaskInput](task.TaskInput, nil)
	if createEnvelopeTaskInput == nil {
		return errors.New("CreateEnvelopeAndSentFromBox:createEnvelopeTaskInput is nil")
	}
	if err := createEnvelopeTaskInput.Verify(); err != nil {
		return err
	}
	contractFolderId, err := c.BoxcontractUsecase.ContractFolderId(task.IncrId)
	if err != nil {
		return err
	}

	//if createEnvelopeTaskInput.ClientEmail == "lialing@foxmail.com" {
	//	c.log.Info("mock email: lialing@foxmail.com, does not really create box sign.")
	//	return nil
	//}

	res, contractId, err := c.BoxUsecase.SignRequests(createEnvelopeTaskInput, contractFolderId)
	if err != nil {
		return err
	}
	str := ""
	if res != nil {
		str = *res
	}
	err = c.ClientEnvelopeUsecase.Add(task.IncrId, EsignVendor_box, contractId, str, Type_FeeContract, 0)
	if err == nil {
		err = c.RollpoingUsecase.Upsert(Rollpoing_Vendor_boxsign, contractId)
	}
	return err
}

func (c *Dag) CreateEnvelopeAndSentFromBoxWithTemplate(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}

	createEnvelopeTaskInput := lib.StringToTDef[*CreateEnvelopeTaskInput](task.TaskInput, nil)
	if createEnvelopeTaskInput == nil {
		return errors.New("CreateEnvelopeAndSentFromBox:createEnvelopeTaskInput is nil")
	}
	if err := createEnvelopeTaskInput.Verify(); err != nil {
		return err
	}
	contractFolderId, err := c.BoxcontractUsecase.ContractFolderId(task.IncrId)
	if err != nil {
		return err
	}
	signFileBytes, err := c.GopdfUsecase.CreateContract(CreateContractVo{
		ClientName:  GenFullName(createEnvelopeTaskInput.ClientFirstName, createEnvelopeTaskInput.ClientLastName),
		ClientEmail: createEnvelopeTaskInput.ClientEmail,
		VsName:      GenFullName(createEnvelopeTaskInput.AgentFirstName, createEnvelopeTaskInput.AgentLastName),
		VsEmail:     createEnvelopeTaskInput.AgentEmail,
	}, createEnvelopeTaskInput.ContractIndex)
	if err != nil {
		c.log.Error(err)
		return err
	}
	folderName := "Contract_" + uuid.UuidWithoutStrike()
	signFolderId, err := c.BoxUsecase.CreateFolder(folderName, contractFolderId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	signFileId, err := c.BoxUsecase.UploadFile(signFolderId, bytes.NewReader(signFileBytes), "Agreement for Consulting Services.pdf")
	if err != nil {
		return err
	}
	res, contractId, err := c.BoxUsecase.SignRequestsWithoutTemplate(createEnvelopeTaskInput, contractFolderId, signFileId)
	if err != nil {
		return err
	}
	str := ""
	if res != nil {
		str = *res
	}
	err = c.ClientEnvelopeUsecase.Add(task.IncrId, EsignVendor_box, contractId, str, Type_FeeContract, 0)
	if err == nil {
		err = c.RollpoingUsecase.Upsert(Rollpoing_Vendor_boxsign, contractId)
	}
	return err
}

func (c *Dag) CreateEnvelopeAndSentFromBoxAm(task *TaskEntity) (err error) {

	if task == nil {
		return errors.New("task is nil")
	}
	caseId := task.IncrId
	return c.DoCreateEnvelopeAndSentFromBoxAm(caseId)
}

func (c *Dag) DoCreateEnvelopeAndSentFromBoxAm(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	tClient, _, err := c.DataComboUsecase.ClientWithCase(*tCase)
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	contractVetVo := GenContractVetVo(*tClient, *tCase)

	attorney, err := c.AmUsecase.DoAttorney(*tCase)
	if err != nil {
		return err
	}
	if attorney == nil {
		return errors.New("attorney is nil")
	}
	contractAttorneyVo := attorney.ToContractAttorneyVo()
	contractTime := time.Now()
	signFileBytes, err := c.GopdfUsecase.CreateContractAm(contractTime, contractVetVo, contractAttorneyVo)
	if err != nil {
		return err
	}

	contractFolderId, err := c.BoxcontractUsecase.ContractFolderId(caseId)
	if err != nil {
		return err
	}

	folderName := "Contract_" + uuid.UuidWithoutStrike()
	signFolderId, err := c.BoxUsecase.CreateFolder(folderName, contractFolderId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	signFileId, err := c.BoxUsecase.UploadFile(signFolderId, bytes.NewReader(signFileBytes), "Your VA Representation Agreement with August Miles.pdf")
	if err != nil {
		return err
	}
	contractAttorneyEmail := ""

	if contractVetVo.Email == "lialing@foxmail.com" {
		contractAttorneyEmail = "glliao@vetbenefitscenter.com"
		//contractAttorneyEmail = "mvplinchen888@gmail.com"
	} else {
		contractAttorneyEmail = contractAttorneyVo.Email
	}
	res, contractId, err := c.BoxUsecase.SignRequestsWithoutTemplateAm(contractVetVo.Email, contractAttorneyEmail, contractAttorneyVo.FullName(), contractFolderId, signFileId)
	if err != nil {
		return err
	}
	str := ""
	if res != nil {
		str = *res
	}
	err = c.ClientEnvelopeUsecase.Add(caseId, EsignVendor_box, contractId, str, Type_AmContract, attorney.ID)
	if err == nil {
		err = c.RollpoingUsecase.Upsert(Rollpoing_Vendor_boxsign, contractId)
	}
	return err
}

func (c *Dag) CreateEnvelopeAndSent(task *TaskEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	var taskInput lib.TypeMap
	taskInput = lib.ToTypeMapByString(task.TaskInput)
	signType := InterfaceToString(taskInput.Get("signType"))
	if signType == Sign_type_box {
		// Current: Support Box Sign only
		return c.CreateEnvelopeAndSentFromBox(task)
	} else if signType == Sign_type_adobe {
		return c.CreateEnvelopeAndSentFromAdobeSign(task)
	}
	docusignTemplateId := InterfaceToString(taskInput.Get("docusignTemplateId"))
	clientName := InterfaceToString(taskInput.Get("clientName"))
	clientEmail := InterfaceToString(taskInput.Get("clientEmail"))
	agentName := InterfaceToString(taskInput.Get("agentName"))
	agentEmail := InterfaceToString(taskInput.Get("agentEmail"))

	if len(docusignTemplateId) == 0 || len(clientEmail) == 0 || len(agentEmail) == 0 {
		return errors.New("Parameters is wrong.")
	}
	envelopeSummary, err := c.DocuSignUsecase.CreateEnvelopeAndSent(docusignTemplateId, clientName, clientEmail, agentName, agentEmail)
	if err != nil {
		return err
	}
	if envelopeSummary == nil {
		return errors.New("envelopeSummary is nil.")
	}

	entity := ClientEnvelopeEntity{
		ClientId:       task.IncrId,
		EnvelopeId:     envelopeSummary.EnvelopeID,
		Uri:            envelopeSummary.URI,
		Status:         envelopeSummary.Status,
		StatusDatetime: envelopeSummary.StatusDateTime.Format(time.RFC3339),
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}
	return c.ClientEnvelopeUsecase.CommonUsecase.DB().Save(&entity).Error
}

func (c *Dag) BuzEmailCustom(MailTaskInput *MailTaskInput) error {

	if MailTaskInput == nil {
		return errors.New("BuzEmailCustom:MailTaskInput is nil")
	}
	return c.MailUsecase.SendEmail(InitMailServiceConfig(), MailTaskInput.MailMessage, MailAttach_No, nil)
}

func (c *Dag) BuzEmail(task *TaskEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	mailTaskInput := lib.StringToTDef[*MailTaskInput](task.TaskInput, nil)
	if mailTaskInput == nil {
		return errors.New("mailTaskInput is nil.")
	}
	if mailTaskInput.Genre == MailGenre_Custom {
		return c.BuzEmailCustom(mailTaskInput)
	}
	tCustomer, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": task.IncrId})
	if err != nil {
		return err
	}

	tpl := mailTaskInput.Genre
	subId := mailTaskInput.SubId
	tTpl, err := c.TUsecase.Data(Kind_email_tpls, Eq{"tpl": tpl, "sub_id": subId})
	if err != nil {
		return err
	}
	if tCustomer != nil && tTpl != nil {
		if tTpl.CustomFields.TextValueByNameBasic("tpl") == MailGenre_PersonalStatementsReadyforYourReview ||
			tTpl.CustomFields.TextValueByNameBasic("tpl") == MailGenre_PleaseReviewYourPersonalStatementsinSharedFolder {

			useNewEmailTpl, err := c.PersonalWebformUsecase.IsUseNewPersonalWebForm(tCustomer.Id())
			if err != nil {
				return err
			}
			if useNewEmailTpl {
				devTpl := tTpl.CustomFields.TextValueByNameBasic("tpl") + "_ForWebForm"
				devTTpl, er := c.TUsecase.Data(Kind_email_tpls, Eq{"tpl": devTpl, "sub_id": subId})
				if er != nil {
					c.log.Error(devTpl, er)
				}
				if devTTpl != nil {
					subject := devTTpl.CustomFields.TextValueByNameBasic("subject")
					body := devTTpl.CustomFields.TextValueByNameBasic("body")
					tTpl.CustomFields.SetTextValueByName("subject", &subject)
					tTpl.CustomFields.SetTextValueByName("body", &body)
				} else {
					c.log.Error("devTTpl is nil")
				}
			}
		}
	}

	err, mailSubject, mailBody, email, senderEmail, senderName := c.MailUsecase.SendEmailWithData(tCustomer, tTpl, mailTaskInput)
	if err != nil {
		return err
	}
	task.TaskStatus = Task_TaskStatus_finish
	return c.CommonUsecase.DB().Save(&EmailLogEntity{
		ClientId:   task.IncrId,
		Email:      email,
		TaskId:     task.ID,
		Tpl:        tpl,
		SubId:      subId,
		SenderMail: senderEmail,
		SenderName: senderName,
		Subject:    mailSubject,
		Body:       mailBody,
	}).Error
}

func (c *Dag) BoxCreateClientContractsAdobe(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	taskInput := lib.ToTypeMapByString(task.TaskInput)
	AdobeWebhookEventId := taskInput.GetInt("AdobeWebhookEventId")
	if AdobeWebhookEventId == 0 {
		return errors.New("AdobeWebhookEventId is 0.")
	}
	var entity *AdobeWebhookEventEntity
	var err error
	entity, err = c.AdobeWebhookEventUsecase.GetByCond(Eq{"id": AdobeWebhookEventId})
	if err != nil {
		return err
	}
	if entity == nil {
		return errors.New("AdobeWebhookEventEntity is nil.")
	}

	clientAgreement, err := c.ClientAgreementUsecase.GetByCond(Eq{"agreement_id": entity.AgreementId})
	if err != nil {
		return err
	}
	if clientAgreement == nil {
		return errors.New("clientAgreement is nil.")
	}

	boxFolderId, err := c.BoxcontractUsecase.ContractFolderId(clientAgreement.ClientId)
	if err != nil {
		return err
	}

	err = c.TaskCreateUsecase.CreateTask(AdobeWebhookEventId,
		map[string]interface{}{
			"type":         Sign_type_adobe,
			"agreement_id": entity.AgreementId,
			"folder_id":    boxFolderId},
		Task_Dag_SaveSignedContractInBox, 0, "", "")

	return err
}

func (c *Dag) BoxCreateClientContracts(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}

	taskInput := lib.ToTypeMapByString(task.TaskInput)
	taskType := taskInput.GetString("type")
	if taskType == Sign_type_adobe {
		return c.BoxCreateClientContractsAdobe(task)
	}

	EnvelopeStatusChangeId := taskInput.GetInt("EnvelopeStatusChangeId")
	if EnvelopeStatusChangeId == 0 {
		return errors.New("EnvelopeStatusChangeId is 0.")
	}
	var entity *EnvelopeStatusChangeEntity
	var err error
	entity, err = c.EnvelopeStatusChangeUsecase.GetByCond(Eq{"id": EnvelopeStatusChangeId})
	if err != nil {
		return err
	}
	if entity == nil {
		return errors.New("EnvelopeStatusChangeEntity is nil.")
	}
	clientEnvelope, err := c.ClientEnvelopeUsecase.GetByCond(And(
		Eq{"envelope_id": entity.EnvelopeId},
		Eq{"esign_vendor": EsignVendor_docusign}))
	if err != nil {
		return err
	}
	if clientEnvelope == nil {
		return errors.New("clientEnvelope is nil.")
	}

	boxFolderId, err := c.BoxcontractUsecase.ContractFolderId(clientEnvelope.ClientId)
	if err != nil {
		return err
	}

	err = c.TaskCreateUsecase.CreateTask(EnvelopeStatusChangeId,
		map[string]interface{}{"envelope_id": entity.EnvelopeId, "folder_id": boxFolderId},
		Task_Dag_SaveSignedContractInBox, 0, "", "")

	return err
}

const EnableClientCaseParentFolder = false

// BoxCreateFolderForNewClient 在Box中，创建客户文件夹
func (c *Dag) BoxCreateFolderForNewClient(task *TaskEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}

	taskInput := lib.ToTypeMapByString(task.TaskInput)
	clientCaseId := taskInput.GetInt("ClientId")
	if clientCaseId == 0 {
		return errors.New("ClientId is 0.")
	}
	return c.BizCreateBoxFolder(clientCaseId)
}

func (c *Dag) BizCreateBoxFolder(clientCaseId int32) error {

	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return err
	}
	if tClientCase == nil {
		return errors.New("tClientCase is nil.")
	}

	_, tContactFields, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return err
	}
	if tContactFields == nil {
		return errors.New("tContactFields is nil.")
	}

	folderName := ClientFolderNameForBox(tContactFields.TextValueByNameBasic("first_name"),
		tContactFields.TextValueByNameBasic("last_name"))

	useVBCActiveFolder, parentFolderId := c.BoxbuzUsecase.GetClientFolderRootId(*tClientCase)
	if EnableClientCaseParentFolder {
		parentKey := fmt.Sprintf("%s%d", Map_ClientBoxFolderIdParentId, clientCaseId)
		pFolderId, err := c.MapUsecase.GetForString(parentKey)
		if err != nil {
			return err
		}
		if pFolderId == "" {
			parentFolderName := fmt.Sprintf("%s #%d", folderName, clientCaseId)
			pFolderId, err = c.BoxUsecase.CreateFolder(parentFolderName, parentFolderId)
			if err != nil {
				return err
			}
			err = c.MapUsecase.Set(parentKey, pFolderId)
		}
		parentFolderId = pFolderId
	}
	newFolderName := fmt.Sprintf("%s #%d", folderName, clientCaseId)
	boxFolderId, _, err := c.BoxUsecase.CopyFolder(c.conf.Box.NewClientFolderStructureId, newFolderName, parentFolderId)
	if err != nil {
		c.log.Error(err, " newFolderName: ", newFolderName)
		return err
	}
	if useVBCActiveFolder {
		er := c.BoxCollaborationBuzUsecase.HandleUseVBCActiveCases(clientCaseId)
		if er != nil {
			c.log.Error(er, " HandleUseVBCActiveCases clientCaseId: ", clientCaseId)
		}
	}
	if len(boxFolderId) == 0 {
		return errors.New("boxFolderId is empty.")
	}
	key := fmt.Sprintf("%s%d", Map_ClientBoxFolderId, clientCaseId)
	err = c.MapUsecase.Set(key, boxFolderId)
	if err != nil {
		c.log.Error(err)
	}

	// 创建共享
	email := tContactFields.TextValueByNameBasic("email")
	if email == "" {
		return errors.New("Email does not exists.")
	}
	_, err = c.BoxUsecase.Collaborations(boxFolderId, email)

	if err != nil {
		return err
	}

	if configs.StoppedZoho {

		caseFileFolderValue := tClientCase.CustomFields.TextValueByNameBasic(FieldName_case_files_folder)
		if caseFileFolderValue == "" {
			row := make(lib.TypeMap)
			key = fmt.Sprintf("%s%d", Map_ClientBoxFolderId, clientCaseId)
			boxFolderId, _ = c.MapUsecase.GetForString(key)
			if boxFolderId != "" {
				row.Set(FieldName_case_files_folder, "https://veteranbenefitscenter.app.box.com/folder/"+boxFolderId)
			}
			if len(row) > 0 {
				row.Set(DataEntry_gid, tClientCase.Gid())
				_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry(row), DataEntry_gid, nil)
			}
		}
	} else {
		// 更新zoho box文件夹
		deal, err := c.ZohoUsecase.GetDeal(tClientCase.CustomFields.TextValueByNameBasic("gid"))
		if err != nil {
			return err
		}
		if deal == nil {
			return errors.New("zoho deal is nil.")
		}
		if deal.GetString("Case_Files_Folder") == "" {
			row := make(lib.TypeMap)
			key = fmt.Sprintf("%s%d", Map_ClientBoxFolderId, clientCaseId)
			boxFolderId, _ = c.MapUsecase.GetForString(key)
			if boxFolderId != "" {
				row.Set("Case_Files_Folder", "https://veteranbenefitscenter.app.box.com/folder/"+boxFolderId)
			}
			if len(row) > 0 {
				row.Set("id", tClientCase.CustomFields.TextValueByNameBasic("gid"))
				c.ZohoUsecase.PutRecordV1(config_zoho.Deals, row)
			}
			c.log.Info("BizCreateBoxFolder Case_Files_Folder: row: ", row)
		}
	}
	c.BoxCollaborationBuzUsecase.DoAddPermissionForBox(clientCaseId)
	return err
}

func (c *Dag) GetEnvelopeDocuments(task *TaskEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	taskInput := lib.ToTypeMapByString(task.TaskInput)
	envelopeId := InterfaceToString(taskInput.Get("envelope_id"))
	if len(envelopeId) == 0 {
		return errors.New("envelopeId不存在")
	}

	cred, err := c.DocuSignUsecase.DocuSignCredential()
	if err != nil {
		return err
	}
	srv := envelopes.New(cred)
	// 641eafbd-1d03-409b-b092-37219af0ae41
	//
	a, err := srv.DocumentsList(envelopeId).Do(context.Background())
	if err != nil {
		return err
	}

	key := Map_EnvelopeDocuments + envelopeId
	return c.MapUsecase.Set(key, InterfaceToString(a))
}

func (c *Dag) SaveSignedContractInBoxAdobe(task *TaskEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	taskInput := lib.ToTypeMapByString(task.TaskInput)
	agreementId := InterfaceToString(taskInput.Get("agreement_id"))
	if len(agreementId) == 0 {
		return errors.New("agreementId不存在")
	}
	folderId := taskInput.GetString("folder_id")
	if len(folderId) == 0 {
		return errors.New("folderId不存在")
	}
	adobeClient, err := c.AdobeSignUsecase.Client()
	if err != nil {
		return err
	}
	contractBytes, err := adobeClient.AgreementService.GetCombinedDocument(context.Background(), agreementId)
	if err != nil {
		return err
	}

	_, errContract := c.BoxUsecase.UploadFile(folderId, strings.NewReader(string(contractBytes)), "Contract.pdf")
	//errCertificate := c.BoxUsecase.UploadContract(folderId, certificateDownload, "Certificate.pdf")
	if errContract != nil {
		return errContract
	}
	//if errCertificate != nil {
	//	return errCertificate
	//}
	return nil
}

func (c *Dag) SaveSignedContractInBox(task *TaskEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	taskInput := lib.ToTypeMapByString(task.TaskInput)
	ty := taskInput.GetString("type")
	if ty == Sign_type_adobe {
		return c.SaveSignedContractInBoxAdobe(task)
	}

	envelopeId := InterfaceToString(taskInput.Get("envelope_id"))
	if len(envelopeId) == 0 {
		return errors.New("envelopeId不存在")
	}
	folderId := InterfaceToString(taskInput.Get("folder_id"))
	if len(folderId) == 0 {
		return errors.New("folderId不存在")
	}

	cred, err := c.DocuSignUsecase.DocuSignCredential()
	if err != nil {
		return err
	}
	srv := envelopes.New(cred)
	// documentId: 1
	// documentId: certificate
	certificateDownload, err := srv.DocumentsGet("certificate", envelopeId).Do(context.Background())
	if err != nil {
		return err
	}

	contractDownload, err := srv.DocumentsGet("1", envelopeId).Do(context.Background())
	if err != nil {
		return err
	}
	_, errContract := c.BoxUsecase.UploadFile(folderId, contractDownload, "Contract.pdf")
	_, errCertificate := c.BoxUsecase.UploadFile(folderId, certificateDownload, "Certificate.pdf")
	if errContract != nil {
		return errContract
	}
	if errCertificate != nil {
		return errCertificate
	}
	return nil
}

func (c *Dag) ReminderMedicalTeamFormsContractSent(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}

	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, task.IncrId)
	if err != nil {
		return err
	}
	if tClientCase == nil {
		return errors.New("tClientCase is nil.")
	}

	var tUser *TData
	if configs.UseOwnerSendingSMS {
		tUser, err = c.UserUsecase.GetByGid(tClientCase.CustomFields.TextValueByNameBasic(FieldName_user_gid))
		if err != nil {
			c.log.Error(err)
			return err
		}
	} else {
		primaryVSFullName := ""

		if tClientCase.CustomFields.TextValueByNameBasic("email") == "lialing@foxmail.com" ||
			tClientCase.CustomFields.TextValueByNameBasic("email") == "liaogling@gmail.com" {
			primaryVSFullName = "Engineering Team"
		} else {
			primaryVSFullName = tClientCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs)
		}

		if primaryVSFullName == "" {
			return errors.New("primaryVSFullName is empty")
		}
		tUser, err = c.UserUsecase.GetByFullName(primaryVSFullName)
		if err != nil {
			return err
		}
	}

	contractSentAt := c.BehaviorUsecase.MedicalTeamFormsContractSentAt(task.IncrId)
	vo, err := c.RemindUsecase.FollowingUpSignMedicalTeamFormsEmailBody(tClientCase, tUser, contractSentAt)
	if err != nil {
		return err
	}
	mailServiceConfig := InitMailServiceConfig()
	mailMessage := &MailMessage{
		To:      vo.Email,
		Subject: vo.Subject,
		Body:    vo.Body,
	}
	if !configs.IsDev() {
		mailMessage.Cc = []string{"info@vetbenefitscenter.com"}
	}
	err = c.MailUsecase.SendEmail(mailServiceConfig, mailMessage, "", nil)

	c.CommonUsecase.DB().Save(&EmailLogEntity{
		ClientId:   task.IncrId,
		Email:      vo.Email,
		TaskId:     task.ID,
		Tpl:        "FollowingUpSignMedicalTeamFormsEmail",
		SubId:      0,
		SenderMail: mailServiceConfig.Username,
		SenderName: mailServiceConfig.Name,
		Subject:    mailMessage.Subject,
		Body:       mailMessage.Body,
	})
	return err
}

func (c *Dag) ContractNonResponsive(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}

	return c.HandleNonResponsive(task.IncrId, false)
}

func (c *Dag) AmContractNonResponsive(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}

	return c.HandleNonResponsive(task.IncrId, true)
}

func (c *Dag) HandleNonResponsive(caseId int32, isAmContract bool) error {

	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tClientCase == nil { // 可能是测试的case，已经删除了
		return nil
	}
	deleted, err := IsDeletedCase(tClientCase)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if deleted {
		return nil
	}

	stages := tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages == config_vbc.Stages_FeeScheduleandContract || stages == config_vbc.Stages_AmContractPending { // 说明客户合同没有签属：设置为：non responsive

		var envelopeEntity *ClientEnvelopeEntity
		if isAmContract {
			envelopeEntity, err = c.ClientEnvelopeUsecase.GetBoxSignByCaseId(caseId, Type_AmContract)
		} else {
			envelopeEntity, err = c.ClientEnvelopeUsecase.GetBoxSignByCaseId(caseId, Type_FeeContract)
		}
		if err != nil {
			c.log.Error(err)
			return err
		}
		if envelopeEntity == nil {
			return errors.New("envelopeEntity is nil")
		}

		cancelResult, err := c.BoxUsecase.CancelSignRequest(envelopeEntity.EnvelopeId)
		if err != nil {
			cancelResultString := InterfaceToString(cancelResult)
			c.log.Error(err, " : ", cancelResultString, "envelopeEntity.ID: ", envelopeEntity.ID)
			return err
		}

		data := make(TypeDataEntry)
		data[DataEntry_gid] = tClientCase.Gid()
		data[FieldName_stages] = config_vbc.Stages_Terminated
		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
		if err != nil {
			c.log.Error(err)
			return err
		}
		behaviorType := BehaviorType_contract_non_responsive
		if isAmContract {
			behaviorType = BehaviorType_am_contract_non_responsive
		}

		return c.BehaviorUsecase.Add(caseId, behaviorType, time.Now(), "")
	} else {

	}
	return nil
}

func (c *Dag) HandleReminder(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}

	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, task.IncrId)
	if err != nil {
		return err
	}
	if tClientCase == nil { // 可能是测试的case，已经删除了
		return nil
	}
	deleted, err := IsDeletedCase(tClientCase)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if deleted {
		return nil
	}

	stages := tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages)

	timeLocation := GetCaseTimeLocation(tClientCase, c.log)
	taskInput := lib.ToTypeMapByString(task.TaskInput)
	reminderType := config_vbc.ReminderType(taskInput.GetString("ReminderType"))

	var mailGenre string
	if config_vbc.IsGroupReminderAmIntakeForm(reminderType) {
		if stages != config_vbc.Stages_AmInformationIntake {
			return nil
		}
		mailGenre = MailGenre_AmIntakeFormReminder
	}
	reminderConfig, err := config_vbc.GetReminderConfigs().GetVo(reminderType)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if mailGenre == "" {
		return errors.New("mailGenre is empty")
	}
	err = c.TaskCreateUsecase.CreateTaskMail(task.IncrId,
		mailGenre,
		reminderConfig.EmailTplSubId,
		nil, 0, "", "")

	if err != nil {
		return err
	}

	if reminderType == config_vbc.AmIntakeFormReminderFirst {
		tReminderConfig, err := config_vbc.GetReminderConfigs().GetVo(config_vbc.AmIntakeFormReminderSecond)
		if err != nil {
			return err
		}
		nextAtTime, err := tReminderConfig.ReminderTime(time.Now(), timeLocation)
		if err != nil {
			return err
		}
		err = c.TaskCreateUsecase.CreateGroupReminderAmIntakeFormTask(task.IncrId, config_vbc.AmIntakeFormReminderSecond, nextAtTime.Unix())

	} else if reminderType == config_vbc.AmIntakeFormReminderSecond {
		tReminderConfig, err := config_vbc.GetReminderConfigs().GetVo(config_vbc.AmIntakeFormReminderThird)
		if err != nil {
			return err
		}
		nextAtTime, err := tReminderConfig.ReminderTime(time.Now(), timeLocation)
		if err != nil {
			return err
		}
		err = c.TaskCreateUsecase.CreateGroupReminderAmIntakeFormTask(task.IncrId, config_vbc.AmIntakeFormReminderThird, nextAtTime.Unix())
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Dag) HandleContractReminder(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}

	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, task.IncrId)
	if err != nil {
		return err
	}
	if tClientCase == nil { // 可能是测试的case，已经删除了
		return nil
	}
	deleted, err := IsDeletedCase(tClientCase)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if deleted {
		return nil
	}

	stages := tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages != config_vbc.Stages_FeeScheduleandContract && stages != config_vbc.Stages_AmContractPending { // 说明客户已经签属合同了
		return nil
	}

	isAmContract := false
	if IsAmContract(*tClientCase) {
		isAmContract = true
	}

	timeLocation := GetCaseTimeLocation(tClientCase, c.log)

	taskInput := lib.ToTypeMapByString(task.TaskInput)
	contractReminderType := config_vbc.ContractReminderType(taskInput.GetString("ContractReminderType"))

	contractReminderConfig, err := config_vbc.GetContractReminderConfigs().GetVo(contractReminderType)
	if err != nil {
		c.log.Error(err)
		return err
	}

	mailGenre := MailGenre_ContractReminder
	if isAmContract {
		mailGenre = MailGenre_AmContractReminder
	}

	err = c.TaskCreateUsecase.CreateTaskMail(task.IncrId,
		mailGenre,
		contractReminderConfig.EmailTplSubId,
		nil, 0, "", "")

	if err != nil {
		return err
	}

	if contractReminderType == config_vbc.ContractReminderFirst {

		tContractReminderConfig, err := config_vbc.GetContractReminderConfigs().GetVo(config_vbc.ContractReminderSecond)
		if err != nil {
			return err
		}
		nextAtTime, err := tContractReminderConfig.ContractReminderTime(time.Now(), timeLocation)
		if err != nil {
			return err
		}
		err = c.TaskCreateUsecase.CreateTask(task.IncrId, map[string]interface{}{
			"ContractReminderType": config_vbc.ContractReminderSecond,
		}, Task_Dag_HandleContractReminder, nextAtTime.Unix(), "", "")
	} else if contractReminderType == config_vbc.ContractReminderSecond {
		tContractReminderConfig, err := config_vbc.GetContractReminderConfigs().GetVo(config_vbc.ContractReminderThird)
		if err != nil {
			return err
		}
		nextAtTime, err := tContractReminderConfig.ContractReminderTime(time.Now(), timeLocation)
		if err != nil {
			return err
		}
		err = c.TaskCreateUsecase.CreateTask(task.IncrId, map[string]interface{}{
			"ContractReminderType": config_vbc.ContractReminderThird,
		}, Task_Dag_HandleContractReminder, nextAtTime.Unix(), "", "")
	} else if contractReminderType == config_vbc.ContractReminderThird {
		tContractReminderConfig, err := config_vbc.GetContractReminderConfigs().GetVo(config_vbc.ContractReminderFourth)
		if err != nil {
			return err
		}
		nextAtTime, err := tContractReminderConfig.ContractReminderTime(time.Now(), timeLocation)
		if err != nil {
			return err
		}
		err = c.TaskCreateUsecase.CreateTask(task.IncrId, map[string]interface{}{
			"ContractReminderType": config_vbc.ContractReminderFourth,
		}, Task_Dag_HandleContractReminder, nextAtTime.Unix(), "", "")
	} else if contractReminderType == config_vbc.ContractReminderFourth {
		nextAtTime, err := utils.CalIntervalDayTime(time.Now(), 7, "08:00", timeLocation)
		if err != nil {
			return err
		}
		if isAmContract {
			err = c.TaskCreateUsecase.CreateTask(task.IncrId, nil,
				Task_Dag_AmContractNonResponsive, nextAtTime.Unix(), "", "")
		} else {
			err = c.TaskCreateUsecase.CreateTask(task.IncrId, nil,
				Task_Dag_ContractNonResponsive, nextAtTime.Unix(), "", "")
		}
	}

	if err != nil {
		return err
	}
	return nil
}

type CronTriggerVo struct {
	HandleSendSMSType HandleSendSMSType
	Params            lib.TypeMap
}

func (c *Dag) CronTrigger(task *TaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	cronTriggerVo := lib.StringToTDef[*CronTriggerVo](task.TaskInput, nil)
	if cronTriggerVo == nil {
		return errors.New("cronTriggerVo is nil")
	}

	err := c.CronTriggerUsecase.Handle(cronTriggerVo.HandleSendSMSType, task.IncrId, *cronTriggerVo)
	if err != nil {
		c.log.Error(err, " task.IncrId: ", task.IncrId)
		return err
	}
	return nil
}
