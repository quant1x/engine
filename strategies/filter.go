package strategies

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/rules"
)

// GeneralFilter 过滤条件
func GeneralFilter(snapshot factors.QuoteSnapshot) error {
	passed, failed, err := rules.Filter(snapshot)
	_ = passed
	_ = failed
	return err
}
