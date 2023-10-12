package cache

const (
	UseGoroutine = false // 更新和修复数据是否启用协程
)

type OpKind int

const (
	OpUpdate OpKind = iota + 1 // 更新
	OpRepair                   // 修复
	OpIncr                     // 增量
)
