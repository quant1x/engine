package dfcf

import (
	"encoding/json"
	"fmt"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"testing"
)

func TestStockNoticeReport(t *testing.T) {
	notices, _, err := StockNotices("603045", "20240101", "20240430", 1)
	if err != nil {
		return
	}
	data, err := json.Marshal(notices)
	fmt.Println(api.Bytes2String(data))
}

func TestNoticeDateForAnnualReport(t *testing.T) {
	code := "603045"
	date := "2024-04-09"
	y, q := NoticeDateForReport(code, date)
	fmt.Println(y, q)
	ys := exchange.DateRange(date, y)
	fmt.Println(len(ys))
	qs := exchange.DateRange(date, q)
	fmt.Println(len(qs))
}
