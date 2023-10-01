package cache

import "fmt"

// XdxrFilename XDXR缓存路径
func XdxrFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	xdxrPath := fmt.Sprintf("%s/%s/%s.csv", GetXdxrPath(), cacheId[:length-3], cacheId)
	return xdxrPath
}

// KLineFilename 基础K线缓存路径
func KLineFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	filepath := fmt.Sprintf("%s/%s/%s.csv", GetDayPath(), cacheId[:length-3], cacheId)
	return filepath
}
