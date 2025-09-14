package base

import (
	"github.com/quant1x/engine/cache"
	"github.com/quant1x/exchange"
	"github.com/quant1x/gotdx"
	"github.com/quant1x/gotdx/quotes"
	"github.com/quant1x/x/api"
	"github.com/quant1x/x/logger"
)

// UpdateXdxrInfo 除权除息数据
func UpdateXdxrInfo(securityCode string) {
	securityCode = exchange.CorrectSecurityCode(securityCode)
	tdxApi := gotdx.GetTdxApi()
	xdxrInfos, err := tdxApi.GetXdxrInfo(securityCode)
	if err != nil {
		logger.Errorf("获取除权除息数据失败", err)
		return
	}
	//slices.SortFunc(xdxrInfos, func(a, b quotes.XdxrInfo) int {
	//	if a.Date == b.Date {
	//		return 0
	//	} else if a.Date > b.Date {
	//		return 1
	//	} else {
	//		return -1
	//	}
	//})
	if len(xdxrInfos) > 0 {
		filename := cache.XdxrFilename(securityCode)
		_ = api.SlicesToCsv(filename, xdxrInfos)
	}
}

// GetCacheXdxrList 获取除权除息的数据列表
func GetCacheXdxrList(securityCode string) []quotes.XdxrInfo {
	securityCode = exchange.CorrectSecurityCode(securityCode)
	filename := cache.XdxrFilename(securityCode)
	var list []quotes.XdxrInfo
	_ = api.CsvToSlices(filename, &list)
	return list
}
