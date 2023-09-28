package base

import (
	"gitee.com/quant1x/engine/factors"
)

// KLine 日K线基础结构
type KLine struct {
	Date   string  `dataframe:"date"`   // 日期
	Open   float64 `dataframe:"open"`   // 开盘价
	Close  float64 `dataframe:"close"`  // 收盘价
	High   float64 `dataframe:"high"`   // 最高价
	Low    float64 `dataframe:"low"`    // 最低价
	Volume float64 `dataframe:"volume"` // 成交量
	Amount float64 `dataframe:"amount"` // 成交金额
	Up     int     `dataframe:"up"`     // 指数类是上涨家数, 个股类是外盘
	Down   int     `dataframe:"down"`   // 指数类是下跌家数, 个股类是内盘
}

func (K *KLine) Init() error {
	return nil
}

func (K *KLine) Kind() factors.FeatureKind {
	return factors.FeatureBaseKLine
}

func (K *KLine) Name() string {
	return "基础K线"
}

func (K *KLine) Update(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (K *KLine) Repair(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (K *KLine) Increase(snapshot factors.Snapshot) {
	//TODO implement me
	panic("implement me")
}
