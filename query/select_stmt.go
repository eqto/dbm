package query

import "strings"

type SelectStmt struct {
	fields []Field
	tables []Table
	wheres []string
	orders []string
	offset int
	count  int
}

//From table name with or without alias.
//Ex:
//Without alias => "Select books.* From books"
//With alias => "Select b.* From books b"
func (s *SelectStmt) From(table string) *SelectStmt {
	t := parseTable(table, s.fields)
	if len(s.tables) == 0 {
		s.tables = append(s.tables, t)
	} else {
		s.tables[0] = t
	}
	return s
}

func (s *SelectStmt) InnerJoin(table, condition string) *SelectStmt {
	t := parseTable(table, s.fields)
	t.Condition = condition
	s.tables = append(s.tables, t)
	return s
}

func (s *SelectStmt) Where(condition string) *SelectStmt {
	s.wheres = append(s.wheres, strings.TrimSpace(condition))
	return s
}

//OrderBy
//query: "title" => Select books.* From books ORDER BY title
//query: "title DESC" => Select books.* From books ORDER BY title DESC
func (s *SelectStmt) OrderBy(order string) *SelectStmt {
	s.orders = append(s.orders, strings.TrimSpace(order))
	return s
}

//Limit used by MySQL. Parameters 'num' can be single int for "LIMIT n" or double for "LIMIT n1, n2"
//Ex:
//SELECT * FROM books LIMIT 1. offset = 0, length = 1
//SELECT * FROM books LIMIT 0, 10. offset = 0, length = 10
func (s *SelectStmt) Limit(num ...int) *SelectStmt {
	if len(num) > 1 {
		s.offset = num[0]
		s.count = num[1]
	} else if len(num) == 1 {
		s.count = num[0]
	}
	return s
}

func Select(query string) *SelectStmt {
	return &SelectStmt{
		fields: parseFields(query),
	}
}
