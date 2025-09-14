package command

import (
	"github.com/quant1x/engine/rules"
	cmder "github.com/spf13/cobra"
)

// CmdRules 规则
var CmdRules = &cmder.Command{
	Use:   "rules",
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

func initRules() {
	//CmdRules.SetFlagErrorFunc(func(cmd *cmder.Command, err error) error {
	//	args := os.Args[1:]
	//	cmd_, flags, err := cmd.Parent().Find(args)
	//	fmt.Println(cmd_, flags, err)
	//	return nil
	//})
}
