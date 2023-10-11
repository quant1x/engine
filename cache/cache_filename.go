package cache

import (
	"fmt"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
)

// XdxrFilename XDXR缓存路径
func XdxrFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	xdxrPath := fmt.Sprintf("%s/%s/%s.csv", GetXdxrPath(), cacheId[:length-3], cacheId)
	return xdxrPath
}

// KLineFilename 基础K线缓存路径
func KLineFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	filepath := fmt.Sprintf("%s/%s/%s.csv", GetDayPath(), cacheId[:length-3], cacheId)
	return filepath
}

// Top10HoldersFilename 前十大流通股股东缓存文件名
func Top10HoldersFilename(code, date string) string {
	idPath := CacheIdPath(code)
	q, _, _ := api.GetQuarterByDate(date)
	filename := fmt.Sprintf("%s/%s/%s.csv", GetHoldingPath(), q, idPath)
	return filename
}

func quarterlyCachePath(date string) string {
	q, _, _ := api.GetQuarterByDate(date)
	path := fmt.Sprintf("%s/%s", GetQuarterlyPath(), q)
	return path
}

// QuarterlyReportFilename 季报存储路径
//
//	info
//	  |-- YYYYQ1
//	        |--  sh600105.report
//	  |-- YYYYQ2
//	  |-- YYYYQ3
//	  |-- YYYYQ4
//	Deprecated: 不推荐使用
func QuarterlyReportFilename(code, date string) string {
	idPath := CacheIdPath(code)
	path := quarterlyCachePath(date)
	filename := fmt.Sprintf("%s/%s.report", path, idPath)
	return filename
}

// ReportsFilename 报告数据文件名
func ReportsFilename(date string) string {
	keyword := "reports"
	path := quarterlyCachePath(date)
	filename := fmt.Sprintf("%s/%s.csv", path, keyword)
	return filename
}

// TickFilename tick文件比较多, 目录结构${tick}/${YYYY}/${YYYYMMDD}/${CacheIdPath}
func TickFilename(code, date string) string {
	date = trading.FixTradeDate(date, TDX_FORMAT_PROTOCOL_DATE)
	cacheId := CacheId(code)
	tickPath := fmt.Sprintf("%s/%s/%s/%s.csv", GetTickPath(), date[0:4], date, cacheId)
	return tickPath
}
