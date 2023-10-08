package base

import (
	"gitee.com/quant1x/gotdx/trading"
	"strconv"
)

const (
	TDX_FORMAT_PROTOCOL_DATE = "20060102"
	TickDefaultStartDate     = "2023-01-01" // 分笔成交最早的日期
)

// 通达信协议日期为YYYYMMDD格式的十进制整型
func toTdxProtocolDate(date string) uint32 {
	protoDate := trading.FixTradeDate(date, TDX_FORMAT_PROTOCOL_DATE)
	transDate, err := strconv.ParseUint(protoDate, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(transDate)
}
