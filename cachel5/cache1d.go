package cachel5

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/concurrent"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/gox/util/treemap"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
	"sync"
)

// Cache1D 缓存所有证券代码的特征组合数据
//
//	每天1个证券代码1条数据
type Cache1D[T factors.Feature] struct {
	once     coroutine.PeriodicOnce
	m        sync.RWMutex
	factory  func(date, securityCode string) T
	cacheKey string // 缓存关键字
	Date     string // 日期
	filename string // 缓存文件名
	//mapCache    map[string]T
	//mapCache    concurrent.ConcurrentHashMap[string, T]
	mapCache    *concurrent.TreeMap[string, T]
	replaceDate string // 替换缓存的日期
	allCodes    []string
	tShadow     T // 泛型T的影子
}

// NewCache1D 创建一个新的C1D对象
//
//	key支持多级相对路径, 比如a/b, 创建的路径是~/.quant1x/a/b.yyyy-mm-dd
func NewCache1D[T factors.Feature](key string, factory func(date, securityCode string) T) *Cache1D[T] {
	d1 := &Cache1D[T]{
		cacheKey: key,
		Date:     "",
		factory:  factory,
		//mapCache:    map[string]T{},
		//mapCache:    concurrent.NewHashMap[string, T](),
		mapCache:    concurrent.NewTreeMap[string, T](),
		replaceDate: "",
		allCodes:    []string{},
	}
	d1.Date = cache.DefaultCanReadDate()
	d1.allCodes = market.GetCodeList()
	//d1.Checkout(d1.Date)
	d1.tShadow = d1.factory(d1.Date, defaultSecurityCode)
	RegisterCacheLoader(key, d1)
	return d1
}

func (this *Cache1D[T]) Factory(date, securityCode string) factors.Feature {
	return this.tShadow.Factory(date, securityCode)
}

func (this *Cache1D[T]) Init(ctx context.Context, date, securityCode string) error {
	_ = ctx
	_ = date
	_ = securityCode
	return nil
}

func (this *Cache1D[T]) Owner() string {
	return this.tShadow.Owner()
}

func (this *Cache1D[T]) Kind() cache.Kind {
	return this.tShadow.Kind()
}

func (this *Cache1D[T]) Key() string {
	return this.tShadow.Key()
}

func (this *Cache1D[T]) Name() string {
	return this.tShadow.Name()
}

func (this *Cache1D[T]) Usage() string {
	return this.tShadow.Usage()
}

// Length 获取长度
func (this *Cache1D[T]) Length() int {
	return len(this.allCodes)
}

// loadCache 加载指定日期的数据
func (this *Cache1D[T]) loadCache(date string) {
	this.allCodes = market.GetCodeList()
	this.Date = trading.FixTradeDate(date)
	this.filename = getCache1DFilepath(this.cacheKey, this.Date)
	var list []T
	err := api.CsvToSlices(this.filename, &list)
	if err != nil || len(list) == 0 {
		logger.Errorf("%s 没有有效数据, error=%+v", this.filename, err)
		return
	}
	for _, v := range list {
		code := v.GetSecurityCode()
		//this.mapCache[code] = v
		//this.mapCache.Set(code, v)
		this.mapCache.Put(code, v)
	}
}

// 加载默认数据, 日期为当前交易中的日期
func (this *Cache1D[T]) loadDefault() {
	this.loadCache(this.Date)
}

// ReplaceCache 替换当前缓存数据
func (this *Cache1D[T]) ReplaceCache() {
	this.mapCache.Clear()
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

func checkoutTable(v any) (headers []string, records [][]string) {
	headers = []string{"字段", "数值"}
	fields := tags.GetHeadersByTags(v)
	values := tags.GetValuesByTags(v)
	num := len(fields)
	if num > len(values) {
		num = len(values)
	}
	for i := 0; i < num; i++ {
		records = append(records, []string{fields[i], strings.TrimSpace(values[i])})
	}
	return
}

func (this *Cache1D[T]) Print(code string, date ...string) {
	securityCode := proto.CorrectSecurityCode(code)
	//name := securities.GetStockName(securityCode)
	//tradeDate = trading.FixTradeDate(tradeDate)
	tradeDate := cache.DefaultCanReadDate()
	if len(date) > 0 {
		tradeDate = trading.FixTradeDate(date[0])
	}
	//fmt.Printf("%s: %s, %s\n", securityCode, name, tradeDate)
	value := this.Get(securityCode, tradeDate)
	if value != nil {
		headers, records := checkoutTable(*value)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(headers)
		table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
		table.AppendBulk(records)
		table.Render()
	}
}

func (this *Cache1D[T]) Check(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

// Get 获取指定证券代码的数据
func (this *Cache1D[T]) Get(securityCode string, date ...string) *T {
	this.Checkout(date...)
	this.once.Do(this.loadDefault)
	t, ok := this.mapCache.Get(securityCode)
	if ok {
		return &t
	}
	return nil
}

// Set 更新map中指定证券代码的数据
func (this *Cache1D[T]) Set(securityCode string, newValue T, date ...string) {
	this.Checkout(date...)
	this.once.Do(this.loadDefault)
	this.mapCache.Put(securityCode, newValue)
}

// Apply 数据合并
//
//	泛型T需要保持一个string类型的Date字段
func (this *Cache1D[T]) Apply(merge func(code string, local *T) (updated bool)) {
	list := []T{}
	for _, securityCode := range this.allCodes {
		v, found := this.mapCache.Get(securityCode)
		if !found && this.factory != nil {
			v = this.factory(this.Date, securityCode)
		}
		if merge != nil {
			ok := merge(securityCode, &v)
			if ok {
				this.mapCache.Put(securityCode, v)
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
		v, found := this.mapCache.Get(securityCode)
		if !found && this.factory != nil {
			v = this.factory(this.Date, securityCode)
		}
		if p != nil {
			tmp, ok := p.Get(securityCode)
			if ok {
				_ = api.CopyWithOption(v, tmp, api.Option{})
				if ok {
					this.mapCache.Put(securityCode, v)
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
