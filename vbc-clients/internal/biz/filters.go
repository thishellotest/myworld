package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

const (
	Filter_FieldName_filter_name = "filter_name"
	Filter_FieldName_content     = "content"
	Filter_FieldName_kind        = "kind"
	Filter_FieldName_user_gid    = "user_gid"
	Filter_FieldName_table_type  = "table_type"
)

type FilterEntity struct {
	ID         int32 `gorm:"primaryKey"`
	Kind       string
	UserGid    string
	FilterName string
	Content    string
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  int64
}

func (FilterEntity) TableName() string {
	return "filters"
}

type FilterVo struct {
	FilterId   int32  `json:"filter_id"`
	FilterName string `json:"filter_name"`
	Content    string `json:"content"`
}

func (c *FilterEntity) ToFilterVo() FilterVo {

	return FilterVo{
		FilterId:   c.ID,
		FilterName: c.FilterName,
		Content:    c.Content,
	}
}

type FilterUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[FilterEntity]
}

func NewFilterUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *FilterUsecase {
	uc := &FilterUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
