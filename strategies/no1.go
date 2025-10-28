package strategies

import (
	"gitee.com/quant1x/data/level1/securities"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/realtime"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/gox/concurrent"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
)

const (
	// MaximumResultDays 结果最大天数
	MaximumResultDays int = 3
)

func init() {
	err := models.Register(ModelNo1{})
	if err != nil {
		logger.Fatalf("%+v", err)
	}
}

// ModelNo1 1号模型
//
//	FormulaNo1 3天内5天线上穿10天线，10天线上穿20天线的个股
//	count(cross(MA(c,5),MA(c,10)),3)>=1 and count(cross(MA(c,10),MA(c,20)),3)>=1
type ModelNo1 struct {
}

func (m ModelNo1) Code() models.ModelKind {
	return models.ModelHousNo1
}

func (m ModelNo1) Name() string {
	return models.MapStrategies[m.Code()].Name
}

func (m ModelNo1) OrderFlag() string {
	return models.OrderFlagTail
}

func (m ModelNo1) Filter(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error {
	return GeneralFilter(ruleParameter, snapshot)
}

func (m ModelNo1) Sort(snapshots []factors.QuoteSnapshot) models.SortedStatus {
	_ = snapshots
	return models.SortDefault
}

func (m ModelNo1) Evaluate(securityCode string, result *concurrent.TreeMap[string, models.ResultInfo]) {
	history := factors.GetL5History(securityCode)
	if history == nil {
		return
	}
	snapshot := models.GetStrategySnapshot(securityCode)
	if snapshot == nil {
		return
	}

	// 取出昨日的数据
	r1MA5 := history.MA5
	r1MA10 := history.MA10
	r1MA20 := history.MA20

	ma5 := realtime.IncrementalMovingAverage(history.MA4, 5, snapshot.Price)
	ma10 := realtime.IncrementalMovingAverage(history.MA9, 10, snapshot.Price)
	ma20 := realtime.IncrementalMovingAverage(history.MA19, 20, snapshot.Price)

	// 组织series
	s5 := pandas.ToSeries(r1MA5, ma5)
	s10 := pandas.ToSeries(r1MA10, ma10)
	s20 := pandas.ToSeries(r1MA20, ma20)

	// 两个金叉
	c1 := CROSS(s5, s10)
	c2 := CROSS(s10, s20)
	// 横向对比
	d := c1.And(c2)
	s := utils.BoolIndexOf(d, -1)
	if s {
		price := snapshot.Price
		date := snapshot.Date
		result.Put(securityCode, models.ResultInfo{Code: securityCode,
			Name:         securities.GetStockName(securityCode),
			Date:         date,
			Rate:         0.00,
			Buy:          price,
			Sell:         price * 1.05,
			StrategyCode: m.Code(),
			StrategyName: m.Name()})
	}
}
