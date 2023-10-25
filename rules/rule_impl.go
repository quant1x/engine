package rules

import (
	"gitee.com/quant1x/engine/models"
)

type RuleImpl struct {
	kind Kind
	name string
	exec func(snapshot models.QuoteSnapshot) error
}

func (r RuleImpl) Kind() Kind {
	return r.kind
}

func (r RuleImpl) Name() string {
	return r.name
}

func (r RuleImpl) Exec(snapshot models.QuoteSnapshot) error {
	return r.exec(snapshot)
}
