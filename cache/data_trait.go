package cache

// Trait 基础的特性
//
//	这也是一个特征, 为啥起这个名字, 自己可以脑补 哈哈~
type Trait interface {
	// Kind 数据类型
	Kind() Kind
	// Key 数据关键词, key与cache落地强关联
	Key() string
	// Desc 数据描述(data description)
	Desc() string
	// Provider 提供者
	Provider() string
}
