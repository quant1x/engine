package cache

import "context"

const (
	KBarIndex = "barIndex"
	KLineMin  = 120 // K线最少记录数
)

// Schema 缓存的概要信息
type Schema interface {
	// Kind 数据类型
	Kind() Kind
	// Owner 提供者
	Owner() string
	// Key 数据关键词, key与cache落地强关联
	Key() string
	// Name 特性名称
	Name() string
	// Usage 控制台参数提示信息, 数据描述(data description)
	Usage() string
}

// Initialization 初始化接口
type Initialization interface {
	// Init 初始化, 接受context, 日期和证券代码作为入参
	Init(ctx context.Context, date string) error
}

// Properties 属性接口
type Properties interface {
	// GetDate 日期
	GetDate() string
	// GetSecurityCode 证券代码
	GetSecurityCode() string
}

// Manifest 提要
type Manifest interface {
	Schema
	Properties
	Initialization
}

// Validator 验证接口
type Validator interface {
	// Check 数据校验
	Check(featureDate string) error
}

// DataFile 基础数据文件接口
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
}

// Future 预备数据的接口
type Future interface {
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
