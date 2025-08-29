package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"net/url"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

type ZohoDealScanJobUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	MapUsecase        *MapUsecase
	ZohoUsecase       *ZohoUsecase
	ZohobuzUsecase    *ZohobuzUsecase
	UsageStatsUsecase *UsageStatsUsecase
}

func NewZohoDealScanJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	ZohoUsecase *ZohoUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
	UsageStatsUsecase *UsageStatsUsecase) *ZohoDealScanJobUsecase {
	uc := &ZohoDealScanJobUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		MapUsecase:        MapUsecase,
		ZohoUsecase:       ZohoUsecase,
		ZohobuzUsecase:    ZohobuzUsecase,
		UsageStatsUsecase: UsageStatsUsecase,
	}
	return uc
}

func (c *ZohoDealScanJobUsecase) RunJob(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ZohoDealScanJobUsecase:RunJob:Done")
				return
			default:
				err := c.BizRunJob()
				if err != nil {
					c.log.Info("ZohoDealScanJobUsecase:", err)
					time.Sleep(60 * time.Second) // 报错了，延时大些
				}
				// 测试环境加快速度， 后续改为 15秒
				time.Sleep(time.Duration(configs.ZohoContactAndDealSyncSlowTimes) * 18 * time.Second)
			}
		}
	}()
	return nil
}

func (c *ZohoDealScanJobUsecase) BizRunJob() error {
	modifiedTime, err := c.MapUsecase.GetForString(Map_ZohoDealHandleLastModifiedTime)
	if err != nil {
		return err
	}
	var lastModifiedTime string
	err = c.BatchHandle(&lastModifiedTime, modifiedTime, 1)
	if err != nil {
		c.log.Info("ZohoDealScanJobUsecase:", err)
		return err
	}
	formatLastModifiedTime, err := lib.TimeParse(lastModifiedTime)
	if err != nil {
		c.log.Info("ZohoDealScanJobUsecase:", err)
		return err
	}
	formatLastModifiedTime = formatLastModifiedTime.Add(-1 * time.Second) // 防同一时刻数据不正确
	nModifiedTime := formatLastModifiedTime.Format(time.RFC3339Nano)
	if modifiedTime == nModifiedTime {
		return nil
	}
	return c.MapUsecase.Set(Map_ZohoDealHandleLastModifiedTime, nModifiedTime)
}

func (c *ZohoDealScanJobUsecase) BatchHandle(lastModifiedTime *string, modifiedTime string, page int) error {

	perPage := 100
	fields := config_zoho.DealLayout().DealApiNames()
	time.Sleep(time.Second)
	params := make(url.Values)
	params.Add("page", InterfaceToString(page))
	params.Add("per_page", InterfaceToString(perPage))
	c.UsageStatsUsecase.Stat(UsageType_GetDealRecords, time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	if err != nil {
		c.log.Info("ZohoDealScanJobUsecase:", err)
		return err
	}
	listMaps := records.GetTypeList("data")
	if len(listMaps) == 0 {
		return nil
	}
	if page == 1 && len(listMaps) > 0 {
		t := listMaps[0].GetString("Modified_Time")
		*lastModifiedTime = t
	}
	err = c.ZohobuzUsecase.SyncClientCases(listMaps, "")
	if err != nil {
		c.log.Info("ZohoDealScanJobUsecase:", err)
		return err
	}

	if len(listMaps) == perPage { // 有下一页
		nModifiedTime := listMaps[perPage-1].GetString("Modified_Time")
		if modifiedTime == "" {
			return c.BatchHandle(lastModifiedTime, modifiedTime, page+1)
		} else {
			divideModifiedTime, _ := lib.TimeParse(modifiedTime)
			newModifiedTime, _ := lib.TimeParse(nModifiedTime)
			if newModifiedTime.Unix() > divideModifiedTime.Unix() {
				return c.BatchHandle(lastModifiedTime, modifiedTime, page+1)
				//return nil
			}
		}
	}
	return nil
}

func (c *ZohoDealScanJobUsecase) Handle() {

}
