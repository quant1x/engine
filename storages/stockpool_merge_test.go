package storages

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/concurrent"
)

type TestModel struct{}

func (m TestModel) Code() models.ModelKind {
	return 0
}

func (m TestModel) Name() string {
	return "0号策略"
}

func (m TestModel) OrderFlag() string {
	return models.OrderFlagTick
}

func (m TestModel) Filter(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error {
	//TODO implement me
	panic("implement me")
}

func (m TestModel) Sort(snapshots []factors.QuoteSnapshot) models.SortedStatus {
	//TODO implement me
	panic("implement me")
}

func (m TestModel) Evaluate(securityCode string, result *concurrent.TreeMap[string, models.ResultInfo]) {
	//TODO implement me
	panic("implement me")
}
