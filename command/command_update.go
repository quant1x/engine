package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/storages"
	cmder "github.com/spf13/cobra"
)

var (
	updateBases    []cmdFlag[bool]
	updateFeatures []cmdFlag[bool]
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
			//handleUpdateBaseData(featureDate)
			keywords := []string{}
			for _, m := range updateBases {
				if m.Value {
					keywords = append(keywords, m.Name)
					break
				}
			}
			handleUpdateBaseDataWithKeywords(cacheDate, featureDate, keywords...)
		} else if flagFeatures.Value {
			keywords := []string{}
			for _, m := range updateFeatures {
				if m.Value {
					keywords = append(keywords, m.Name)
					break
				}
			}
			handleUpdateFeaturesWithKeywords(cacheDate, featureDate, keywords...)
		}
	},
}

func initUpdate() {
	commandInit(CmdUpdate, &flagAll)
	commandInit(CmdUpdate, &flagBaseData)
	commandInit(CmdUpdate, &flagFeatures)

	plugins := cache.Plugins(cache.PluginMaskBaseData)
	updateBases = make([]cmdFlag[bool], len(plugins))
	for i, plugin := range plugins {
		key := plugin.Key()
		usage := plugin.Usage()
		updateBases[i] = cmdFlag[bool]{Name: key, Usage: plugin.Owner() + ": " + usage, Value: false}
		//CmdUpdate.Flags().BoolVar(&(updateFeatures[i].Value), updateFeatures[i].Name, updateFeatures[i].Value, updateFeatures[i].Usage)
		commandInit(CmdUpdate, &updateBases[i])
	}

	plugins = cache.Plugins(cache.PluginMaskFeature)
	updateFeatures = make([]cmdFlag[bool], len(plugins))
	for i, plugin := range plugins {
		key := plugin.Key()
		usage := plugin.Usage()
		updateFeatures[i] = cmdFlag[bool]{Name: key, Usage: plugin.Owner() + ": " + usage, Value: false}
		//CmdUpdate.Flags().BoolVar(&(updateFeatures[i].Value), updateFeatures[i].Name, updateFeatures[i].Value, updateFeatures[i].Usage)
		commandInit(CmdUpdate, &updateFeatures[i])
	}
}

// 全部更新
func handleUpdateAll(cacheDate, featureDate string) {
	handleUpdateBaseData(featureDate)
	handleUpdateFeaturesWithKeywords(cacheDate, featureDate)
}

// 更新基础数据
func handleUpdateBaseData(date string) {
	// 1. 获取全部注册的数据集插件
	mask := cache.PluginMaskBaseData
	plugins := cache.Plugins(mask)
	// 2. 执行操作
	storages.BaseDataUpdate(barIndex, date, plugins, cache.OpUpdate)
}

// 更新基础数据
func handleUpdateBaseDataWithKeywords(cacheDate, featureDate string, keywords ...string) {
	plugins := cache.PluginsWithName(cache.PluginMaskBaseData, keywords...)
	if len(plugins) == 0 {
		// 1. 获取全部注册的数据集插件
		mask := cache.PluginMaskBaseData
		//dataSetList := flash.DataSetList()
		plugins = cache.Plugins(mask)
	}
	storages.BaseDataUpdate(barIndex, featureDate, plugins, cache.OpUpdate)
}

// 更新特征组合
func handleUpdateFeaturesWithKeywords(cacheDate, featureDate string, keywords ...string) {
	plugins := cache.PluginsWithName(cache.PluginMaskFeature, keywords...)
	if len(plugins) == 0 {
		// 1. 获取全部注册的数据集插件
		mask := cache.PluginMaskFeature
		//dataSetList := flash.DataSetList()
		plugins = cache.Plugins(mask)
	}
	storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpUpdate)
}
