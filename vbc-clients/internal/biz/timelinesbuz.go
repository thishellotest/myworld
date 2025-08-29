package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type TimelinesVo struct {
	Timelines []TimelinesItemVo `json:"timelines"`
}

type TimelinesItemVo struct {
	Gid           string                        `json:"gid"`
	CreatedBy     *ResponseUser                 `json:"created_by"`
	RelatedRecord *ResponseRelatedRecord        `json:"related_record"`
	Record        TimelinesRecordVo             `json:"record"`
	Action        string                        `json:"action"`
	CreatedAt     int32                         `json:"created_at"` // 创建时间
	FieldHistory  []TimelinesFieldHistoryItemVo `json:"field_history"`
}

type TimelinesFieldHistoryItemVo struct {
	FieldName  string `json:"field_name"`
	FieldType  string `json:"field_type"`
	FieldLabel string `json:"field_label"`
	//Value        TimelinesFieldHistoryItemValueVo        `json:"value"`
	//DisplayValue TimelinesFieldHistoryItemDisplayValueVo `json:"display_value"`

	NewMultiValues TFieldMultiValues `json:"new_multi_values"`
	OldMultiValues TFieldMultiValues `json:"old_multi_values"`

	OldValue *TimelinesFieldHistoryItemValue `json:"old_value"`
	NewValue *TimelinesFieldHistoryItemValue `json:"new_value"`
}

type TimelinesFieldHistoryItemValue struct {
	Label string `json:"label"` // DisplayValue
	Value string `json:"value"` // Value
}

//type TimelinesFieldHistoryItemValueVo struct {
//	NewValue string `json:"new_value"`
//	OldValue string `json:"old_value"`
//}

//type TimelinesFieldHistoryItemDisplayValueVo struct {
//	NewValue string `json:"new_value"`
//	OldValue string `json:"old_value"`
//}

//type TimelinesCreatedByVo struct {
//	Gid  string `json:"gid"`
//	Name string `json:"name"` // 一般用于显示
//}
//
//type TimelinesRelatedRecordVo struct {
//	ModuleName  string `json:"module_name"`
//	ModuleLabel string `json:"module_label"`
//	Gid         string `json:"gid"`
//	Name        string `json:"name"` // 一般用于显示
//}

type TimelinesRecordVo struct {
	ModuleName  string `json:"module_name"`
	ModuleLabel string `json:"module_label"`
	Gid         string `json:"gid"`
	Name        string `json:"name"` // 一般用于显示
}

type TimelinesbuzUsecase struct {
	log                    *log.Helper
	CommonUsecase          *CommonUsecase
	conf                   *conf.Data
	KindUsecase            *KindUsecase
	TUsecase               *TUsecase
	FieldUsecase           *FieldUsecase
	FieldOptionUsecase     *FieldOptionUsecase
	TRelaUsecase           *TRelaUsecase
	FieldPermissionUsecase *FieldPermissionUsecase
}

func NewTimelinesbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	KindUsecase *KindUsecase,
	TUsecase *TUsecase,
	FieldUsecase *FieldUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
	TRelaUsecase *TRelaUsecase,
	FieldPermissionUsecase *FieldPermissionUsecase) *TimelinesbuzUsecase {
	uc := &TimelinesbuzUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		KindUsecase:            KindUsecase,
		TUsecase:               TUsecase,
		FieldUsecase:           FieldUsecase,
		FieldOptionUsecase:     FieldOptionUsecase,
		TRelaUsecase:           TRelaUsecase,
		FieldPermissionUsecase: FieldPermissionUsecase,
	}
	return uc
}

func (c *TimelinesbuzUsecase) GetCaches(list TDataList) (lib.Cache[*TData], error) {

	tRelaMap := TRelaMapInit()
	//caseKindEntity, _ := c.KindUsecase.GetByKind(Kind_client_cases)
	//if caseKindEntity == nil {
	//	return timelinesVo, total, errors.New("caseKindEntity is nil")
	//}
	//clientKindEntity, _ := c.KindUsecase.GetByKind(Kind_clients)
	//if clientKindEntity == nil {
	//	return timelinesVo, total, errors.New("clientKindEntity is nil")
	//}

	for _, tTimeline := range list {
		timelineKind := tTimeline.CustomFields.TextValueByNameBasic(Timeline_FieldName_kind)
		gid := tTimeline.CustomFields.TextValueByNameBasic(Timeline_FieldName_kind_gid)
		if gid == "" {
			continue
		}
		if timelineKind == Kind_notes {

		} else { // 支持cases / clients
			//timelineKindEntity, err := c.KindUsecase.GetByKind(timelineKind)
			//if err != nil {
			//	c.log.Error(err)
			//	continue
			//}
			//if timelineKindEntity == nil {
			//	c.log.Error("timelineKindEntity is nil")
			//	continue
			//}
			recordKindFields, err := c.FieldUsecase.CacheStructByKind(timelineKind)
			if err != nil {
				c.log.Error(err)
				continue
			}
			note := lib.StringToTDef[TimelineFieldHistoryNotes](tTimeline.CustomFields.TextValueByNameBasic(Timeline_FieldName_notes), TimelineFieldHistoryNotes{})
			for _, v := range note.FieldHistory {
				//var timelinesFieldHistoryItemVo TimelinesFieldHistoryItemVo
				// todo:lgl 权限过虑
				fieldEntity := recordKindFields.GetByFieldName(v.FieldName)
				if fieldEntity != nil {
					//timelinesFieldHistoryItemVo.FieldType = fieldEntity.FieldType
					//timelinesFieldHistoryItemVo.FieldName = fieldEntity.FieldName
					//timelinesFieldHistoryItemVo.FieldLabel = fieldEntity.FieldLabel

					if fieldEntity.FieldType == FieldType_multilookup {
						val := v.OldValue
						if val != "" {
							vals := strings.Split(val, ",")
							tRelaMap.Set(fieldEntity.RelaKind, fieldEntity.RelaName, vals)
						}
						nVal := v.NewValue
						if nVal != "" {
							nVals := strings.Split(nVal, ",")
							tRelaMap.Set(fieldEntity.RelaKind, fieldEntity.RelaName, nVals)
						}
					} else if fieldEntity.FieldType == FieldType_lookup {
						val := v.OldValue
						if val != "" {
							tRelaMap.Set(fieldEntity.RelaKind, fieldEntity.RelaName, []string{val})
						}
						nVal := v.NewValue
						if nVal != "" {
							tRelaMap.Set(fieldEntity.RelaKind, fieldEntity.RelaName, []string{nVal})
						}
					}
				}
			}
		}
	}
	caches, _ := c.TUsecase.GetRelaCachesByRelaMap(tRelaMap)
	return caches, nil
}

func (c *TimelinesbuzUsecase) List(kind string, kindGid string, userFacade *UserFacade, page int, pageSize int) (timelinesVo TimelinesVo, total int64, err error) {

	noteKindEntity, _ := c.KindUsecase.GetByKind(Kind_timelines)
	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return timelinesVo, total, err
	}
	if kindEntity.NoTimelines == NoChangeHistory_Yes {
		return timelinesVo, total, errors.New(kind + " is no timelines")
	}
	if kindGid == "" {
		return timelinesVo, total, errors.New("Gid is empty")
	}
	var conds []Cond
	conds = append(conds, Or(Eq{Timeline_FieldName_kind: kind, Timeline_FieldName_kind_gid: kindGid},
		Eq{Timeline_FieldName_related_kind: kind, Timeline_FieldName_related_kind_gid: kindGid}))

	list, err := c.TUsecase.ListByCondWithPaging(*noteKindEntity, And(conds...), "id desc", page, pageSize)
	if err != nil {
		return timelinesVo, total, err
	}

	total, err = c.TUsecase.TotalByCond(*noteKindEntity, And(conds...))
	if err != nil {
		return timelinesVo, total, err
	}

	//caches := lib.CacheInit[*TData]()
	caches, _ := c.GetCaches(list)

	for k, v := range list {
		var timelinesItemVo TimelinesItemVo
		timelinesItemVo.Gid = v.Gid()
		timelinesItemVo.Action = v.CustomFields.TextValueByNameBasic(Timeline_FieldName_action)
		timelinesItemVo.CreatedAt = v.CustomFields.NumberValueByNameBasic(DataEntry_created_at)
		createdByGid := v.CustomFields.TextValueByNameBasic(DataEntry_created_by)
		if createdByGid != "" {
			timelinesItemVo.CreatedBy = &ResponseUser{
				Gid:  v.CustomFields.TextValueByNameBasic(DataEntry_created_by),
				Name: v.CustomFields.DisplayValueByNameBasic(DataEntry_created_by),
			}
		}

		record, fieldHistory, err := c.TimelinesItemRecord(&caches, list[k], userFacade)
		if err != nil {
			return TimelinesVo{}, 0, err
		}
		timelinesItemVo.Record = record
		timelinesItemVo.FieldHistory = fieldHistory

		timelinesVo.Timelines = append(timelinesVo.Timelines, timelinesItemVo)

	}

	//lib.DPrintln("list:", list, err)
	//lib.DPrintln("total:", total)

	return timelinesVo, total, nil
}

func (c *TimelinesbuzUsecase) TimelinesItemRecord(caches *lib.Cache[*TData], timelineData TData, userFacade *UserFacade) (timelinesRecordVo TimelinesRecordVo, fieldHistory []TimelinesFieldHistoryItemVo, err error) {

	timelineFields := timelineData.CustomFields
	timelineKind := timelineFields.TextValueByNameBasic(Timeline_FieldName_kind)

	kindEntity, err := c.KindUsecase.GetByKind(timelineKind)
	if err != nil {
		return timelinesRecordVo, fieldHistory, err
	}
	if kindEntity == nil {
		return timelinesRecordVo, fieldHistory, errors.New("kindEntity is nil")
	}
	if timelineKind == Kind_notes {
		timelinesRecordVo.ModuleName = KindConvertToModule(Kind_notes)
		timelinesRecordVo.ModuleLabel = kindEntity.Label
		timelinesRecordVo.Gid = timelineFields.TextValueByNameBasic(Timeline_FieldName_kind_gid)
		note := lib.StringToTDef[TimelineForNotes](timelineFields.TextValueByNameBasic(Timeline_FieldName_notes), TimelineForNotes{})
		timelinesRecordVo.Name = note.Content
	} else {

		timelinesRecordVo.ModuleName = KindConvertToModule(timelineKind)
		timelinesRecordVo.ModuleLabel = kindEntity.Label
		timelinesRecordVo.Gid = timelineFields.TextValueByNameBasic(Timeline_FieldName_kind_gid)

		recordKindFields, err := c.FieldUsecase.CacheStructByKind(timelineKind)
		if err != nil {
			c.log.Error(err)
		} else {
			var fieldPermissionCenter FieldPermissionCenter
			if userFacade != nil {
				fieldPermissionCenter, err = c.FieldPermissionUsecase.CacheFieldPermissionCenter(timelineKind, userFacade.ProfileGid())
				if err != nil {
					c.log.Error(err)
				}
				//fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(FieldName_amount)
				//if err != nil {
				//	return nil, err
				//}
				//if fieldPermissionVo.CanShow() {
				//
				//}
			}

			_, err := c.FieldOptionUsecase.CacheStructByKind(timelineKind)
			// todo:lgl Option颜色处理
			if err != nil {
				c.log.Error(err)
			} else {
				note := lib.StringToTDef[TimelineFieldHistoryNotes](timelineFields.TextValueByNameBasic(Timeline_FieldName_notes), TimelineFieldHistoryNotes{})

				for _, v := range note.FieldHistory {
					var timelinesFieldHistoryItemVo TimelinesFieldHistoryItemVo
					if userFacade != nil {
						fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v.FieldName)
						if err != nil {
							//c.log.Error(err)
							continue
						}
						if !fieldPermissionVo.CanShow() {
							continue
						}
					}

					fieldEntity := recordKindFields.GetByFieldName(v.FieldName)
					if fieldEntity != nil {
						timelinesFieldHistoryItemVo.FieldType = fieldEntity.FieldType
						timelinesFieldHistoryItemVo.FieldName = fieldEntity.FieldName
						timelinesFieldHistoryItemVo.FieldLabel = fieldEntity.FieldLabel

						if fieldEntity.FieldType == FieldType_multilookup {

							timelinesFieldHistoryItemVo.OldMultiValues = c.TUsecase.GenTFieldMultiValues(caches, *fieldEntity, v.OldValue)
							timelinesFieldHistoryItemVo.NewMultiValues = c.TUsecase.GenTFieldMultiValues(caches, *fieldEntity, v.NewValue)

						} else {
							if v.OldValue != "" {
								oldValueTField := c.TUsecase.GenTField(caches, timelineKind, *fieldEntity, v.OldValue)
								timelinesFieldHistoryItemVo.OldValue = &TimelinesFieldHistoryItemValue{
									Label: InterfaceToString(oldValueTField.DisplayValue),
									Value: v.OldValue,
								}
							}
							if v.NewValue != "" {
								newValueTField := c.TUsecase.GenTField(caches, timelineKind, *fieldEntity, v.NewValue)
								timelinesFieldHistoryItemVo.NewValue = &TimelinesFieldHistoryItemValue{
									Label: InterfaceToString(newValueTField.DisplayValue),
									Value: v.NewValue,
								}
							}

							//timelinesFieldHistoryItemVo.Value = TimelinesFieldHistoryItemValueVo{
							//	NewValue: v.NewValue,
							//	OldValue: v.OldValue,
							//}
							//newValueTField := c.TUsecase.GenTField(caches, timelineKind, *fieldEntity, v.NewValue)
							//newDisplayValue := InterfaceToString(newValueTField.DisplayValue)
							//
							//oldValueTField := c.TUsecase.GenTField(caches, timelineKind, *fieldEntity, v.OldValue)
							//oldDisplayValue := InterfaceToString(oldValueTField.DisplayValue)
							//
							//timelinesFieldHistoryItemVo.DisplayValue = TimelinesFieldHistoryItemDisplayValueVo{
							//	NewValue: newDisplayValue,
							//	OldValue: oldDisplayValue,
							//}
						}
						fieldHistory = append(fieldHistory, timelinesFieldHistoryItemVo)
					}
				}
			}
		}
	}
	return timelinesRecordVo, fieldHistory, nil
}
