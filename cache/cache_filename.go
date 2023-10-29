package cache

import (
	"fmt"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
)

const (
	FilenameDate = "20060102" // 缓存文件相关的日期格式
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

// FeatureFilename 特征数据缓存路径
func FeatureFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	filepath := fmt.Sprintf("%s/%s/%s.csv", GetFeaturesPath(), cacheId[:length-3], cacheId)
	return filepath
}

func MinuteFilename(code, date string) string {
	date = trading.FixTradeDate(date, FilenameDate)
	cacheId := CacheId(code)
	filename := fmt.Sprintf("%s/%s/%s/%s.csv", GetMinutePath(), date[0:4], date, cacheId)
	return filename
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

func quarterlyFilename(date, keyword string) string {
	path := quarterlyCachePath(date)
	filename := fmt.Sprintf("%s/%s.csv", path, keyword)
	return filename
}

// ReportsFilename 报告数据文件名
func ReportsFilename(date string) string {
	return quarterlyFilename(date, "reports")
}

// PreviewReportFilename 业绩预告文件名
func PreviewReportFilename(date string) string {
	return quarterlyFilename(date, "preview")
}

// TickFilename tick文件比较多, 目录结构${tick}/${YYYY}/${YYYYMMDD}/${CacheIdPath}
func TickFilename(code, date string) string {
	date = trading.FixTradeDate(date, FilenameDate)
	cacheId := CacheId(code)
	tickPath := fmt.Sprintf("%s/%s/%s/%s.csv", GetTickPath(), date[0:4], date, cacheId)
	return tickPath
}

// FundFlowFilename 通过证券代码获取资金流向的缓存文件路径
func FundFlowFilename(securityCode string) string {
	cacheId := CacheId(securityCode)
	length := len(cacheId)
	filepath := fmt.Sprintf("%s/%s/%s.csv", GetFundFlowPath(), cacheId[:length-3], cacheId)
	return filepath
}

// SnapshotFilename 快照数据文件
func SnapshotFilename(securityCode string) string {
	cacheId := CacheId(securityCode)
	length := len(cacheId)
	filepath := fmt.Sprintf("%s/%s/%s.csv", GetSnapshotPath(), cacheId[:length-3], cacheId)
	return filepath
}
