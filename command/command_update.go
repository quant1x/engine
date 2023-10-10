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
		fmt.Println()
		currentDate := cache.DefaultCanUpdateDate()
		cacheDate, featureDate := cache.CorrectDate(currentDate)
		if flagAll.Value {
			// 全部更新
			handleUpdateAll(cacheDate, featureDate)
		} else if flagBaseData.Value {
			handleUpdateBaseData(cacheDate, featureDate)
		} else if flagFeatures.Value {
			handleUpdateFeatures(cacheDate, featureDate)
		}
	},
}

func init() {
	commandInit(CmdUpdate, &flagAll)
	commandInit(CmdUpdate, &flagFeatures)
}

// 全部更新
func handleUpdateAll(cacheDate, featureDate string) {
	handleUpdateBaseData(cacheDate, featureDate)
	handleUpdateFeatures(cacheDate, featureDate)
}

// 更新基础数据
func handleUpdateBaseData(cacheDate, featureDate string) {
	storages.UpdateBaseData(&barIndex, cacheDate, featureDate)
}

// 更新特征组合
func handleUpdateFeatures(cacheDate, featureDate string) {
	storages.UpdateFeatures(&barIndex, cacheDate, featureDate)
}
