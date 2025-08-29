package biz

//
//import (
//	"github.com/go-kratos/kratos/v2/log"
//	"vbc/internal/conf"
//	"vbc/internal/vbc_config"
//	"vbc/lib"
//)
//
//type EventClientSyncCaseUsecase struct {
//	log                      *log.Helper
//	conf                     *conf.Data
//	EventBus                 *EventBus
//	ClientCaseSyncbuzUsecase *ClientCaseSyncbuzUsecase
//	TUsecase                 *TUsecase
//}
//
//func NewEventClientSyncCaseUsecase(logger log.Logger,
//	conf *conf.Data,
//	EventBus *EventBus,
//	ClientCaseSyncbuzUsecase *ClientCaseSyncbuzUsecase,
//	TUsecase *TUsecase,
//) *EventClientSyncCaseUsecase {
//	uc := &EventClientSyncCaseUsecase{
//		log:                      log.NewHelper(logger),
//		conf:                     conf,
//		EventBus:                 EventBus,
//		ClientCaseSyncbuzUsecase: ClientCaseSyncbuzUsecase,
//		TUsecase:                 TUsecase,
//	}
//	// 有顺序问题，移入到queue处理
//	uc.EventBus.Subscribe(EventBus_AfterHandleUpdate, uc.HandleEventClientSyncCase)
//	uc.EventBus.Subscribe(EventBus_AfterInsertData, uc.HandleEventClientSyncCaseByInsertData)
//	return uc
//}
//
//func (c *EventClientSyncCaseUsecase) HandleEventClientSyncCaseByInsertData(kindEntity KindEntity,
//	structField *TypeFieldStruct,
//	recognizeFieldName string,
//	dataEntryOperResult DataEntryOperResult,
//	sourceData TypeDataEntryList,
//	modifiedBy string) {
//	if lib.IsProd() { // todo:lgl线上环境不能开放，因为使用了zoho，后续zoho下架，就需要开放
//		return
//	}
//	//lib.DPrintln("dataEntryOperResult:", dataEntryOperResult, modifiedBy)
//	//lib.DPrintln("sourceData:", sourceData, modifiedBy)
//	//if kindEntity.Kind == Kind_client_cases && recognizeFieldName == DataEntry_gid {
//	//	for gid, v := range dataEntryOperResult {
//	//		if v.IsNewRecord {
//	//			//c.UpdateByGid(gid, modifiedBy)
//	//		}
//	//	}
//	//}
//}
//
//func (c *EventClientSyncCaseUsecase) HandleEventClientSyncCase(kindEntity KindEntity,
//	structField *TypeFieldStruct,
//	recognizeFieldName string,
//	dataEntryOperResult DataEntryOperResult,
//	sourceData TypeDataEntryList,
//	modifiedBy string) {
//
//	if lib.IsProd() { // todo:lgl线上环境不能开放，因为使用了zoho，后续zoho下架，就需要开放
//		return
//	}
//
//	if kindEntity.Kind == Kind_client_cases && recognizeFieldName == DataEntry_gid {
//		operUser, _ := c.TUsecase.DataByGid(Kind_users, modifiedBy)
//		fieldNames := vbc_config.SyncFieldNamesForCase()
//		for gid, v := range dataEntryOperResult {
//			var clientCaseSyncList ClientCaseSyncList
//			if v.IsUpdated {
//				for fieldName, v1 := range v.DataEntryModifyDataMap {
//					if lib.InArray(fieldName, fieldNames) { // 此处须在FieldName_stages前面操作
//						newVal := InterfaceToString(v1.NewVal)
//						clientCaseSyncList = append(clientCaseSyncList, ClientCaseSyncVo{
//							FieldName:  fieldName,
//							FieldValue: newVal,
//						})
//					}
//				}
//			}
//			err := c.ClientCaseSyncbuzUsecase.CaseToClient(gid, clientCaseSyncList, operUser)
//			if err != nil {
//				c.log.Error(err)
//			}
//		}
//	}
//}
