package rules

import (
	"gitee.com/quant1x/engine/models"
)

type RuleSubNew struct{}

func (r RuleSubNew) Kind() Kind {
	return RuleSubNewStock
}

func (r RuleSubNew) Name() string {
	return "次新股"
}

func (r RuleSubNew) Exec(snapshot models.QuoteSnapshot) error {
	return ErrExecuteFailed
}

func init() {
	err := Register(RuleSubNew{})
	if err != nil {
		panic(err)
	}
}
