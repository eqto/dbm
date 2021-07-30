package query

type SelectStmt struct {
	fields      []Field
	tableStmt   *TableStmt
	whereStmt   *WhereStmt
	orderByStmt *OrderByStmt
	groupByStmt *GroupByStmt
	offset      int
	count       int
}

//From table name with or without alias.
//Ex:
//Without alias => "Select books.* From books"
//With alias => "Select b.* From books b"
func (s *SelectStmt) From(table string) *TableStmt {
	s.tableStmt = parseTable(s, table)
	return s.tableStmt
}

func Select(query string) *SelectStmt {
	return &SelectStmt{
		fields: parseFields(query),
	}
}
