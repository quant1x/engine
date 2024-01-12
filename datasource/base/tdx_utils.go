package base

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/logger"
	"strconv"
)

// 通达信协议日期为YYYYMMDD格式的十进制整型
func toTdxProtocolDate(date string) (uint32, error) {
	protoDate := exchange.FixTradeDate(date, cache.TDX_FORMAT_PROTOCOL_DATE)
	transDate, err := strconv.ParseUint(protoDate, 10, 32)
	if err != nil {
		logger.Errorf("转换日期[%s]到uint32异常, error=%+v", date, err)
		return 0, err
	}
	return uint32(transDate), nil
}
