package stmt

type Delete struct {
	whereStatement
	table  string
	wheres []WhereParam
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
