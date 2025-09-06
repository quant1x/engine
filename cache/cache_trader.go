package cache

import (
	"fmt"
	"strings"
	"sync"

	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/gox/api"
)

const (
	keyOrdersPath = "qmt" // QMT订单缓存路径
)

var (
	qmtOnce      sync.Once
	qmtOrderPath = "" // QMT订单路径
)

// QMT默认路径
func defaultQmtCachePath() string {
	path := fmt.Sprintf("%s/%s", GetRootPath(), keyOrdersPath)
	return path
}

func lazyInitQmt() {
	traderParameter := &config.GlobalConfig.Trader
	traderParameter.OrderPath = strings.TrimSpace(traderParameter.OrderPath)
	if len(traderParameter.OrderPath) > 0 && api.CheckFilepath(traderParameter.OrderPath, true) == nil {
		// 如果配置了路径且有效
		qmtOrderPath = traderParameter.OrderPath
	} else {
		qmtOrderPath = defaultQmtCachePath()
	}
	traderParameter.OrderPath = qmtOrderPath
}

func initMiniQmt() {
	qmtOnce.Do(lazyInitQmt)
}

// GetQmtCachePath QMT订单文件路径
func GetQmtCachePath() string {
	qmtOnce.Do(lazyInitQmt)
	return qmtOrderPath
}
