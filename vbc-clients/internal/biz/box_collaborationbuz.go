package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	. "vbc/lib/builder"
)

func NeedHandleBoxCollaboration(boxUserId string) bool {
	if boxUserId == config_box.YN_BoxUserId ||
		boxUserId == config_box.ED_BoxUserId ||
		boxUserId == config_box.VBCTeam_BoxUserId ||
		boxUserId == config_box.Donald_BoxUserId ||
		boxUserId == config_box.Victoria_BoxUserId ||
		boxUserId == config_box.LILI_BoxUserId {
		return false
	}
	return true
}

type BoxCollaborationBuzUsecase struct {
	log                     *log.Helper
	conf                    *conf.Data
	CommonUsecase           *CommonUsecase
	BoxCollaborationUsecase *BoxCollaborationUsecase
	BoxUsecase              *BoxUsecase
	TUsecase                *TUsecase
	BoxUserUsecase          *BoxUserUsecase
	UserUsecase             *UserUsecase
	BoxbuzUsecase           *BoxbuzUsecase
	MapUsecase              *MapUsecase
}

func NewBoxCollaborationBuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	BoxCollaborationUsecase *BoxCollaborationUsecase,
	BoxUsecase *BoxUsecase,
	TUsecase *TUsecase,
	BoxUserUsecase *BoxUserUsecase,
	UserUsecase *UserUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	MapUsecase *MapUsecase,
) *BoxCollaborationBuzUsecase {
	uc := &BoxCollaborationBuzUsecase{
		log:                     log.NewHelper(logger),
		CommonUsecase:           CommonUsecase,
		conf:                    conf,
		BoxCollaborationUsecase: BoxCollaborationUsecase,
		BoxUsecase:              BoxUsecase,
		TUsecase:                TUsecase,
		BoxUserUsecase:          BoxUserUsecase,
		UserUsecase:             UserUsecase,
		BoxbuzUsecase:           BoxbuzUsecase,
		MapUsecase:              MapUsecase,
	}

	return uc
}

type UserRelatedBox struct {
	TUser   TData
	BoxUser *BoxUserEntity
}

func GetBoxUserByUser(boxUsers []*BoxUserEntity, tUser TData) (boxUser *BoxUserEntity) {
	userBoxUserId := tUser.CustomFields.TextValueByNameBasic(UserFieldName_box_user_id)
	userEmail := tUser.CustomFields.TextValueByNameBasic(UserFieldName_email)
	for k, v := range boxUsers {
		if userBoxUserId != "" {
			if userBoxUserId == v.BoxUserId {
				return boxUsers[k]
			}
		} else {
			if userEmail == v.Login {
				return boxUsers[k]
			}
		}
	}
	return nil
}

type RelatedBoxUserIdsMap map[string]UserRelatedBox

func (c RelatedBoxUserIdsMap) GetBoxUserIdByUserGid(userGid string) (boxUserId string) {

	if c == nil {
		return ""
	}
	if _, ok := c[userGid]; ok {
		if c[userGid].BoxUser != nil {
			return c[userGid].BoxUser.BoxUserId
		}
	}
	return ""
}

func (c *BoxCollaborationBuzUsecase) GetRelatedBoxUserIds(userGids []string) (res RelatedBoxUserIdsMap, err error) {
	if len(userGids) == 0 {
		return
	}
	users, err := c.TUsecase.ListByCond(Kind_users, In("gid", userGids))
	if err != nil {
		return nil, err
	}
	boxUsers, err := c.BoxUserUsecase.AllByCond(Eq{"deleted_at": 0})
	if err != nil {
		return nil, err
	}
	res = make(RelatedBoxUserIdsMap)
	for k, v := range users {
		res[v.Gid()] = UserRelatedBox{
			TUser:   *users[k],
			BoxUser: GetBoxUserByUser(boxUsers, *users[k]),
		}
	}
	return res, nil
}

const (
	HandleCollaborationFromCase_BizType_ClientFolder = "ClientFolder"
	HandleCollaborationFromCase_BizType_DCFolder     = "DCFolder"
)

func (c *BoxCollaborationBuzUsecase) HandleUseVBCActiveCases(caseId int32) error {
	key := MapKeyUseVBCActiveCasesFolder(caseId)
	return c.MapUsecase.Set(key, "1")
}

func (c *BoxCollaborationBuzUsecase) UseVBCActiveCases(caseId int32) (use bool, err error) {
	key := MapKeyUseVBCActiveCasesFolder(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return false, err
	}
	if val == "1" {
		return true, nil
	}
	return false, nil
}

func (c *BoxCollaborationBuzUsecase) DoAddPermissionForBox(caseId int32) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
	}
	if tCase != nil {
		err = c.DoHandleCollaborationFromCaseWithDCFolder(*tCase)
		if err != nil {
			c.log.Error(err, " DoHandleCollaborationFromCaseWithDCFolder: ", caseId)
		}
		err = c.DoHandleCollaborationFromCaseWithClientFolder(*tCase)
		if err != nil {
			c.log.Error(err, " DoHandleCollaborationFromCaseWithDCFolder: ", caseId)
		}
	}
}

func (c *BoxCollaborationBuzUsecase) DoHandleCollaborationFromCaseWithClientFolder(tCase TData) error {

	val := tCase.CustomFields.TextValueByNameBasic(FieldName_case_files_folder)
	if val != "" {
		useVBCActiveFolder, err := c.UseVBCActiveCases(tCase.Id())
		if err != nil {
			return err
		}
		if !useVBCActiveFolder {
			return nil
		}
		clientBoxFolderId, err := c.BoxbuzUsecase.GetClientBoxFolderId(&tCase)
		if err != nil {
			return err
		}
		if clientBoxFolderId != "" {
			return c.HandleCollaborationFromCase(tCase, HandleCollaborationFromCase_BizType_ClientFolder)
		}
	}
	return nil
}

func (c *BoxCollaborationBuzUsecase) DoHandleCollaborationFromCaseWithDCFolder(tCase TData) error {

	val := tCase.CustomFields.TextValueByNameBasic(FieldName_data_collection_folder)
	if val != "" {

		useVBCActiveFolder, err := c.UseVBCActiveCases(tCase.Id())
		if err != nil {
			return err
		}
		if !useVBCActiveFolder {
			return nil
		}

		dcBoxFolderId, err := c.BoxbuzUsecase.GetDCFolderId(tCase.Id())
		if err != nil {
			return err
		}
		if dcBoxFolderId != "" {
			return c.HandleCollaborationFromCase(tCase, HandleCollaborationFromCase_BizType_DCFolder)
		}
	}
	return nil
}

func (c *BoxCollaborationBuzUsecase) HandleCollaborationFromCase(tCase TData, bizType string) error {

	if bizType != HandleCollaborationFromCase_BizType_ClientFolder &&
		bizType != HandleCollaborationFromCase_BizType_DCFolder {
		return errors.New("HandleCollaborationFromCase bizType is wrong")
	}

	var destBoxFolderId string
	if bizType == HandleCollaborationFromCase_BizType_ClientFolder {
		clientBoxFolderId, err := c.BoxbuzUsecase.GetClientBoxFolderId(&tCase)
		if err != nil {
			return err
		}
		if clientBoxFolderId == "" {
			return errors.New("clientBoxFolderId is empty")
		}
		destBoxFolderId = clientBoxFolderId
	} else if bizType == HandleCollaborationFromCase_BizType_DCFolder {
		dcBoxFolderId, err := c.BoxbuzUsecase.GetDCFolderId(tCase.Id())
		if err != nil {
			return err
		}
		if dcBoxFolderId == "" {
			return errors.New("dcBoxFolderId is empty")
		}
		destBoxFolderId = dcBoxFolderId
	}

	var allowedUserGidsForBoxFolder []string
	userGid := tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid)
	if userGid != "" {
		allowedUserGidsForBoxFolder = append(allowedUserGidsForBoxFolder, userGid)
	}
	leadVS, _ := c.UserUsecase.GetUserByLeadVS(&tCase)
	leadCP, _ := c.UserUsecase.GetUserByLeadCP(&tCase)
	if leadVS != nil {
		allowedUserGidsForBoxFolder = append(allowedUserGidsForBoxFolder, leadVS.Gid())
	}
	if leadCP != nil {
		allowedUserGidsForBoxFolder = append(allowedUserGidsForBoxFolder, leadCP.Gid())
	}
	supportCP, _ := c.UserUsecase.GetUserBySupportCP(&tCase)
	if supportCP != nil {
		allowedUserGidsForBoxFolder = append(allowedUserGidsForBoxFolder, supportCP.Gid())
	}
	leadCO, _ := c.UserUsecase.GetUserByLeadCO(&tCase)
	if leadCO != nil {
		allowedUserGidsForBoxFolder = append(allowedUserGidsForBoxFolder, leadCO.Gid())
	}

	relatedBoxUserIdsMap, err := c.GetRelatedBoxUserIds(allowedUserGidsForBoxFolder)
	if err != nil {
		return err
	}
	if userGid != "" {
		desUserGid := userGid
		boxUserId := relatedBoxUserIdsMap.GetBoxUserIdByUserGid(desUserGid)
		if boxUserId == "" {
			c.log.Error("desUserGid: ", desUserGid, " The corresponding BoxUserID was not found ")
		} else {
			err = c.HandleAddCollaboration(Box_collaboration_ow, destBoxFolderId, boxUserId, tCase.Id(), desUserGid)
			if err != nil {
				c.log.Error("HandleAddCollaboration error: ", " destBoxFolderId: ", destBoxFolderId,
					" boxUserId: ", boxUserId, " caseId: ", tCase.Id(), " desUserGid: ", desUserGid)
			}
		}
	}
	if leadVS != nil {
		desUserGid := leadVS.Gid()
		boxUserId := relatedBoxUserIdsMap.GetBoxUserIdByUserGid(desUserGid)
		if boxUserId == "" {
			c.log.Error("desUserGid: ", desUserGid, " The corresponding BoxUserID was not found ")
		} else {
			err = c.HandleAddCollaboration(Box_collaboration_vs, destBoxFolderId, boxUserId, tCase.Id(), desUserGid)
			if err != nil {
				c.log.Error("HandleAddCollaboration error: ", " destBoxFolderId: ", destBoxFolderId,
					" boxUserId: ", boxUserId, " caseId: ", tCase.Id(), " desUserGid: ", desUserGid)
			}
		}
	}

	if leadCP != nil {
		desUserGid := leadCP.Gid()
		boxUserId := relatedBoxUserIdsMap.GetBoxUserIdByUserGid(desUserGid)
		if boxUserId == "" {
			c.log.Error("desUserGid: ", desUserGid, " The corresponding BoxUserID was not found ")
		} else {
			err = c.HandleAddCollaboration(Box_collaboration_cp, destBoxFolderId, boxUserId, tCase.Id(), desUserGid)
			if err != nil {
				c.log.Error("HandleAddCollaboration error: ", " destBoxFolderId: ", destBoxFolderId,
					" boxUserId: ", boxUserId, " caseId: ", tCase.Id(), " desUserGid: ", desUserGid)
			}
		}
	}
	if supportCP != nil {
		desUserGid := supportCP.Gid()
		boxUserId := relatedBoxUserIdsMap.GetBoxUserIdByUserGid(desUserGid)
		if boxUserId == "" {
			c.log.Error("desUserGid: ", desUserGid, " The corresponding BoxUserID was not found ")
		} else {
			err = c.HandleAddCollaboration(Box_collaboration_support_cp, destBoxFolderId, boxUserId, tCase.Id(), desUserGid)
			if err != nil {
				c.log.Error("HandleAddCollaboration error: ", " destBoxFolderId: ", destBoxFolderId,
					" boxUserId: ", boxUserId, " caseId: ", tCase.Id(), " desUserGid: ", desUserGid)
			}
		}
	}

	if leadCO != nil {
		desUserGid := leadCO.Gid()
		boxUserId := relatedBoxUserIdsMap.GetBoxUserIdByUserGid(desUserGid)
		if boxUserId == "" {
			c.log.Error("desUserGid: ", desUserGid, " The corresponding BoxUserID was not found ")
		} else {
			err = c.HandleAddCollaboration(Box_collaboration_lead_co, destBoxFolderId, boxUserId, tCase.Id(), desUserGid)
			if err != nil {
				c.log.Error("HandleAddCollaboration error: ", " destBoxFolderId: ", destBoxFolderId,
					" boxUserId: ", boxUserId, " caseId: ", tCase.Id(), " desUserGid: ", desUserGid)
			}
		}
	}

	return nil
}

func (c *BoxCollaborationBuzUsecase) HandleAddCollaboration(permissionSource string, boxFolderId string, boxUserId string, caseId int32, userGid string) error {

	if !NeedHandleBoxCollaboration(boxUserId) {
		return nil
	}

	records, err := c.BoxCollaborationUsecase.AllByCond(Eq{"box_folder_id": boxFolderId,
		"box_user_id": boxUserId,
		"deleted_at":  0,
	})
	if err != nil {
		return err
	}
	var boxCollaborationId string

	hasPermission := false
	if len(records) == 0 {
		boxCollaborationId, _, err = c.BoxUsecase.CollaborationsByBoxUserId(boxFolderId, boxUserId)
		if err != nil {
			return err
		}
	} else {

		for _, v := range records {
			if v.PermissionSource == permissionSource && v.UserGid == userGid { // 已经加过权限了
				hasPermission = true
				break
			}
			// 使用其它来源的权限
			boxCollaborationId = v.BoxCollaborationId
		}
	}
	if !hasPermission {
		newEntity := &BoxCollaborationEntity{
			CaseId:             caseId,
			BoxFolderId:        boxFolderId,
			PermissionSource:   permissionSource,
			BoxUserId:          boxUserId,
			UserGid:            userGid,
			BoxCollaborationId: boxCollaborationId,
			CreatedAt:          time.Now().Unix(),
			UpdatedAt:          time.Now().Unix(),
		}
		return c.CommonUsecase.DB().Save(&newEntity).Error
	}
	return nil
}

func (c *BoxCollaborationBuzUsecase) RunHandleDeleteCollaborationByUserFullName(permissionSource string, caseId int32, userFullName string) error {

	tUser, err := c.UserUsecase.GetByFullName(userFullName)
	if err != nil {
		return err
	}
	if tUser == nil {
		c.log.Error("tUser is nil; UserFullName:" + userFullName)
		return errors.New("tUser is nil")
	}
	return c.RunHandleDeleteCollaborationByCaseId(permissionSource, caseId, tUser.Gid())
}

func (c *BoxCollaborationBuzUsecase) RunHandleDeleteCollaborationByCaseId(permissionSource string, caseId int32, userGid string) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	return c.RunHandleDeleteCollaboration(permissionSource, *tCase, userGid)
}

// RunHandleDeleteCollaboration userGid:需要删除权限的用户
func (c *BoxCollaborationBuzUsecase) RunHandleDeleteCollaboration(permissionSource string, tCase TData, userGid string) error {

	useVBCActiveFolder, err := c.UseVBCActiveCases(tCase.Id())
	if err != nil {
		return err
	}
	if !useVBCActiveFolder {
		return nil
	}

	clientBoxFolderId, err := c.BoxbuzUsecase.GetClientBoxFolderId(&tCase)
	if err != nil {
		c.log.Error(err)
	} else if clientBoxFolderId != "" {
		err = c.DoHandleDeleteCollaboration(permissionSource, clientBoxFolderId, userGid)
		if err != nil {
			c.log.Error(err, " DoHandleDeleteCollaboration: ", permissionSource, " clientBoxFolderId:", clientBoxFolderId, " userGid:", userGid)
		}
	}

	dcFolderId, err := c.BoxbuzUsecase.GetDCFolderId(tCase.Id())
	if err != nil {
		c.log.Error(err)
	} else if dcFolderId != "" {
		err = c.DoHandleDeleteCollaboration(permissionSource, dcFolderId, userGid)
		if err != nil {
			c.log.Error(err, " DoHandleDeleteCollaboration: ", permissionSource, " dcFolderId:", dcFolderId, " userGid:", userGid)
		}
	}
	return nil
}

func (c *BoxCollaborationBuzUsecase) DoHandleDeleteCollaboration(permissionSource string, boxFolderId string, userGid string) error {
	relatedBoxUserIdsMap, err := c.GetRelatedBoxUserIds([]string{userGid})
	if err != nil {
		return err
	}
	boxUserId := relatedBoxUserIdsMap.GetBoxUserIdByUserGid(userGid)
	if boxUserId == "" {
		c.log.Error("desUserGid: ", userGid, " The corresponding BoxUserID was not found ")
		return nil
	}
	return c.HandleDeleteCollaboration(permissionSource, boxFolderId, boxUserId, userGid)
}

// HandleDeleteCollaboration 处理删除权限
func (c *BoxCollaborationBuzUsecase) HandleDeleteCollaboration(permissionSource string, boxFolderId string, boxUserId string, userGid string) error {

	if !NeedHandleBoxCollaboration(boxUserId) {
		return nil
	}

	records, err := c.BoxCollaborationUsecase.AllByCond(Eq{"box_folder_id": boxFolderId,
		"box_user_id": boxUserId,
		"deleted_at":  0,
	})
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	needInvokeApi := false

	var collaborationId string
	for k, v := range records {
		if v.PermissionSource == permissionSource && v.UserGid == userGid {
			entity := records[k]
			entity.DeletedAt = time.Now().Unix()
			entity.UpdatedAt = time.Now().Unix()
			err = c.CommonUsecase.DB().Save(&entity).Error
			if err != nil {
				return err
			}
			if len(records) == 1 {
				needInvokeApi = true
			}
			collaborationId = entity.BoxCollaborationId
		}
	}
	if needInvokeApi && collaborationId != "" {
		_, er := c.BoxUsecase.DeleteCollaborations(collaborationId)
		if er != nil {
			c.log.Error("BoxUsecase.DeleteCollaborations: ", collaborationId, er)
		}
	}
	return nil
}
