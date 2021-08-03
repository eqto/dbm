package query

func assignWhere(stmt interface{}, where *Where) {
	switch stmt := stmt.(type) {
	case *SelectStmt:
		stmt.whereStmt = where
	}
}

func assignGroupBy(stmt interface{}, groupBy *GroupBy) {
	switch stmt := stmt.(type) {
	case *SelectStmt:
		stmt.groupByStmt = groupBy
	}
}

func assignOrderBy(stmt interface{}, orderBy *OrderBy) {
	switch stmt := stmt.(type) {
	case *SelectStmt:
		stmt.orderByStmt = orderBy
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
