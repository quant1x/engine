package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/storages"
	flags "github.com/spf13/cobra"
)

// CmdUpdate 更新数据
var CmdUpdate = &flags.Command{
	Use:     "update",
	Example: Application + " update --all",
	//Args:    args.MinimumNArgs(0),
	Args: func(cmd *flags.Command, args []string) error {
		return nil
	},
	Short: "更新股市数据",
	Long:  `更新股市数据`,
	Run: func(cmd *flags.Command, args []string) {
		if flagHistory.Value {
			handleUpdateAll()
		}
	},
}

func init() {
	flagHistory.init(CmdUpdate)
}

func handleUpdateAll() {
	fmt.Println()
	currentDate := cachel5.DefaultCanUpdateDate()
	cacheDate, featureDate := cachel5.CorrectDate(currentDate)

	//storages.UpdateBaseCache(&barIndex, cacheDate, featureDate)
	storages.UpdateBaseData(&barIndex, cacheDate, featureDate)
	_ = cacheDate
	_ = featureDate
}
