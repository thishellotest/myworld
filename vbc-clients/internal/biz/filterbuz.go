package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type FilterbuzUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	FilterUsecase *FilterUsecase
}

func NewFilterbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	FilterUsecase *FilterUsecase,
) *FilterbuzUsecase {
	uc := &FilterbuzUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		FilterUsecase: FilterUsecase,
	}

	return uc
}

func (c *FilterbuzUsecase) FilterDelete(userGid string, filterIds []int32) error {
	return c.FilterUsecase.UpdatesByCond(map[string]interface{}{"deleted_at": time.Now().Unix()},
		And(Eq{"user_gid": userGid},
			In("id", filterIds), Eq{"deleted_at": 0}))
}

func (c *FilterbuzUsecase) FilterList(userGid string, kind string, tableType string) ([]FilterVo, error) {
	records, err := c.FilterUsecase.AllByCondWithOrderBy(And(Eq{"kind": kind, "user_gid": userGid, "biz_deleted_at": 0, "deleted_at": 0, "table_type": tableType}), "id desc", 1000)
	if err != nil {
		return nil, err
	}
	var filters []FilterVo
	for _, v := range records {
		filters = append(filters, v.ToFilterVo())
	}
	return filters, nil

}

func (c *FilterbuzUsecase) BizFilterList(userGid string, kind string, tableType string) (lib.TypeMap, error) {

	return nil, nil
}

func (c *FilterbuzUsecase) BizFilterSave(userGid string, kind string, content string, filterName string) (lib.TypeMap, error) {

	entity := &FilterEntity{
		FilterName: filterName,
		Kind:       kind,
		UserGid:    userGid,
		Content:    content,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	err := c.FilterUsecase.DB.Save(&entity).Error
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	data.Set("filter", entity.ToFilterVo())
	return data, nil

}
