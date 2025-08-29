package biz

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
)

const (
	UserFieldName_mobile              = "mobile"
	UserFieldName_dialpad_userid      = "dialpad_userid"
	UserFieldName_dialpad_phonenumber = "dialpad_phonenumber"
	UserFieldName_first_name          = "first_name"
	UserFieldName_last_name           = "last_name"
	UserFieldName_title               = "title"
	UserFieldName_fullname            = "full_name"
	UserFieldName_status              = "status"
	UserFieldName_email               = "email"
	UserFieldName_role_gid            = "role_gid"
	UserFieldName_gender              = "gender"
	User_FieldName_profile_gid        = "profile_gid"
	User_FieldName_timezone_id        = "timezone_id"
	UserFieldName_MailSender          = "mail_username"
	UserFieldName_MailPassword        = "mail_password"
	UserFieldName_permissions         = "permissions"
	UserFieldName_box_user_id         = "box_user_id"
)

type UserEntity struct {
	ID int32 `gorm:"primaryKey"`
	//Phone string
	//UserName  string
	Email     string
	Password  string
	FullName  string
	FirstName string
	LastName  string
	Title     string
	Mobile    string
	PicUrl    string
	UpdatedAt int64
}

func (UserEntity) TableName() string {
	return "users"
}

func UserToInfoApi(uf *UserFacade, userProfile TData, userRole TData, TimezonesUsecase *TimezonesUsecase) lib.TypeMap {
	if uf == nil {
		return nil
	}
	info := make(lib.TypeMap)
	info.Set("id", uf.CustomFields.NumberValueByNameBasic("id"))
	info.Set("gid", uf.Gid())
	info.Set("full_name", uf.CustomFields.TextValueByNameBasic("full_name"))
	info.Set("email", uf.CustomFields.TextValueByNameBasic("email"))
	info.Set("pic_url", uf.CustomFields.TextValueByNameBasic("pic_url"))
	info.Set("is_admin", userProfile.CustomFields.NumberValueByNameBasic(Profile_FieldName_is_admin))
	//info.Set("is_user_manager", userRole.CustomFields.NumberValueByNameBasic(RoleFieldName_is_user_manager))

	timezoneEntity, _ := uf.GetTimezonesEntity(TimezonesUsecase)
	if timezoneEntity != nil {
		info.Set("timezone", timezoneEntity.CodeValue)
		info.Set("timezone_title", timezoneEntity.Title)
	} else {
		info.Set("timezone", Default_Timezones_CodeValue)
		info.Set("timezone_title", "Pacific Time (PT)")
	}
	return info
}
func (c *UserEntity) ToInfo() lib.TypeMap {
	info := make(lib.TypeMap)
	info.Set("id", c.ID)
	info.Set("full_name", c.FullName)
	info.Set("email", c.Email)
	info.Set("pic_url", c.PicUrl)
	return info
}

func UserToRelaApi(tUser *TData) lib.TypeMap {
	if tUser == nil {
		return nil
	}
	data := make(lib.TypeMap)
	data.Set("gid", tUser.CustomFields.TextValueByNameBasic("gid"))
	data.Set("full_name", tUser.CustomFields.TextValueByNameBasic("full_name"))
	data.Set("email", tUser.CustomFields.TextValueByNameBasic("email"))
	data.Set("pic_url", tUser.CustomFields.TextValueByNameBasic("pic_url"))
	return data

}

func (c *UserEntity) UserToRelaApi() lib.TypeMap {
	info := make(lib.TypeMap)
	info.Set("id", c.ID)
	info.Set("full_name", c.FullName)
	info.Set("email", c.Email)
	info.Set("pic_url", c.PicUrl)
	return info
}

type UserUsecase struct {
	CommonUsecase *CommonUsecase
	TUsecase      *TUsecase
	DBUsecase[UserEntity]
	BUsecase *BUsecase
}

func NewUserUsecase(CommonUsecase *CommonUsecase, TUsecase *TUsecase, BUsecase *BUsecase) *UserUsecase {
	uc := &UserUsecase{
		CommonUsecase: CommonUsecase,
		TUsecase:      TUsecase,
		BUsecase:      BUsecase,
	}
	uc.DB = CommonUsecase.DB()
	return uc
}

func (c *UserUsecase) GetUserFacadesByGids(gids []string) (map[string]UserFacade, error) {

	userFacades := make(map[string]UserFacade)
	records, err := c.TUsecase.ListByCond(Kind_users, In(DataEntry_gid, gids))
	if err != nil {
		return userFacades, err
	}
	for k, v := range records {
		userFacades[v.Gid()] = UserFacade{
			TData: *records[k],
		}
	}
	return userFacades, nil
}

func (c *UserUsecase) GetUserFacadeByGid(gid string) (*UserFacade, error) {
	tUser, err := c.TUsecase.DataByGid(Kind_users, gid)
	if err != nil {
		return nil, err
	}
	if tUser == nil {
		return nil, nil
	}
	userFacade := UserFacade{
		TData: *tUser,
	}
	return &userFacade, nil
}

func (c *UserUsecase) GetUserFacadeById(id int32) (*UserFacade, error) {
	tUser, err := c.TUsecase.DataById(Kind_users, id)
	if err != nil {
		return nil, err
	}
	if tUser == nil {
		return nil, nil
	}
	userFacade := UserFacade{
		TData: *tUser,
	}
	return &userFacade, nil
}

func (c *UserUsecase) FetchByUserId(userId int32) (*UserEntity, error) {
	var entity UserEntity
	err := c.CommonUsecase.DB().Where("id=?", userId).Take(&entity).Error
	if err == nil {
		return &entity, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (c *UserUsecase) GetUserByLeadVS(tCase *TData) (tUser *TData, err error) {

	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	primaryVs := tCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs)
	if primaryVs == "" {
		return nil, errors.New("primaryVs is empty")
	}
	tUser, err = c.GetByFullName(primaryVs)
	if err != nil {
		return nil, err
	}
	if tUser == nil {
		return nil, errors.New("tUser is nil")
	}
	return tUser, nil
}

func (c *UserUsecase) GetUserByLeadCP(tCase *TData) (tUser *TData, err error) {

	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	primaryCP := tCase.CustomFields.TextValueByNameBasic(FieldName_primary_cp)
	if primaryCP == "" {
		return nil, errors.New("primaryCP is empty")
	}
	tUser, err = c.GetByFullName(primaryCP)
	if err != nil {
		return nil, err
	}
	if tUser == nil {
		return nil, errors.New("tUser is nil")
	}
	return tUser, nil
}

func (c *UserUsecase) GetUserBySupportCP(tCase *TData) (tUser *TData, err error) {

	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	supportCP := tCase.CustomFields.TextValueByNameBasic(FieldName_support_cp)
	if supportCP == "" {
		return nil, errors.New("supportCP is empty")
	}
	tUser, err = c.GetByFullName(supportCP)
	if err != nil {
		return nil, err
	}
	if tUser == nil {
		return nil, errors.New("tUser is nil")
	}
	return tUser, nil
}

func (c *UserUsecase) GetUserByLeadCO(tCase *TData) (tUser *TData, err error) {

	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	leadCO := tCase.CustomFields.TextValueByNameBasic(FieldName_lead_co)
	if leadCO == "" {
		return nil, errors.New("leadCO is empty")
	}
	tUser, err = c.GetByFullName(leadCO)
	if err != nil {
		return nil, err
	}
	if tUser == nil {
		return nil, errors.New("tUser is nil")
	}
	return tUser, nil
}

func (c *UserUsecase) GetByFullName(fullName string) (*TData, error) {
	return c.TUsecase.Data(Kind_users, Eq{"full_name": fullName})
}

func (c *UserUsecase) GetByEmail(email string) (*TData, error) {
	return c.TUsecase.Data(Kind_users, Eq{"email": email, "biz_deleted_at": 0})
}

func (c *UserUsecase) GetByMicrosoftEmail(email string) (*TData, error) {
	return c.TUsecase.Data(Kind_users, Eq{"microsoft_email": email, "biz_deleted_at": 0})
}

func (c *UserUsecase) GetByGid(userGid string) (*TData, error) {
	return c.TUsecase.Data(Kind_users, Eq{"gid": userGid})
}

func (c *UserUsecase) InitPassword(email string) error {

	user, err := c.GetByCond(Eq{"email": email})
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user is nil")
	}
	pwd := lib.GeneratePassword(16)

	encryptPwd, err := c.EncryptPassword(pwd)
	if err != nil {
		return err
	}
	user.Password = encryptPwd
	user.UpdatedAt = time.Now().Unix()

	return c.CommonUsecase.DB().Save(&user).Error
}

func (c *UserUsecase) EncryptPassword(password string) (encryptPassword string, err error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (c *UserUsecase) VerifyPassword(originPassword string, encryptPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(originPassword))
	if err != nil {
		return false
	}
	return true
}

func (c *UserUsecase) GetUserWithCache(caches lib.Cache[*TData], gid string) (*TData, error) {
	key := gid
	entity, exists := caches.Get(key)
	if exists {
		return entity, nil
	}

	entity, err := c.TUsecase.DataByGid(Kind_users, gid)
	if err != nil {
		return nil, err
	}
	caches.Set(key, entity)
	return entity, nil
}

func (c *UserUsecase) GetProfile(tUser *TData) (tProfile *TData, err error) {
	if tUser == nil {
		return nil, errors.New("GetProfile: tUser is nil")
	}
	return tUser.RelaData(c.BUsecase, User_FieldName_profile_gid)
}

func (c *UserUsecase) VSTeamUsers() ([]*TData, error) {
	return c.TUsecase.ListByCond(Kind_users, Eq{"deleted_at": 0, "biz_deleted_at": 0, "role_gid": config_vbc.Role_VSTeam_gid})
}
