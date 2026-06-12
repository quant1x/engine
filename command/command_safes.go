package command

import (
	"fmt"

	"github.com/quant1x/engine/trader"
	cmder "github.com/spf13/cobra"
)

var (
	safesList         bool       // 输出列表
	safesSecureType   int    = 0 // 名单类型
	safesSecurityCode string     // 证券代码
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
			if safesList {
				trader.GetBlackAndWhiteList()
			} else {
				if len(safesSecurityCode) == 0 {
					fmt.Println("证券代码不能为空")
					return
				}
				trader.AddCodeToBlackList(safesSecurityCode, trader.SecureType(safesSecureType))
			}
		},
	}
	CmdSafes.Flags().BoolVar(&safesList, "list", false, "显示黑白名单列表")
	CmdSafes.Flags().StringVar(&safesSecurityCode, "code", "", "证券代码")
	CmdSafes.Flags().IntVar(&safesSecureType, "type", 0, trader.UsageOfSecureType())
}
