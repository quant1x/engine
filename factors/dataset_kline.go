package factors

import (
	"context"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/gotdx/quotes"
)

const (
	klineFreq = "1d"
)

type DataKLine struct {
	cache.DataSummary
	Date string
	Code string
}

func init() {
	summary := __mapDataSets[BaseKLine]
	_ = cache.Register(&DataKLine{DataSummary: summary})
}

func (k *DataKLine) Clone(date, code string) DataSet {
	summary := __mapDataSets[BaseKLine]
	var dest = DataKLine{
		DataSummary: summary,
		Date:        date,
		Code:        code,
	}
	return &dest
}

func (k *DataKLine) GetDate() string {
	return k.Date
}

func (k *DataKLine) GetSecurityCode() string {
	return k.Code
}

func (k *DataKLine) Init(ctx context.Context, date string) error {
	//_ = barIndex
	_ = ctx
	_ = date
	return nil
}

func (k *DataKLine) Checkout(securityCode, date string) {
	//TODO implement me
	_ = securityCode
	_ = date
	panic("implement me")
}

func (k *DataKLine) Check(cacheDate, featureDate string) error {
	//TODO implement me
	_ = cacheDate
	_ = featureDate
	panic("implement me")
}

func (k *DataKLine) Filename(date, code string) string {
	//TODO implement me
	_ = code
	_ = date
	panic("implement me")
}

func (k *DataKLine) Print(code string, date ...string) {
	//TODO implement me
	_ = code
	_ = date
	panic("implement me")
}

func (k *DataKLine) Update(date string) error {
	base.UpdateAllBasicKLine(k.GetSecurityCode())
	_ = date
	return nil
}

func (k *DataKLine) Repair(date string) error {
	base.UpdateAllBasicKLine(k.GetSecurityCode())
	_ = date
	return nil
}

func (k *DataKLine) Increase(snapshot quotes.Snapshot) error {
	//TODO K线增量更新数据的条件是缓存的数据最晚的日期是上一个交易日, 否则会缺失缓存数据中最后1条数据和当日之间的数据, 所以只能按照K线的更新方法, 不适合用快照更新
	// 第一步: 取出最后一条数据的记录
	// 第二步: 检查时间的有效性
	// 第三步: 用快照组织K线结构
	// 第四步: 如果不符合快照更新, 则忽略
	_ = snapshot
	panic("implement me")
}
