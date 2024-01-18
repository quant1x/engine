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

const (
	// 状态文件时间格式
	timeLayoutOfState = "150405"
)

func stateFilename(date, timestamp string) string {
	date = exchange.FixTradeDate(date)
	t, _ := time.ParseInLocation(exchange.CN_SERVERTIME_FORMAT, timestamp, time.Local)
	timestamp = t.Format(timeLayoutOfState)
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
	pattern := filepath.Join(statePath, "update.*")
	filePaths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, filePath := range filePaths {
		_ = os.Remove(filePath)
	}
	return nil
}
