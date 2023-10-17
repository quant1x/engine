package cachel5

import (
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gox/util/treemap"
	"sync"
)

// CacheAdapter 缓存加载器
type CacheAdapter interface {
	// Name 名称
	//Name() string

	// Checkout 加载指定日期的缓存
	Checkout(date ...string)
	// Merge 合并数据
	Merge(p *treemap.Map)
	// Factory 工厂
	Factory(date, securityCode string) factors.Feature
}

var (
	__loadMutex sync.Mutex
	__mapLoader = map[string]CacheAdapter{}
)

func RegisterCacheLoader(key string, loader CacheAdapter) {
	__loadMutex.Lock()
	defer __loadMutex.Unlock()
	__mapLoader[key] = loader
}

// SwitchDate 统一切换数据的缓存日期
func SwitchDate(date string) {
	__loadMutex.Lock()
	defer __loadMutex.Unlock()
	for _, v := range __mapLoader {
		v.Checkout(date)
	}
}

func Get(key string) CacheAdapter {
	return __mapLoader[key]
}
