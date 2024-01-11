package cache

import (
	"gitee.com/quant1x/exchange"
	"time"
)

const (
	TDX_FORMAT_PROTOCOL_DATE = "20060102"   // 通达信协议的日期字符串格式
	CACHE_DATE               = "20060102"   // 缓存日期
	INDEX_DATE               = "2006-01-02" // 索引日期格式
	TDX_DATE                 = "20060102"   // 通达信日期
	YearOnly                 = "2006"       // 仅年份
	TimeStampMilli           = "2006-01-02 15:04:05.000"
	TimeStampMicro           = "2006-01-02 15:04:05.000000"
	TimeStampNano            = "2006-01-02 15:04:05.000000000"
)

func Today() string {
	now := time.Now()
	return now.Format(CACHE_DATE)
}

//// CorrectDate 矫正日期, 统一格式: 20060102
//func CorrectDate(date string) string {
//	dt, err := api.ParseTime(date)
//	if err != nil {
//		return Today()
//	}
//	date = dt.Format(CACHE_DATE)
//	return date
//}

// CorrectDate 校正日期
func CorrectDate(date string) (cacheDate, resourcesDate string) {
	cacheDate = exchange.FixTradeDate(date)
	dates := exchange.LastNDate(cacheDate, 1)
	if len(dates) == 0 {
		// TODO: 存在一定概率的坑, 比如1990-12-19执行
		resourcesDate = cacheDate
	} else {
		resourcesDate = dates[0]
	}
	cacheDate = exchange.NextTradeDate(resourcesDate)
	return
}

// DefaultCanReadDate 获取默认可以读缓存文件的日期
func DefaultCanReadDate() string {
	dateOfReadingData := exchange.GetCurrentDate()
	return dateOfReadingData
}

// DefaultCanUpdateDate 获取默认可以更新缓存文件的日期
func DefaultCanUpdateDate() string {
	currentDate := exchange.GetCurrentDate()
	dateOfUpdatingData := exchange.NextTradeDate(currentDate)
	return dateOfUpdatingData
}

// DefaultCanUpdateDate 获取默认可以更新缓存文件的日期
func testDefaultCanUpdateDate() string {
	dateOfUpdatingData := exchange.GetCurrentlyDay()
	return dateOfUpdatingData
}
