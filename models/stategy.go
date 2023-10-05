package models

import "gitee.com/quant1x/gox/util/treemap"

// ModelKind 做多64个策略
type ModelKind = int

const (
	ModelZero ModelKind = 0 // 0号策略
)

const (
	ModelHousNo1 ModelKind = 1 << iota // 1号策略
	ModelTail                          // 尾盘策略
)

const (
	DefaultStrategy = ModelHousNo1
	KLineMin        = 89 // K线最少记录数
)

// Strategy 策略/公式指标(features)接口
type Strategy interface {
	// Code 策略编号
	Code() ModelKind
	// Name 策略名称
	Name() string
	// Evaluate 评估 日线数据
	Evaluate(securityCode string, result *treemap.Map)
}

type StrategyWrap struct {
	Type ModelKind
	Name string
}

var (
	mapStrategies = map[ModelKind]StrategyWrap{
		ModelZero:    {Type: ModelZero, Name: "0号策略"},
		ModelHousNo1: {Type: ModelHousNo1, Name: "1号策略"},
		ModelTail:    {Type: ModelTail, Name: "尾盘策略"},
	}
)

// ResultInfo 策略结果
type ResultInfo struct {
	Code           string  `name:"证券代码" dataframe:"code"`
	Name           string  `name:"证券名称" dataframe:"name"`
	Date           string  `name:"信号日期" dataframe:"date"`
	TurnZ          float64 `name:"开盘换手Z" dataframe:"turn_z"`
	Rate           float64 `name:"涨跌幅%" dataframe:"rate"`
	Buy            float64 `name:"委托价格" dataframe:"buy"`
	Sell           float64 `name:"目标价格" dataframe:"sell"`
	StrategyCode   int     `name:"策略编码" dataframe:"strategy_code"`
	StrategyName   string  `name:"策略名称" dataframe:"strategy_name"`
	BlockType      string  `name:"板块类型" dataframe:"block_type"`
	BlockCode      string  `name:"板块代码" dataframe:"block_code"`
	BlockName      string  `name:"板块名称" dataframe:"block_name"`
	BlockRate      float64 `name:"板块涨幅%" dataframe:"block_rate"`
	BlockTop       int     `name:"板块排名" dataframe:"block_top"`
	BlockZhangTing string  `name:"板块涨停数" dataframe:"block_zhangting"`
	BlockDescribe  string  `name:"涨/跌/平" dataframe:"block_describe"`
	BlockRank      int     `name:"个股排名" dataframe:"block_rank"`
	BlockTopCode   string  `name:"领涨股代码" dataframe:"block_top_code"`
	BlockTopName   string  `name:"领涨股名称" dataframe:"block_top_name"`
	BlockTopRate   float64 `name:"领涨股涨幅%" dataframe:"block_top_rate"`
	Tendency       string  `name:"短线趋势" dataframe:"tendency"`
}
