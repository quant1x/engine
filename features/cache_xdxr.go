package features

import "gitee.com/quant1x/gotdx/quotes"

type DataXdxr struct {
	DataCache
}

func (x *DataXdxr) Kind() FeatureKind {
	return FeatureBaseXdxr
}

func (x *DataXdxr) Name() string {
	return mapFeatures[x.Kind()].Name
}

func (x *DataXdxr) Key() string {
	return mapFeatures[x.Kind()].Key
}

func (x *DataXdxr) Filename(date, code string) string {
	//TODO implement me
	panic("implement me")
}

func (x *DataXdxr) Update(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (x *DataXdxr) Repair(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (x *DataXdxr) Increase(snapshot quotes.Snapshot) {
	//TODO implement me
	panic("implement me")
}

func (x *DataXdxr) Clone(date string, code string) DataSet {
	var dest = DataXdxr{DataCache{Date: date, Code: code}}
	return &dest
}
