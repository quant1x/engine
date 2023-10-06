package base

import (
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/quotes"
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
