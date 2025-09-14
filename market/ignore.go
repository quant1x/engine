package market

import (
	"strings"

	"gitee.com/quant1x/gotdx/securities"
)

// IsNeedIgnore 需要忽略的个股
//
//	检测需要的剔除ST、退市和摘牌的个股
func IsNeedIgnore(securityCode string) bool {
	securityInfo, ok := securities.CheckoutSecurityInfo(securityCode)
	if !ok {
		// 没找到, 忽略
		return true
	}
	name := strings.ToUpper(securityInfo.Name)
	if strings.Contains(name, "ST") {
		// ST标志, 忽略
		return true
	}
	if strings.Contains(name, "退") {
		// 退市标志, 忽略
		return true
	}
	if strings.Contains(name, "摘牌") {
		// 摘牌标志, 忽略
		return true
	}
	return false
}
