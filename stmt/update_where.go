package stmt

type UpdateWhere struct {
	stmt *Update
}

func (u *UpdateWhere) And(condition string) *UpdateWhere {
	u.stmt.wheres = append(u.stmt.wheres, WhereParam{condition, false})
	return u
}

func (u *UpdateWhere) Or(condition string) *UpdateWhere {
	u.stmt.wheres = append(u.stmt.wheres, WhereParam{condition, true})
	return u
}
