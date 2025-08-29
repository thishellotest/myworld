package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
)

const (
	BlobComment_Type_Default  = 0
	BlobComment_Type_OnlyText = 1
)

type BlobCommentEntity struct {
	ID          int32 `gorm:"primaryKey"`
	Gid         string
	BlobGid     string
	Content     string
	JsonData    string
	Page        int32
	Type        int
	UserGid     string
	HaReportGid string
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64
}

func (BlobCommentEntity) TableName() string {
	return "blob_comments"
}

func (c *BlobCommentEntity) BlobCommentToApi(userCaches lib.Cache[*TData], UserUsecase *UserUsecase, log *log.Helper) lib.TypeMap {

	typeMap := make(lib.TypeMap)
	typeMap.Set("id", c.ID)
	typeMap.Set("gid", c.Gid)
	typeMap.Set("blob_gid", c.BlobGid)
	typeMap.Set("content", c.Content)
	typeMap.Set("json_data", c.JsonData)
	typeMap.Set("page", c.Page)
	typeMap.Set("type", c.Type)
	typeMap.Set("created_at", c.CreatedAt)
	typeMap.Set("updated_at", c.UpdatedAt)
	if c.UserGid != "" {
		user, err := UserUsecase.GetUserWithCache(userCaches, c.UserGid)
		if err != nil {
			log.Error(err)
		}
		if user != nil {
			typeMap.Set(Fab_User, UserToRelaApi(user))
		}
	}
	return typeMap
}

type BlobCommentUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[BlobCommentEntity]
}

func NewBlobCommentUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *BlobCommentUsecase {
	uc := &BlobCommentUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}
