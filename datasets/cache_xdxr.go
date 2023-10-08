package datasets

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx/quotes"
)

type DataXdxr struct {
	DataCache
}

func (x *DataXdxr) Kind() DataKind {
	return BaseXdxr
}

func (x *DataXdxr) Name() string {
	return mapDataSets[x.Kind()].Name
}

func (x *DataXdxr) Key() string {
	return mapDataSets[x.Kind()].Key
}

func (x *DataXdxr) Filename(date, code string) string {
	x.filename = cache.XdxrFilename(x.Code)
	return x.filename
}

func (x *DataXdxr) Update(cacheDate, featureDate string) {
	base.UpdateXdxrInfo(x.Code)
}

func (x *DataXdxr) Repair(cacheDate, featureDate string) {
	base.UpdateXdxrInfo(x.Code)
}

func (x *DataXdxr) Increase(snapshot quotes.Snapshot) {
	// 除权除息没有增量计算的逻辑
}

func (x *DataXdxr) Clone(date string, code string) DataSet {
	var dest = DataXdxr{DataCache{Date: date, Code: code}}
	return &dest
}