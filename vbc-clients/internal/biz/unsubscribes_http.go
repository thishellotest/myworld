package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type UnsubscribesHttpUsecase struct {
	log                    *log.Helper
	conf                   *conf.Data
	JWTUsecase             *JWTUsecase
	UnsubscribesUsecase    *UnsubscribesUsecase
	ClientCaseUsecase      *ClientCaseUsecase
	ClientUsecase          *ClientUsecase
	TUsecase               *TUsecase
	TimezonesUsecase       *TimezonesUsecase
	LogUsecase             *LogUsecase
	UnsubscribesbuzUsecase *UnsubscribesbuzUsecase
}

func NewUnsubscribesHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	UnsubscribesUsecase *UnsubscribesUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	ClientUsecase *ClientUsecase,
	TUsecase *TUsecase,
	TimezonesUsecase *TimezonesUsecase,
	LogUsecase *LogUsecase,
	UnsubscribesbuzUsecase *UnsubscribesbuzUsecase,
) *UnsubscribesHttpUsecase {
	return &UnsubscribesHttpUsecase{
		log:                    log.NewHelper(logger),
		conf:                   conf,
		JWTUsecase:             JWTUsecase,
		UnsubscribesUsecase:    UnsubscribesUsecase,
		ClientCaseUsecase:      ClientCaseUsecase,
		ClientUsecase:          ClientUsecase,
		TUsecase:               TUsecase,
		TimezonesUsecase:       TimezonesUsecase,
		LogUsecase:             LogUsecase,
		UnsubscribesbuzUsecase: UnsubscribesbuzUsecase,
	}
}

type UnsubscribesHttpChangeStatusRequest struct {
	Action string `json:"action"`
	Id     int32  `json:"id"`
}

const (
	UnsubscribesHttpChangeStatusRequest_Opt_in  = "opt-in"
	UnsubscribesHttpChangeStatusRequest_Opt_out = "opt-out"
)

func (c *UnsubscribesHttpUsecase) Delete(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizDelete(userFacade, body)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *UnsubscribesHttpUsecase) BizDelete(userFacade UserFacade, body lib.TypeMap) (lib.TypeMap, error) {

	id := body.GetInt("id")
	if id <= 0 {
		return nil, errors.New("parameter incorrect")
	}

	entity, err := c.UnsubscribesUsecase.GetByCond(Eq{"id": id})
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, errors.New("parameter incorrect")
	} else {
		entity.BizDeletedAt = time.Now().Unix()
		err = c.UnsubscribesUsecase.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return nil, err
		}
	}
	c.LogUsecase.SaveLog(entity.ID, "UnsubscribesHttp:Delete", map[string]interface{}{
		"ID":      entity.ID,
		"userGid": userFacade.Gid(),
	})

	return nil, nil
}

type UnsubscribesHttpSaveRequest struct {
	Phone string `json:"phone"`
}

func (c *UnsubscribesHttpUsecase) Save(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	request := lib.BytesToTDef[UnsubscribesHttpSaveRequest](rawData, UnsubscribesHttpSaveRequest{})
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizSave(userFacade, request)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *UnsubscribesHttpUsecase) BizSave(userFacade UserFacade, request UnsubscribesHttpSaveRequest) (lib.TypeMap, error) {

	if !IsValidUSAPhoneNumber(request.Phone) {
		return nil, errors.New("The phone format is incorrect")
	}
	phone, err := FormatUSAPhoneHandle(request.Phone)
	if err != nil {
		return nil, err
	}

	entity, err := c.UnsubscribesUsecase.GetByCond(Eq{"contact_phone_number": phone})
	if err != nil {
		return nil, err
	}

	if entity == nil {
		entity = &UnsubscribesEntity{
			CreatedAt:          time.Now().Unix(),
			ContactPhoneNumber: phone,
		}
		entity.LatestFromId = ""
	} else {
		if entity.BizDeletedAt == 0 {
			return nil, errors.New("The phone number is already in the list, please operate in the list.")
		}
		entity.BizDeletedAt = 0

	}
	entity.Status = Unsubscribes_Status_Yes
	entity.UpdatedAt = time.Now().Unix()

	err = c.UnsubscribesUsecase.CommonUsecase.DB().Save(&entity).Error
	if err != nil {
		return nil, err
	}
	c.LogUsecase.SaveLog(entity.ID, "UnsubscribesHttp:Save", map[string]interface{}{
		"ID":      entity.ID,
		"userGid": userFacade.Gid(),
	})
	return nil, nil
}

func (c *UnsubscribesHttpUsecase) ChangeStatus(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	request := lib.BytesToTDef[UnsubscribesHttpChangeStatusRequest](rawData, UnsubscribesHttpChangeStatusRequest{})
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizChangeStatus(userFacade, request)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *UnsubscribesHttpUsecase) BizChangeStatus(userFacade UserFacade, request UnsubscribesHttpChangeStatusRequest) (lib.TypeMap, error) {

	if request.Action != UnsubscribesHttpChangeStatusRequest_Opt_in &&
		request.Action != UnsubscribesHttpChangeStatusRequest_Opt_out {
		return nil, errors.New("Parameter Error")
	}
	if request.Id <= 0 {
		return nil, errors.New("Parameter Error")
	}

	entity, _ := c.UnsubscribesUsecase.GetByCond(Eq{"id": request.Id})
	if entity == nil {
		return nil, errors.New("Record does not exist")
	}
	needSave := false
	if request.Action == UnsubscribesHttpChangeStatusRequest_Opt_in {
		if entity.Status == Unsubscribes_Status_Yes {
			entity.Status = Unsubscribes_Status_No
			entity.UpdatedAt = time.Now().Unix()
			needSave = true
		}
	}
	if request.Action == UnsubscribesHttpChangeStatusRequest_Opt_out {
		if entity.Status == Unsubscribes_Status_No {
			entity.Status = Unsubscribes_Status_Yes
			entity.UpdatedAt = time.Now().Unix()
			needSave = true
		}
	}
	if needSave {
		err := c.UnsubscribesUsecase.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return nil, err
		}
		c.LogUsecase.SaveLog(entity.ID, "UnsubscribesHttp:ChangeStatus", map[string]interface{}{
			"Modified status": entity.Status,
			"userGid":         userFacade.Gid(),
		})
	}
	return nil, nil
}

func (c *UnsubscribesHttpUsecase) List(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	page := HandlePage(ctx.Query("page"))
	pageSize := HandlePageSize(ctx.Query("page_size"))
	data, err := c.BizList(userFacade, page, pageSize)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *UnsubscribesHttpUsecase) BizList(userFacade UserFacade, page int, pageSize int) (lib.TypeMap, error) {

	cond := Eq{"biz_deleted_at": 0}
	records, err := c.UnsubscribesUsecase.ListByCondWithPaging(cond, "id desc", page, pageSize)
	if err != nil {
		return nil, err
	}

	phones := make(lib.TypeMap)
	for _, v := range records {
		p1, p2, p3, _ := FormatPhoneNumber(v.ContactPhoneNumber)
		if p1 != "" {
			phones.Set(p1, 1)
		}
		if p3 != "" {
			phones.Set(p2, 1)
		}
		if p3 != "" {
			phones.Set(p3, 1)
		}
	}
	var phonesList []string
	for k, _ := range phones {
		phonesList = append(phonesList, k)
	}

	clients := make(map[string][]*TData)
	clientCases := make(map[string][]*TData)
	if len(phonesList) > 0 {
		res, err := c.TUsecase.ListByCond(Kind_clients, And(In(FieldName_phone, phonesList), Eq{FieldName_biz_deleted_at: 0}))
		if err != nil {
			return nil, err
		}
		for k, v := range res {
			phone := v.CustomFields.TextValueByNameBasic(FieldName_phone)
			destPhone, _ := FormatUSAPhoneHandle(phone)
			clients[destPhone] = append(clients[destPhone], res[k])
		}

		res1, err := c.TUsecase.ListByCond(Kind_client_cases, And(In(FieldName_phone, phonesList), Eq{FieldName_biz_deleted_at: 0}))
		if err != nil {
			return nil, err
		}
		for k, v := range res1 {
			phone := v.CustomFields.TextValueByNameBasic(FieldName_phone)
			destPhone, _ := FormatUSAPhoneHandle(phone)
			clientCases[destPhone] = append(clientCases[destPhone], res1[k])
		}
	}

	var destRecords lib.TypeList
	for _, v := range records {
		destRecords = append(destRecords, v.ToApi(&userFacade, c.TimezonesUsecase, clients, clientCases))
	}

	//lib.DPrintln("clientCases:", clientCases)
	//lib.DPrintln("clients:", clients)

	total, err := c.UnsubscribesUsecase.Total(cond)
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)
	data.Set(Fab_TRecords, destRecords)
	data.Set(Fab_TTotal, int32(total))
	data.Set(Fab_TPage, page)
	data.Set(Fab_TPageSize, pageSize)

	return data, nil
}
