package factors

import (
	"fmt"
	"strings"
	"sync"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gox/util/treemap"
)

const (
	// 闪存路径
	cache1dPrefix = "flash"
	// 默认证券代码为上证指数
	defaultSecurityCode = "sh000001"
)

// 返回文件路径, 指定关键字和日期
func getCache1DFilepath(key, date string) string {
	cachePath, key, found := strings.Cut(key, "/")
	if !found {
		key = cachePath
		cachePath = cache1dPrefix
	}
	cachePath = cache.GetRootPath() + "/" + cachePath
	year := date[:4]
	filename := fmt.Sprintf("%s/%s/%s.%s", cachePath, year, key, date)
	return filename
}

// FeatureRotationAdapter 特征缓存日旋转适配器
//
//	一天一个特征组合缓存文件
type FeatureRotationAdapter interface {
	// Kind 获取缓存类型
	Kind() cache.Kind
	// Name 名称
	Name() string
	// Element 获取指定日期的特征
	Element(securityCode string, date ...string) Feature
	// Checkout 加载指定日期的缓存
	Checkout(date ...string)
	// Merge 合并数据
	Merge(p *treemap.Map)
	// Factory 工厂
	Factory(date, securityCode string) Feature
}

var (
	__mutexFeatureRotationAdapters sync.Mutex
	__mapFeatureRotationAdapters   = map[string]FeatureRotationAdapter{}
)

func RegisterFeatureRotationAdapter(key string, adapter FeatureRotationAdapter) {
	__mutexFeatureRotationAdapters.Lock()
	defer __mutexFeatureRotationAdapters.Unlock()
	__mapFeatureRotationAdapters[key] = adapter
}

// SwitchDate 统一切换数据的缓存日期
func SwitchDate(date string) {
	__mutexFeatureRotationAdapters.Lock()
	defer __mutexFeatureRotationAdapters.Unlock()
	for _, v := range __mapFeatureRotationAdapters {
		v.Checkout(date)
	}
}

func Get(key string) FeatureRotationAdapter {
	__mutexFeatureRotationAdapters.Lock()
	defer __mutexFeatureRotationAdapters.Unlock()
	return __mapFeatureRotationAdapters[key]
}
