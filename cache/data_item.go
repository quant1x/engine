package cache

import "context"

// DataItem 单行数据
type DataItem interface {
	// Init 初始化, 接受context, 日期和证券代码作为入参
	Init(ctx context.Context, date, securityCode string)
	// GetDate 得到日期
	GetDate() string // 日期
	// GetSecurityCode 得到证券代码
	GetSecurityCode() string // 证券代码
	// GetSecurityName 获取证券名称
	GetSecurityName() string // 证券名称
}
