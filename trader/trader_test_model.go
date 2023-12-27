package trader

import (
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/concurrent"
)

type TestModel struct{}

func (TestModel) Code() models.ModelKind {
	return 82
}

func (s TestModel) Name() string {
	//TODO implement me
	panic("implement me")
}

func (s TestModel) OrderFlag() string {
	//TODO implement me
	panic("implement me")
}

func (s TestModel) Filter(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error {
	//TODO implement me
	panic("implement me")
}

func (s TestModel) Sort(snapshots []factors.QuoteSnapshot) models.SortedStatus {
	//TODO implement me
	panic("implement me")
}

func (s TestModel) Evaluate(securityCode string, result *concurrent.TreeMap[string, models.ResultInfo]) {
	//TODO implement me
	panic("implement me")
}
