package stmt

func StatementOf(stmt interface{}) interface{} {
	switch stmt := stmt.(type) {
	case *Where:
		return stmt.stmt
	case *UpdateFields:
		return stmt.stmt
	}
	return stmt
}

func FieldsOf(stmt interface{}) Fields {
	switch stmt := stmt.(type) {
	case *Select:
		return stmt.fields.value
	case *Insert:
		return stmt.fields
	}
	return Fields{}
}

func TableOf(stmt interface{}) string {
	switch stmt := stmt.(type) {
	case *Insert:
		return stmt.table
	case *Update:
		return stmt.table
	}
	return ``
}

func NameValueOf(stmt *Update) []string {
	return stmt.namevalue
}

func TablesOf(stmt *Select) Tables {
	return stmt.tables
}

func WheresOf(stmt interface{}) []WhereParam {
	switch stmt := stmt.(type) {
	case *Select:
		return stmt.wheres
	case *Update:
		return stmt.wheres
	}
	return nil
}

func GroupByOf(stmt *Select) string {
	return stmt.groupBy
}

func OrderByOf(stmt *Select) string {
	return stmt.orderBy
}

func LimitOf(stmt *Select) (int, int) {
	return stmt.offset, stmt.count
}
