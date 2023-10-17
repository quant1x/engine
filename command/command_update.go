package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/storages"
	cmder "github.com/spf13/cobra"
)

var (
	updateModules = []cmdFlag[bool]{}
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
			keywords := []string{}
			for _, m := range updateModules {
				if m.Value {
					keywords = append(keywords, m.Name)
					break
				}
			}
			handleUpdateFeatures(cacheDate, featureDate, keywords...)
		}
	},
}

func initUpdate() {
	commandInit(CmdUpdate, &flagAll)
	commandInit(CmdUpdate, &flagBaseData)
	commandInit(CmdUpdate, &flagFeatures)

	plugins := cache.Plugins(cache.PluginMaskFeature)
	updateModules = make([]cmdFlag[bool], len(plugins))
	for i, plugin := range plugins {
		key := plugin.Key()
		usage := plugin.Usage()
		updateModules[i] = cmdFlag[bool]{Name: key, Usage: plugin.Owner() + ": " + usage, Value: false}
		CmdUpdate.Flags().BoolVar(&(updateModules[i].Value), updateModules[i].Name, updateModules[i].Value, updateModules[i].Usage)
	}
}

// 全部更新
func handleUpdateAll(cacheDate, featureDate string) {
	handleUpdateBaseData(cacheDate, featureDate)
	handleUpdateFeatures(cacheDate, featureDate)
}

// 更新基础数据
func handleUpdateBaseData(cacheDate, featureDate string) {
	// 1. 获取全部注册的数据集插件
	mask := cache.PluginMaskBaseData
	plugins := cache.Plugins(mask)
	// 2. 执行操作
	storages.BaseDataUpdate(barIndex, cacheDate, featureDate, plugins, cache.OpUpdate)
}

// 更新特征组合
func handleUpdateFeatures(cacheDate, featureDate string, keywords ...string) {
	plugins := cache.PluginsWithName(cache.PluginMaskFeature, keywords...)
	if len(plugins) == 0 {
		// 1. 获取全部注册的数据集插件
		mask := cache.PluginMaskFeature
		//dataSetList := flash.DataSetList()
		plugins = cache.Plugins(mask)
	}
	storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpUpdate)
}
