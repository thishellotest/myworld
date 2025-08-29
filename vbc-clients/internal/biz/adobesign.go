package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib/adobesign"
)

type AdobeSignUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	Oauth2TokenUsecase *Oauth2TokenUsecase
}

func NewAdobeSignUsecase(logger log.Logger, CommonUsecase *CommonUsecase, conf *conf.Data, Oauth2TokenUsecase *Oauth2TokenUsecase) *AdobeSignUsecase {
	return &AdobeSignUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		Oauth2TokenUsecase: Oauth2TokenUsecase,
	}
}

func (c *AdobeSignUsecase) Client() (*adobesign.Client, error) {
	token, err := c.Oauth2TokenUsecase.GetByAppId(Oauth2_AppId_adobesign)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("token is nil")
	}
	client := adobesign.NewClient(token.AccessToken, "na3", "")
	return client, nil
}

type CreateAgreementMember struct {
	Email     string
	FirstName string
	LastName  string
}

// CreateAgreement libraryDocumentId:CBJCHBCAABAAxqRFmyaz9biyPYcmsoRYRwOxIAW8nifY
func (c *AdobeSignUsecase) CreateAgreement(agreementName string, libraryDocumentId string, client CreateAgreementMember, vs CreateAgreementMember) (agreementId string, err error) {

	adobesignClient, err := c.Client()
	if err != nil {
		return "", err
	}

	memberInfo := adobesign.MemberInfo{
		Email: client.Email,
	}
	memberInfo.SecurityOption.NameInfo.FirstName = client.FirstName
	memberInfo.SecurityOption.NameInfo.LastName = client.LastName

	memberInfo1 := adobesign.MemberInfo{
		Email: vs.Email,
	}
	memberInfo1.SecurityOption.NameInfo.FirstName = vs.FirstName
	memberInfo1.SecurityOption.NameInfo.LastName = vs.LastName

	participantSetInfo := adobesign.ParticipantSetInfo{
		MemberInfos: []adobesign.MemberInfo{
			memberInfo,
		},
		Order: 1,                                // the order in which the signer appears on the document
		Role:  adobesign.ParticipantRole.Signer, // the role of the member
		//Name:  "Ling Liao",
		//Id:    "client_01",
	}

	participantSetInfo1 := adobesign.ParticipantSetInfo{
		MemberInfos: []adobesign.MemberInfo{
			memberInfo1,
		},
		Order: 2,                                // the order in which the signer appears on the document
		Role:  adobesign.ParticipantRole.Signer, // the role of the member
		//Name:  "Geng Ling",
		//Id:    "vs_01",
	}

	agreement, err := adobesignClient.AgreementService.CreateAgreement(context.Background(), adobesign.Agreement{
		FileInfos: []adobesign.FileInfo{
			{LibraryDocumentId: libraryDocumentId},
		},
		Name: agreementName,
		ParticipantSetsInfo: []adobesign.ParticipantSetInfo{
			participantSetInfo,
			participantSetInfo1,
		},
		SignatureType: adobesign.SignatureType.Esign,
		State:         adobesign.AgreementState.InProcess,
		MergeFieldInfo: []adobesign.MergeFieldInfo{
			{
				DefaultValue: GenFullName(client.FirstName, client.LastName),
				FieldName:    "ClientName",
			},
			{
				DefaultValue: GenFullName(client.FirstName, client.LastName),
				FieldName:    "ClientNameSign",
			},
			{
				DefaultValue: GenFullName(vs.FirstName, vs.LastName),
				FieldName:    "VSNameSign",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return agreement.Id, nil
}

func GenFullName(firstName string, lastName string) string {
	if lastName == "" {
		return firstName
	}
	return firstName + " " + lastName
}

func GenFullNameWithMiddleName(firstName string, middleName string, lastName string) string {
	if middleName != "" {
		firstName += " " + middleName
	}
	if lastName != "" {
		firstName += " " + lastName
	}
	return firstName
}

func (c *AdobeSignUsecase) GetAgreement(ctx context.Context, agreementId string) (*adobesign.Agreement, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}
	return client.AgreementService.GetAgreement(ctx, agreementId)
}
