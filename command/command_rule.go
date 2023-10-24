package command

import (
	"fmt"
	"gitee.com/quant1x/engine/rules"
	cmder "github.com/spf13/cobra"
	"os"
)

// CmdRule 规则
var CmdRule = &cmder.Command{
	Use:   "rule",
	Short: "规则",
	Args: func(cmd *cmder.Command, args []string) error {
		return nil
	},
	//ValidArgsFunction: func(cmd *cmder.Command, args []string, toComplete string) ([]string, cmder.ShellCompDirective) {
	//
	//},
	PreRun: func(cmd *cmder.Command, args []string) {

	},
	Run: func(cmd *cmder.Command, args []string) {
		rules.PrintRuleList()
	},
}

func initRule() {
	CmdRule.SetFlagErrorFunc(func(cmd *cmder.Command, err error) error {
		args := os.Args[1:]
		cmd_, flags, err := cmd.Parent().Find(args)
		fmt.Println(cmd_, flags, err)
		return nil
	})
}
