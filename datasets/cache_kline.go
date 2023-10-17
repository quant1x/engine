package datasets

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx/quotes"
)

type DataKLine struct {
	DataCache
}

func init() {
	_ = cache.Register(&DataKLine{})
}

func (k *DataKLine) Init(ctx context.Context, date, securityCode string) error {
	//_ = barIndex
	_ = ctx
	_ = date
	_ = securityCode
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

func (k *DataKLine) Kind() cache.Kind {
	return BaseKLine
}

func (k *DataKLine) Key() string {
	return mapDataSets[k.Kind()].Key()
}

func (k *DataKLine) Name() string {
	return mapDataSets[k.Kind()].Name()
}

func (k *DataKLine) Owner() string {
	return mapDataSets[k.Kind()].Owner()
}

func (k *DataKLine) Usage() string {
	return mapDataSets[k.Kind()].Name()
}

func (k *DataKLine) Filename(date, code string) string {
	k.filename = cache.KLineFilename(code)
	_ = date
	return k.filename
}

func (k *DataKLine) Print(code string, date ...string) {
	//TODO implement me
	_ = code
	_ = date
	panic("implement me")
}

func (k *DataKLine) Update(cacheDate, featureDate string) {
	base.UpdateAllBasicKLine(k.Code)
	_ = cacheDate
	_ = featureDate
}

func (k *DataKLine) Repair(cacheDate, featureDate string) {
	base.UpdateAllBasicKLine(k.Code)
	_ = cacheDate
	_ = featureDate
}

func (k *DataKLine) Increase(snapshot quotes.Snapshot) {
	//TODO K线增量更新数据的条件是缓存的数据最晚的日期是上一个交易日, 否则会缺失缓存数据中最后1条数据和当日之间的数据, 所以只能按照K线的更新方法, 不适合用快照更新
	// 第一步: 取出最后一条数据的记录
	// 第二步: 检查时间的有效性
	// 第三步: 用快照组织K线结构
	// 第四步: 如果不符合快照更新, 则忽略
	_ = snapshot
	panic("implement me")
}

func (k *DataKLine) Clone(date, code string) DataSet {
	var dest = DataKLine{DataCache{Date: date, Code: code}}
	return &dest
}
