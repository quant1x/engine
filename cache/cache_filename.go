package cache

import "fmt"

// XdxrFilename XDXR缓存路径
func XdxrFilename(code string) string {
	cacheId := CacheId(code)
	length := len(cacheId)
	xdxrPath := fmt.Sprintf("%s/%s/%s.csv", GetXdxrPath(), cacheId[:length-3], cacheId)
	return xdxrPath
}
