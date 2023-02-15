package stmt

type UpdateWhere struct {
	stmt whereStatement
}

func (u *UpdateWhere) And(condition string) *UpdateWhere {
	u.stmt.where(WhereParam{condition, false})
	return u
}

func (u *UpdateWhere) Or(condition string) *UpdateWhere {
	u.stmt.where(WhereParam{condition, true})
	return u
}
