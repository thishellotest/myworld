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

type ZohoTaskScanJobUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	MapUsecase        *MapUsecase
	ZohoUsecase       *ZohoUsecase
	ZohobuzUsecase    *ZohobuzUsecase
	UsageStatsUsecase *UsageStatsUsecase
}

func NewZohoTaskScanJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	ZohoUsecase *ZohoUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
	UsageStatsUsecase *UsageStatsUsecase) *ZohoTaskScanJobUsecase {
	uc := &ZohoTaskScanJobUsecase{
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

func (c *ZohoTaskScanJobUsecase) RunJob(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ZohoTaskScanJobUsecase:RunJob:Done")
				return
			default:
				err := c.BizRunJob()
				if err != nil {
					c.log.Error(err)
					time.Sleep(60 * time.Second) // 报错了，延时大些
				}
				//fmt.Println("ZohoTaskScanJobUsecase: RunJob")
				// 测试环境加快速度， 后续改为 15秒
				time.Sleep(time.Duration(configs.ZohoContactAndDealSyncSlowTimes) * 10 * time.Second)
			}
		}
	}()
	return nil
}

func (c *ZohoTaskScanJobUsecase) BizRunJob() error {
	modifiedTime, err := c.MapUsecase.GetForString(Map_ZohoTaskHandleLastModifiedTime)
	if err != nil {
		return err
	}
	var lastModifiedTime string
	err = c.BatchHandle(&lastModifiedTime, modifiedTime, 1)
	if err != nil {
		return err
	}
	formatLastModifiedTime, err := lib.TimeParse(lastModifiedTime)
	if err != nil {
		return err
	}
	formatLastModifiedTime = formatLastModifiedTime.Add(-1 * time.Second) // 防同一时刻数据不正确
	nModifiedTime := formatLastModifiedTime.Format(time.RFC3339Nano)
	if modifiedTime == nModifiedTime {
		return nil
	}
	return c.MapUsecase.Set(Map_ZohoTaskHandleLastModifiedTime, nModifiedTime)
}

func (c *ZohoTaskScanJobUsecase) BatchHandle(lastModifiedTime *string, modifiedTime string, page int) error {

	perPage := 100
	fields := config_zoho.TaskLayout().TaskApiNames()
	time.Sleep(time.Second)
	params := make(url.Values)
	params.Add("page", InterfaceToString(page))
	params.Add("per_page", InterfaceToString(perPage))
	//lib.DPrintln(params)
	c.UsageStatsUsecase.Stat(UsageType_GetTaskRecords, time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Tasks, fields, params)
	if err != nil {
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
	err = c.ZohobuzUsecase.SyncTasks(listMaps)
	if err != nil {
		return err
	}

	//lib.DPrintln("ZohoTaskScanJobUsecase", len(listMaps), "===", perPage)
	if len(listMaps) == perPage { // 有下一页
		//lib.DPrintln("ZohoTaskScanJobUsecase", "----2")
		nModifiedTime := listMaps[perPage-1].GetString("Modified_Time")
		if modifiedTime == "" {
			//lib.DPrintln("ZohoTaskScanJobUsecase", "----1")
			return c.BatchHandle(lastModifiedTime, modifiedTime, page+1)
		} else {
			divideModifiedTime, _ := lib.TimeParse(modifiedTime)
			newModifiedTime, _ := lib.TimeParse(nModifiedTime)
			//lib.DPrintln("ZohoTaskScanJobUsecase", "newModifiedTime:", nModifiedTime, "||", "divideModifiedTime:", modifiedTime, "page:", page)
			//lib.DPrintln("ZohoTaskScanJobUsecase", "newModifiedTime.Unix():", newModifiedTime.Unix(), "divideModifiedTime.Unix():", divideModifiedTime.Unix())
			if newModifiedTime.Unix() > divideModifiedTime.Unix() {
				return c.BatchHandle(lastModifiedTime, modifiedTime, page+1)
				//return nil
			}
		}
	}
	return nil
}

func (c *ZohoTaskScanJobUsecase) Handle() {

}
