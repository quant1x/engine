package factors

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/factors/pb"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/pkg/yaml"
	"google.golang.org/protobuf/proto"
	"os"
)

type DataChip struct {
	cache.DataSummary
	Date string
	Code string
}

//func init() {
//	summary := __mapDataSets[BaseChipDistribution]
//	_ = cache.Register(&DataChip{DataSummary: summary})
//}

// Print 控制台输出指定日期的数据
func (d *DataChip) Print(code string, date ...string) {
	//TODO implement me
	panic("implement me")
}

func (d *DataChip) Clone(date, code string) BaseData {
	//TODO implement me
	panic("implement me")
}

func (d *DataChip) Update(date, code string) error {
	//TODO implement me
	panic("implement me")
}

func (d *DataChip) Repair(date, code string) error {
	//TODO implement me
	panic("implement me")
}

func (d *DataChip) Increase(snapshot quotes.Snapshot) error {
	//TODO implement me
	panic("implement me")
}

type tt struct {
	CD map[float64]float64 `dataframe:"CD"`
}

// 更新筹码分布
func v1updateChipDistribution(date, code string) error {
	securityCode := exchange.CorrectSecurityCode(code)
	cacheDate := exchange.FixTradeDate(date)
	trans := base.CheckoutTransactionData(securityCode, cacheDate, true)
	tmp := map[float64]float64{}
	for _, v := range trans {
		price := v.Price
		vol := v.Vol
		tmp[price] = tmp[price] + float64(vol)
	}
	filename := "t1.yaml"
	//var list []tt
	//t1 := tt{CD: tmp}
	//list = append(list, t1)
	//err := api.SlicesToCsv(fn, list)
	//if err != nil {
	//	fmt.Println(err)
	//}
	data, err := yaml.Marshal(tmp)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	return err
}

func updateChipDistribution(date, code string) error {
	securityCode := exchange.CorrectSecurityCode(code)
	cacheDate := exchange.FixTradeDate(date)
	trans := base.CheckoutTransactionData(securityCode, cacheDate, true)
	chips := pb.Chips{
		Date: cacheDate,
		Dist: map[int32]float64{},
	}
	for _, v := range trans {
		price := int32(100 * v.Price)
		vol := v.Vol
		chips.Dist[price] = chips.Dist[price] + float64(vol)
	}
	cd := pb.ChipDistribution{Data: make(map[string]*pb.Chips)}
	cd.Data[cacheDate] = &chips
	filename := "t1.bin"
	data, err := proto.Marshal(&cd)
	err = os.WriteFile(filename, data, 0644)
	return err
}
