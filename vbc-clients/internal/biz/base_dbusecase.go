package biz

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	. "vbc/lib/builder"
)

type DBUsecaseInterface interface {
	TableName() string
}

type DBUsecase[T DBUsecaseInterface] struct {
	DB *gorm.DB
}

func (c *DBUsecase[T]) GetByCond(cond Cond) (*T, error) {
	condSql, err := ToBoundSQL(cond)
	if err != nil {
		return nil, err
	}
	var entity T
	err = c.DB.Where(condSql).Take(&entity).Error
	if err == nil {
		return &entity, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

//
//func (c *DBUsecase[T]) GetByRawSql(rawSql string) (*T, error) {
//	var record T
//	err := c.DB.Raw(rawSql).Scan(&record).Error
//	return &record, err
//}

// GetByCondWithOrderBy order: id desc
func (c *DBUsecase[T]) GetByCondWithOrderBy(cond Cond, order interface{}) (*T, error) {
	condSql, err := ToBoundSQL(cond)
	if err != nil {
		return nil, err
	}
	var entity T
	err = c.DB.Where(condSql).Order(order).Take(&entity).Error
	if err == nil {
		return &entity, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (c *DBUsecase[T]) AllByCond(cond Cond) ([]*T, error) {
	var records []*T
	query := c.DB
	if cond != nil {
		condSql, err := ToBoundSQL(cond)
		if err != nil {
			return nil, err
		}
		query = query.Where(condSql)
	}
	err := query.Find(&records).Error
	if err == nil {
		return records, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (c *DBUsecase[T]) AllByCondWithOrderBy(cond Cond, order interface{}, limit int) ([]*T, error) {
	var records []*T
	condSql, err := ToBoundSQL(cond)
	if err != nil {
		return nil, err
	}
	query := c.DB.Where(condSql).Order(order)
	if limit > 0 {
		query.Limit(limit)
	}
	err = query.Find(&records).Error
	if err == nil {
		return records, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (c *DBUsecase[T]) AllByCondWithOrderBySelect(selects []string, cond Cond, order interface{}, limit int) ([]*T, error) {
	var records []*T
	condSql, err := ToBoundSQL(cond)
	if err != nil {
		return nil, err
	}
	query := c.DB.Select(selects).Where(condSql).Order(order)
	if limit > 0 {
		query.Limit(limit)
	}
	err = query.Find(&records).Error
	if err == nil {
		return records, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (c *DBUsecase[T]) AllByRawSql(rawSql string) ([]*T, error) {
	var records []*T
	err := c.DB.Raw(rawSql).Scan(&records).Error
	return records, err
}

func (c *DBUsecase[T]) Total(cond Cond) (total int64, err error) {
	var tt T
	query := c.DB.Model(tt)
	if cond != nil {
		condSql, err := ToBoundSQL(cond)
		if err != nil {
			return 0, err
		}
		query = query.Where(condSql)
	}
	err = query.Count(&total).Error
	return
}
func (c *DBUsecase[T]) ListByCondWithPaging(cond Cond, order interface{}, page int, pageSize int) ([]*T, error) {

	var records []*T
	query := c.DB
	if cond != nil {
		condSql, err := ToBoundSQL(cond)
		if err != nil {
			return nil, err
		}
		query = query.Where(condSql)
	}
	if order != "" {
		query = query.Order(order)
	}

	query = query.Offset(HandleOffset(page, pageSize)).Limit(pageSize)
	err := query.Find(&records).Error
	if err == nil {
		return records, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (c *DBUsecase[T]) UpdatesByCond(values map[string]interface{}, cond Cond) error {
	if values == nil {
		return errors.New("UpdateByCond:data is nil")
	}
	if cond == nil {
		return errors.New("UpdateByCond:cond is nil")
	}
	var entity T
	condSql, err := ToBoundSQL(cond)
	if err != nil {
		return err
	}
	err = c.DB.Table(entity.TableName()).Where(condSql).Updates(values).Error
	return err
}

type Facade interface {
	GetTData() TData // 使用方法代替直接嵌入
	SetTData(tData TData)
	New() Facade // 定义一个构造函数
}

//
//func GetFacadesByGids[T Facade](TUsecase *TUsecase, kind string, gids []string) (map[string]T, error) {
//
//	facades := make(map[string]T)
//	records, err := TUsecase.ListByCond(kind, In(DataEntry_gid, gids))
//	if err != nil {
//		return facades, err
//	}
//	for k, v := range records {
//		obj := T.New()
//		//obj :=  // 这里会导致错误，因为T不是指针
//		obj.SetTData(*records[k])
//		//obj.SetTData(*records[k])
//		facades[v.Gid()] = obj
//	}
//	return facades, nil
//}
