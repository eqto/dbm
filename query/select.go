package query

type SelectStmt struct {
	fields      []Field
	tableStmt   *TableStmt
	whereStmt   *Where
	orderByStmt *OrderBy
	groupByStmt *GroupBy
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
