package biz

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
	"vbc/lib"
)

/*

CREATE TABLE `long_maps` (
  `mkey` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `mval` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` int NOT NULL DEFAULT '0',
  `updated_at` int NOT NULL DEFAULT '0',
  UNIQUE KEY `uniq_k` (`mkey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='long_maps';

*/

const (
	LongMapKey_sck1 = "sck1"
	LongMapKey_sck2 = "sck2"
)

type LongMapEntity struct {
	Mkey      string
	Mval      string
	CreatedAt int64
	UpdatedAt int64
}

func (LongMapEntity) TableName() string {
	return "long_maps"
}

type LongMapUsecase struct {
	CommonUsecase *CommonUsecase
	DBUsecase[LongMapEntity]
}

func NewLongMapUsecase(CommonUsecase *CommonUsecase) *LongMapUsecase {

	uc := &LongMapUsecase{
		CommonUsecase: CommonUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *LongMapUsecase) SetInt(mkey string, mval int) error {
	return c.Set(mkey, InterfaceToString(mval))
}

func (c *LongMapUsecase) Set(mkey string, mval string) error {
	currentTime := time.Now()
	mval = lib.SqlBindValue(mval)
	sql := fmt.Sprintf("INSERT INTO %s (mkey, mval ,created_at, updated_at) VALUES(\"%s\", %s, %d, %d ) ON DUPLICATE KEY UPDATE updated_at=values(updated_at),mval=values(mval)",
		LongMapEntity{}.TableName(), mkey, mval, currentTime.Unix(), currentTime.Unix())
	return c.CommonUsecase.DB().Exec(sql).Error
}

func (c *LongMapUsecase) GetForString(mkey string) (string, error) {
	var entity LongMapEntity
	err := c.CommonUsecase.DB().Where("mkey=?", mkey).
		Take(&entity).
		Error
	if err == nil {
		return entity.Mval, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	return "", err
}

func (c *LongMapUsecase) GetForInt(mkey string) (int32, error) {
	val, err := c.GetForString(mkey)
	if err != nil {
		return 0, err
	}
	i, _ := strconv.ParseInt(val, 10, 32)
	return int32(i), nil
}
