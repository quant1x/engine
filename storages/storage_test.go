package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/concurrent"
	"path/filepath"
	"testing"
)

type TestModel82 struct{}

func (TestModel82) Code() models.ModelKind {
	return 0
}

func (TestModel82) Name() string {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) OrderFlag() string {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) Filter(ruleParameter config.RuleParameter, snapshot factors.QuoteSnapshot) error {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) Sort(snapshots []factors.QuoteSnapshot) models.SortedStatus {
	//TODO implement me
	panic("implement me")
}

func (TestModel82) Evaluate(securityCode string, result *concurrent.TreeMap[string, models.ResultInfo]) {
	//TODO implement me
	panic("implement me")
}

func TestFilePathClean(t *testing.T) {
	s := "d:\\quant\\data/qmt/var/20231225/20231225-8881479758-s0-b-*.done"
	fmt.Println(filepath.Clean(s))
}
