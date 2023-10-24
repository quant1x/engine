package rules

import (
	"gitee.com/quant1x/engine/models"
)

type RSubNew struct{}

func (r RSubNew) Kind() Kind {
	return RuleSubnew
}

func (r RSubNew) Name() string {
	return "次新股"
}

func (r RSubNew) Exec(snapshot models.QuoteSnapshot) error {
	return ErrExecuteFailed
}

func init() {
	err := Register(RSubNew{})
	if err != nil {
		panic(err)
	}
}
