package cachel5

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/features"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/util/treemap"
	"strings"
	"sync"
)

// 返回文件路径, 指定关键字和日期
func getCache1DFilepath(key, date string) string {
	cachePath, key, found := strings.Cut(key, "/")
	if !found {
		key = cachePath
		cachePath = "cache"
	}
	cachePath = cache.GetRootPath() + "/" + cachePath
	year := date[:4]
	filename := fmt.Sprintf("%s/%s/%s.%s", cachePath, year, key, date)
	return filename
}

// Cache1D 每天1个证券代码1条数据
type Cache1D[T features.Feature] struct {
	once        coroutine.RollingMutex
	m           sync.RWMutex
	factory     func(date, securityCode string) T
	Key         string // 缓存关键字
	Date        string // 日期
	filename    string // 缓存文件名
	mapCache    map[string]T
	replaceDate string // 替换缓存的日期
	allCodes    []string
}

// NewCache1D 创建一个新的C1D对象
//
//	key支持多级相对路径, 比如a/b, 创建的路径是~/.quant1x/a/b.yyyy-mm-dd
func NewCache1D[T features.Feature](key string, factory func(date, securityCode string) T) *Cache1D[T] {
	d1 := Cache1D[T]{
		Key:         key,
		Date:        "",
		factory:     factory,
		mapCache:    map[string]T{},
		replaceDate: "",
		allCodes:    []string{},
	}
	d1.Date = cache.DefaultCanReadDate()
	d1.allCodes = market.GetCodeList()
	RegisterCacheLoader(key, &d1)
	return &d1
}

func (this *Cache1D[T]) Factory(date, securityCode string) features.Feature {
	return this.factory(date, securityCode)
}

func (this *Cache1D[T]) Name() string {
	var t T
	return t.FeatureName()
}

// Length 获取长度
func (this *Cache1D[T]) Length() int {
	return len(this.allCodes)
}

// loadCache 加载指定日期的数据
func (this *Cache1D[T]) loadCache(date string) {
	this.allCodes = market.GetCodeList()
	this.Date = trading.FixTradeDate(date)
	this.filename = getCache1DFilepath(this.Key, this.Date)
	var list []T
	err := api.CsvToSlices(this.filename, &list)
	if err != nil || len(list) == 0 {
		logger.Errorf("%s 没有有效数据, error=%+v", this.filename, err)
		return
	}
	for _, v := range list {
		code := v.GetSecurityCode()
		this.mapCache[code] = v
	}
}

// 加载默认数据, 日期为当前交易中的日期
func (this *Cache1D[T]) loadDefault() {
	//this.Date = DefaultCanReadDate()
	this.loadCache(this.Date)
}

// ReplaceCache 替换当前缓存数据
func (this *Cache1D[T]) ReplaceCache() {
	clear(this.mapCache)
	this.loadCache(this.replaceDate)
}

func (this *Cache1D[T]) Checkout(date ...string) {
	if len(date) > 0 {
		this.m.Lock()
		destDate := trading.FixTradeDate(date[0])
		if this.Date != destDate {
			this.replaceDate = destDate
		}
		this.m.Unlock()
	}
	if len(this.replaceDate) == 0 || this.Date == this.replaceDate {
		this.once.Do(this.loadDefault)
	} else {
		// 重置once锁计数器为0
		this.once.Reset()
		this.once.Do(this.ReplaceCache)
	}
}

// Get 获取指定证券代码的数据
func (this *Cache1D[T]) Get(securityCode string, date ...string) *T {
	this.Checkout(date...)
	this.once.Do(this.loadDefault)
	t, ok := this.mapCache[securityCode]
	if ok {
		return &t
	}
	return nil
}

// Set 更新map中指定证券代码的数据
func (this *Cache1D[T]) Set(securityCode string, newValue T, date ...string) {
	this.Checkout(date...)
	this.once.Do(this.loadDefault)
	this.mapCache[securityCode] = newValue
}

// Apply 数据合并
//
//	泛型T需要保持一个string类型的Date字段
func (this *Cache1D[T]) Apply(merge func(code string, local *T) (updated bool)) {
	list := []T{}
	for _, securityCode := range this.allCodes {
		v, found := this.mapCache[securityCode]
		if !found && this.factory != nil {
			v = this.factory(this.Date, securityCode)
		}
		if merge != nil {
			ok := merge(securityCode, &v)
			if ok {
				this.mapCache[securityCode] = v
			}
		}
		list = append(list, v)
	}
	if len(list) > 0 {
		err := api.SlicesToCsv(this.filename, list)
		if err != nil {
			logger.Errorf("刷新%s异常:%+v", this.filename, err)
		}
	}
}

func (this *Cache1D[T]) Merge(p *treemap.Map) {
	list := []T{}
	for _, securityCode := range this.allCodes {
		v, found := this.mapCache[securityCode]
		if !found && this.factory != nil {
			v = this.factory(this.Date, securityCode)
		}
		if p != nil {
			tmp, ok := p.Get(securityCode)
			if ok {
				_ = api.CopyWithOption(v, tmp, api.Option{})
				if ok {
					this.mapCache[securityCode] = v
				}
			}
		}
		list = append(list, v)
	}
	if len(list) > 0 {
		err := api.SlicesToCsv(this.filename, list)
		if err != nil {
			logger.Errorf("刷新%s异常:%+v", this.filename, err)
		}
	}
}
