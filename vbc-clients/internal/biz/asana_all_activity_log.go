package biz

type AsanaAllActivityLogEntity struct {
	ID                            int32 `gorm:"primaryKey"`
	UserGid                       string
	UserResourceType              string
	Action                        string
	AsanaCreatedAt                string
	ResourceGid                   string
	ResourceType                  string
	ResourceSubtype               string
	ParentGid                     string
	ParentResourceType            string
	ParentResourceSubtype         string
	ChangeField                   string
	ChangeAction                  string
	ChangeNewValueGid             string
	ChangeNewValueResourceType    string
	ChangeNewValueResourceSubtype string
	CreatedAt                     int64
}

func (AsanaAllActivityLogEntity) TableName() string {
	return "asana_all_activity_log"
}
