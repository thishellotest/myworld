package biz

import (
	"strings"
	"vbc/internal/conf"
	. "vbc/lib/builder"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	AttorneyFieldName_first_name           = "first_name"
	AttorneyFieldName_last_name            = "last_name"
	AttorneyFieldName_accreditation_date   = "accreditation_date"
	AttorneyFieldName_accreditation_number = "accreditation_number"
	AttorneyFieldName_email                = "email"
	AttorneyFieldName_zip_code             = "zip_code"
	AttorneyFieldName_province             = "province"
	AttorneyFieldName_city                 = "city"
	AttorneyFieldName_street               = "street"
	AttorneyFieldName_company_name         = "company_name"
	AttorneyFieldName_ro_email             = "ro_email"
	AttorneyFieldName_status               = "status"
)

/*
CREATE TABLE `attorneys` (

	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`first_name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`last_name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`accreditation_date` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`accreditation_number` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`email` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`zip_code` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`province` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`city` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`street` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`created_at` int(11) NOT NULL DEFAULT '0',
	`updated_at` int(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='attorneys';
*/
type AttorneyEntity struct {
	ID                  int32 `gorm:"primaryKey"`
	Gid                 string
	FullName            string
	FirstName           string
	LastName            string
	AccreditationDate   string
	AccreditationNumber string
	Email               string
	ZipCode             string
	Province            string
	City                string
	Street              string
	CompanyName         string
	RoEmail             string
	ForwardEmails       string
	CreatedAt           int64
	UpdatedAt           int64
	DeletedAt           int64
}

func (c *AttorneyEntity) ToContractAttorneyVo() (vo ContractAttorneyVo) {

	vo.Street = c.Street
	vo.City = c.City
	vo.Province = c.Province
	vo.ZipCode = c.ZipCode
	vo.FirstName = c.FirstName
	vo.LastName = c.LastName
	vo.AccreditationDate = c.AccreditationDate
	vo.AccreditationNumber = c.AccreditationNumber
	vo.Email = c.Email
	return
}

func (AttorneyEntity) TableName() string {
	return "attorneys"
}

type AttorneyUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[AttorneyEntity]
}

func NewAttorneyUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *AttorneyUsecase {
	uc := &AttorneyUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *AttorneyUsecase) GetByGid(gid string) (*AttorneyEntity, error) {

	return c.GetByCond(Eq{"gid": gid})
}

func (c *AttorneyUsecase) GetByName(name string) (*AttorneyEntity, error) {
	parts := strings.Fields(name)
	if len(parts) < 2 {
		return nil, nil
	}
	first := strings.Join(parts[:len(parts)-1], " ")
	last := parts[len(parts)-1]
	return c.GetByCond(And(
		Expr("LOWER(first_name) = LOWER(?)", first),
		Expr("LOWER(last_name) = LOWER(?)", last),
	))
}
