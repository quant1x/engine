package realtime

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"time"
)

// TradeSessionHasEnd 是否收盘
func TradeSessionHasEnd(date string) bool {
	tm, _ := time.Parse(cache.INDEX_DATE, date)
	_, status := exchange.CanUpdateInRealtime(tm)
	return status == exchange.ExchangeClosing
}
