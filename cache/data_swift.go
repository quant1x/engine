package cache

// Swift 快速接口
//
//	securityCode 证券代码, 2位交易所缩写+6位数字代码, 例如sh600600, 代表上海市场的青岛啤酒
//	cacheDate 缓存日期
//	featureDate 特征数据的日期
type Swift interface {
	// Init 初始化, 加载配置信息
	Init(barIndex *int, date string) error
	// GetDate 日期
	GetDate() string
	// GetSecurityCode 证券代码
	GetSecurityCode() string
	// Check 数据校验
	Check(cacheDate, featureDate string)
	// Checkout 捡出指定日期的缓存数据
	Checkout(securityCode, date string)
	// Update 更新数据
	//	whole 是否完整的数据, false是加工成半成品数据, 为了配合Increase
	Update(securityCode, cacheDate, featureDate string, whole bool)
	// Repair 回补数据
	Repair(securityCode, cacheDate, featureDate string, whole bool)
	// Increase 增量计算, 用快照增量计算特征
	//Increase(securityCode string, snapshot quotes.Snapshot)
	//Assign(target *cachel5.CacheAdapter) bool   // 赋值
}
