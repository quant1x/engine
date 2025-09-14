package strategies

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/rules"
)

// GeneralFilter 过滤条件
//
//	执行所有在册的规则
func GeneralFilter(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error {
	passed, failed, err := rules.Filter(ruleParameter, snapshot)
	_ = passed
	_ = failed
	return err
}
