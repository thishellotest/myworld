package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type RecordbuzSearchUsecase struct {
	log                         *log.Helper
	conf                        *conf.Data
	CommonUsecase               *CommonUsecase
	TFilterUsecase              *TFilterUsecase
	KindUsecase                 *KindUsecase
	BUsecase                    *BUsecase
	PermissionDataFilterUsecase *PermissionDataFilterUsecase
	TUsecase                    *TUsecase
}

func NewRecordbuzSearchUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TFilterUsecase *TFilterUsecase,
	KindUsecase *KindUsecase,
	BUsecase *BUsecase,
	PermissionDataFilterUsecase *PermissionDataFilterUsecase,
	TUsecase *TUsecase,
) *RecordbuzSearchUsecase {
	uc := &RecordbuzSearchUsecase{
		log:                         log.NewHelper(logger),
		CommonUsecase:               CommonUsecase,
		conf:                        conf,
		TFilterUsecase:              TFilterUsecase,
		KindUsecase:                 KindUsecase,
		BUsecase:                    BUsecase,
		PermissionDataFilterUsecase: PermissionDataFilterUsecase,
		TUsecase:                    TUsecase,
	}
	return uc
}

func (c *RecordbuzSearchUsecase) NewRecordbuzSearchCls(NormalUserOnlyOwner bool) *RecordbuzSearchCls {
	return CreateRecordbuzSearchCls(c.TFilterUsecase,
		c.KindUsecase,
		c.BUsecase,
		c.PermissionDataFilterUsecase, c.TUsecase, NormalUserOnlyOwner)
}
