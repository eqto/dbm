package stmt

import "errors"

func Copy(dest, src interface{}) error {
	src = StatementOf(src)
	if dest, ok := dest.(*Select); ok {
		if src, ok := src.(*Select); ok {
			dest.fields = src.fields
			dest.tables = src.tables
			dest.wheres = src.wheres
			dest.groupBy = src.groupBy
			dest.orderBy = src.orderBy
			dest.offset = src.offset
			dest.count = src.count
			return nil
		}
	}
	return errors.New(`copy failed`)
}

func StatementOf(s interface{}) interface{} {
	switch s := s.(type) {
	case *SelectWhere:
		return s.stmt
	case *GroupBy:
		return s.stmt
	case *OrderBy:
		return s.stmt
	case *UpdateWhere:
		return s.stmt
	case *UpdateFields:
		return s.stmt
	}
	return s
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
