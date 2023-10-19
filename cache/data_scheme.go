package cache

// Scheme 特征的基础数据
type Scheme struct {
	date  string
	code  string
	kind  Kind
	owner string
	key   string
	name  string
	usage string
}

func DataScheme(date, code string, summary DataSummary) Scheme {
	return Scheme{
		date:  date,
		code:  code,
		kind:  summary.Kind(),
		owner: summary.Owner(),
		key:   summary.Key(),
		name:  summary.Name(),
		usage: summary.Usage(),
	}
}

func (d Scheme) Kind() Kind {
	return d.kind
}

func (d Scheme) Owner() string {
	return d.owner
}

func (d Scheme) Key() string {
	return d.key
}

func (d Scheme) Name() string {
	return d.name
}

func (d Scheme) Usage() string {
	return d.usage
}

func (d Scheme) GetDate() string {
	return d.date
}

func (d Scheme) GetSecurityCode() string {
	return d.code
}
