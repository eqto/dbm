package sqlserver

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"

	_ "github.com/denisenkom/go-mssqldb"
	db "github.com/eqto/dbqoo"
	"github.com/eqto/dbqoo/query"
)

func init() {
	db.Register(`sqlserver`, &Driver{})
}

type Driver struct {
	db.Driver
}

func (Driver) Name() string {
	return `sqlserver`
}

func (Driver) Query(stmt interface{}) string {
	stmt = query.StatementOf(stmt)
	switch stmt := stmt.(type) {
	case *query.SelectStmt:
		return querySelect(stmt)
	case *query.InsertStmt:
		return queryInsert(stmt)
	case *query.UpdateStmt:
		return queryUpdate(stmt)
	}
	return ``
}

func (Driver) DataSourceName(hostname string, port int, username, password, name string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local`,
		username, password,
		hostname, port,
		name,
	)
}
func (Driver) IsDuplicate(msg string) bool {
	return regexp.MustCompile(`^Duplicate entry.*`).MatchString(msg)
}

func (Driver) BuildContents(colTypes []*sql.ColumnType) ([]interface{}, error) {
	vals := make([]interface{}, len(colTypes))
	for idx, colType := range colTypes {
		scanType := colType.ScanType()
		switch scanType.Kind() {
		case reflect.Int8:
			vals[idx] = new(int8)
		case reflect.Uint8:
			vals[idx] = new(uint8)
		case reflect.Int16:
			vals[idx] = new(int16)
		case reflect.Uint16:
			vals[idx] = new(uint16)
		case reflect.Int32:
			vals[idx] = new(int32)
		case reflect.Uint32:
			vals[idx] = new(uint32)
		case reflect.Int64:
			vals[idx] = new(int64)
		case reflect.Uint64:
			vals[idx] = new(uint64)
		case reflect.Float32:
			vals[idx] = new(float32)
		case reflect.Float64:
			vals[idx] = new(float64)
		case reflect.Slice:
			switch colType.DatabaseTypeName() {
			case `DECIMAL`:
				if null, ok := colType.Nullable(); null || ok {
					vals[idx] = new(sql.NullFloat64)
				} else {
					vals[idx] = new(float64)
				}
			default:
				vals[idx] = new([]byte)
			}
		case reflect.Struct:
			switch scanType.Name() {
			case `NullInt64`:
				vals[idx] = new(sql.NullInt64)
			case `NullFloat64`:
				vals[idx] = new(sql.NullFloat64)
			case `NullTime`:
				vals[idx] = new(sql.NullTime)
			}
		}
		if vals[idx] == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		}
	}
	return vals, nil
}

// func (*driver) BuildQuery(param db.QueryParameter) string {
// 	s := strings.Builder{}
// 	println(param.Mode())
// 	switch param.Mode() {
// 	case db.ModeInsert:
// 		values := []string{}
// 		s.WriteString(`INSERT INTO ` + param.Table())
// 		if len(param.Keys()) > 0 {
// 			s.WriteString(`(` + strings.Join(param.Keys(), `, `) + `)`)
// 			values = append(values, `?`)
// 		}
// 		s.WriteString(fmt.Sprintf(` VALUES(%s)`, strings.Join(values, `, `)))
// 	case db.ModeSelect:
// 		strFields := strings.Join(param.Fields(), `, `)
// 		if strFields == `` {
// 			strFields = `*`
// 		}
// 		s.WriteString(fmt.Sprintf(`SELECT %s FROM %s`, strFields, param.Table()))
// 		if len(param.Wheres()) > 0 {
// 			s.WriteString(fmt.Sprintf(` WHERE %s`, strings.Join(param.Wheres(), ` AND `)))
// 		}
// 		if param.Count() > 0 {
// 			s.WriteString(fmt.Sprintf(` LIMIT %d, %d`, param.Start(), param.Count()))
// 		}
// 	}
// 	return s.String()
// }

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

// func init() {
// 	driver.Register(`sqlserver`, &sqlserver{})
// }

// type sqlserver struct {
// }

// func (*sqlserver) BuildContents(colTypes []*sql.ColumnType) ([]interface{}, error) {
// 	vals := make([]interface{}, len(colTypes))
// 	for idx, colType := range colTypes {
// 		scanType := colType.ScanType()
// 		switch scanType.Kind() {
// 		case reflect.Int64:
// 			vals[idx] = new(*int64)
// 		case reflect.Bool:
// 			vals[idx] = new(*bool)
// 		case reflect.String:
// 			vals[idx] = new(*string)
// 		case reflect.Struct:
// 			switch scanType.Name() {
// 			case `Time`:
// 				vals[idx] = new(*time.Time)
// 			}
// 		}
// 		if vals[idx] == nil {
// 			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
// 		}
// 	}
// 	return vals, nil
// }

// func (d *sqlserver) ConnectionString(hostname string, port int, username, password, name string) string {
// 	u := url.URL{
// 		Scheme:   `sqlserver`,
// 		User:     url.UserPassword(username, password),
// 		Host:     fmt.Sprintf("%s:%d", hostname, port),
// 		RawQuery: name,
// 	}
// 	return u.String()
// }

// func (*sqlserver) InsertQuery(tableName string, fields []string) string {
// 	values := make([]string, len(fields))
// 	for i := range values {
// 		values[i] = fmt.Sprintf(`@p%d`, i+1)
// 	}
// 	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
// 		tableName,
// 		strings.Join(fields, `, `),
// 		strings.Join(values, `, `))
// }

// func (*sqlserver) IsDuplicate(msg string) bool {
// 	return regexp.MustCompile(`.*Cannot insert duplicate key.*`).MatchString(msg)
// }
