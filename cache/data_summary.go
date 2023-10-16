package cache

// DataSummary 数据概要
type DataSummary struct {
	kind  Kind   // 类型
	key   string // 关键字
	name  string // 描述
	owner string // 拥有者
}

func Summary(kind Kind, key, name, owner string) DataSummary {
	return DataSummary{
		kind:  kind,
		key:   key,
		name:  name,
		owner: owner,
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
