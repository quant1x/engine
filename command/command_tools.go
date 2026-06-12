package command

import (
	"fmt"
	"slices"
	"strings"

	"github.com/quant1x/engine/tools"
	"github.com/quant1x/num"
	"github.com/quant1x/pkg/tools/tail"
	cmder "github.com/spf13/cobra"
)

const (
	toolsCommand     = "tool"
	toolsDescription = "工具"
)

var (
	CmdTools    *cmder.Command = nil // 工具集合
	cmdToolTail *cmder.Command = nil // tail工具
)

func initTools() {
	CmdTools = &cmder.Command{
		Use:     toolsCommand,
		Example: Application + " " + toolsCommand + " --help",
		Short:   toolsDescription,
		Run: func(cmd *cmder.Command, args []string) {

		},
	}
	toolsInitTail()
	CmdTools.AddCommand(cmdToolTail)
}

func toolsInitTail() {
	var taiConfig tail.Config
	var n int
	cmdToolTail = &cmder.Command{
		Use:                "tail",
		Example:            Application + " tool tail -f runtime.log",
		Short:              "文件末端阅览",
		DisableFlagParsing: true,
		Run: func(cmd *cmder.Command, args []string) {
			for i := 0; i < len(args); i++ {
				args[i] = strings.TrimSpace(args[i])
			}
			if len(args) != 2 || slices.Contains(args, "--help") || slices.Contains(args, "-h") {
				_ = cmd.Usage()
				return
			}
			if args[0] == "-f" {
				taiConfig.Follow = true
				taiConfig.Poll = true
			} else if args[0][:1] == "-" {
				n = int(num.AnyToInt64(args[0]))
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
	cmdToolTail.SetUsageFunc(func(command *cmder.Command) error {
		fmt.Println("Usage:\n" + Application + " tool tail [-f] [-n #] [file ...]")
		return nil
	})
	cmdToolTail.Flags().BoolVarP(&taiConfig.Follow, commandDefaultLongFlag, "f", false, "一直等待新数据添加到文件")
}
