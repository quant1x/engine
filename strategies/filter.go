package strategies

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/rules"
)

// GeneralFilter 过滤条件
func GeneralFilter(snapshot factors.QuoteSnapshot) bool {
	passed, failed, err := rules.Filter(snapshot)
	if failed != rules.Pass {
		return false
	}
	_ = passed
	_ = err
	return true
}
