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
	// Init 初始化, 接受context, 日期作为入参
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

// FactorSignalEvaluator 因子验证接口
type FactorSignalEvaluator interface {
	// Check 对指定featureDate进行特征数据验证。
	// 返回值含义:
	//   hasSignal == true : 该日期产生了有效信号(用于统计信号覆盖率)
	//   err == nil        : 校验通过(即数据有效/逻辑正确)
	//   err != nil        : 校验失败, 需要记录或输出错误信息, hasSignal值在此场景下不参与胜率统计
	// 统计建议:
	//   Signals 计数: hasSignal==true 且 err==nil
	//   Passed  计数: err==nil (通过的数据样本)
	//   WinRate = PassedSignals / Signals (若需要区分“通过且有信号”)
	//   当前接口保持最小语义, 扩展需求可引入 ValidationResult 结构
	// NOTE: 不将 err==nil 自动推导为 hasSignal=true, 两者独立
	//
	Check(featureDate string) (hasSignal bool, err error)
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

// Future 预备数据的接口
type Future interface {
	// Update 更新数据
	// 	whole 是否完整的数据, false是加工成半成品数据, 为了配合Increase
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
