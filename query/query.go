package query

func TableStmtOf(stmt *SelectStmt) *TableStmt {
	return stmt.tableStmt
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

func JoinOf(stmt *TableStmt) *TableStmt {
	return stmt.join
}

func FieldsOf(stmt interface{}) []Field {
	stmt = StatementOf(stmt)
	if stmt, ok := stmt.(*SelectStmt); ok {
		return stmt.fields
	}
	return nil
}

func StatementOf(q interface{}) interface{} {
	switch q := q.(type) {
	case *TableStmt:
		return q.stmt
	case *WhereStmt:
		return q.ConditionStmt.stmt
	case *OrderByStmt:
		return q.table.stmt
	case *ConditionStmt:
		return q.stmt
	}
	return q
}

func WhereOf(stmt interface{}) []string {
	switch stmt := stmt.(type) {
	case *SelectStmt:
		if stmt.where != nil {
			return stmt.where.conditions
		}
	case *UpdateStmt:
		if stmt.condition != nil {
			return stmt.condition.conditions
		}
	}
	return nil
}

func OrderByOf(stmt interface{}) []string {
	if stmt, ok := stmt.(*SelectStmt); ok && stmt.orderBy != nil {
		return stmt.orderBy.orders
	}
	return nil
}

func LimitOf(stmt interface{}) (int, int) {
	if stmt, ok := stmt.(*SelectStmt); ok {
		return stmt.offset, stmt.count
	}
	return 0, 0
}

func ValueOf(stmt interface{}) []string {
	if stmt, ok := stmt.(*InsertStmt); ok {
		return stmt.values
	}
	return nil
}

func assignWhere(stmt interface{}, where *WhereStmt) {
	switch stmt := stmt.(type) {
	case *SelectStmt:
		stmt.where = where
	}
}

func assignOrderBy(stmt interface{}, orderBy *OrderByStmt) {
	switch stmt := stmt.(type) {
	case *SelectStmt:
		stmt.orderBy = orderBy
	}
}

func assignLimit(stmt interface{}, num ...int) {
	switch stmt := stmt.(type) {
	case *SelectStmt:
		switch len(num) {
		case 1:
			stmt.count = num[0]
		case 2:
			stmt.offset, stmt.count = num[0], num[1]
		}
	}
}
