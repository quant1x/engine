package datasets

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx/quotes"
)

type DataXdxr struct {
	Manifest
}

func init() {
	_ = cache.Register(&DataXdxr{Manifest: Manifest{Kind_: BaseXdxr}})
}

func (x *DataXdxr) Clone(date string, code string) DataSet {
	manifest := Manifest{Date: date, Code: code, Kind_: BaseXdxr}
	var dest = DataXdxr{Manifest: manifest}
	return &dest
}

func (x *DataXdxr) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (x *DataXdxr) Checkout(securityCode, date string) {
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

func (x *DataXdxr) Update(date string) {
	base.UpdateXdxrInfo(x.Code)
	_ = date
}

func (x *DataXdxr) Repair(date string) {
	base.UpdateXdxrInfo(x.Code)
	_ = date
}

func (x *DataXdxr) Increase(snapshot quotes.Snapshot) {
	// 除权除息没有增量计算的逻辑
	_ = snapshot
}
