package cache

import (
	"context"
)

const (
	KBarIndex = "barIndex"
)

// Initialize 初始化接口
type Initialize interface {
	// Init 初始化, 接受context, 日期和证券代码作为入参
	Init(ctx context.Context, date, securityCode string) error
}

//// DataItem 单行数据
//type DataItem interface {
//	// Initialize 初始化接口
//	Initialize
//	// GetDate 得到日期
//	GetDate() string // 日期
//	// GetSecurityCode 得到证券代码
//	GetSecurityCode() string // 证券代码
//	// GetSecurityName 获取证券名称
//	GetSecurityName() string // 证券名称
//}

// Trait 基础的特性
//
//	这也是一个特征, 为啥起这个名字, 自己可以脑补 哈哈~
type Trait interface {
	// Kind 数据类型
	Kind() Kind
	// Key 数据关键词, key与cache落地强关联
	Key() string
	// Name 特性名称
	Name() string
	// Usage 控制台参数提示信息, 数据描述(data description)
	Usage() string
	// Owner 提供者
	Owner() string
	// Initialize 初始化
	//Initialize
	// Init 初始化, 接受context, 日期和证券代码作为入参
	Init(ctx context.Context, date, securityCode string) error
}

type DataFile interface {
	// Checkout 捡出指定日期的缓存数据
	Checkout(securityCode, date string)
	// Filename 缓存文件名
	//	接受两个参数 日期和证券代码
	// 	文件名为空不缓存
	Filename(date, securityCode string) string
	// Check 数据校验
	Check(cacheDate, featureDate string) error
	// Print 控制台输出指定日期的数据
	//Print(code string, date ...string)
}

// Swift 快速接口
//
//	securityCode 证券代码, 2位交易所缩写+6位数字代码, 例如sh600600, 代表上海市场的青岛啤酒
//	cacheDate 缓存日期
//	featureDate 特征数据的日期
type Swift interface {
	// GetDate 日期
	GetDate() string
	// GetSecurityCode 证券代码
	GetSecurityCode() string
	// Update 更新数据
	//	whole 是否完整的数据, false是加工成半成品数据, 为了配合Increase
	Update(securityCode, cacheDate, featureDate string, whole bool)
	// Repair 回补数据
	Repair(securityCode, cacheDate, featureDate string, whole bool)
	// Increase 增量计算, 用快照增量计算特征
	//Increase(securityCode string, snapshot quotes.Snapshot)
	//Assign(target *cachel5.CacheAdapter) bool   // 赋值
}

// Operator 缓存操作接口
//
//	数据操作, 包含初始化和拉取两个接口
type Operator interface {
	// Pull 拉取数据
	Pull(date, securityCode string) Operator
}
