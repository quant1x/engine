package models

import (
	"errors"
	"fmt"
	"gitee.com/quant1x/gox/util/treemap"
	"golang.org/x/exp/maps"
	"slices"
	"sync"
)

// ModelKind 做多64个策略
type ModelKind = int

const (
	ModelZero ModelKind = 0 // 0号策略
)

const (
	ModelHousNo1 ModelKind = 1 << iota // 1号策略
	ModelTail                          // 尾盘策略
	ModelTick                          // 盘中实时策略
)

const (
	DefaultStrategy = ModelHousNo1
	KLineMin        = 89 // K线最少记录数
)

const (
	OrderFlagHead = "head" // 早盘订单标志
	OrderFlagTick = "tick" // 实时订单标志
	OrderFlagTail = "tail" // 尾盘订单标志
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
	Filter(snapshot QuoteSnapshot) bool
	// Sort 排序
	Sort([]QuoteSnapshot) SortedStatus
	// Evaluate 评估 日线数据
	Evaluate(securityCode string, result *treemap.Map)
}

var (
	_mutexStrategies sync.Mutex
	_mapStrategies   = map[ModelKind]Strategy{}
	ErrAlreadyExists = errors.New("strategy is already exists") // 已经存在
	ErrNotFound      = errors.New("strategy not found")         // 不存在
)

// Register 注册策略
func Register(strategy Strategy) error {
	_mutexStrategies.Lock()
	defer _mutexStrategies.Unlock()
	_, ok := _mapStrategies[strategy.Code()]
	if ok {
		return ErrAlreadyExists
	}
	_mapStrategies[strategy.Code()] = strategy
	return nil
}

// CheckoutStrategy 捡出策略对象
func CheckoutStrategy(strategyNumber int) (Strategy, error) {
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
	kinds := maps.Keys(_mapStrategies)
	slices.Sort(kinds)
	usage := ""
	for _, kind := range kinds {
		if rule, ok := _mapStrategies[kind]; ok {
			usage += fmt.Sprintf("%d: %s\n", rule.Code(), rule.Name())
		}
	}
	return usage
}

type StrategyWrap struct {
	Type ModelKind
	Name string
}

var (
	MapStrategies = map[ModelKind]StrategyWrap{
		ModelZero:    {Type: ModelZero, Name: "0号策略"},
		ModelHousNo1: {Type: ModelHousNo1, Name: "1号策略"},
		ModelTail:    {Type: ModelTail, Name: "尾盘策略"},
	}
)
