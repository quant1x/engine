package command

import (
	"fmt"

	cmder "github.com/spf13/cobra"
)

// CmdVersion 版本
var CmdVersion = &cmder.Command{
	Use:   "version",
	Short: "显示版本号",
	Run: func(cmd *cmder.Command, args []string) {
		fmt.Println(MinVersion)
	},
}
