package sqlserver

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/eqto/go-db/driver"
)

func init() {
	driver.Register(`sqlserver`, &sqlserver{})
}

type sqlserver struct {
}

func (*sqlserver) BuildContents(colTypes []*sql.ColumnType) ([]interface{}, error) {
	vals := make([]interface{}, len(colTypes))
	for idx, colType := range colTypes {
		scanType := colType.ScanType()
		switch scanType.Kind() {
		case reflect.Int64:
			vals[idx] = new(*int64)
		case reflect.Bool:
			vals[idx] = new(*bool)
		case reflect.String:
			vals[idx] = new(*string)
		case reflect.Struct:
			switch scanType.Name() {
			case `Time`:
				vals[idx] = new(*time.Time)
			}
		}
		if vals[idx] == nil {
			return nil, fmt.Errorf(`not supported type %s:%s as kind %s`, colType.Name(), colType.DatabaseTypeName(), scanType.Kind().String())
		}
	}
	return vals, nil
}

func (d *sqlserver) ConnectionString(hostname string, port int, username, password, name string) string {
	u := url.URL{
		Scheme:   `sqlserver`,
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%d", hostname, port),
		RawQuery: name,
	}
	return u.String()
}

func (*sqlserver) InsertQuery(tableName string, fields []string) string {
	values := make([]string, len(fields))
	for i := range values {
		values[i] = fmt.Sprintf(`@p%d`, i+1)
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		tableName,
		strings.Join(fields, `, `),
		strings.Join(values, `, `))
}

func (*sqlserver) IsDuplicate(msg string) bool {
	return regexp.MustCompile(`.*Cannot insert duplicate key.*`).MatchString(msg)
}
