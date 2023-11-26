package storages

import (
	"fmt"
)

// StockPool 股票池
type StockPool struct {
	Status         StrategyStatus `name:"策略状态" dataframe:"status"`
	Date           string         `name:"信号日期" dataframe:"date"`
	Code           string         `name:"证券代码" dataframe:"code"`
	Name           string         `name:"证券名称" dataframe:"name"`
	TurnZ          float64        `name:"开盘换手Z" dataframe:"turn_z"`
	Rate           float64        `name:"涨跌幅%" dataframe:"rate"`
	Buy            float64        `name:"委托价格" dataframe:"buy"`
	Sell           float64        `name:"目标价格" dataframe:"sell"`
	StrategyCode   int            `name:"策略编码" dataframe:"strategy_code"`
	StrategyName   string         `name:"策略名称" dataframe:"strategy_name"`
	Rules          uint64         `name:"规则" dataframe:"rules"`
	BlockType      string         `name:"板块类型" dataframe:"block_type"`
	BlockCode      string         `name:"板块代码" dataframe:"block_code"`
	BlockName      string         `name:"板块名称" dataframe:"block_name"`
	BlockRate      float64        `name:"板块涨幅%" dataframe:"block_rate"`
	BlockTop       int            `name:"板块排名" dataframe:"block_top"`
	BlockRank      int            `name:"个股排名" dataframe:"block_rank"`
	BlockZhangTing string         `name:"板块涨停数" dataframe:"block_zhangting"`
	BlockDescribe  string         `name:"涨/跌/平" dataframe:"block_describe"`
	BlockTopCode   string         `name:"领涨股代码" dataframe:"block_top_code"`
	BlockTopName   string         `name:"领涨股名称" dataframe:"block_top_name"`
	BlockTopRate   float64        `name:"领涨股涨幅%" dataframe:"block_top_rate"`
	Tendency       string         `name:"短线趋势" dataframe:"tendency"`
	OrderStatus    int            `name:"订单状态" dataframe:"order_status"`
	Active         int            `name:"活跃度" dataframe:"active"`
	Speed          float64        `name:"涨速" dataframe:"speed"`
	CreateTime     string         `name:"创建时间" dataframe:"create_time"`
	UpdateTime     string         `name:"更新时间" dataframe:"update_time"`
}

// Key 索引字段: 日期/策略代码/证券代码
func (sp StockPool) Key() string {
	return fmt.Sprintf("%s/%d/%s", sp.Date, sp.StrategyCode, sp.Code)
}

type StrategyStatus int

const (
	StrategyMiss          StrategyStatus = 0x0000 // 策略 - 未命中
	StrategyHit           StrategyStatus = 0x0001 // 策略 - 命中
	StrategyCancel        StrategyStatus = 0x0002 // 策略 - 召回
	StrategyPassed        StrategyStatus = 0x0004 // 策略 - 成功
	StrategyOrdered       StrategyStatus = 0x0008 // 策略 - 已下单
	StrategySucceeded     StrategyStatus = 0x0010 // 策略 - 订单已成功
	StrategyJunk          StrategyStatus = 0x0080 // 策略 - 作废
	StrategyAlreadyExists StrategyStatus = 0x8000 // 已存在
)

var (
	mapStrategiesOfOrder = map[StrategyStatus]string{
		StrategyMiss:          "未命中",
		StrategyHit:           "命中",
		StrategyCancel:        "召回",
		StrategyPassed:        "通过",
		StrategyJunk:          "作废",
		StrategyAlreadyExists: "已存在",
	}
)

func (s *StrategyStatus) test(other StrategyStatus) bool {
	return (*s & other) == other
}

// Set 设置状态
func (s *StrategyStatus) Set(other StrategyStatus, on bool) {
	if on {
		*s |= other
	} else {
		*s &= ^other
	}
}

// IsHit 是否命中
func (s *StrategyStatus) IsHit() bool {
	return s.test(StrategyHit)
}

// IsCancel 是否召回/撤销
func (s *StrategyStatus) IsCancel() bool {
	return s.test(StrategyCancel)
}

func (s *StrategyStatus) IsPassed() bool {
	return s.test(StrategyPassed)
}
