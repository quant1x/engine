package storages

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
)

const (
	// 状态文件扩展名
	orderStateFileExtension = ".done"
)

// 获取状态机路径
func state_file_path(stateDate string) string {
	flagPath := filepath.Join(cache.GetQmtCachePath(), "var", stateDate)
	return flagPath
}

// 获取状态文件前缀
func state_file_prefix(stateDate, quantStrategyName string, direction trader.Direction) string {
	quantStrategyName = strings.ToLower(quantStrategyName)
	prefix := fmt.Sprintf("%s-%s-%s-%s", stateDate, traderConfig.AccountId, quantStrategyName, direction.Flag())
	return prefix
}

// 分拣订单状态字段
func order_state_fields(date, quantStrategyName string, direction trader.Direction) (flagPath, filenamePrefix string) {
	stateDate := exchange.FixTradeDate(date, cache.CACHE_DATE)
	flagPath = state_file_path(stateDate)
	filenamePrefix = state_file_prefix(stateDate, quantStrategyName, direction)
	return
}

// 从策略分拣订单状态字段
func order_state_fields_from_strategy(date string, model models.Strategy, direction trader.Direction) (flagPath, filenamePrefix string) {
	quantStrategyName := models.QmtStrategyName(model)
	flagPath, filenamePrefix = order_state_fields(date, quantStrategyName, direction)
	return
}

// 获得订单标识文件名
func order_state_filename(date string, model models.Strategy, direction trader.Direction, code string) string {
	orderFlagPath, filenamePrefix := order_state_fields_from_strategy(date, model, direction)
	securityCode := exchange.CorrectSecurityCode(code)
	filename := fmt.Sprintf("%s-%s.done", filenamePrefix, securityCode)
	state_filename := path.Join(orderFlagPath, filename)
	return state_filename
}

// CheckOrderState 检查订单执行状态
func CheckOrderState(date string, model models.Strategy, code string, direction trader.Direction) bool {
	filename := order_state_filename(date, model, direction, code)
	return api.FileExist(filename)
}

// PushOrderState 推送订单完成状态
func PushOrderState(date string, model models.Strategy, code string, direction trader.Direction) error {
	filename := order_state_filename(date, model, direction, code)
	return api.Touch(filename)
}

// 捡出策略订单状态文件列表
func checkoutStrategyOrderFiles(date string, model models.Strategy, direction trader.Direction) []string {
	orderFlagPath, filenamePrefix := order_state_fields_from_strategy(date, model, direction)
	pattern := filepath.Join(orderFlagPath, filenamePrefix+"-*"+orderStateFileExtension)
	files, err := filepath.Glob(pattern)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return files
}

// CountStrategyOrders 统计策略订单数
func CountStrategyOrders(date string, model models.Strategy, direction trader.Direction) int {
	files := checkoutStrategyOrderFiles(date, model, direction)
	return len(files)
}

// FetchListForFirstPurchase 获取指定日期交易的个股列表
func FetchListForFirstPurchase(date, quantStrategyName string, direction trader.Direction) []string {
	orderFlagPath, filenamePrefix := order_state_fields(date, quantStrategyName, direction)
	var list []string
	prefix := filepath.Join(orderFlagPath, filenamePrefix+"-")
	pattern := prefix + "*" + orderStateFileExtension
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		return list
	}
	for _, filename := range files {
		after, found := strings.CutPrefix(filename, prefix)
		if !found {
			continue
		}
		before, found := strings.CutSuffix(after, orderStateFileExtension)
		if !found {
			continue
		}
		list = append(list, before)
	}
	return list
}
