package command

import (
	"gitee.com/quant1x/engine/datasets/base"
	cmder "github.com/spf13/cobra"
	"unsafe"
)

var (
	flagAll       = cmdFlag[bool]{Name: "all", Value: false, Usage: "全部"}
	flagDataSet   = cmdFlag[bool]{Name: "dataset", Value: false, Usage: "数据集"}
	flagHistory   = cmdFlag[bool]{Name: "history", Value: false, Usage: "历史特征数据"}
	flagStartDate = cmdFlag[string]{Name: "start", Value: base.TickDefaultStartDate, Usage: "开始日期"}
	flagEndDate   = cmdFlag[string]{Name: "end", Value: "", Usage: "结束日期"}
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
