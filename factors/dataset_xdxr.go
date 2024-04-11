package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/gotdx/quotes"
)

// DataXdxr 除权除息
type DataXdxr struct {
	cache.DataSummary
	Date string
	Code string
}

func init() {
	summary := __mapDataSets[BaseXdxr]
	_ = cache.Register(&DataXdxr{DataSummary: summary})
}

func (x *DataXdxr) Clone(date string, code string) DataSet {
	summary := __mapDataSets[BaseXdxr]
	var dest = DataXdxr{DataSummary: summary, Date: date, Code: code}
	return &dest
}

func (x *DataXdxr) GetDate() string {
	return x.Date
}

func (x *DataXdxr) GetSecurityCode() string {
	return x.Code
}

func (x *DataXdxr) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

//	func (x *DataXdxr) Checkout(securityCode, date string) {
//		//TODO implement me
//		panic("implement me")
//	}
//
//	func (x *DataXdxr) Check(cacheDate, featureDate string) error {
//		//TODO implement me
//		panic("implement me")
//	}

func (x *DataXdxr) Print(code string, date ...string) {
	//TODO implement me
	panic("implement me")
}

func (x *DataXdxr) Filename(date, code string) string {
	//TODO implement me
	_ = code
	_ = date
	panic("implement me")
}

func (x *DataXdxr) Update(date string) {
	base.UpdateXdxrInfo(x.GetSecurityCode())
	_ = date
}

func (x *DataXdxr) Repair(date string) {
	base.UpdateXdxrInfo(x.GetSecurityCode())
	_ = date
}

func (x *DataXdxr) Increase(snapshot quotes.Snapshot) {
	// 除权除息没有增量计算的逻辑
	_ = snapshot
}
