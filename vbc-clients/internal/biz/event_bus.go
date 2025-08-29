package biz

import "sync"

const (
	EventBus_AfterInsertData   = "AfterInsertData"
	EventBus_AfterHandleUpdate = "AfterHandleUpdate"
)

type EventBus struct {
	subscribers map[string][]func(kindEntity KindEntity,
		structField *TypeFieldStruct,
		recognizeFieldName string,
		dataEntryOperResult DataEntryOperResult,
		sourceData TypeDataEntryList,
		modifiedBy string)
	lock sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]func(kindEntity KindEntity,
			structField *TypeFieldStruct,
			recognizeFieldName string,
			dataEntryOperResult DataEntryOperResult,
			sourceData TypeDataEntryList,
			modifiedBy string)),
	}
}

// Subscribe 添加事件监听器
func (e *EventBus) Subscribe(event string, handler func(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList,
	modifiedBy string)) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.subscribers[event] = append(e.subscribers[event], handler)
}

// Publish 发布事件
func (e *EventBus) Publish(event string, kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList, modifiedBy string) {
	//event := EventBus_AfterHandleUpdate
	e.lock.RLock()
	defer e.lock.RUnlock()
	if handlers, ok := e.subscribers[event]; ok {
		for _, handler := range handlers {
			handler(kindEntity, structField, recognizeFieldName, dataEntryOperResult, sourceData, modifiedBy) // 异步调用
		}
	}
}
