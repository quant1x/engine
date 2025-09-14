package factors

import (
	"context"

	"github.com/quant1x/engine/cache"
	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/gotdx/quotes"
)

// DataMinutes 分时数据
type DataMinutes struct {
	Manifest
}

func init() {
	summary := __mapDataSets[BaseMinutes]
	_ = cache.Register(&DataMinutes{Manifest: Manifest{DataSummary: summary}})
}

func (this *DataMinutes) Clone(date string, code string) DataSet {
	summary := __mapDataSets[BaseMinutes]
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

func (this *DataMinutes) Update(date string) error {
	base.UpdateMinutes(this.GetSecurityCode(), date)
	return nil
}

func (this *DataMinutes) Repair(date string) error {
	this.Update(date)
	return nil
}

func (this *DataMinutes) Increase(snapshot quotes.Snapshot) error {
	_ = snapshot
	//TODO implement me
	panic("implement me")

}

func (this *DataMinutes) Print(code string, date ...string) {
	_ = code
	_ = date
	//TODO implement me
	panic("implement me")
}
