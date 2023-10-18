package factors

import "gitee.com/quant1x/engine/cache"

// Manifest 特征的基础数据
type Manifest struct {
	Date     string // 日期
	Code     string // 证券代码
	filename string // 文件名
	kind     cache.Kind
}

func (d Manifest) Kind() cache.Kind {
	return d.kind
}

func (d Manifest) Owner() string {
	return mapFeatures[d.Kind()].Owner()
}

func (d Manifest) Key() string {
	return mapFeatures[d.Kind()].Key()
}

func (d Manifest) Name() string {
	return mapFeatures[d.Kind()].Name()
}

func (d Manifest) Usage() string {
	return mapFeatures[d.Kind()].Usage()
}

func (d Manifest) GetDate() string {
	return d.Date
}

func (d Manifest) GetSecurityCode() string {
	return d.Code
}
