package datasets

import (
	"context"
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

func (r *TransactionRecord) Print(code string, date ...string) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Kind() cache.Kind {
	return BaseTransaction
}

func (r *TransactionRecord) Key() string {
	return mapDataSets[r.Kind()].Key()
}

func (r *TransactionRecord) Name() string {
	return mapDataSets[r.Kind()].Name()
}

func (r *TransactionRecord) Owner() string {
	return mapDataSets[r.Kind()].Owner()
}

func (r *TransactionRecord) Usage() string {
	return mapDataSets[r.Kind()].Name()
}

func (r *TransactionRecord) Init(ctx context.Context, date, securityCode string) error {
	_ = ctx
	_ = date
	_ = securityCode
	return nil
}

func (r *TransactionRecord) Checkout(securityCode, date string) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Check(cacheDate, featureDate string) error {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Filename(date, code string) string {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Update(cacheDate, featureDate string) {
	base.UpdateTickStartDate(cacheDate)
	base.GetTickAll(r.Code)
	_ = featureDate
}

func (r *TransactionRecord) Repair(cacheDate, featureDate string) {
	//base.GetTickAll(r.Code)
	base.GetTickData(r.Code, cacheDate)
	_ = featureDate
}

func (r *TransactionRecord) Increase(snapshot quotes.Snapshot) {
	//TODO implement me
	panic("implement me")
}

func (r *TransactionRecord) Clone(date string, code string) DataSet {
	var dest = TransactionRecord{DataCache: DataCache{Date: date, Code: code}}
	return &dest
}
