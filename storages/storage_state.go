package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	orderConfig = config.OrderConfig()
)

// Touch 创建一个空文件
func Touch(filename string) error {
	_ = api.CheckFilepath(filename, true)
	return os.WriteFile(filename, nil, 0644)
}

// 获取状态机路径
func state_filepath(state_date string) string {
	flagPath := fmt.Sprintf("%s/var/%s", getQmtCachePath(), state_date)
	return flagPath
}

// 订单状态文件前缀
func order_state_prefix(state_date string, model models.Strategy, direction trader.Direction) string {
	qmtStrategyName := models.QmtStrategyName(model)
	qmtStrategyName = strings.ToLower(qmtStrategyName)
	prefix := fmt.Sprintf("%s-%s-%s-%s", state_date, orderConfig.AccountId, qmtStrategyName, direction.Flag())
	return prefix
}

// 获得订单标识文件名
func order_state_filename(date string, model models.Strategy, code string, direction trader.Direction) string {
	state_date := trading.FixTradeDate(date, cache.CACHE_DATE)
	orderFlagPath := state_filepath(state_date)
	prefix := order_state_prefix(date, model, direction)
	securityCode := proto.CorrectSecurityCode(code)
	filename := fmt.Sprintf("%s-%s.done", prefix, securityCode)
	state_filename := path.Join(orderFlagPath, filename)
	return state_filename
}

// CheckOrderState 检查订单执行状态
func CheckOrderState(date string, model models.Strategy, code string, direction trader.Direction) bool {
	filename := order_state_filename(date, model, code, direction)
	return api.FileExist(filename)
}

// PushOrderState 推送订单完成状态
func PushOrderState(date string, model models.Strategy, code string, direction trader.Direction) error {
	filename := order_state_filename(date, model, code, direction)
	return Touch(filename)
}

// CountStrategyOrders 统计策略订单数
func CountStrategyOrders(date string, model models.Strategy, direction trader.Direction) int {
	stateDate := trading.FixTradeDate(date, cache.CACHE_DATE)
	orderFlagPath := state_filepath(stateDate)
	prefix := order_state_prefix(stateDate, model, direction)
	files, err := filepath.Glob(orderFlagPath + "/" + prefix + "-*.done")
	if err != nil {
		logger.Error(err)
		return 0
	}
	return len(files)
}
