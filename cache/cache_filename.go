package cache

import (
	"fmt"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"path/filepath"
)

const (
	FilenameDate = "20060102" // 缓存文件相关的日期格式
)

const (
	// 使用时：cacheID + chipFileSuffix
	chipFileSuffix = ".bin"
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
	filename := fmt.Sprintf("%s/%s/%s.csv", GetDayPath(), cacheId[:length-3], cacheId)
	return filename
}

// WideFilename 宽表据缓存路径
func WideFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	filename := fmt.Sprintf("%s/%s/%s.csv", GetWidePath(), cacheId[:length-3], cacheId)
	return filename
}

func MinuteFilename(code, date string) string {
	date = exchange.FixTradeDate(date, FilenameDate)
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

// TransFilename 历史成交数据文件比较多, 目录结构${tick}/${YYYY}/${YYYYMMDD}/${CacheIdPath}
func TransFilename(code, date string) string {
	date = exchange.FixTradeDate(date, FilenameDate)
	cacheId := CacheId(code)
	tickPath := fmt.Sprintf("%s/%s/%s/%s.csv", GetTransPath(), date[0:4], date, cacheId)
	return tickPath
}

// FundFlowFilename 通过证券代码获取资金流向的缓存文件路径
func FundFlowFilename(securityCode string) string {
	cacheId := CacheId(securityCode)
	length := len(cacheId)
	filename := fmt.Sprintf("%s/%s/%s.csv", GetFundFlowPath(), cacheId[:length-3], cacheId)
	return filename
}

// SnapshotFilename 快照数据文件
func SnapshotFilename(securityCode string, date string) string {
	date = exchange.FixTradeDate(date, FilenameDate)
	cacheId := CacheId(securityCode)
	//length := len(cacheId)
	filename := fmt.Sprintf("%s/%s/%s/%s.csv", GetSnapshotPath(), date[0:4], date, cacheId)
	return filename
}

// ChipsFilename 筹码分布文件
func ChipsFilename(securityCode string) string {
	idCode := CacheId(securityCode)
	idPath := CacheIdPath(idCode)
	filename := filepath.Join(GetChipsPath(), idPath+chipFileSuffix)
	return filename
}
