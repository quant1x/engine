package cachel5

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/util/treemap"
	"testing"
)

func TestCache1D_Apply(t *testing.T) {
	c1d := NewCache1D[*FeatureN01](n01, NewFeatureNo1)
	m := treemap.NewWithStringComparator()
	c1d.Merge(m)
}

const (
	n01 = "n01"
)

type FeatureN01 struct {
	Date string
	Code string
}

func NewFeatureNo1(date string, code string) *FeatureN01 {
	v := FeatureN01{
		Date: date,
		Code: code,
	}
	return &v
}

func (f *FeatureN01) Factory(date string, code string) factors.Feature {
	v := FeatureN01{
		Date: date,
		Code: code,
	}
	return &v
}

func (f *FeatureN01) Kind() factors.FeatureKind {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) FeatureName() string {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Key() string {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Init(barIndex *int, date string) error {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) GetDate() string {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) GetSecurityCode() string {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) FromHistory(history factors.History) factors.Feature {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Update(code, cacheDate, featureDate string, complete bool) {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Repair(code, cacheDate, featureDate string, complete bool) {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Increase(snapshot quotes.Snapshot) factors.Feature {
	//TODO implement me
	panic("implement me")
}
