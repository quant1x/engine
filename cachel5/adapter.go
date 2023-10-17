package cachel5

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gox/util/treemap"
	"strings"
	"sync"
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

// CacheAdapter 缓存适配器
//
//	一天一个特征组合缓存文件
type CacheAdapter interface {
	// Name 名称
	Name() string
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
