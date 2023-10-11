package factors

// Increase 增量数据计算接口
//
//	deprecated: 不推荐
type Increase[T any] interface {
	Add(data T) T
}
