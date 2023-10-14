package rules

// RuleKind 规则类型
type RuleKind int

const (
	RuleMiss RuleKind = 0 // 规则未命中
)

const (
	RuleHit    RuleKind = 1 << iota // 命中
	RuleCancel                      // 撤回
	RulePassed                      // 成功
	RuleFailed                      // 失败
)

// Rule 规则接口
type Rule interface {
	Kind()
	Name()
}
