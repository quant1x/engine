package cache

import "gitee.com/quant1x/gotdx/trading"

// CorrectDate 校正日期
func CorrectDate(date string) (cacheDate, resourcesDate string) {
	cacheDate = trading.FixTradeDate(date)
	dates := trading.LastNDate(cacheDate, 1)
	if len(dates) == 0 {
		// TODO: 存在一定概率的坑, 比如1990-12-19执行
		resourcesDate = cacheDate
	} else {
		resourcesDate = dates[0]
	}
	cacheDate = trading.NextTradeDate(resourcesDate)
	return
}

// DefaultCanReadDate 获取默认可以读缓存文件的日期
func DefaultCanReadDate() string {
	dateOfReadingData := trading.GetCurrentDate()
	return dateOfReadingData
}

// DefaultCanUpdateDate 获取默认可以更新缓存文件的日期
func DefaultCanUpdateDate() string {
	currentDate := trading.GetCurrentDate()
	dateOfUpdatingData := trading.NextTradeDate(currentDate)
	return dateOfUpdatingData
}

// DefaultCanUpdateDate 获取默认可以更新缓存文件的日期
func testDefaultCanUpdateDate() string {
	dateOfUpdatingData := trading.GetCurrentlyDay()
	return dateOfUpdatingData
}
