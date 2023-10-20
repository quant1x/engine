package cachel5

import (
	"context"
	"gitee.com/quant1x/engine/cache"
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

func (f *FeatureN01) Key() string {
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

func (f *FeatureN01) Update(securityCode, cacheDate, featureDate string, whole bool) {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Repair(securityCode, cacheDate, featureDate string, whole bool) {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) FromHistory(history factors.History) factors.Feature {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Increase(snapshot quotes.Snapshot) factors.Feature {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Kind() cache.Kind {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Owner() string {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Name() string {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Usage() string {
	//TODO implement me
	panic("implement me")
}

func (f *FeatureN01) Init(ctx context.Context, date string) error {
	//TODO implement me
	panic("implement me")
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
