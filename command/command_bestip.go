package command

import (
	"gitee.com/quant1x/data/level1/quotes"
	cmder "github.com/spf13/cobra"
)

// CmdBestIP 检测服务器地址
var CmdBestIP = &cmder.Command{
	Use:   "bestip",
	Short: "检测服务器网速",
	Run: func(cmd *cmder.Command, args []string) {
		quotes.BestIP()
	},
}
