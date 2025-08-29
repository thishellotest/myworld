package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

const (
	Unsubscribes_Status_No  = 0 // 没有退订
	Unsubscribes_Status_Yes = 1 // 已经退订
)

type UnsubscribesEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	LatestFromId       string
	ContactPhoneNumber string
	Status             int
	BizDeletedAt       int64
	CreatedAt          int64
	UpdatedAt          int64
}

func (c *UnsubscribesEntity) SysStatus() string {
	if c.Status == Unsubscribes_Status_Yes {
		return "Opt-out SMS"
	}
	if c.Status == Unsubscribes_Status_No {
		return "Opt-in SMS"
	}
	return ""
}

func (UnsubscribesEntity) TableName() string {
	return "unsubscribes"
}

func (c *UnsubscribesEntity) ToApi(userFacde *UserFacade, TimezonesUsecase *TimezonesUsecase, clients map[string][]*TData, clientCases map[string][]*TData) lib.TypeMap {
	data := make(lib.TypeMap)
	data.Set("id", c.ID)
	phone, _, _, _ := FormatPhoneNumber(c.ContactPhoneNumber)
	data.Set("phone", phone)

	if client, ok := clients[c.ContactPhoneNumber]; ok {
		var rowClients lib.TypeList
		for _, v := range client {
			rowClients = append(rowClients, lib.TypeMap{
				DataEntry_gid:       v.Gid(),
				FieldName_full_name: v.CustomFields.TextValueByNameBasic(FieldName_full_name),
			})
		}
		data.Set("clients", rowClients)
	}
	if tCase, ok := clientCases[c.ContactPhoneNumber]; ok {
		var rowCases lib.TypeList
		for _, v := range tCase {
			rowCases = append(rowCases, lib.TypeMap{
				DataEntry_gid:       v.Gid(),
				FieldName_deal_name: v.CustomFields.TextValueByNameBasic(FieldName_deal_name),
			})
		}
		data.Set("client_cases", rowCases)
	}
	createdTime, _ := TimestampToStringByUserFacade(userFacde, TimezonesUsecase, c.CreatedAt)
	data.Set("created_time", createdTime)

	updatedTime, _ := TimestampToStringByUserFacade(userFacde, TimezonesUsecase, c.UpdatedAt)
	data.Set("updated_time", updatedTime)
	data.Set("status", c.Status)
	data.Set("sys__status", c.SysStatus())

	return data
}

type UnsubscribesUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[UnsubscribesEntity]
}

func NewUnsubscribesUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *UnsubscribesUsecase {
	uc := &UnsubscribesUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

/*

SELECT REPLACE(REPLACE(REPLACE(REPLACE(phone, '-', ''), '(', ''), ')', ''), ' ', '')  as aa
FROM client_cases
WHERE REPLACE(REPLACE(REPLACE(REPLACE(phone, '-', ''), '(', ''), ')', ''), ' ', '') LIKE '%415%';
*/

// CanSendSms 格式支持USA：+19044157090 (402) 215-6064， (402) 215-6064
func (c *UnsubscribesUsecase) CanSendSms(phone string) (bool, error) {
	var newPhone string
	var err error
	if len(phone) == 12 {
		if strings.Index(phone, "+1") == 0 {
			newPhone = phone
		} else {
			return false, errors.New(phone + " Phone format is wrong")
		}
	} else {
		newPhone, err = USAPhoneHandle(phone)
		if err != nil {
			return false, err
		}
		newPhone = "+1" + newPhone
	}
	a, err := c.GetByCond(Eq{"contact_phone_number": newPhone, "biz_deleted_at": 0})
	if err != nil {
		return false, err
	}
	if a != nil {
		if a.Status == Unsubscribes_Status_Yes {
			return false, nil
		}
	}
	return true, nil
}
