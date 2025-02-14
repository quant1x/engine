package command

import (
	"fmt"
	"gitee.com/quant1x/engine/trader"
	cmder "github.com/spf13/cobra"
)

var (
	safesSecureType   int = 0
	safesSecurityCode string
)

var (
	// CmdSafes 安全类-黑白名单
	CmdSafes *cmder.Command = nil
)

func initSafes() {
	CmdSafes = &cmder.Command{
		Use:     "safes",
		Example: Application + " safes --code=sh000001 --type=1",
		Short:   "黑白名单",
		Run: func(cmd *cmder.Command, args []string) {
			if len(safesSecurityCode) == 0 {
				fmt.Println("证券代码不能为空")
				return
			}
			trader.AddCodeToBlackList(safesSecurityCode, trader.SecureType(safesSecureType))
		},
	}

	CmdSafes.Flags().StringVar(&safesSecurityCode, "code", "", "证券代码")
	CmdSafes.Flags().IntVar(&safesSecureType, "type", 0, trader.UsageOfSecureType())
}
