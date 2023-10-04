package base

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

// UpdateXdxrInfo 除权除息数据
func UpdateXdxrInfo(securityCode string) {
	securityCode = proto.CorrectSecurityCode(securityCode)
	tdxApi := gotdx.GetTdxApi()
	xdxrInfos, err := tdxApi.GetXdxrInfo(securityCode)
	if err != nil {
		logger.Errorf("获取除权除息数据失败", err)
		return
	}
	if len(xdxrInfos) > 0 {
		filename := cache.XdxrFilename(securityCode)
		_ = api.SlicesToCsv(filename, xdxrInfos)
	}
}

// GetCacheXdxrList 获取除权除息的数据列表
func GetCacheXdxrList(securityCode string) []quotes.XdxrInfo {
	securityCode = proto.CorrectSecurityCode(securityCode)
	filename := cache.XdxrFilename(securityCode)
	list := []quotes.XdxrInfo{}
	_ = api.CsvToSlices(filename, &list)
	return list
}
