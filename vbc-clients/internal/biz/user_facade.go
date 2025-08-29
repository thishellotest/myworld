package biz

import (
	"github.com/pkg/errors"
	"sync"
	"time"
	"vbc/lib/builder"
)

var _ Facade = (*UserFacade)(nil)

type UserFacade struct {
	TData
	timezonesEntity        *TimezonesEntity
	getTimezonesEntityOnce sync.Once
}

func (c *UserFacade) New() Facade {
	return &UserFacade{}
}

func (c *UserFacade) GetTData() TData {
	return c.TData
}

func (c *UserFacade) SetTData(tData TData) {
	c.TData = tData
}

func (c *UserFacade) ProfileGid() string {
	return c.CustomFields.TextValueByNameBasic(User_FieldName_profile_gid)
}

func (c *UserFacade) TimezoneId() string {
	timezoneId := c.CustomFields.TextValueByNameBasic(User_FieldName_timezone_id)
	if timezoneId == "" {
		timezoneId = Default_Timezones_CodeValue
	}
	return timezoneId
}

func (c *UserFacade) ToFabUser() FabUser {
	return FabUser{
		Gid:      c.Gid(),
		FullName: c.CustomFields.TextValueByNameBasic(UserFieldName_fullname),
		Email:    c.CustomFields.TextValueByNameBasic(UserFieldName_email),
	}
}

func (c *UserFacade) GetTimezonesEntity(TimezonesUsecase *TimezonesUsecase) (TimezonesEntity *TimezonesEntity, err error) {

	if TimezonesUsecase == nil {
		return nil, errors.New("TimezonesUsecase is nil")
	}
	c.getTimezonesEntityOnce.Do(func() {
		timezoneId := c.TimezoneId()
		c.timezonesEntity, err = TimezonesUsecase.GetByCond(builder.Eq{"code_value": timezoneId})
	})
	return c.timezonesEntity, err
}

func (c *UserFacade) GetTimeLocation(TimezonesUsecase *TimezonesUsecase) (*time.Location, error) {
	entity, err := c.GetTimezonesEntity(TimezonesUsecase)
	if err != nil {
		return nil, err
	}
	return time.LoadLocation(entity.CodeValue)
}
