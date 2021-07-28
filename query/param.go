package query

type Param struct {
	b *Builder
}

func (p *Param) Mode() int {
	return p.b.mode
}

func (p *Param) Table() string {
	return p.b.table
}

func (p *Param) Keys() []string {
	return p.b.keys
}

func (p *Param) Fields() []string {
	return p.b.fields
}

func (p *Param) Values() []interface{} {
	return p.b.values
}

func (p *Param) Wheres() []string {
	return p.b.wheres
}

func (p *Param) Start() int {
	return p.b.start
}

func (p *Param) Count() int {
	return p.b.count
}
