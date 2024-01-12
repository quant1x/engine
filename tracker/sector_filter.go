package tracker

import (
	"slices"
)

// 策略常量
const (
	StoreSnapshotTimeBegin = "09:25:00" // 缓存快照数据的开始时间
	StoreSnapshotTimeEnd   = "09:29:59" // 缓存快照数据的结束时间
	SectorMinChangeRate    = 0.10       // 板块开盘最小涨幅
	SectorMinVolume        = 1e8        // 板块开盘最低成交金额
)

var (
	blockIgnoreList = []string{"sh880516"} // ST板块
)

// 板块过滤规则
func sectorFilter(info SectorInfo) bool {
	if slices.Contains(blockIgnoreList, info.Code) {
		return false
	}
	if info.OpenChangeRate <= SectorMinChangeRate {
		return false
	}
	if info.OpenAmount <= SectorMinVolume {
		return false
	}
	return true
}

// 盘中过滤规则
func tickSectorFilter(info SectorInfo) bool {
	if slices.Contains(blockIgnoreList, info.Code) {
		return false
	}
	if info.OpenChangeRate <= SectorMinChangeRate {
		return false
	}
	if info.OpenAmount <= SectorMinVolume {
		return false
	}
	if info.ChangeRate < 0.00 {
		return false
	}
	return true
}
