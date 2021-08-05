package stmt

type WhereParam struct {
	Condition string
	Or        bool //OR, default: false = AND
}

type Where struct {
	stmt interface{}
}

type SelectWhere struct {
	Where
}

func (w *Where) And(condition string) *Where {
	switch stmt := w.stmt.(type) {
	case *Select:
		stmt.wheres = append(stmt.wheres, WhereParam{condition, false})
	case *Update:
		stmt.wheres = append(stmt.wheres, WhereParam{condition, false})
	}
	return w
}

func (w *Where) Or(condition string) *Where {
	switch stmt := w.stmt.(type) {
	case *Select:
		stmt.wheres = append(stmt.wheres, WhereParam{condition, true})
	case *Update:
		stmt.wheres = append(stmt.wheres, WhereParam{condition, false})
	}
	return w
}

func (s *SelectWhere) GroupBy(groupBy string) *GroupBy {
	return s.stmt.(*Select).GroupBy(groupBy)
}

func (s *SelectWhere) OrderBy(orderBy string) *OrderBy {
	return s.stmt.(*Select).OrderBy(orderBy)
}

func (s *SelectWhere) Limit(num ...int) *Select {
	return s.stmt.(*Select).Limit(num...)
}
