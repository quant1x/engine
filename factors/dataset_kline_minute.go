package factors

import (
	"context"
	"log"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/num"
)

const (
	cacheKeyKLine = "kline"
)

var (
	paramsKLine        = config.GetDataConfig().Cache[cacheKeyKLine]
	minuteKLine        = "1min"
	minuteKLineEnabled = false
)

type KLineMinute struct {
	cache.DataSummary
	Date string
	Code string
}

func init() {
	// 初始化data.cache.kline参数
	n := len(paramsKLine)
	if n > 1 {
		log.Fatal("分钟级K线配置只能是1个")
	}
	for k, v := range paramsKLine {
		enabled, ok := v.(bool)
		if ok {
			_, err := num.ParseFreq(k)
			if err == nil {
				minuteKLineEnabled = enabled
				minuteKLine = k
			}
		}
	}
	if minuteKLineEnabled {
		summary := __mapDataSets[BaseKLineMinute]
		_ = cache.Register(&KLineMinute{DataSummary: summary})
	}
}

func (k *KLineMinute) Clone(date, code string) DataSet {
	summary := __mapDataSets[BaseKLineMinute]
	var dest = KLineMinute{
		DataSummary: summary,
		Date:        date,
		Code:        code,
	}
	return &dest
}

func (k *KLineMinute) GetDate() string {
	return k.Date
}

func (k *KLineMinute) GetSecurityCode() string {
	return k.Code
}

func (k *KLineMinute) Init(ctx context.Context, date string) error {
	//_ = barIndex
	_ = ctx
	_ = date
	return nil
}

func (k *KLineMinute) Checkout(securityCode, date string) {
	//TODO implement me
	_ = securityCode
	_ = date
	panic("implement me")
}

func (k *KLineMinute) Check(cacheDate, featureDate string) error {
	//TODO implement me
	_ = cacheDate
	_ = featureDate
	panic("implement me")
}

func (k *KLineMinute) Filename(date, code string) string {
	//TODO implement me
	_ = code
	_ = date
	panic("implement me")
}

func (k *KLineMinute) Print(code string, date ...string) {
	//TODO implement me
	_ = code
	_ = date
	panic("implement me")
}

func (k *KLineMinute) Update(date string) error {
	base.UpdateAllKLine(k.GetSecurityCode(), minuteKLine)
	_ = date
	return nil
}

func (k *KLineMinute) Repair(date string) error {
	base.UpdateAllKLine(k.GetSecurityCode(), minuteKLine)
	_ = date
	return nil
}

func (k *KLineMinute) Increase(snapshot quotes.Snapshot) error {
	//TODO K线增量更新数据的条件是缓存的数据最晚的日期是上一个交易日, 否则会缺失缓存数据中最后1条数据和当日之间的数据, 所以只能按照K线的更新方法, 不适合用快照更新
	// 第一步: 取出最后一条数据的记录
	// 第二步: 检查时间的有效性
	// 第三步: 用快照组织K线结构
	// 第四步: 如果不符合快照更新, 则忽略
	_ = snapshot
	panic("implement me")
}
