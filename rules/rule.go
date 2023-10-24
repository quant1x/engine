package rules

import (
	"errors"
	"gitee.com/quant1x/engine/models"
	bitmap "github.com/bits-and-blooms/bitset"
	"golang.org/x/exp/maps"
	"slices"
	"sync"
)

// Kind 规则类型
type Kind = uint

const (
	Pass Kind = 0
)

const (
	RuleMiss   Kind = iota //规则未命中
	RuleHit                // 命中
	RuleCancel             // 撤回
	RulePassed             // 成功
	RuleFailed             // 失败
)

// Rule 规则接口
type Rule interface {
	Kind() Kind                               // 类型
	Name() string                             // 名称
	Exec(snapshot models.QuoteSnapshot) error // 执行
}

var (
	mutex    sync.RWMutex
	mapRules = map[Kind]Rule{}
)

var (
	ErrAlreadyExists = errors.New("rule is already exists") // 规则已经存在
	ErrExecuteFailed = errors.New("rule execute failed")    // 规则执行失败
)

// Register 注册规则
func Register(rule Rule) error {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := mapRules[rule.Kind()]
	if ok {
		return ErrAlreadyExists
	}
	mapRules[rule.Kind()] = rule
	return nil
}

// Each 遍历所有规则
func Each(snapshot models.QuoteSnapshot) (passed []uint64, failed Kind) {
	mutex.RLock()
	defer mutex.RUnlock()
	var bitset bitmap.BitSet
	// 规则排序
	kinds := maps.Keys(mapRules)
	slices.Sort(kinds)
	for _, kind := range kinds {
		if rule, ok := mapRules[kind]; ok {
			err := rule.Exec(snapshot)
			if err != nil {
				failed = rule.Kind()
				break
			}
			bitset.Set(rule.Kind())
		}
	}
	return bitset.Bytes(), failed
}
