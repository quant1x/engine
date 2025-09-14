package factors

import (
	"time"

	"github.com/quant1x/engine/cache"
	"github.com/quant1x/gotdx/securities"
)

// GetTimestamp 时间戳
//
//	格式: YYYY-MM-DD hh:mm:ss.SSS
func GetTimestamp() string {
	now := time.Now()
	return now.Format(cache.TimeStampMilli)
}

// PriceDigits 获取证券标的价格保留小数点后几位
//
//	默认范围2, 即小数点后2位
func PriceDigits(securityCode string) int {
	securityInfo, ok := securities.CheckoutSecurityInfo(securityCode)
	if !ok {
		return 2
	}
	return int(securityInfo.DecimalPoint)
}
