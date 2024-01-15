package factors

import (
	"gitee.com/quant1x/engine/cache"
	"sync"
)

var (
	__l5Once sync.Once
	// 历史数据
	__l5History *Cache1D[*History] = nil
	// 基本面F10
	__l5F10 *Cache1D[*F10] = nil
	// 扩展交易特征
	__l5Misc *Cache1D[*Misc] = nil
	// 平台
	__l5Box *Cache1D[*Box] = nil
)

func init() {
	__l5Once.Do(lazyInitFeatures)
}

func lazyInitFeatures() {
	// 历史数据
	__l5History = NewCache1D[*History](cacheL5KeyHistory, NewHistory)
	err := cache.Register(__l5History)
	if err != nil {
		panic(err)
	}
	// 基本面F10
	__l5F10 = NewCache1D[*F10](cacheL5KeyF10, NewF10)
	err = cache.Register(__l5F10)
	if err != nil {
		panic(err)
	}
	// 扩展信息
	__l5Misc = NewCache1D[*Misc](cacheL5KeyMisc, NewMisc)
	err = cache.Register(__l5Misc)
	if err != nil {
		panic(err)
	}
	// 平台
	__l5Box = NewCache1D[*Box](cacheL5KeyBox, NewBox)
	err = cache.Register(__l5Box)
	if err != nil {
		panic(err)
	}
}

func GetL5History(securityCode string, date ...string) *History {
	__l5Once.Do(lazyInitFeatures)
	data := __l5History.Get(securityCode, date...)
	if data == nil {
		return nil
	}
	return *data
}

func GetL5F10(securityCode string, date ...string) *F10 {
	__l5Once.Do(lazyInitFeatures)
	data := __l5F10.Get(securityCode, date...)
	if data == nil {
		return nil
	}
	return *data
}

// GetL5Misc 获取扩展信息
func GetL5Misc(securityCode string, date ...string) (exchange *Misc) {
	__l5Once.Do(lazyInitFeatures)
	v := __l5Misc.Get(securityCode, date...)
	if v == nil {
		return nil
	}
	return *v
}

// UpdateL5Misc 更新当日exchange
func UpdateL5Misc(misc *Misc) {
	__l5Once.Do(lazyInitFeatures)
	__l5Misc.Set(misc.Code, misc, cache.DefaultCanReadDate())
}

// RefreshL5Misc 刷新缓存
func RefreshL5Misc() {
	__l5Once.Do(lazyInitFeatures)
	__l5Misc.Apply(nil)
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
