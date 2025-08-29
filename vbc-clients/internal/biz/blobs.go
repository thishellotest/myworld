package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

const (
	BlobType_pdf           = "pdf"
	Blob_Status_processing = "0" // 等待OCR完成
	Blob_Status_ready      = "1" // 已经好了
)

const (
	BlobFieldName_status        = "status"
	BlobFieldName_blob_type     = "blob_type"
	BlobFieldName_uniqblob      = "uniqblob"
	BlobFieldName_case_gid      = "case_gid"
	BlobFieldName_user_gid      = "user_gid"
	BlobFieldName_file_blobname = "file_blobname"
)

//
//type BlobEntity struct {
//	ID                 int32 `gorm:"primaryKey"`
//	Gid                string
//	HandleStatus       int
//	HandleResult       int
//	HandleResultDetail string
//	Status             int
//	Name               string
//	BlobType           string
//	Uniqblob           string
//	CaseId             int32
//	UserId             int32
//	CreatedAt          int64
//	UpdatedAt          int64
//	DeletedAt          int64
//}
//
//func (BlobEntity) TableName() string {
//	return "blobs"
//}
//
//func (c *BlobEntity) AppendHandleResultDetail(str string) {
//	if c.HandleResultDetail == "" {
//		c.HandleResultDetail = time.Now().Format(time.RFC3339) + " " + str
//	} else {
//		c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
//	}
//}

func GetUniqblobFileInfo(uniqblob string) (fileId string, versionId string) {
	aa := strings.Split(uniqblob, "_")
	if len(aa) == 2 {
		return aa[0], aa[1]
	}
	return "", ""
}

//
//func (c *BlobEntity) GetFileInfo() (fileId string, versionId string) {
//	aa := strings.Split(c.Uniqblob, "_")
//	if len(aa) == 2 {
//		return aa[0], aa[1]
//	}
//	return "", ""
//}

func BlobToApi(blob *TData, userCaches lib.Cache[*TData], UserUsecase *UserUsecase,
	caseCaches lib.Cache[*TData], ClientCaseUsecase *ClientCaseUsecase,
	AzstorageUsecase *AzstorageUsecase,
	log *log.Helper) lib.TypeMap {
	if blob == nil {
		return nil
	}

	apiMap := blob.CustomFields.ToApiMap()
	blobName := apiMap.GetString(BlobFieldName_file_blobname + ".text_value")
	if blobName != "" {
		blobUrl, err := AzstorageUsecase.SasReadUrl(blobName, nil)
		if err != nil {
			log.Error(err)
		} else {
			apiMap.Set(BlobFieldName_file_blobname+".text_value", blobUrl)
			apiMap.Set(BlobFieldName_file_blobname+".display_value", blobUrl)
		}
	}

	data := make(lib.TypeMap)
	data.Set("custom_fields", apiMap)

	user, err := UserUsecase.GetUserWithCache(userCaches, blob.CustomFields.TextValueByNameBasic("gid"))
	if err != nil {
		log.Error(err)
	}
	if user != nil {
		data.Set(Fab_User, UserToRelaApi(user))
	}
	tCase, err := ClientCaseUsecase.GetCaseWithCache(caseCaches, blob.CustomFields.NumberValueByNameBasic("case_id"))
	if err != nil {
		log.Error(err)
	}
	if tCase != nil {
		data.Set(Fab_Case, CaseToRelaApi(tCase, log))
	}
	return data
}

//
//func (c *BlobEntity) BlobToApi(userCaches lib.Cache[*UserEntity], UserUsecase *UserUsecase,
//	caseCaches lib.Cache[*TData], ClientCaseUsecase *ClientCaseUsecase,
//	log *log.Helper) lib.TypeMap {
//
//	data := make(lib.TypeMap)
//	data.Set("gid", c.Gid)
//	data.Set("created_at", c.CreatedAt)
//	data.Set("updated_at", c.UpdatedAt)
//	data.Set("name", c.Name)
//
//	user, err := UserUsecase.GetUserWithCache(userCaches, c.UserId)
//	if err != nil {
//		log.Error(err)
//	}
//	if user != nil {
//		data.Set(Fab_User, user.UserToApi())
//	}
//	tCase, err := ClientCaseUsecase.GetCaseWithCache(caseCaches, c.CaseId)
//	if err != nil {
//		log.Error(err)
//	}
//	if tCase != nil {
//		data.Set(Fab_Case, CaseToApi(tCase, log))
//	}
//	return data
//}

type BlobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	//DBUsecase[BlobEntity]
	TUsecase *TUsecase
}

func NewBlobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase) *BlobUsecase {
	uc := &BlobUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
	}

	//uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *BlobUsecase) GetByUniqblob(uniqblob string) (*TData, error) {
	return c.TUsecase.Data(Kind_blobs, Eq{"uniqblob": uniqblob, "deleted_at": 0})
}

func (c *BlobUsecase) GetByGid(gid string) (*TData, error) {
	return c.TUsecase.Data(Kind_blobs, Eq{"gid": gid, "deleted_at": 0})
}
