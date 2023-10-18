package datasets

import "gitee.com/quant1x/engine/cache"

// dataManifest 基础的数据缓存
type dataManifest struct {
	Date     string // 日期
	Code     string // 证券代码
	filename string // 文件名
	kind     cache.Kind
}

func (d dataManifest) Kind() cache.Kind {
	return d.kind
}

func (d dataManifest) Owner() string {
	return mapDataSets[d.Kind()].Owner()
}

func (d dataManifest) Key() string {
	return mapDataSets[d.Kind()].Key()
}

func (d dataManifest) Name() string {
	return mapDataSets[d.Kind()].Name()
}

func (d dataManifest) Usage() string {
	return mapDataSets[d.Kind()].Usage()
}

func (d dataManifest) GetDate() string {
	return d.Date
}

func (d dataManifest) GetSecurityCode() string {
	return d.Code
}
