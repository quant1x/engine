package cache

// Data 数据接口
type Data interface {
	// Kind 数据类型
	Kind() Kind
	// Key 数据关键词, key与cache落地强关联
	Key() string
	// Desc 数据描述(data description)
	Desc() string
	// FileName 缓存文件名
	//	接受两个参数 证券代码和日期
	FileName(securityCode, date string) string
	// IsBaseData 是否基础数据
	IsBaseData() bool
	// IsFeature 是否特征
	IsFeature() bool
	// Init 初始化, 加载配置信息
	Init(barIndex *int, date string) error
	// Check 数据校验
	Check(cacheDate, featureDate string)
	// Print 控制台输出指定日期的数据
	Print(code string, date ...string)
}
