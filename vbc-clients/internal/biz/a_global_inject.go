package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

/*
解决一些像Event Bus，业务没有用的实例无法注入的问题
*/
type GlobalInjectUsecase struct {
	log *log.Helper
}

func NewGlobalInjectUsecase(logger log.Logger,
	GlobalEventBusBuzUsecase *GlobalEventBusBuzUsecase,
	ConditionCategoryUsecase *ConditionCategoryUsecase,
	EventBusReferrerUsecase *EventBusReferrerUsecase,
) *GlobalInjectUsecase {
	uc := &GlobalInjectUsecase{
		log: log.NewHelper(logger),
	}

	return uc
}
