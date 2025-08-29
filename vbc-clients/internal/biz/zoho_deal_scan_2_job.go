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

type ZohoDealScan2JobUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	MapUsecase        *MapUsecase
	ZohoUsecase       *ZohoUsecase
	ZohobuzUsecase    *ZohobuzUsecase
	UsageStatsUsecase *UsageStatsUsecase
}

func NewZohoDealScan2JobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	ZohoUsecase *ZohoUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
	UsageStatsUsecase *UsageStatsUsecase) *ZohoDealScan2JobUsecase {
	uc := &ZohoDealScan2JobUsecase{
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

func (c *ZohoDealScan2JobUsecase) RunJob(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ZohoDealScan2JobUsecase:RunJob:Done")
				return
			default:
				//fmt.Println("ZohoDealScan2JobUsecase:RunJob 1")
				err := c.BizRunJob()
				if err != nil {
					c.log.Info("ZohoDealScan2JobUsecase:", err)
					time.Sleep(60 * time.Second) // 报错了，延时大些
				}
				// 测试环境加快速度， 后续改为 15秒
				time.Sleep(time.Duration(configs.ZohoContactAndDealSyncSlowTimes) * 15 * time.Second)
				//fmt.Println("ZohoDealScan2JobUsecase:RunJob 2")
			}
		}
	}()
	return nil
}

func (c *ZohoDealScan2JobUsecase) BizRunJob() error {
	modifiedTime, err := c.MapUsecase.GetForString(Map_ZohoDealHandleLastModifiedTime2)
	if err != nil {
		return err
	}
	var lastModifiedTime string
	err = c.BatchHandle(&lastModifiedTime, modifiedTime, 1)
	if err != nil {
		c.log.Info("ZohoDealScan2JobUsecase:", err)
		return err
	}
	formatLastModifiedTime, err := lib.TimeParse(lastModifiedTime)
	if err != nil {
		c.log.Info("ZohoDealScan2JobUsecase:", err)
		return err
	}
	formatLastModifiedTime = formatLastModifiedTime.Add(-1 * time.Second) // 防同一时刻数据不正确
	nModifiedTime := formatLastModifiedTime.Format(time.RFC3339Nano)
	if modifiedTime == nModifiedTime {
		return nil
	}
	return c.MapUsecase.Set(Map_ZohoDealHandleLastModifiedTime2, nModifiedTime)
}

func (c *ZohoDealScan2JobUsecase) BatchHandle(lastModifiedTime *string, modifiedTime string, page int) error {

	perPage := 100
	fields := config_zoho.DealLayout().DealApiNames2()
	time.Sleep(time.Second)
	params := make(url.Values)
	params.Add("page", InterfaceToString(page))
	params.Add("per_page", InterfaceToString(perPage))
	//params.Add("ids", "6159272000009972111")
	//lib.DPrintln(params)
	c.UsageStatsUsecase.Stat(UsageType_GetDealRecords, time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	if err != nil {
		c.log.Info("ZohoDealScan2JobUsecase:", err)
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
	err = c.ZohobuzUsecase.SyncClientCases(listMaps, "job2")
	if err != nil {
		c.log.Info("ZohoDealScan2JobUsecase:", err)
		return err
	}

	//lib.DPrintln("ZohoDealScan2JobUsecase", len(listMaps), "===", perPage)
	if len(listMaps) == perPage { // 有下一页
		//lib.DPrintln("ZohoDealScan2JobUsecase", "----2")
		nModifiedTime := listMaps[perPage-1].GetString("Modified_Time")
		if modifiedTime == "" {
			//lib.DPrintln("ZohoDealScan2JobUsecase", "----1")
			return c.BatchHandle(lastModifiedTime, modifiedTime, page+1)
		} else {
			divideModifiedTime, _ := lib.TimeParse(modifiedTime)
			newModifiedTime, _ := lib.TimeParse(nModifiedTime)
			//lib.DPrintln("ZohoDealScan2JobUsecase", "newModifiedTime:", nModifiedTime, "||", "divideModifiedTime:", modifiedTime, "page:", page)
			//lib.DPrintln("ZohoDealScan2JobUsecase", "newModifiedTime.Unix():", newModifiedTime.Unix(), "divideModifiedTime.Unix():", divideModifiedTime.Unix())
			if newModifiedTime.Unix() > divideModifiedTime.Unix() {
				return c.BatchHandle(lastModifiedTime, modifiedTime, page+1)
				//return nil
			}
		}
	}
	return nil
}

func (c *ZohoDealScan2JobUsecase) Handle() {

}
