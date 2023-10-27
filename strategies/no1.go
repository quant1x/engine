package strategies

import (
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/util/treemap"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
	"gitee.com/quant1x/pandas/stat"
)

const (
	// MaximumResultDays 结果最大天数
	MaximumResultDays int = 3
)

// ModelNo1 1号模型
//
//	FormulaNo1 3天内5天线上穿10天线，10天线上穿20天线的个股
//	count(cross(MA(c,5),MA(c,10)),3)>=1 and count(cross(MA(c,10),MA(c,20)),3)>=1
type ModelNo1 struct {
}

func (m *ModelNo1) OrderFlag() string {
	return models.OrderFlagTail
}

func (m *ModelNo1) Filter(snapshot models.QuoteSnapshot) bool {
	return RuleFilter(snapshot)
}

func (m *ModelNo1) Code() models.ModelKind {
	return models.ModelHousNo1
}

func (m *ModelNo1) Name() string {
	return models.MapStrategies[m.Code()].Name
}

func (m *ModelNo1) v1Evaluate(securityCode string, result *treemap.Map) {
	lastDate := trading.LastTradeDate()
	klines := base.CheckoutKLines(securityCode, lastDate)
	if len(klines) < models.KLineMin {
		return
	}
	df := pandas.LoadStructs(klines)
	if df.Err != nil {
		return
	}
	var (
		DATE  = df.Col("date")
		CLOSE = df.ColAsNDArray("close")
	)
	days := CLOSE.Len()
	if days < 1 {
		return
	}

	// 取5、10、20日均线
	ma5 := MA(CLOSE, 5)
	ma10 := MA(CLOSE, 10)
	ma20 := MA(CLOSE, 20)
	if ma5.Len() != days || ma10.Len() != days || ma20.Len() != days {
		logger.Errorf("[%s]: 均线, 数据没对齐", m.Name())
		return
	}
	// 两个金叉
	c1 := CROSS(ma5, ma10)
	c2 := CROSS(ma10, ma20)
	// 两个统计
	N := MaximumResultDays
	r1 := COUNT(c1.Bools(), N)
	r2 := COUNT(c2.Bools(), N)
	// 横向对比
	d := r1.And(r2)
	s := factors.BoolIndexOf(d, -1)
	if s {
		price := factors.SeriesIndexOf(CLOSE, -1)
		result.Put(securityCode, models.ResultInfo{Code: securityCode,
			Name:         securities.GetStockName(securityCode),
			Date:         factors.StringIndexOf(DATE, -1),
			Rate:         0.00,
			Buy:          price,
			Sell:         price * 1.05,
			StrategyCode: m.Code(),
			StrategyName: m.Name()})
	}
}

func (m *ModelNo1) Evaluate(securityCode string, result *treemap.Map) {
	history := smart.GetL5History(securityCode)
	if history == nil {
		return
	}
	snapshot := models.GetQuoteSnapshot(securityCode)
	if snapshot == nil {
		return
	}

	// 取出昨日的数据
	lastNo1 := history.Last.No1
	r1MA5 := lastNo1.MA5
	r1MA10 := lastNo1.MA10
	r1MA20 := lastNo1.MA20
	// 取出今日的半成品数据
	today := history.Payloads.No1.Increase(*snapshot).(*factors.HousNo1)
	ma5 := today.MA5
	ma10 := today.MA10
	ma20 := today.MA20

	// 组织series
	s5 := stat.NewSeries(r1MA5, ma5)
	s10 := stat.NewSeries(r1MA10, ma10)
	s20 := stat.NewSeries(r1MA20, ma20)

	// 两个金叉
	c1 := CROSS(s5, s10)
	c2 := CROSS(s10, s20)
	// 横向对比
	d := c1.And(c2)
	s := factors.BoolIndexOf(d, -1)
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
