package factors

import "gitee.com/quant1x/gotdx/quotes"

// QuoteSnapshot 即时行情快照(副本)
type QuoteSnapshot struct {
	Date                  string               // 交易日期
	ServerTime            string               // 时间
	SecurityCode          string               // 证券代码
	ExchangeState         quotes.ExchangeState // 交易状态
	State                 quotes.TradeState    // 上市公司状态
	Market                uint8                // 市场ID
	Code                  string               `name:"证券代码"`  // 代码
	Name                  string               `name:"证券名称"`  // 证券名称
	Active                uint16               `name:"活跃度"`   // 活跃度
	LastClose             float64              `name:"昨收"`    // 昨收
	Open                  float64              `name:"开盘价"`   // 开盘
	OpeningChangeRate     float64              `name:"开盘涨幅%"` // 开盘
	Price                 float64              `name:"现价"`    // 现价
	ChangeRate            float64              `name:"涨跌幅%"`  // 涨跌幅
	PremiumRate           float64              `name:"溢价率%"`  // 集合竞价买入溢价, hedgeRatio
	High                  float64              // 最高
	Low                   float64              // 最低
	Vol                   int                  // 总量
	CurVol                int                  // 现量
	Amount                float64              // 总金额
	SVol                  int                  // 内盘
	BVol                  int                  // 外盘
	IndexOpenAmount       int                  // 指数-集合竞价成交金额=开盘成交金额
	StockOpenAmount       int                  // 个股-集合竞价成交金额=开盘成交金额
	OpenVolume            int                  `name:"开盘量"` // 集合竞价-开盘量, 单位是股
	CloseVolume           int                  `name:"收盘量"` /// 集合竞价-收盘量, 单位是股
	IndexUp               int                  // 指数有效-上涨数
	IndexUpLimit          int                  // 指数有效-涨停数
	IndexDown             int                  // 指数有效-下跌数
	IndexDownLimit        int                  // 指数有效-跌停数
	OpenBiddingDirection  int                  `name:"开盘竞价" dataframe:"开盘竞价"` // 竞价方向, 交易当日集合竞价开盘时更新
	OpenVolumeDirection   int                  `name:"开盘竞量" dataframe:"开盘竞量"` // 委托量差, 交易当日集合竞价开盘时更新
	CloseBiddingDirection int                  `name:"收盘竞价" dataframe:"收盘竞价"` // 竞价方向, 交易当日集合竞价收盘时更新
	CloseVolumeDirection  int                  `name:"收盘竞量" dataframe:"收盘竞量"` // 委托量差, 交易当日集合竞价收盘时更新
	Rate                  float64              // 涨速
	TopNo                 int                  // 板块排名
	TopCode               string               // 领涨个股
	TopName               string               // 领涨个股名称
	TopRate               float64              // 领涨个股涨幅
	ZhanTing              int                  // 涨停数
	Ling                  int                  // 平盘数
	Count                 int                  // 总数
	Capital               float64              `name:"流通盘"`    // 流通盘
	FreeCapital           float64              `name:"自由流通股本"` // 自由流通股本
	OpenTurnZ             float64              `name:"开盘换手Z%"` // 开盘换手
	OpenQuantityRatio     float64              `name:"开盘量比"`
	QuantityRatio         float64              `name:"量比"`
	ChangePower           float64              `name:"涨跌力度"` // 开盘金额除以开盘涨幅
	AverageBiddingVolume  int                  `name:"委托均量"` // 委托均量
	NextOpen              float64              // 仅回测有效: 下一个交易日开盘价
	NextClose             float64              // 仅回测有效: 下一个交易日收盘价
	NextHigh              float64              // 仅回测有效: 下一个交易日最高价
	NextLow               float64              // 仅回测有效: 下一个交易日最低价
	UpdateTime            string               // 本地时间戳
}

// ExistUpwardGap 是否存在向上跳空缺口
func (q QuoteSnapshot) ExistUpwardGap() bool {
	history := GetL5History(q.SecurityCode, q.Date)
	if history == nil {
		return false
	}
	return history.HIGH < q.Low
}

// ExistDownwardGap 是否存在向下跳空缺口
func (q QuoteSnapshot) ExistDownwardGap() bool {
	history := GetL5History(q.SecurityCode, q.Date)
	if history == nil {
		return false
	}
	return history.LOW > q.High
}

// BoxUpwardGap K线实体位于昨日K线实体上方
func (q QuoteSnapshot) BoxUpwardGap() bool {
	history := GetL5History(q.SecurityCode, q.Date)
	if history == nil {
		return false
	}
	lastBoxHigh := max(history.OPEN, history.CLOSE)
	boxLow := min(q.Open, q.Price)
	return lastBoxHigh <= boxLow
}

// BoxDownwardGap K线实体位于昨日K线实体下方
func (q QuoteSnapshot) BoxDownwardGap() bool {
	history := GetL5History(q.SecurityCode, q.Date)
	if history == nil {
		return false
	}
	lastBoxLow := min(history.OPEN, history.CLOSE)
	boxHigh := max(q.Open, q.Price)
	return lastBoxLow >= boxHigh
}

// KLineWeaknessToStrength K线弱转强, 毕竟还是弱, 下一个交易日有低吸机会
func (q QuoteSnapshot) KLineWeaknessToStrength() bool {
	history := GetL5History(q.SecurityCode, q.Date)
	if history == nil {
		return false
	}
	// 昨日收盘于均价线之下
	c1 := history.CLOSE <= history.AveragePrice
	// 当日开盘不低于昨日均价线
	c2 := q.Open >= history.AveragePrice
	return c1 && c2
}

// WeaknessToStrength 弱转强
func (q QuoteSnapshot) WeaknessToStrength() bool {
	history := GetL5History(q.SecurityCode, q.Date)
	if history == nil {
		return false
	}
	// 1. 昨日转强
	// 昨日阴线
	c11 := history.CLOSE < history.OPEN
	// 昨日收盘于均价线之上
	c12 := history.CLOSE >= history.AveragePrice
	// 当日开盘不低于昨日收盘
	c13 := q.Open >= history.CLOSE
	frontToStrength := c11 && c12 && c13
	// 2. 隔日转强
	// 昨日收盘于均价线之下
	c21 := history.CLOSE <= history.AveragePrice
	// 当日开盘不低于昨日均价线
	c22 := q.Open >= history.AveragePrice
	nextToStrength := c21 && c22

	return frontToStrength || nextToStrength
}
