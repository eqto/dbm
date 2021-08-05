package stmt

type SelectWhere struct {
	stmt *Select
}

func (s *SelectWhere) And(condition string) *SelectWhere {
	s.stmt.wheres = append(s.stmt.wheres, WhereParam{condition, false})
	return s
}

func (s *SelectWhere) Or(condition string) *SelectWhere {
	s.stmt.wheres = append(s.stmt.wheres, WhereParam{condition, true})
	return s
}

func (s *SelectWhere) GroupBy(groupBy string) *GroupBy {
	return s.stmt.GroupBy(groupBy)
}

func (s *SelectWhere) OrderBy(orderBy string) *OrderBy {
	return s.stmt.OrderBy(orderBy)
}

func (s *SelectWhere) Limit(num ...int) *Select {
	return s.stmt.Limit(num...)
}
