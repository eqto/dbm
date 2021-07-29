package query

func TableOf(stmt interface{}) []Table {
	if stmt, ok := stmt.(*SelectStmt); ok {
		return stmt.tables
	}
	return nil
}

func WhereOf(stmt interface{}) []string {
	if stmt, ok := stmt.(*SelectStmt); ok {
		return stmt.wheres
	}
	return nil
}

func OrderByOf(stmt interface{}) []string {
	if stmt, ok := stmt.(*SelectStmt); ok {
		return stmt.orders
	}
	return nil
}

func LimitOf(stmt interface{}) (int, int) {
	if stmt, ok := stmt.(*SelectStmt); ok {
		return stmt.offset, stmt.count
	}
	return 0, 0
}
