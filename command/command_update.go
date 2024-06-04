package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/global"
	"gitee.com/quant1x/engine/storages"
	cmder "github.com/spf13/cobra"
)

const (
	updateCommand     = "update"
	updateDescription = "更新数据"
)

var (
	// CmdUpdate 更新数据
	CmdUpdate *cmder.Command = nil
	barIndex                 = 1
)

func initUpdate(variables *global.Variables) {
	CmdUpdate = &cmder.Command{
		Use:     updateCommand,
		Example: Application + " " + updateCommand + " --all",
		//Args:    args.MinimumNArgs(0),
		Args: func(cmd *cmder.Command, args []string) error {
			return nil
		},
		Short: updateDescription,
		Long:  updateDescription,
		Run: func(cmd *cmder.Command, args []string) {
			fmt.Println()
			currentDate := cache.DefaultCanUpdateDate()
			cacheDate, featureDate := cache.CorrectDate(currentDate)
			if flagAll.Value {
				// 全部更新
				handleUpdateAll(cacheDate, featureDate, variables)
			} else if len(flagBaseData.Value) > 0 {
				all, keywords := parseFields(flagBaseData.Value)
				if all || len(keywords) == 0 {
					clear(keywords)
				}
				handleUpdateBaseDataWithKeywords(cacheDate, featureDate, variables, keywords...)
			} else if len(flagFeatures.Value) > 0 {
				all, keywords := parseFields(flagFeatures.Value)
				if all || len(keywords) == 0 {
					clear(keywords)
				}
				handleUpdateFeaturesWithKeywords(cacheDate, featureDate, variables, keywords...)
			} else {
				fmt.Println("Error: 非全部更新, 必须携带--features或--base")
				_ = cmd.Usage()
			}
		},
	}
	commandInit(CmdUpdate, &flagAll)

	// 1. 基础数据
	plugins := cache.Plugins(cache.PluginMaskBaseData)
	flagBaseData.Usage = getPluginsUsage(plugins)
	commandInit(CmdUpdate, &flagBaseData)

	// 2. 特征数据
	plugins = cache.Plugins(cache.PluginMaskFeature)
	flagFeatures.Usage = getPluginsUsage(plugins)
	commandInit(CmdUpdate, &flagFeatures)

	//// 3. 处理异常
	//CmdUpdate.SetFlagErrorFunc(func(cmd *cmder.Command, err error) error {
	//	return nil
	//})
}

// 全部更新
func handleUpdateAll(cacheDate, featureDate string, variables *global.Variables) {
	handleUpdateBaseData(featureDate, variables)
	handleUpdateFeaturesWithKeywords(cacheDate, featureDate, variables)
}

// 更新基础数据
func handleUpdateBaseData(date string, variables *global.Variables) {
	// 1. 获取全部注册的数据集插件
	mask := cache.PluginMaskBaseData
	plugins := cache.Plugins(mask)
	// 2. 执行操作
	storages.BaseDataUpdate(barIndex, date, plugins, cache.OpUpdate, *variables.MarketData)
}

// 更新基础数据
func handleUpdateBaseDataWithKeywords(cacheDate, featureDate string, variables *global.Variables, keywords ...string) {
	plugins := cache.PluginsWithName(cache.PluginMaskBaseData, keywords...)
	if len(plugins) == 0 {
		// 1. 获取全部注册的数据集插件
		mask := cache.PluginMaskBaseData
		plugins = cache.Plugins(mask)
	}
	storages.BaseDataUpdate(barIndex, featureDate, plugins, cache.OpUpdate, *variables.MarketData)
	_ = cacheDate
}

// 更新特征组合
func handleUpdateFeaturesWithKeywords(cacheDate, featureDate string, variables *global.Variables, keywords ...string) {
	plugins := cache.PluginsWithName(cache.PluginMaskFeature, keywords...)
	if len(plugins) == 0 {
		// 1. 获取全部注册的数据集插件
		mask := cache.PluginMaskFeature
		plugins = cache.Plugins(mask)
	}
	storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpUpdate, *variables.MarketData)
}
