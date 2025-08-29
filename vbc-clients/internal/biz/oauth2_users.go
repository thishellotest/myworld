package biz

type Oauth2UserEntity struct {
	ID        int32 `gorm:"primaryKey"`
	AppId     string
	Account   string
	Email     string
	Name      string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

func (Oauth2UserEntity) TableName() string {
	return "oauth2_users"
}

type Oauth2UserUsecase struct {
}

func NewOauth2UserUsecase() *Oauth2UserUsecase {
	return &Oauth2UserUsecase{}
}
