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
	Code         string  `name:"证券代码" json:"code" csv:"code" array:"0"`
	Name         string  `name:"证券名称" json:"name" csv:"name" array:"1"`
	Date         string  `name:"信号日期" json:"date" csv:"date" array:"2"`
	TurnZ        float64 `name:"开盘换手Z" json:"turn_z" csv:"turn_z" arrar:"3"`
	Rate         float64 `name:"涨跌幅%" json:"rate" csv:"rate"`
	Buy          float64 `name:"委托价格" json:"buy" csv:"buy" array:"3"`
	Sell         float64 `name:"目标价格" json:"sell" csv:"sell" array:"4"`
	StrategyCode int     `name:"策略编码" json:"strategy_code" csv:"strategy_code" array:"5"`
	StrategyName string  `name:"策略名称" json:"strategy_name" csv:"strategy_name" array:"6"`
	//BlockCode    string  `name:"板块代码" json:"block_code" csv:"block_code" array:"7"`
	BlockType      string  `name:"板块类型"`
	BlockName      string  `name:"板块名称" json:"block_name" csv:"block_name" array:"7"`
	BlockRate      float64 `name:"板块涨幅%" json:"block_rate" csv:"block_rate" array:"8"`
	BlockTop       int     `name:"板块排名" json:"block_top" csv:"block_top" array:"9"`
	BlockRank      int     `name:"个股排名" json:"block_rank" csv:"block_top" array:"10"`
	BlockZhangTing string  `name:"板块涨停数" json:"block_zhangting" csv:"block_zhangting" array:"11"`
	BlockDescribe  string  `name:"上涨/下跌/平盘" json:"block_describe" csv:"block_describe" array:"12"`
	//BlockTopCode string  `name:"领涨股代码" json:"block_top_code" csv:"block_top_code" array:"12"`
	BlockTopName string  `name:"领涨股名称" json:"block_top_name" csv:"block_top_name" array:"13"`
	BlockTopRate float64 `name:"领涨股涨幅%" json:"block_top_rate" csv:"block_top_rate" array:"14"`
	Tendency     string  `name:"短线趋势" json:"tendency" csv:"tendency" array:"15"`
	//Open         float64 `name:"预计开盘" json:"open" csv:"open" array:"16"`
	//CLOSE        float64 `name:"预计收盘" json:"close" csv:"close" array:"17"`
	//High         float64 `name:"预计最高" json:"high" csv:"high" array:"18"`
	//Low          float64 `name:"预计最低" json:"low" csv:"low" array:"19"`
}
