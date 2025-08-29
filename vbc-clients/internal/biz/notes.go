package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

const (
	Notes_FieldName_content  = "content"
	Notes_FieldName_kind     = "kind"
	Notes_FieldName_kind_gid = "kind_gid"
)

type NotesUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	DataEntryUsecase *DataEntryUsecase
	TUsecase         *TUsecase
	TimelineUsecase  *TimelineUsecase
}

func NewNotesUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	DataEntryUsecase *DataEntryUsecase,
	TUsecase *TUsecase,
	TimelineUsecase *TimelineUsecase) *NotesUsecase {
	uc := &NotesUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		DataEntryUsecase: DataEntryUsecase,
		TUsecase:         TUsecase,
		TimelineUsecase:  TimelineUsecase,
	}

	return uc
}

func (c *NotesUsecase) GetNoteFacadesByGids(gids []string) (map[string]NoteFacade, error) {

	userFacades := make(map[string]NoteFacade)
	records, err := c.TUsecase.ListByCond(Kind_notes, In(DataEntry_gid, gids))
	if err != nil {
		return userFacades, err
	}
	for k, v := range records {
		userFacades[v.Gid()] = NoteFacade{
			TData: *records[k],
		}
	}
	return userFacades, nil
}

// Save 保存Notes, 如果gid为空，说明是新建
func (c *NotesUsecase) Save(gid string, kind string, kindGid string, content string, operUser *TData) (notesGid string, err error) {

	var tNotes *TData
	isUpdating := false
	if gid == "" {
		notesGid = uuid.UuidWithoutStrike()
	} else {
		notesGid = gid
		isUpdating = true
		tNotes, err = c.TUsecase.Data(Kind_notes, Eq{DataEntry_gid: notesGid})
		if err != nil {
			return "", err
		}
		if tNotes == nil {
			return "", errors.New("tNotes does not exist")
		}
		if tNotes.CustomFields.TextValueByNameBasic(Notes_FieldName_kind_gid) != kindGid {
			return "", errors.New("Disable this operation")
		}
	}

	data := make(TypeDataEntry)
	data[DataEntry_gid] = notesGid
	data[Notes_FieldName_content] = content
	data[Notes_FieldName_kind] = kind
	data[Notes_FieldName_kind_gid] = kindGid
	result, err := c.DataEntryUsecase.HandleOne(Kind_notes, data, DataEntry_gid, operUser)
	if err != nil {
		return "", err
	}
	lib.DPrintln("result:", result)
	if isUpdating {
		r := result.GetByValue(notesGid)
		if r.IsUpdated {

			timelineForNotes := TimelineForNotes{
				Content: content,
			}
			_, err := c.TimelineUsecase.Create(Kind_notes, notesGid, Timeline_action_updated, kind, kindGid, InterfaceToString(timelineForNotes), operUser)
			if err != nil {
				return "", err
			}
		}
	}
	return notesGid, nil
}
