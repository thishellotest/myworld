package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"net/url"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

type ZohoNoteScanJobUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	MapUsecase        *MapUsecase
	ZohoUsecase       *ZohoUsecase
	ZohobuzUsecase    *ZohobuzUsecase
	UsageStatsUsecase *UsageStatsUsecase
	InvokeLogUsecase  *InvokeLogUsecase
}

func NewZohoNoteScanJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	ZohoUsecase *ZohoUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
	UsageStatsUsecase *UsageStatsUsecase,
	InvokeLogUsecase *InvokeLogUsecase) *ZohoNoteScanJobUsecase {
	uc := &ZohoNoteScanJobUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		MapUsecase:        MapUsecase,
		ZohoUsecase:       ZohoUsecase,
		ZohobuzUsecase:    ZohobuzUsecase,
		UsageStatsUsecase: UsageStatsUsecase,
		InvokeLogUsecase:  InvokeLogUsecase,
	}
	return uc
}

func (c *ZohoNoteScanJobUsecase) RunJob(ctx context.Context) error {
	c.log.Info("ZohoNoteScanJobUsecase:RunJob start")
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.log.Info("ZohoNoteScanJobUsecase:RunJob:Done")
				return
			default:
				err := c.BizRunJob()
				if err != nil {
					c.log.Error(err)
					time.Sleep(60 * time.Second) // 报错了，延时大些
				}
				// 测试环境加快速度， 后续改为 15秒
				time.Sleep(time.Duration(configs.ZohoContactAndDealSyncSlowTimes) * 15 * time.Second)
			}
		}
	}()
	return nil
}

func (c *ZohoNoteScanJobUsecase) BizRunJob() error {
	modifiedTime, err := c.MapUsecase.GetForString(Map_ZohoNoteHandleLastModifiedTime)
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
	return c.MapUsecase.Set(Map_ZohoNoteHandleLastModifiedTime, nModifiedTime)
}

func (c *ZohoNoteScanJobUsecase) BatchHandle(lastModifiedTime *string, modifiedTime string, page int) error {

	c.log.Info("ZohoNoteScanJobUsecase:BatchHandle:", page)
	perPage := 200
	fields := config_zoho.NotesLayout().NoteApiNames()
	time.Sleep(time.Second)
	params := make(url.Values)
	params.Add("page", InterfaceToString(page))
	params.Add("per_page", InterfaceToString(perPage))
	c.log.Info("params:", params)
	c.UsageStatsUsecase.Stat(UsageType_GetNoteRecords, time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Notes, fields, params)
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
	for k, _ := range listMaps {
		err = c.SyncRow(listMaps[k])
		if err != nil {
			return err
		}
	}

	//lib.DPrintln(listMaps)

	//lib.DPrintln("ZohoNoteScanJobUsecase", len(listMaps), "===", perPage)
	//if len(listMaps) == perPage { // 有下一页
	//	//lib.DPrintln("ZohoNoteScanJobUsecase", "----2")
	//	nModifiedTime := listMaps[perPage-1].GetString("Modified_Time")
	//	if modifiedTime == "" {
	//		//lib.DPrintln("ZohoNoteScanJobUsecase", "----1")
	//		return c.BatchHandle(lastModifiedTime, modifiedTime, page+1)
	//	} else {
	//		divideModifiedTime, _ := lib.TimeParse(modifiedTime)
	//		newModifiedTime, _ := lib.TimeParse(nModifiedTime)
	//
	//		//lib.DPrintln("ZohoNoteScanJobUsecase", "newModifiedTime:", nModifiedTime, "||", "divideModifiedTime:", modifiedTime, "page:", page)
	//		//lib.DPrintln("ZohoNoteScanJobUsecase", "newModifiedTime.Unix():", newModifiedTime.Unix(), "divideModifiedTime.Unix():", divideModifiedTime.Unix())
	//		if newModifiedTime.Unix() > divideModifiedTime.Unix() {
	//			return c.BatchHandle(lastModifiedTime, modifiedTime, page+1)
	//		}
	//	}
	//}
	return nil
}

func (c *ZohoNoteScanJobUsecase) SyncRow(row lib.TypeMap) error {

	content := row
	mTime := row.GetString("Modified_Time")
	cTime := row.GetString("Created_Time")
	rId := row.GetString("id")
	_, err := c.InvokeLogUsecase.Upsert(rId, mTime, cTime, content)
	if err != nil {
		return err
	}
	return nil
}

func (c *ZohoNoteScanJobUsecase) SyncAll(pageToken string, times int) error {

	params := make(url.Values)
	if pageToken != "" {
		params.Add("page_token", InterfaceToString(pageToken))
	}
	c.log.Info("ZohoNoteScanJobUsecase SyncAll:", params, times)
	fields := config_zoho.NotesLayout().NoteApiNames()
	time.Sleep(time.Second)
	c.UsageStatsUsecase.Stat(UsageType_GetNoteRecords, time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Notes, fields, params)
	if err != nil {
		return err
	}
	data := records.GetTypeList("data")
	c.log.Info("more_records:", records.GetString("info.more_records"), records.GetString("info.next_page_token"))

	for k, _ := range data {
		err = c.SyncRow(data[k])
		if err != nil {
			return err
		}
	}
	if records.GetString("info.more_records") == "true" {
		if records.GetString("info.next_page_token") != "" {
			times += 1
			err = c.SyncAll(records.GetString("info.next_page_token"), times)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
