package db

type QueryParameter struct {
	qb *QueryBuilder
}

func (p *QueryParameter) Mode() int {
	return p.qb.mode
}

func (p *QueryParameter) Table() Table {
	return p.qb.table
}

func (p *QueryParameter) Keys() []string {
	return p.qb.keys
}

func (p *QueryParameter) Fields() []Field {
	return p.qb.fields
}

func (p *QueryParameter) Values() []interface{} {
	return p.qb.values
}

func (p *QueryParameter) Wheres() []string {
	return p.qb.wheres
}

func (p *QueryParameter) Start() int {
	return p.qb.start
}

func (p *QueryParameter) Count() int {
	return p.qb.count
}
