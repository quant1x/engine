package command

import (
	"fmt"
	"gitee.com/quant1x/engine/tools"
	"gitee.com/quant1x/pandas/stat"
	"gitee.com/quant1x/pkg/tools/tail"
	cmder "github.com/spf13/cobra"
	"slices"
	"strings"
)

var CmdTools = &cmder.Command{
	Use:     "tool",
	Example: Application + " tool --help",
	Short:   "工具",
	Run: func(cmd *cmder.Command, args []string) {

	},
}

func initTools() {
	toolsInitTail()

	CmdTools.AddCommand(cmdToolTail)
}

var (
	taiConfig tail.Config
	n         int
)

// tail工具
var cmdToolTail = &cmder.Command{
	Use:                "tail",
	Example:            Application + " tool tail -f runtime.log",
	Short:              "tail",
	DisableFlagParsing: true,
	Run: func(cmd *cmder.Command, args []string) {
		for i := 0; i < len(args); i++ {
			args[i] = strings.TrimSpace(args[i])
		}
		//fmt.Println(args)
		if len(args) != 2 || slices.Contains(args, "--help") || slices.Contains(args, "-h") {
			_ = cmd.Usage()
			return
		}
		if args[0] == "-f" {
			taiConfig.Follow = true
			taiConfig.Poll = true
		} else if args[0][:1] == "-" {
			n = int(stat.AnyToInt64(args[0]))
			n = -n
		}
		name := args[1]
		if n > 0 {
			tools.TailFileWithNumber(name, taiConfig, n)
		} else {
			done := make(chan bool)
			tools.TailFile(name, taiConfig, done)
			<-done
		}
	},
	//PreRunE: func(cmd *cmder.Command, args []string) error {
	//	//fmt.Println(args)
	//	if slices.Contains(args, "--help") || slices.Contains(args, "-h") {
	//		cmd.Usage()
	//	}
	//	return nil
	//},
}

func toolsInitTail() {
	cmdToolTail.SetUsageFunc(func(command *cmder.Command) error {
		fmt.Println("Usage:\n" + Application + " tool tail [-f] [-n #] [file ...]")
		return nil
	})
	//cmdToolTail.SetFlagErrorFunc(func(command *cmder.Command, err error) error {
	//	errText := err.Error()
	//	if !strings.HasPrefix(errText, "unknown shorthand flag:") {
	//		return err
	//	}
	//	arr := strings.Fields(errText)
	//	param := arr[len(arr)-1]
	//	fmt.Println("hh", param)
	//	command.Run(command, []string{param})
	//	return nil
	//})
	//cmdToolTail.Flags().IntVarP(&n, commandDefaultLongFlag, "n", 0, "最后的n行")
	cmdToolTail.Flags().BoolVarP(&taiConfig.Follow, commandDefaultLongFlag, "f", false, "一直等待新数据添加到文件")
}
