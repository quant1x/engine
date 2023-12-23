package factors

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/gotdx/quotes"
)

// DataMinutes 分时数据
type DataMinutes struct {
	Manifest
}

func init() {
	summary := mapDataSets[BaseMinutes]
	_ = cache.Register(&DataMinutes{Manifest: Manifest{DataSummary: summary}})
}

func (this *DataMinutes) Clone(date string, code string) DataSet {
	summary := mapDataSets[BaseMinutes]
	var dest = DataMinutes{
		Manifest: Manifest{
			DataSummary: summary,
			Date:        date,
			Code:        code,
		},
	}
	return &dest
}

func (this *DataMinutes) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (this *DataMinutes) Update(date string) {
	base.UpdateMinutes(this.GetSecurityCode(), date)
}

func (this *DataMinutes) Repair(date string) {
	this.Update(date)
}

func (this *DataMinutes) Increase(snapshot quotes.Snapshot) {
	//TODO implement me
	panic("implement me")
}

func (this *DataMinutes) Print(code string, date ...string) {
	//TODO implement me
	panic("implement me")
}