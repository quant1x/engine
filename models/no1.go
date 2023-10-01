package models

import (
	"gitee.com/quant1x/engine/features"
	"gitee.com/quant1x/engine/features/base"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/util/treemap"
	"gitee.com/quant1x/pandas"
	. "gitee.com/quant1x/pandas/formula"
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

func (m *ModelNo1) Code() ModelKind {
	return ModelHousNo1
}

func (m *ModelNo1) Name() string {
	return mapStrategies[m.Code()].Name
}

func (m *ModelNo1) Evaluate(securityCode string, result *treemap.Map) {
	lastDate := trading.LastTradeDate()
	klines := base.CheckoutKLines(securityCode, lastDate)
	if len(klines) < KLineMin {
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
	s := features.BoolIndexOf(d, -1)
	if s {
		price := features.SeriesIndexOf(CLOSE, -1)
		result.Put(securityCode, ResultInfo{Code: securityCode,
			Name:         securities.GetStockName(securityCode),
			Date:         features.StringIndexOf(DATE, -1),
			Rate:         0.00,
			Buy:          price,
			Sell:         price * 1.05,
			StrategyCode: m.Code(),
			StrategyName: m.Name()})
	}
}
