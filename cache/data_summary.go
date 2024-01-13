package cache

// DataSummary 数据概要
type DataSummary struct {
	kind  Kind   // 类型
	key   string // 关键字
	name  string // 名称
	owner string // 拥有者
	usage string // 用法
}

func Summary(kind Kind, key, name, owner string, usage ...string) DataSummary {
	var description string
	if len(usage) > 0 {
		description = usage[0]
	}
	return DataSummary{
		kind:  kind,
		key:   key,
		name:  name,
		owner: owner,
		usage: description,
	}
}

func (d DataSummary) Kind() Kind {
	return d.kind
}

func (d DataSummary) Key() string {
	return d.key
}

func (d DataSummary) Name() string {
	return d.name
}

func (d DataSummary) Owner() string {
	return d.owner
}

func (d DataSummary) Usage() string {
	return d.usage
}
