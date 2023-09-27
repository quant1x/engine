package factors

// Factor 因子
type Factor interface {
	Init()                                       // 初始化
	Name() string                                // 因子名称
	Weight() uint64                              // 权重
	Execute(securityCode, date string) (ok bool) // 执行
}
