package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"os"
	"path"
)

var (
	orderConfig = config.OrderConfig()
)

// Touch 创建一个空文件
func Touch(filename string) error {
	_ = api.CheckFilepath(filename, true)
	return os.WriteFile(filename, nil, 0644)
}

// 获取订单标识路径
func getOrderFlagPath(statusDate string) string {
	flagPath := fmt.Sprintf("%s/var/%s", getQmtCachePath(), statusDate)
	return flagPath
}

// 获得订单标识文件名
func getOrderFlagFilename(date string, code string, direction trader.Direction) string {
	statusDate := trading.FixTradeDate(date, cache.CACHE_DATE)
	securityCode := proto.CorrectSecurityCode(code)
	orderFlag := direction.Flag()
	filename := fmt.Sprintf("%s-%s-%s-%s.done", statusDate, orderConfig.AccountId, securityCode, orderFlag)
	stock_order_path := path.Join(getOrderFlagPath(statusDate), filename)
	return stock_order_path
}

// 检查订单执行状态
func checkOrderStatus(date string, code string, direction trader.Direction) bool {
	filename := getOrderFlagFilename(date, code, direction)
	return api.FileExist(filename)
}

// 推送订单完成状态
func pushOrderStatus(date string, code string, direction trader.Direction) error {
	filename := getOrderFlagFilename(date, code, direction)
	return Touch(filename)
}
