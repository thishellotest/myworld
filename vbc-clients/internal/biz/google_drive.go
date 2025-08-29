package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"vbc/internal/conf"
)

/*
https://drive.google.com/drive/folders/1KZo9l1sOENV4s8DJbL2vIx8FBiBIznww

https://developers.google.com/drive/api/guides/search-files?hl=zh-cn#java

https://developers.google.com/drive/api/reference/rest/v3/files/listLabels?hl=zh-cn

*/

type GoogleDriveUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	Oauth2TokenUsecase *Oauth2TokenUsecase
}

func NewGoogleDriveUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	Oauth2TokenUsecase *Oauth2TokenUsecase) *GoogleDriveUsecase {
	uc := &GoogleDriveUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		Oauth2TokenUsecase: Oauth2TokenUsecase,
	}

	return uc
}

func (c *GoogleDriveUsecase) Service(ctx context.Context) (*drive.Service, error) {

	a, err := c.Oauth2TokenUsecase.GetByAppId(Oauth2_AppId_google)
	if err != nil {
		return nil, err
	}

	tokenSource, err := a.TokenSourceMod()
	if err != nil {
		return nil, err
	}
	srv, err := drive.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

// ListItemsInFolder id: 1xwolHdQbmw-_uVz3Xxi-AaBlUoJpxgGR
func (c *GoogleDriveUsecase) ListItemsInFolder(ctx context.Context, id string) (*drive.FileList, error) {
	srv, err := c.Service(ctx)
	if err != nil {
		return nil, err
	}
	a := srv.Files.List()
	return a.Q("'" + id + "' in parents").Do()
}

// CreateFolder id: 1xwolHdQbmw-_uVz3Xxi-AaBlUoJpxgGR
// 返回结果：{"id":"1rCY79_oPSDhB1N8LUKeJV7RnKmXPzdm4","kind":"drive#file","mimeType":"application/vnd.google-apps.folder","name":"test folder"}
// google drive的文件夹名称是可以重复的
// 当没有设置文件夹名称时，默认提供：New Folder
func (c *GoogleDriveUsecase) CreateFolder(ctx context.Context, parentFolderId string, folderName string) (*drive.File, error) {
	srv, err := c.Service(ctx)
	if err != nil {
		return nil, err
	}
	//parentFolderId := c.conf.GoogleDrive.PaymentsFolderId
	file := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentFolderId},
	}
	a := srv.Files.Create(file)
	return a.Do()
}

func (c *GoogleDriveUsecase) UploadFile(ctx context.Context, parentFolderId string, fileName string, fileReader io.Reader) (*drive.File, error) {
	srv, err := c.Service(ctx)
	if err != nil {
		return nil, err
	}
	//parentFolderId := "1OVaJtb0x-xO0iKrQnnIsw56u6Q3j32hf"
	file := &drive.File{
		Name:    fileName,
		Parents: []string{parentFolderId},
	}
	a := srv.Files.Create(file)
	a.Media(fileReader)
	return a.Do()
}
