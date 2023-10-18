package factors

import "gitee.com/quant1x/engine/cache"

// featureManifest 特征的基础数据
type featureManifest struct {
	Date     string // 日期
	Code     string // 证券代码
	filename string // 文件名
	kind     cache.Kind
}

func (d featureManifest) Kind() cache.Kind {
	return d.kind
}

func (d featureManifest) Owner() string {
	return mapFeatures[d.Kind()].Owner()
}

func (d featureManifest) Key() string {
	return mapFeatures[d.Kind()].Key()
}

func (d featureManifest) Name() string {
	return mapFeatures[d.Kind()].Name()
}

func (d featureManifest) Usage() string {
	return mapFeatures[d.Kind()].Usage()
}

func (d featureManifest) GetDate() string {
	return d.Date
}

func (d featureManifest) GetSecurityCode() string {
	return d.Code
}
