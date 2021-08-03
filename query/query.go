package query

func Select(query string) *SelectStmt {
	return &SelectStmt{
		fields: parseFields(query),
	}
}

func InsertInto(table string, fields string) *InsertStmt {
	return &InsertStmt{table: Table{Name: table, Fields: parseFields(fields)}}
}

func Update(table string) *UpdateStmt {
	return NewUpdate(table)
}

func StatementOf(q interface{}) interface{} {
	switch q := q.(type) {
	case *TableStmt:
		return q.stmt
	case *Where:
		return q.Condition.stmt
	case *OrderBy:
		return q.table.stmt
	case *Condition:
		return q.stmt
	}
	return q
}

func TableOf(stmt interface{}) Table {
	switch stmt := stmt.(type) {
	case *TableStmt:
		return stmt.table
	case *InsertStmt:
		return stmt.table
	case *UpdateStmt:
		return stmt.table
	}
	return Table{}
}

func ValueOf(stmt interface{}) []string {
	if stmt, ok := stmt.(*InsertStmt); ok {
		return stmt.values
	}
	return nil
}

func FieldsOf(stmt *SelectStmt) []Field {
	return stmt.fields
}

func JoinOf(stmt *TableStmt) *TableStmt {
	return stmt.joinTo
}

func WhereOf(stmt interface{}) []string {
	switch stmt := stmt.(type) {
	case *SelectStmt:
		if stmt.whereStmt != nil {
			return stmt.whereStmt.conditions
		}
	case *UpdateStmt:
		if stmt.condition != nil {
			return stmt.condition.conditions
		}
	}
	return nil
}

func TableStmtOf(stmt *SelectStmt) *TableStmt {
	return stmt.tableStmt
}

func GroupByOf(stmt *SelectStmt) []string {
	if stmt.groupByStmt != nil {
		return stmt.groupByStmt.groups
	}
	return nil
}

func OrderByOf(stmt *SelectStmt) []string {
	if stmt.orderByStmt != nil {
		return stmt.orderByStmt.orders
	}
	return nil
}

func LimitOf(stmt *SelectStmt) (int, int) {
	return stmt.offset, stmt.count
}
