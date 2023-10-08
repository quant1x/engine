package cache

const (
	cacheMetaPath     = "meta"     // 元数据缓存路径
	cacheDayPath      = "day"      // 日线路径
	cacheInfoPath     = "info"     // 信息路径
	cacheTickPath     = "tick"     // tick路径
	cacheXdxrPath     = "xdxr"     // 除权除息路径
	cacheFeaturePath  = "features" // 特征路径
	cacheFinancePath  = "finance"  // 财务信息路径
	cacheSnapshotPath = "snapshot" // 快照数据路径
	cacheHoldingPath  = "holding"  // 流通股东数据路径
	cacheFundFlowPath = "fund"     // 资金流向
)

// GetMetaPath 元数据路径
func GetMetaPath() string {
	return GetRootPath() + "/" + cacheMetaPath
}

// GetDayPath 历史数据-日线缓存路径
func GetDayPath() string {
	return GetRootPath() + "/" + cacheDayPath
}

// GetFeaturesPath 获取特征路径
func GetFeaturesPath() string {
	return GetRootPath() + "/" + cacheFeaturePath
}

// GetXdxrPath 除权除息文件存储路径
func GetXdxrPath() string {
	return GetRootPath() + "/" + cacheXdxrPath
}

// GetTickPath tick数据路径
func GetTickPath() string {
	return GetRootPath() + "/" + cacheTickPath
}

// GetHoldingPath 十大流通股股东数据路径
func GetHoldingPath() string {
	return GetRootPath() + "/" + cacheHoldingPath
}

// GetInfoPath 信息路径
func GetInfoPath() string {
	return GetRootPath() + "/" + cacheInfoPath
}

// GetQuarterlyPath 季报路径
//
//	Deprecated: 不推荐
func GetQuarterlyPath() string {
	return GetRootPath() + "/" + cacheInfoPath + "q"
}

// GetSnapshotPath 快照路径
func GetSnapshotPath() string {
	return GetRootPath() + "/" + cacheSnapshotPath
}

// GetFundFlowPath 资金流向目录
func GetFundFlowPath() string {
	return GetRootPath() + "/" + cacheFundFlowPath
}
