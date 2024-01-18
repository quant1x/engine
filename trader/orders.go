package trader

import (
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"path/filepath"
)

// GetOrderFilename 获得订单文件名
//
//	qmt/账户id/orders.yyyy-mm-dd
func GetOrderFilename(date ...string) string {
	var tradeDate string
	if len(date) > 0 {
		tradeDate = exchange.FixTradeDate(date[0])
	} else {
		tradeDate = exchange.LastTradeDate()
	}
	filename := filepath.Join(traderQmtOrderPath, "orders."+tradeDate)
	return filename
}

// GetOrderList 获取指定日期的订单列表
func GetOrderList(date string) []OrderDetail {
	filename := GetOrderFilename(date)
	var list []OrderDetail
	_ = api.CsvToSlices(filename, &list)
	return list
}
