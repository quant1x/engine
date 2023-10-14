package cache

// DataCommand 数据子命令接口
type DataCommand interface {
	// Usage 控制台参数提示信息
	Usage() string
}
