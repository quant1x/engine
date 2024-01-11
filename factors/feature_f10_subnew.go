package factors

import (
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/api"
)

const (
	SubNewStockYears = 1 // 次新股几年内
)

// 检查在date之前是否存在除权除息
func checkXdxr(list []quotes.XdxrInfo, date string) *quotes.XdxrInfo {
	for _, v := range list {
		if v.Category == 1 && date >= v.Date {
			return &v
		}
	}
	return nil
}

//// IsSubNewStock 检查是否次新股
//func IsSubNewStock(securityCode, date string) bool {
//	date = trading.FixTradeDate(date)
//	securityCode = proto.CorrectSecurityCode(securityCode)
//	f10 := flash.GetL5F10(securityCode, date)
//	if f10 == nil {
//		return false
//	}
//	return IsSubNewStockByIpoDate(securityCode, f10.IpoDate, date)
//}

// IsSubNewStockByIpoDate 检查是否次新股
func IsSubNewStockByIpoDate(securityCode, ipoDate, date string) bool {
	ipoDate = exchange.FixTradeDate(ipoDate)
	date = exchange.FixTradeDate(date)
	listingDate, err := api.ParseTime(ipoDate)
	if err != nil {
		return false
	}
	tm := listingDate.AddDate(SubNewStockYears, 0, 0)
	after := tm.Format(exchange.TradingDayDateFormat)
	if date >= after {
		return false
	}
	//xdxrs := base.GetCacheXdxrList(securityCode)
	//if len(xdxrs) == 0 {
	//	return false
	//}
	//api.SliceSort(xdxrs, func(a, b quotes.XdxrInfo) bool {
	//	return a.Date > b.Date
	//})
	//xdxrInfo := checkXdxr(xdxrs, after)
	//if xdxrInfo == nil {
	//	return true
	//}

	return true
}
