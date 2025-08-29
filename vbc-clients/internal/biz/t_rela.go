package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
)

type TRelaUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
}

func NewTRelaUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *TRelaUsecase {
	uc := &TRelaUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

type TRelaMap map[string]TRelaVo

func (c TRelaMap) Set(kind string, relaName string, gids []string) {
	if _, ok := c[kind]; !ok {
		c[kind] = TRelaVo{
			RelaName: relaName,
			Gids:     make(map[string]bool),
		}
	}
	for _, v := range gids {
		if v == "" {
			continue
		}
		c[kind].Gids[v] = true
	}
}

//func (c TRelaMap) GetGids(kind string) (gids []string) {
//	if _, ok := c[kind]; ok {
//		for k, _ := range c[kind].Gids {
//			gids = append(gids, k)
//		}
//	}
//	return
//}

type TRelaVo struct {
	RelaName string
	Gids     map[string]bool
}

func (c *TRelaVo) GetGids() (gids []string) {
	for k, _ := range c.Gids {
		gids = append(gids, k)
	}
	return
}

func TRelaMapInit() TRelaMap {
	return make(TRelaMap)
}

func (c *TRelaUsecase) GetRelaMap(fields TypeFieldList, records []map[string]interface{}) (tRelaMap TRelaMap) {

	tRelaMap = TRelaMapInit()
	for k, _ := range records {
		for _, v1 := range fields {
			if v1.FieldType == FieldType_lookup {
				val := lib.InterfaceToString(records[k][v1.FieldName])
				if val != "" {
					tRelaMap.Set(v1.RelaKind, v1.RelaName, []string{val})
				}
			} else if v1.FieldType == FieldType_multilookup {
				val := lib.InterfaceToString(records[k][v1.FieldName])
				if val != "" {
					vals := strings.Split(val, ",")
					for _, gid := range vals {
						if gid != "" {
							tRelaMap.Set(v1.RelaKind, v1.RelaName, []string{gid})
						}
					}
				}
			}
		}
	}

	return tRelaMap
}
