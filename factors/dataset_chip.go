package factors

import (
	"context"
	"os"

	"gitee.com/quant1x/data/exchange"
	"gitee.com/quant1x/data/level1/quotes"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/factors/pb"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/util/homedir"
	"google.golang.org/protobuf/proto"
)

func init() {
	summary := __mapDataSets[BaseChipDistribution]
	_ = cache.Register(&DataChip{Manifest{DataSummary: summary}})
}

type DataChip struct {
	Manifest
}

func (d *DataChip) Clone(date, code string) DataSet {
	summary := __mapDataSets[BaseChipDistribution]
	var dest = DataChip{
		Manifest: Manifest{
			DataSummary: summary,
			Date:        date,
			Code:        code,
		},
	}
	return &dest
}

func (d *DataChip) Print(code string, date ...string) {
	//TODO implement me
	panic("implement me")
}

func (d *DataChip) Init(ctx context.Context, date string) error {
	_ = ctx
	_ = date
	return nil
}

func (d *DataChip) Update(featureDate string) error {
	cacheFilename := cache.ChipsFilename(d.GetSecurityCode())
	filepath, err := homedir.Expand(cacheFilename)
	if err != nil {
		return err
	}
	// 检查目录, 不存在就创建
	_ = api.CheckFilepath(filepath, true)
	cd := pb.ChipDistribution{}
	dataBytes, err := os.ReadFile(cacheFilename)
	if err == nil && len(dataBytes) > 0 {
		err = proto.Unmarshal(dataBytes, &cd)
		if err != nil {
			return err
		}
	}
	if cd.Data == nil {
		cd.Data = make(map[string]*pb.Chips)
	}
	chips := updateChipDistribution(d.GetDate(), d.GetSecurityCode())
	cd.Data[featureDate] = chips
	dataBytes, err = proto.Marshal(&cd)
	if err != nil {
		return err
	}
	err = os.WriteFile(cacheFilename, dataBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataChip) Repair(featureDate string) error {
	return d.Update(featureDate)
}

func (d *DataChip) Increase(snapshot quotes.Snapshot) error {
	//TODO implement me
	panic("implement me")
}

// 更新筹码分布
func updateChipDistribution(featureDate, code string) *pb.Chips {
	securityCode := exchange.CorrectSecurityCode(code)
	cacheDate := exchange.FixTradeDate(featureDate)
	trans := base.CheckoutTransactionData(securityCode, cacheDate, true)
	if len(trans) == 0 {
		return nil
	}
	chips := &pb.Chips{
		Date: cacheDate,
		Dist: map[int32]float64{},
	}
	for _, v := range trans {
		price := int32(100 * v.Price)
		vol := v.Vol
		chips.Dist[price] = chips.Dist[price] + float64(vol)
	}
	//cd := pb.ChipDistribution{Data: make(map[string]*pb.Chips)}
	//cd.Data[cacheDate] = &chips
	//filename := "t1.bin"
	//data, err := proto.Marshal(&cd)
	//err = os.WriteFile(filename, data, 0644)
	return chips
}
