package flash

import (
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/datasets"
	"gitee.com/quant1x/engine/features"
)

var (
	// 历史数据
	__l5History = cachel5.NewCache1D[*features.History](features.CacheL5KeyHistory, features.NewHistory)
)

// CacheList 缓存列表
func CacheList() []cachel5.CacheAdapter {
	list := []cachel5.CacheAdapter{__l5History}
	return list
}

func GetHistory() *cachel5.Cache1D[*features.History] {
	return __l5History
}

func GetL5History(securityCode string, date ...string) *features.History {
	data := __l5History.Get(securityCode, date...)
	if data == nil {
		return nil
	}
	return *data
}

// DataSetList 数据集 列表
func DataSetList() []datasets.DataSet {
	list := []datasets.DataSet{
		new(datasets.DataXdxr),  // 除权除息
		new(datasets.DataKLine), // 基础K线
	}
	return list
}
