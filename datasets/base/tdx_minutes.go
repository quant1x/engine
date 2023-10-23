package base

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
)

// GetMinutes 获取分时数据
func GetMinutes(securityCode, date string) (list []quotes.HistoryMinuteTime) {
	tdxApi := gotdx.GetTdxApi()
	hs, err := tdxApi.GetHistoryMinuteTimeData(securityCode, toTdxProtocolDate(date))
	if err != nil || hs.Count == 0 {
		return
	}
	list = append(list, hs.List...)
	_ = hs
	return
}

// UpdateMinutes 更新指定日期的个股分时数据
func UpdateMinutes(securityCode, date string) {
	list := GetMinutes(securityCode, date)
	if len(list) > 0 {
		filename := cache.MinuteFilename(securityCode, date)
		_ = api.SlicesToCsv(filename, list)
	}
}
