package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type NotesHttpUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	NotesUsecase     *NotesUsecase
	JWTUsecase       *JWTUsecase
	TUsecase         *TUsecase
	KindUsecase      *KindUsecase
	DataEntryUsecase *DataEntryUsecase
	TimelineUsecase  *TimelineUsecase
}

func NewNotesHttpUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	NotesUsecase *NotesUsecase,
	JWTUsecase *JWTUsecase,
	TUsecase *TUsecase,
	KindUsecase *KindUsecase,
	DataEntryUsecase *DataEntryUsecase,
	TimelineUsecase *TimelineUsecase) *NotesHttpUsecase {
	uc := &NotesHttpUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		NotesUsecase:     NotesUsecase,
		JWTUsecase:       JWTUsecase,
		TUsecase:         TUsecase,
		KindUsecase:      KindUsecase,
		DataEntryUsecase: DataEntryUsecase,
		TimelineUsecase:  TimelineUsecase,
	}
	return uc
}

type NotesHttpSaveRequestVo struct {
	Gid     string `json:"gid"` // 值存在就是更新
	Content string `json:"content"`
}

func (c *NotesHttpUsecase) List(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	// 通过notes关联的module
	moduleName := ctx.Param("module_name")
	kindGid := ctx.Param("kind_gid")
	tUser, _ := c.JWTUsecase.JWTUser(ctx)
	page := HandlePage(ctx.Query("page"))
	pageSize := HandlePageSize(ctx.Query("page_size"))

	data, err := c.BizList(ModuleConvertToKind(moduleName), kindGid, rawData, &tUser, page, pageSize)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type NotesResponseList []NotesResponseItem

type NotesResponseItem struct {
	Gid           string                 `json:"gid"`
	CreatedAt     int32                  `json:"created_at"`
	UpdatedAt     int32                  `json:"updated_at"`
	CreatedBy     *ResponseUser          `json:"created_by"`
	ModifiedBy    *ResponseUser          `json:"modified_by"`
	Content       string                 `json:"content"`
	RelatedRecord *ResponseRelatedRecord `json:"related_record"`
}

func (c *NotesHttpUsecase) BizList(kind string, kindGid string, rawBytes []byte, operUser *TData, page int, pageSize int) (lib.TypeMap, error) {
	relatedKind, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return nil, err
	}
	if relatedKind == nil {
		return nil, errors.New("Module incorrect parameter")
	}

	noteKindEntity, err := c.KindUsecase.GetByKind(Kind_notes)
	if err != nil {
		return nil, err
	}
	if kindGid == "" {
		return nil, errors.New("Incorrect parameter")
	}

	cond := Eq{
		Notes_FieldName_kind_gid: kindGid,
		Notes_FieldName_kind:     kind,
	}
	list, err := c.TUsecase.ListByCondWithPaging(*noteKindEntity, cond, "created_at desc", page, pageSize)
	if err != nil {
		return nil, err
	}
	total, err := c.TUsecase.TotalByCond(*noteKindEntity, cond)
	if err != nil {
		return nil, err
	}

	var notesResponseList NotesResponseList
	//caches := lib.CacheInit[*TData]()

	var caseGids []string
	var clientGids []string
	for _, v := range list {
		if v.CustomFields.TextValueByNameBasic(Notes_FieldName_kind) == Kind_client_cases {
			caseGids = append(caseGids, v.CustomFields.TextValueByNameBasic(Notes_FieldName_kind_gid))
		} else if v.CustomFields.TextValueByNameBasic(Notes_FieldName_kind) == Kind_clients {
			clientGids = append(clientGids, v.CustomFields.TextValueByNameBasic(Notes_FieldName_kind_gid))
		}
	}
	caseGids = lib.RemoveDuplicates(caseGids)
	clientGids = lib.RemoveDuplicates(clientGids)

	relaMap := make(TRelaMap)
	relaMap.Set(Kind_client_cases, FieldName_deal_name, caseGids)
	relaMap.Set(Kind_clients, FieldName_full_name, clientGids)
	caches, err := c.TUsecase.GetRelaCachesByRelaMap(relaMap)
	if err != nil {
		return nil, err
	}

	for _, v := range list {
		relatedRecordData := c.TUsecase.DataByGidWithCaches(&caches, relatedKind.Kind, v.CustomFields.TextValueByNameBasic(Notes_FieldName_kind_gid))
		item := NotesResponseItem{
			Gid:           v.Gid(),
			CreatedAt:     v.CreatedAt(),
			UpdatedAt:     v.UpdatedAt(),
			CreatedBy:     v.CustomFields.ToResponseUser(DataEntry_created_by),
			ModifiedBy:    v.CustomFields.ToResponseUser(DataEntry_modified_by),
			Content:       v.CustomFields.TextValueByNameBasic(Notes_FieldName_content),
			RelatedRecord: v.CustomFields.ToResponseRelatedRecord(*relatedKind, Notes_FieldName_kind_gid, relatedRecordData),
		}
		notesResponseList = append(notesResponseList, item)
	}

	data := make(lib.TypeMap)
	data.Set("notes", notesResponseList)
	data.Set(Fab_TTotal, int32(total))
	data.Set(Fab_TPage, page)
	data.Set(Fab_TPageSize, pageSize)

	if int64(page*pageSize) >= total {
		data.Set(Fab_HasMore, false)
	} else {
		data.Set(Fab_HasMore, true)
	}

	return data, nil
}

func (c *NotesHttpUsecase) Save(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	// 通过notes关联的module
	moduleName := ctx.Param("module_name")
	kindGid := ctx.Param("kind_gid")
	lib.DPrintln(moduleName)

	tUser, _ := c.JWTUsecase.JWTUser(ctx)

	data, err := c.BizSave(ModuleConvertToKind(moduleName), kindGid, rawData, &tUser)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *NotesHttpUsecase) BizSave(kind string, kindGid string, rawBytes []byte, operUser *TData) (lib.TypeMap, error) {

	notesHttpSaveRequestVo := lib.BytesToTDef[*NotesHttpSaveRequestVo](rawBytes, nil)
	if notesHttpSaveRequestVo == nil {
		return nil, errors.New("Incorrect parameter")
	}
	if notesHttpSaveRequestVo.Content == "" {
		return nil, errors.New("Incorrect parameter")
	}

	relaData, err := c.TUsecase.DataByGid(kind, kindGid)
	if err != nil {
		return nil, err
	}
	if relaData == nil {
		return nil, errors.New("relaData does not exist")
	}

	motesGid, err := c.NotesUsecase.Save(notesHttpSaveRequestVo.Gid, kind, kindGid, notesHttpSaveRequestVo.Content, operUser)
	if err != nil {
		return nil, err
	}
	tNotes, err := c.TUsecase.DataByGid(Kind_notes, motesGid)
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	data.Set("note", lib.TypeMap{
		"gid":     tNotes.Gid(),
		"content": tNotes.CustomFields.TextValueByNameBasic(Notes_FieldName_content),
	})

	return data, nil
}

type NotesRequestDelete struct {
	Gids []string `json:"gids"`
}

func (c *NotesHttpUsecase) Delete(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	// 通过notes关联的module

	tUser, _ := c.JWTUsecase.JWTUser(ctx)

	data, err := c.BizDelete(rawData, &tUser)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *NotesHttpUsecase) BizDelete(rawBytes []byte, operUser *TData) (lib.TypeMap, error) {

	notesRequestDelete := lib.BytesToTDef[*NotesRequestDelete](rawBytes, nil)
	if notesRequestDelete == nil {
		return nil, errors.New("Incorrect parameter")
	}
	if len(notesRequestDelete.Gids) == 0 {
		return nil, errors.New("Incorrect parameter")
	}

	list, err := c.TUsecase.ListByCond(Kind_notes, In("gid", notesRequestDelete.Gids))
	if err != nil {
		return nil, err
	}

	for _, gid := range notesRequestDelete.Gids {
		isOk := false
		for _, row := range list {
			if row.Gid() == gid {
				isOk = true
				break
			}
		}
		if !isOk {
			return nil, errors.New("Records have been deleted or have no permissions")
		}
	}

	noteKind, err := c.KindUsecase.GetByKind(Kind_notes)
	if err != nil {
		return nil, err
	}

	err = c.DataEntryUsecase.Delete(*noteKind, lib.ConvertToInterfaceSlice(notesRequestDelete.Gids), DataEntry_gid)
	if err != nil {
		return nil, err
	}

	for _, v := range list {
		timelineForNotes := TimelineForNotes{
			Content: v.CustomFields.TextValueByNameBasic(Notes_FieldName_content),
		}
		_, err := c.TimelineUsecase.Create(Kind_notes, v.Gid(),
			Timeline_action_deleted, v.CustomFields.TextValueByNameBasic(Notes_FieldName_kind),
			v.CustomFields.TextValueByNameBasic(Notes_FieldName_kind_gid),
			InterfaceToString(timelineForNotes),
			operUser)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
