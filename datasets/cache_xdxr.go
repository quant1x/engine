package datasets

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx/quotes"
)

type DataXdxr struct {
	DataCache
}

func init() {
	_ = cache.Register(&DataXdxr{})
}

func (x *DataXdxr) Init(ctx context.Context, date, securityCode string) error {
	_ = ctx
	_ = date
	_ = securityCode
	return nil
}

func (x *DataXdxr) Kind() cache.Kind {
	return BaseXdxr
}

func (x *DataXdxr) Key() string {
	return mapDataSets[x.Kind()].Key()
}

func (x *DataXdxr) Name() string {
	return mapDataSets[x.Kind()].Name()
}

func (x *DataXdxr) Owner() string {
	return mapDataSets[x.Kind()].Owner()
}

func (x *DataXdxr) Usage() string {
	//TODO implement me
	panic("implement me")
}

func (x *DataXdxr) Check(cacheDate, featureDate string) error {
	//TODO implement me
	panic("implement me")
}

func (x *DataXdxr) Print(code string, date ...string) {
	//TODO implement me
	panic("implement me")
}

func (x *DataXdxr) Filename(date, code string) string {
	x.filename = cache.XdxrFilename(x.Code)
	return x.filename
}

func (x *DataXdxr) Update(cacheDate, featureDate string) {
	base.UpdateXdxrInfo(x.Code)
	_ = cacheDate
	_ = featureDate
}

func (x *DataXdxr) Repair(cacheDate, featureDate string) {
	base.UpdateXdxrInfo(x.Code)
	_ = cacheDate
	_ = featureDate
}

func (x *DataXdxr) Increase(snapshot quotes.Snapshot) {
	// 除权除息没有增量计算的逻辑
	_ = snapshot
}

func (x *DataXdxr) Clone(date string, code string) DataSet {
	var dest = DataXdxr{DataCache{Date: date, Code: code}}
	return &dest
}
