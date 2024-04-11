package factors

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/dfcf"
	"gitee.com/quant1x/exchange"
)

var (
	__mapMarginTradingTargets = map[string]dfcf.SecurityMarginTrading{}
)

// MarginTradingTargetInit 一次性缓存两融数据, 交易日9点后更新上一个交易的两融数据
func MarginTradingTargetInit(date string) {
	clear(__mapMarginTradingTargets)
	_, featureDate := cache.CorrectDate(date)
	list := dfcf.GetMarginTradingList(featureDate)
	for _, v := range list {
		securityCode := exchange.CorrectSecurityCode(v.SecuCode)
		__mapMarginTradingTargets[securityCode] = v
	}
}

// GetMarginTradingTarget 获取两融数据
func GetMarginTradingTarget(code string) (dfcf.SecurityMarginTrading, bool) {
	securityCode := exchange.CorrectSecurityCode(code)
	v, ok := __mapMarginTradingTargets[securityCode]
	return v, ok
}
