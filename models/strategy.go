package models

import (
	"errors"
	"fmt"
	"slices"
	"sync"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/concurrent"
)

// ModelKind 模型类型编码, 整型
type ModelKind = uint64

const (
	ModelZero                ModelKind = 0          // 0号策略, 是一个特殊策略, 不允许覆盖
	ModelHousNo1             ModelKind = 1          // 1号策略, 经典的默认策略, 不允许覆盖
	ModelNo2                 ModelKind = 2          // 2号策略, 不允许覆盖
	ModelNo3                 ModelKind = 3          // 3号策略, 不允许覆盖
	ModelNo4                 ModelKind = 4          // 4号策略, 不允许覆盖
	ModelNo5                 ModelKind = 5          // 5号策略, 不允许覆盖
	ModelNo6                 ModelKind = 6          // 6号策略, 不允许覆盖
	ModelNo7                 ModelKind = 7          // 7号策略, 不允许覆盖
	ModelNo8                 ModelKind = 8          // 8号策略, 不允许覆盖
	ModelNo9                 ModelKind = 9          // 9号策略, 不允许覆盖
	Model89K                 ModelKind = 89         // 89号策略, 89K策略, 不允许覆盖
	ModelOneSizeFitsAllSells ModelKind = 117        // 卖出策略: 一刀切(Panic sell, cookie-cutter, One size fits all sales)
	ModelNoShareHolding      ModelKind = 861        // 卖出策略: 不留了
	ModelForceOverwrite      ModelKind = 0x80000000 // 强制覆盖
)

var (
	// ReserveStrategyNumberRanges Quant1X保留的策略编码范围
	ReserveStrategyNumberRanges = []ModelKind{
		ModelZero,
		ModelHousNo1,
		ModelNo2,
		ModelNo3,
		ModelNo4,
		ModelNo5,
		ModelNo6,
		ModelNo7,
		ModelNo8,
		ModelNo9,
		Model89K,
		ModelOneSizeFitsAllSells,
		ModelNoShareHolding,
	}
)

const (
	DefaultStrategy = ModelHousNo1
	KLineMin        = cache.KLineMin // K线最少记录数
)

const (
	OrderFlagHead = "head" // 早盘订单标志
	OrderFlagTick = "tick" // 实时订单标志
	OrderFlagTail = "tail" // 尾盘订单标志
	OrderFlagSell = "sell" // 卖出订单标志
)

// SortedStatus 排序模式
//
//	这个状态机
type SortedStatus int

const (
	SortNotExecuted = iota // 没有排序, 由engine自行决定
	SortFinished           // 排序已完成
	SortDefault            // 默认排序, 由engine自行决定
	SortNotRequired        // 无需排序, 保持自然顺序, 这个顺序是证券代码的顺序
)

// Strategy 策略/公式指标(features)接口
type Strategy interface {
	// Code 策略编号
	Code() ModelKind
	// Name 策略名称
	Name() string
	// OrderFlag 订单标志
	OrderFlag() string
	// Filter 过滤
	Filter(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error
	// Sort 排序
	Sort([]factors.QuoteSnapshot) SortedStatus
	// Evaluate 评估 日线数据
	Evaluate(securityCode string, result *concurrent.TreeMap[string, ResultInfo])
}

// QmtStrategyName 获取用于QMT系统的策略名称
func QmtStrategyName(model Strategy) string {
	id := model.Code()
	return config.QmtStrategyNameFromId(id)
}

// QmtOrderRemark 获取用于QMT系统的订单备注
func QmtOrderRemark(model Strategy) string {
	return model.OrderFlag()
}

var (
	_mutexStrategies        sync.Mutex
	_mapStrategies          = map[ModelKind]Strategy{}
	_mapStrategiesOverwrite = map[ModelKind]bool{}                      //  强制覆盖缓存
	ErrAlreadyExists        = errors.New("the strategy already exists") // 已经存在
	ErrNotFound             = errors.New("the strategy was not found")  // 不存在
)

// Register 注册策略
func Register(strategy Strategy) error {
	_mutexStrategies.Lock()
	defer _mutexStrategies.Unlock()
	strategyCode := strategy.Code()
	// 检查是否存在覆盖策略的情况
	_, overwritten := _mapStrategiesOverwrite[strategyCode]
	if overwritten {
		return nil
	}
	if strategyCode < ModelForceOverwrite {
		_, ok := _mapStrategies[strategyCode]
		if ok {
			return ErrAlreadyExists
		}
	} else {
		strategyCode = strategyCode &^ ModelForceOverwrite
		_mapStrategiesOverwrite[strategyCode] = true
	}
	_mapStrategies[strategyCode] = strategy
	return nil
}

// CheckoutStrategy 捡出策略对象
func CheckoutStrategy(strategyNumber uint64) (Strategy, error) {
	_mutexStrategies.Lock()
	defer _mutexStrategies.Unlock()
	strategy, ok := _mapStrategies[strategyNumber]
	if ok {
		return strategy, nil
	}

	return nil, ErrNotFound
}

// UsageStrategyList 输出策略列表
func UsageStrategyList() string {
	// 规则按照kind排序
	kinds := api.Keys(_mapStrategies)
	slices.Sort(kinds)
	usage := ""
	for _, kind := range kinds {
		if rule, ok := _mapStrategies[kind]; ok {
			usage += fmt.Sprintf("%d: %s\n", kind, rule.Name())
		}
	}
	return usage
}

type StrategySummary struct {
	Type ModelKind
	Name string
}

var (
	MapStrategies = map[ModelKind]StrategySummary{
		ModelZero:    {Type: ModelZero, Name: "0号策略"},
		ModelHousNo1: {Type: ModelHousNo1, Name: "1号策略"},
	}
)
