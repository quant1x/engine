package rules

import (
	"errors"
	"fmt"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/runtime"
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
	engineBaseRule  Kind = 1
	RuleBaseF10          = engineBaseRule + 0 // 基础规则
	RuleSubNewStock      = engineBaseRule + 1 // 次新股
)

const (
	errorRuleBase = 1000 // 基础规则错误码
)

// Rule 规则接口
type Rule interface {
	// Kind 类型
	Kind() Kind
	// Name 名称
	Name() string
	// Exec 执行, 返回nil即为成功
	Exec(snapshot models.QuoteSnapshot) error
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

// RegisterFunc 注册规则回调函数
func RegisterFunc(kind Kind, name string, cb func(snapshot models.QuoteSnapshot) error) error {
	rule := RuleImpl{kind: kind, name: name, exec: cb}
	return Register(rule)
}

// Filter 遍历所有规则
func Filter(snapshot models.QuoteSnapshot) (passed []uint64, failed Kind, err error) {
	mutex.RLock()
	defer mutex.RUnlock()
	if len(mapRules) == 0 {
		return
	}
	var bitset bitmap.BitSet
	// 规则按照kind排序
	kinds := maps.Keys(mapRules)
	slices.Sort(kinds)
	for _, kind := range kinds {
		if rule, ok := mapRules[kind]; ok {
			err = rule.Exec(snapshot)
			if err != nil {
				failed = rule.Kind()
				break
			}
			bitset.Set(rule.Kind())
		}
	}
	return bitset.Bytes(), failed, err
}

func PrintRuleList() {
	fmt.Println("规则总数:", len(mapRules))
	// 规则按照kind排序
	kinds := maps.Keys(mapRules)
	slices.Sort(kinds)
	for _, kind := range kinds {
		if rule, ok := mapRules[kind]; ok {
			fmt.Printf("kind: %d, name: %s, method: %s\n", rule.Kind(), rule.Name(), runtime.FuncName(rule.Exec))
		}
	}
}
