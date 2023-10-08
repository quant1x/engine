package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/storages"
	cmder "github.com/spf13/cobra"
)

// CmdUpdate 更新数据
var CmdUpdate = &cmder.Command{
	Use:     "update",
	Example: Application + " update --all",
	//Args:    args.MinimumNArgs(0),
	Args: func(cmd *cmder.Command, args []string) error {
		return nil
	},
	Short: "更新股市数据",
	Long:  `更新股市数据`,
	Run: func(cmd *cmder.Command, args []string) {
		if flagAll.Value {
			// 全部更新
			handleUpdateAll()
		} else if flagHistory.Value {
			//handleUpdateAll()
		}
	},
}

func init() {
	commandInit(CmdUpdate, &flagAll)
	commandInit(CmdUpdate, &flagHistory)
}

func handleUpdateAll() {
	fmt.Println()
	currentDate := cache.DefaultCanUpdateDate()
	cacheDate, featureDate := cache.CorrectDate(currentDate)

	//storages.UpdateBaseCache(&barIndex, cacheDate, featureDate)
	storages.UpdateBaseData(&barIndex, cacheDate, featureDate)
	storages.UpdateFeature(&barIndex, cacheDate, featureDate)
	_ = cacheDate
	_ = featureDate
}
