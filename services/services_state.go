package services

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
	"os"
	"path/filepath"
	"time"
)

func stateFilename(date, timestamp string) string {
	date = exchange.FixTradeDate(date)
	t, _ := time.ParseInLocation(exchange.CN_SERVERTIME_FORMAT, timestamp, time.Local)
	//timestamp = t.Format(trading.CN_SERVERTIME_SHORT_FORMAT)
	timestamp = t.Format("150405")
	tm := date + "T" + timestamp
	filename := fmt.Sprintf("%s/update.%s", cache.GetVariablePath(), tm)
	return filename
}

// 状态文件不存在则可更新, 反之不可更新
func checkUpdateState(date, timestamp string) bool {
	filename := stateFilename(date, timestamp)
	return !api.FileExist(filename)
}

// 确定完成当期更新状态
func doneUpdate(date, timestamp string) {
	filename := stateFilename(date, timestamp)
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	api.CloseQuietly(file)
}

// 清理过期的状态文件
func cleanExpiredStateFiles() error {
	statePath := cache.GetVariablePath()
	filePaths, err := filepath.Glob(statePath + "/update.*")
	if err != nil {
		return err
	}
	for _, filePath := range filePaths {
		_ = os.Remove(filePath)
	}
	return nil
}
