package rules

import (
	"errors"
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/runtime"
	bitmap "github.com/bits-and-blooms/bitset"
	"slices"
	"sync"
)

// Kind 规则类型
type Kind = uint

const (
	Pass Kind = 0
)

const (
	engineBaseRule Kind = 1
	KRuleF10            = engineBaseRule + 0 // 基础规则
	KRuleBase           = engineBaseRule + 1 // 基础规则
)

// 规则错误码, 每一组规则错误拟1000个错误码
const (
	errorRuleF10  = (iota + 1) * 1000 // F10错误码
	errorRuleBase                     // 基础规则错误码
)

// Rule 规则接口封装
type Rule struct {
	kind Kind
	name string
	Exec func(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error
}

func (this Rule) Kind() Kind {
	return this.kind
}

func (this Rule) Name() string {
	return this.name
}

var (
	mutex    sync.RWMutex
	mapRules = map[Kind]Rule{}
)

var (
	ErrAlreadyExists = errors.New("the rule already exists") // 规则已经存在
)

// RegisterFunc 注册规则回调函数
func RegisterFunc(kind Kind, name string, cb func(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error) error {
	rule := Rule{kind: kind, name: name, Exec: cb}
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := mapRules[rule.Kind()]
	if ok {
		return ErrAlreadyExists
	}
	mapRules[rule.Kind()] = rule
	return nil
}

// Filter 遍历所有规则
func Filter(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) (passed []uint64, failed Kind, err error) {
	mutex.RLock()
	defer mutex.RUnlock()
	if len(mapRules) == 0 {
		return
	}
	var bitset bitmap.BitSet
	// 规则按照kind排序
	kinds := api.Keys(mapRules)
	slices.Sort(kinds)
	for _, kind := range kinds {
		if rule, ok := mapRules[kind]; ok {
			err = rule.Exec(ruleParameter, snapshot)
			if err != nil {
				failed = rule.Kind()
				break
			}
			bitset.Set(rule.Kind())
		}
	}
	return bitset.Bytes(), failed, err
}

// PrintRuleList 输出规则列表
func PrintRuleList() {
	fmt.Println("规则总数:", len(mapRules))
	// 规则按照kind排序
	kinds := api.Keys(mapRules)
	slices.Sort(kinds)
	for _, kind := range kinds {
		if rule, ok := mapRules[kind]; ok {
			fmt.Printf("kind: %d, name: %s, method: %s\n", rule.Kind(), rule.Name(), runtime.FuncName(rule.Exec))
		}
	}
}
