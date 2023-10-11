package datasets

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx/quotes"
)

// TransactionRecord 成交记录
//
//	最短3秒内的合并统计数据, 与行情数据保持一致
//	不可以当作tick数据来使用
type TransactionRecord struct {
	DataCache
}

func init() {
	_ = cache.Register(&TransactionRecord{})
}

func (r *TransactionRecord) Get(code string, date ...string) any {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Kind() DataKind {
	return BaseTransaction
}

func (r *TransactionRecord) Name() string {
	return mapDataSets[r.Kind()].Name
}

func (r *TransactionRecord) Key() string {
	return mapDataSets[r.Kind()].Key
}

func (r *TransactionRecord) Init(barIndex *int, date string) error {
	_ = barIndex
	_ = date
	return nil
}

func (r *TransactionRecord) Filename(date, code string) string {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Update(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Repair(cacheDate, featureDate string) {
	base.GetTickAll(r.Code)
}

func (r *TransactionRecord) Increase(snapshot quotes.Snapshot) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Clone(date string, code string) DataSet {
	var dest = TransactionRecord{DataCache: DataCache{Date: date, Code: code}}
	return &dest
}
