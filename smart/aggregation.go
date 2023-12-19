package smart

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/factors"
	"sync"
)

var (
	__l5Once sync.Once
	// 历史数据
	__l5History *cachel5.Cache1D[*factors.History] = nil
	// 基本面F10
	__l5F10 *cachel5.Cache1D[*factors.F10] = nil
)

func init() {
	__l5Once.Do(lazyInitFeatures)
}

func lazyInitFeatures() {
	// 历史数据
	__l5History = cachel5.NewCache1D[*factors.History](factors.CacheL5KeyHistory, factors.NewHistory)
	err := cache.Register(__l5History)
	if err != nil {
		panic(err)
	}
	// 基本面F10
	__l5F10 = cachel5.NewCache1D[*factors.F10](factors.CacheL5KeyF10, factors.NewF10)
	err = cache.Register(__l5F10)
	if err != nil {
		panic(err)
	}
}

func CacheHistory() *cachel5.Cache1D[*factors.History] {
	__l5Once.Do(lazyInitFeatures)
	return __l5History
}

func GetL5History(securityCode string, date ...string) *factors.History {
	__l5Once.Do(lazyInitFeatures)
	data := __l5History.Get(securityCode, date...)
	if data == nil {
		return nil
	}
	return *data
}

func CacheF10() *cachel5.Cache1D[*factors.F10] {
	__l5Once.Do(lazyInitFeatures)
	return __l5F10

}

func GetL5F10(securityCode string, date ...string) *factors.F10 {
	__l5Once.Do(lazyInitFeatures)
	data := __l5F10.Get(securityCode, date...)
	if data == nil {
		return nil
	}
	return *data
}
