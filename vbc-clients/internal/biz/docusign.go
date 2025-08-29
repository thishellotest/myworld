package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	. "vbc/lib/builder"
	"vbc/lib/esign"
	"vbc/lib/esign/v2.1/envelopes"
	"vbc/lib/esign/v2.1/model"
)

type DocuSignCredential struct {
	oauth2Token *Oauth2TokenEntity
	accountId   string // 60e36fd1-0481-40c3-b7ca-c1ca4776bd87
}

func (c *DocuSignCredential) AuthDo(ctx context.Context, op *esign.Op) (*http.Response, error) {

	if op.Version == nil {
		return nil, errors.New("no api version set for op")
	}
	req, err := op.CreateRequest()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.oauth2Token.TokenType+" "+c.oauth2Token.AccessToken)
	var rawUrl string

	/*
		demoHost = "demo.docusign.net"
		baseHost = "www.docusign.net"
	*/
	//fmt.Println(op.Path, op.QueryOpts)
	if configs.IsDev() {
		rawUrl = fmt.Sprintf("https://demo.docusign.net/restapi/v2.1/accounts/%s/%s", c.accountId, op.Path)
	} else {
		// use demo todo:
		rawUrl = fmt.Sprintf("https://demo.docusign.net/restapi/v2.1/accounts/%s/%s", c.accountId, op.Path)
	}
	//fmt.Println("rawUrl:", rawUrl)

	//a, err := lib.Request("GET", rawUrl, nil, map[string]string{
	//	"Authorization": c.oauth2Token.TokenType + " " + c.oauth2Token.AccessToken,
	//})
	//fmt.Println(c.oauth2Token.TokenType, c.oauth2Token.AccessToken)
	//fmt.Println(*a, err)

	// /restapi/v2.1/accounts/{accountId}/envelopes/{envelopeId}/documents

	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	//fmt.Println(u.Host, u.Path, u.RawPath)
	req.URL.Scheme = u.Scheme
	req.URL.Host = u.Host
	req.URL.Path = u.Path
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("resp is nil.")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		msg := fmt.Sprintf("response status code %v", resp.StatusCode)
		bs, _ := io.ReadAll(resp.Body)
		res := string(bs)
		defer resp.Body.Close()
		if res == "" {
			return nil, errors.New(msg)
		} else {
			return nil, errors.New(msg + ":" + res)
		}
	}
	return resp, nil
}

func NewDocuSignCredential(oauth2Token *Oauth2TokenEntity, accountId string) *DocuSignCredential {
	return &DocuSignCredential{
		oauth2Token: oauth2Token,
		accountId:   accountId,
	}
}

type DocuSignUsecase struct {
	log                         *log.Helper
	conf                        *conf.Data
	Oauth2TokenUsecase          *Oauth2TokenUsecase
	Oauth2ClientUsecase         *Oauth2ClientUsecase
	EnvelopeStatusChangeUsecase *EnvelopeStatusChangeUsecase
	CommonUsecase               *CommonUsecase
	MapUsecase                  *MapUsecase
	PricingVersionUsecase       *PricingVersionUsecase
}

func NewDocuSignUsecase(logger log.Logger, conf *conf.Data, Oauth2TokenUsecase *Oauth2TokenUsecase, Oauth2ClientUsecase *Oauth2ClientUsecase,
	EnvelopeStatusChangeUsecase *EnvelopeStatusChangeUsecase,
	CommonUsecase *CommonUsecase,
	MapUsecase *MapUsecase,
	PricingVersionUsecase *PricingVersionUsecase) *DocuSignUsecase {
	return &DocuSignUsecase{
		log:                         log.NewHelper(logger),
		conf:                        conf,
		Oauth2TokenUsecase:          Oauth2TokenUsecase,
		Oauth2ClientUsecase:         Oauth2ClientUsecase,
		EnvelopeStatusChangeUsecase: EnvelopeStatusChangeUsecase,
		CommonUsecase:               CommonUsecase,
		MapUsecase:                  MapUsecase,
		PricingVersionUsecase:       PricingVersionUsecase,
	}
}

func (c *DocuSignUsecase) RunDosuSignEnvelopeChangeStatusJob(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := c.HandleEnvelopeChangeStatus()
				if err != nil {
					c.log.Error(err)
				}
				time.Sleep(120 * time.Second)
			}
		}
	}()
	return nil
}

func (c *DocuSignUsecase) DocuSignCredential() (*DocuSignCredential, error) {
	token, err := c.Oauth2TokenUsecase.GetByAppId(Oauth2_AppId_docusign)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("token is nil.")
	}
	return NewDocuSignCredential(token, c.conf.Docusign.AppAccountId), nil
}

func (c *DocuSignUsecase) HandleEnvelopeChangeStatus() error {
	cred, err := c.DocuSignCredential()
	if err != nil {
		return err
	}
	srv := envelopes.New(cred)
	now := time.Now()
	now = now.AddDate(0, 0, -20)
	data, err := srv.ListStatusChanges().FromDate(now).Do(context.Background())
	if err != nil {
		return err
	}
	if data != nil {
		for _, v := range data.Envelopes {
			// builder.Eq{"ssss": v.EnvelopeID}
			entity, err := c.EnvelopeStatusChangeUsecase.
				GetByCond(And(Eq{"envelope_id": v.EnvelopeID},
					Eq{"status": v.Status},
					Eq{"status_changed_datetime": v.StatusChangedDateTime.Format(time.RFC3339)}))
			if err != nil {
				c.log.Error(err)
			} else {
				if entity == nil {
					var SenderAccountId, SenderEmail, SenderUserId, SenderUsername string
					if v.Sender != nil {
						SenderAccountId = v.Sender.AccountID
						SenderEmail = v.Sender.Email
						SenderUserId = v.Sender.UserID
						SenderUsername = v.Sender.UserName
					}
					entity = &EnvelopeStatusChangeEntity{
						EnvelopeId:                  v.EnvelopeID,
						CreatedDatetime:             v.CreatedDateTime.Format(time.RFC3339),
						AttachmentsUri:              v.AttachmentsURI,
						CertificateUri:              v.CertificateURI,
						CustomFieldsUri:             v.CustomFieldsURI,
						DocumentsCombinedUri:        v.DocumentsCombinedURI,
						EnvelopeLocation:            v.EnvelopeLocation,
						DocumentsUri:                v.DocumentsURI,
						IsSignatureProviderEnvelope: string(v.IsSignatureProviderEnvelope),
						LastModifiedDatetime:        v.LastModifiedDateTime.Format(time.RFC3339),
						NotificationUri:             v.NotificationURI,
						PurgeState:                  v.PurgeState,
						RecipientsUri:               v.RecipientsURI,
						SenderAccountId:             SenderAccountId,
						SenderEmail:                 SenderEmail,
						SenderUserId:                SenderUserId,
						SenderUsername:              SenderUsername,
						SentDatetime:                v.SentDateTime.Format(time.RFC3339),
						SigningLocation:             v.SigningLocation,
						Status:                      v.Status,
						StatusChangedDatetime:       v.StatusChangedDateTime.Format(time.RFC3339),
						TemplatesUri:                v.TemplatesURI,
						EmailBlurb:                  v.EmailBlurb,
						EmailSubject:                v.EmailSubject,
						CreatedAt:                   time.Now().Unix(),
						UpdatedAt:                   time.Now().Unix(),
					}
				} else {
					entity.UpdatedAt = time.Now().Unix()
				}
				err = c.CommonUsecase.DB().Save(&entity).Error
				if err != nil {
					c.log.Error(err)
				}
			}
		}
	}

	return nil
}

func (c *DocuSignUsecase) ContractTemplateKey(tClientCase *TData) (string, error) {
	if tClientCase == nil {
		return "", errors.New("tClientCase is nil.")
	}
	var index int
	if tClientCase.CustomFields.TextValueByNameBasic("active_duty") == "Yes" {
		index = -1
	} else {
		index = config_vbc.SignContractIndexV1(tClientCase.CustomFields.NumberValueByNameBasic("effective_current_rating"))
	}

	if c.conf.UseBoxSign {
		return fmt.Sprintf("%s%d", Map_boxsignTpl, index), nil
	} else if c.conf.UseAdobeSign {
		return fmt.Sprintf("%s%d", Map_adobesignTpl, index), nil
	}
	return fmt.Sprintf("%s%d", Map_docusignTpl, index), nil
}

func (c *DocuSignUsecase) ContractTemplateId(tClientCase *TData) (contractIndex int, templateId string, pricingVersion string, err error) {

	//if lib.EnabledDBPricingVersion {
	return c.ContractTemplateIdByDB(tClientCase)
	//}
	//templateId, err = c.ContractTemplateIdOld(tClientCase)
	//return templateId, DefaultPricingVersion, err
}

func (c *DocuSignUsecase) ContractTemplateIdByDB(tClientCase *TData) (contractIndex int, templateId string, pricingVersion string, err error) {

	if tClientCase == nil {
		return 0, "", "", errors.New("tClientCase is nil")
	}

	var index int
	if tClientCase.CustomFields.TextValueByNameBasic("active_duty") == "Yes" {
		index = -1
	} else {
		index = config_vbc.SignContractIndexV1(tClientCase.CustomFields.NumberValueByNameBasic("effective_current_rating"))
	}
	contractIndex = index
	c.log.Debug("index:", index)
	var config *PricingVersionConfig
	var versionEntity *PricingVersionEntity
	// 这是测试数据
	if tClientCase.CustomFields.TextValueByNameBasic("client_gid") == "6159272000003710001_close" {
		config, versionEntity, err = c.PricingVersionUsecase.ConfigByPricingVersion("v20240620")
	} else {
		config, versionEntity, err = c.PricingVersionUsecase.CurrentVersionConfig()
	}
	if err != nil {
		c.log.Error(err)
		return index, "", "", err
	}
	if config == nil {
		return index, "", "", errors.New("config is nil")
	}
	if versionEntity == nil {
		return index, "", "", errors.New("versionEntity is nil")
	}
	templateId, err = config.GetBoxSignTpl(InterfaceToString(index))
	if err != nil {
		c.log.Error(err)
		return index, "", "", err
	}
	pricingVersion = versionEntity.Version
	return
}

func (c *DocuSignUsecase) ContractTemplateIdOld(tClientCase *TData) (string, error) {

	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}

	key, err := c.ContractTemplateKey(tClientCase)
	if err != nil {
		return "", err
	}
	return c.MapUsecase.GetForString(key)
}

func (c *DocuSignUsecase) CreateEnvelopeAndSent(templateId string,
	clientName string,
	clientEmail string,
	agentName string,
	agentEmail string) (*model.EnvelopeSummary, error) {

	token, err := c.Oauth2TokenUsecase.GetByAppId(Oauth2_AppId_docusign)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("token is nil.")
	}
	cred := NewDocuSignCredential(token, c.conf.Docusign.AppAccountId)
	srv := envelopes.New(cred)

	templateRole := model.TemplateRole{}
	templateRole.Name = clientName
	templateRole.RoleName = "Client"
	templateRole.Email = clientEmail

	templateRole1 := model.TemplateRole{}
	templateRole1.Name = agentName
	templateRole1.RoleName = "Agent"
	templateRole1.Email = agentEmail

	//EnvelopeDefinition
	envlopeDef := model.EnvelopeDefinition{}
	envlopeDef.TemplateRoles = append(envlopeDef.TemplateRoles, templateRole, templateRole1)
	envlopeDef.Status = "sent"
	envlopeDef.TemplateID = templateId

	return srv.Create(&envlopeDef).Do(context.TODO())

}
