package cache

// DataSummary 数据概要
type DataSummary struct {
	kind  Kind   // 类型
	key   string // 关键字
	desc  string // 描述
	owner string // 拥有者
}

func Summary(kind Kind, key, desc string) DataSummary {
	return DataSummary{
		kind: kind,
		key:  key,
		desc: desc,
	}
}

func (d DataSummary) Kind() Kind {
	return d.kind
}

func (d DataSummary) Key() string {
	return d.key
}

func (d DataSummary) Desc() string {
	return d.desc
}

func (d DataSummary) Owner() string {
	return d.owner
}
