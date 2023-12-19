package features

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"sync"
)

var (
	__l5Once sync.Once
	// 扩展交易特征
	__l5Exchange *cachel5.Cache1D[*Exchange] = nil
	// 平台
	__l5Box *cachel5.Cache1D[*Box] = nil
)

func init() {
	__l5Once.Do(lazyInitFeatures)
}

func lazyInitFeatures() {
	// 扩展信息
	__l5Exchange = cachel5.NewCache1D[*Exchange](CacheL5KeyExchange, NewExchange)
	err := cache.Register(__l5Exchange)
	if err != nil {
		panic(err)
	}
	// 平台
	__l5Box = cachel5.NewCache1D[*Box](CacheL5KeyBox, NewBox)
	err = cache.Register(__l5Box)
	if err != nil {
		panic(err)
	}

}

func CacheExchange() *cachel5.Cache1D[*Exchange] {
	__l5Once.Do(lazyInitFeatures)
	return __l5Exchange
}

// GetL5Exchange 获取扩展信息
func GetL5Exchange(securityCode string, date ...string) (exchange *Exchange) {
	__l5Once.Do(lazyInitFeatures)
	v := __l5Exchange.Get(securityCode, date...)
	if v == nil {
		return nil
	}
	return *v
}

// UpdateL5Exchange 更新当日exchange
func UpdateL5Exchange(extension *Exchange) {
	__l5Once.Do(lazyInitFeatures)
	__l5Exchange.Set(extension.Code, extension, cache.DefaultCanReadDate())
}

// RefreshL5Exchange 刷新缓存
func RefreshL5Exchange() {
	__l5Once.Do(lazyInitFeatures)
	__l5Exchange.Apply(nil)
}

func CacheBox() *cachel5.Cache1D[*Box] {
	__l5Once.Do(lazyInitFeatures)
	return __l5Box
}

// GetL5Box 获取平台数据
func GetL5Box(securityCode string, date ...string) *Box {
	__l5Once.Do(lazyInitFeatures)
	v := __l5Box.Get(securityCode, date...)
	if v == nil {
		return nil
	}
	return *v
}
