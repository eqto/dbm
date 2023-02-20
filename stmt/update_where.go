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

func (u *UpdateWhere) OrderBy(orderBy string) *UpdateOrderBy {
	u.stmt.orderBy(orderBy)
	return &UpdateOrderBy{stmt: u.stmt}
}

func (u *UpdateWhere) Limit(count int) *Update {
	u.stmt.limit(count)
	return u.stmt.(*Update)
}

type UpdateOrderBy struct {
	stmt whereStatement
}

func (u *UpdateOrderBy) Limit(count int) *Update {
	u.stmt.limit(count)
	return u.stmt.(*Update)
}
