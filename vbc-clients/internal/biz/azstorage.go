package biz

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	"vbc/lib/to"
)

const (
	Azstorage_ContainerName_VBC = "vbc"
)

type AzstorageUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func NewAzstorageUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *AzstorageUsecase {
	uc := &AzstorageUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

func (c *AzstorageUsecase) GetClient(accountName string, accountKey string) (*azblob.Client, error) {
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	return azblob.NewClientWithSharedKeyCredential(url, cred, nil)
}

func (c *AzstorageUsecase) Client() (*azblob.Client, error) {
	return c.GetClient(c.GetAccountInfo())
}

func (c *AzstorageUsecase) GetAccountInfo() (accountName, accountKey string) {
	return configs.EnvAzureStorageAccountName(), configs.EnvAzureStorageAccountKey()
}

func (c *AzstorageUsecase) Credential() *azblob.SharedKeyCredential {
	credential, _ := azblob.NewSharedKeyCredential(c.GetAccountInfo())
	return credential
}

func (c *AzstorageUsecase) SasReadUrl(fileBlobname string, ExpiryTime *time.Time) (url string, err error) {
	accountName, accountKey := c.GetAccountInfo()
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", err
	}

	blobSignValue := sas.BlobSignatureValues{
		Protocol:   sas.ProtocolHTTPS,
		StartTime:  time.Now().UTC(),
		ExpiryTime: time.Now().UTC().Add(48 * time.Hour),
		//ExpiryTime: ExpiryTime,
		//Permissions:   to.Ptr(sas.BlobPermissions{Read: true, Create: true, Write: true, Tag: true}).String(),
		Permissions: to.Ptr(sas.BlobPermissions{Read: true}).String(),
		//ContentType:   "application/pdf",
		ContainerName: Azstorage_ContainerName_VBC,
	}
	if ExpiryTime != nil {
		blobSignValue.ExpiryTime = *ExpiryTime
	}
	_, suffix := lib.FileExt(fileBlobname, true)
	if suffix == "pdf" {
		blobSignValue.ContentType = "application/pdf"
	} else if suffix == "jpg" {
		blobSignValue.ContentType = "image/jpeg"
	} else if suffix == "webp" {
		blobSignValue.ContentType = "image/webp"
	} else if suffix == "json" {
		blobSignValue.ContentType = "application/json"
	}

	sasQueryParams, err := blobSignValue.SignWithSharedKey(credential)
	if err != nil {
		return "", err
	}

	//fileName := url.QueryEscape("STR Full_1.pdf")

	// https://storageblobseu2.blob.core.windows.net/vbc/tmp/STR%20Full_1.pdf
	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/vbc/%s?%s", accountName, fileBlobname, sasQueryParams.Encode())
	return sasURL, nil
}

func (c *AzstorageUsecase) UploadStream(ctx context.Context, blobName string, body io.Reader) (azblob.UploadStreamResponse, error) {

	client, err := c.Client()
	if err != nil {
		return azblob.UploadStreamResponse{}, err
	}

	return client.UploadStream(ctx, Azstorage_ContainerName_VBC, blobName, body, &azblob.UploadStreamOptions{})
}

func (c *AzstorageUsecase) DownloadStream(ctx context.Context, blobName string) (response azblob.DownloadStreamResponse, err error) {
	client, err := c.Client()
	if err != nil {
		return response, err
	}
	return client.DownloadStream(ctx, Azstorage_ContainerName_VBC, blobName, &azblob.DownloadStreamOptions{})
}

func (c *AzstorageUsecase) DeleteBlob(ctx context.Context, blobName string) error {

	client, err := c.Client()
	if err != nil {
		return err
	}
	_, err = client.DeleteBlob(ctx, Azstorage_ContainerName_VBC, blobName, nil)
	return err
}
