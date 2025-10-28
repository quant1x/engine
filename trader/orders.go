package trader

import (
	"path/filepath"
	"strings"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/gox/api"
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

// GetLocalOrderDates 获取本地订单日期列表
func GetLocalOrderDates() (list []string) {
	prefix := "orders."
	pattern := filepath.Join(traderQmtOrderPath, prefix+"*")
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		return list
	}
	for _, filename := range files {
		arr := strings.Split(filename, prefix)
		if len(arr) != 2 {
			continue
		}
		list = append(list, arr[1])
	}
	return
}
