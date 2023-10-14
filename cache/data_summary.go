package cache

// DataSummary 数据概要
type DataSummary struct {
	kind Kind
	key  string
	desc string
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
