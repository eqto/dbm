package stmt

type Delete struct {
	table  string
	wheres []WhereParam
	order  string
	count  int
}

func (d *Delete) where(param WhereParam) {
	d.wheres = append(d.wheres, param)
}

func (d *Delete) Where(conditions ...string) *UpdateWhere {
	if len(conditions) > 0 {
		d.wheres = []WhereParam{}
		for _, c := range conditions {
			d.wheres = append(d.wheres, WhereParam{c, false})
		}
	}
	return &UpdateWhere{d}
}

func (d *Delete) orderBy(orderBy string) {
	d.order = orderBy
}

func (d *Delete) limit(count int) {
	d.count = count
}
