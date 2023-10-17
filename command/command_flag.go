package command

import (
	"gitee.com/quant1x/engine/datasets/base"
	cmder "github.com/spf13/cobra"
	"unsafe"
)

var (
	flagAll       = cmdFlag[bool]{Name: "all", Value: false, Usage: "全部"}
	flagBaseData  = cmdFlag[bool]{Name: "base", Value: false, Usage: "基础数据"}
	flagFeatures  = cmdFlag[bool]{Name: "features", Value: false, Usage: "特征数据"}
	flagStartDate = cmdFlag[string]{Name: "start", Value: base.TickDefaultStartDate, Usage: "开始日期"}
	flagEndDate   = cmdFlag[string]{Name: "end", Value: "", Usage: "结束日期"}
	flagDate      = cmdFlag[string]{Name: "date", Value: "", Usage: "日期"}
	//dataF10       = factors.GetDataDescript(factors.FeatureF10)
	//flagF10       = cmdFlag[string]{Name: dataF10.Key, Value: "", Usage: dataF10.Name}
	//dataTrans = datasets.GetDataDescript(datasets.BaseTransaction)
	//flagTrans = cmdFlag[bool]{Name: dataTrans.Key(), Value: false, Usage: dataTrans.Name()}
)

type Command = cmder.Command

type cmdFlag[T ~int | ~bool | ~string] struct {
	Name  string
	Usage string
	Value T
}

func (cf *cmdFlag[T]) init(cmd *cmder.Command) {
	switch v := any(cf.Value).(type) {
	case bool:
		cmd.Flags().BoolVar((*bool)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case int:
		cmd.Flags().IntVar((*int)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case string:
		cmd.Flags().StringVar((*string)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	}
}

func commandInit[T ~int | ~bool | ~string](cmd *cmder.Command, cf *cmdFlag[T]) {
	switch v := any(cf.Value).(type) {
	case bool:
		cmd.Flags().BoolVar((*bool)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case int:
		cmd.Flags().IntVar((*int)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case string:
		cmd.Flags().StringVar((*string)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	}
}
