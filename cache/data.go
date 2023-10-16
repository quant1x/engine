package cache

import "context"

// Data 数据接口
type Data interface {
	// Kind 数据类型
	Kind() Kind
	// Key 数据关键词, key与cache落地强关联
	Key() string
	// Desc 数据描述(data description)
	Desc() string
	// Filename 缓存文件名
	//	接受两个参数 日期和证券代码
	// 	文件名为空不缓存
	Filename(date, securityCode string) string
	// IsBaseData 是否基础数据
	//IsBaseData() bool
	// IsFeature 是否特征
	//IsFeature() bool
	// Init 初始化, 加载配置信息
	Init(barIndex *int, date string) error
	// Check 数据校验
	Check(cacheDate, featureDate string)
	// Print 控制台输出指定日期的数据
	Print(code string, date ...string)
}

// Operator 缓存操作接口
//
//	数据操作, 包含初始化和拉取两个接口
type Operator interface {
	// Init 初始化, 接受context, 日期和证券代码作为入参
	Init(ctx context.Context, date, securityCode string)
	// Pull 拉取数据
	Pull(date, securityCode string) Operator
}
