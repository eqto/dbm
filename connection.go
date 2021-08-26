package dbm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

//Connection ...
type Connection struct {
	db *sql.DB

	hostname string
	port     int
	username string
	password string
	name     string
	driver   Driver
}

//Connect ...
func (c *Connection) Connect() error {
	db, e := sql.Open(c.driver.Name(), c.driver.DataSourceName(c.hostname, c.port, c.username, c.password, c.name))
	if e != nil {
		return e
	}
	if e := db.Ping(); e != nil {
		return e
	}

	c.db = db
	return nil
}

func (c *Connection) Driver() Driver {
	return c.driver
}

//Ping ...
func (c *Connection) Ping() error {
	return c.db.Ping()
}

//SetConnMaxLifetime ...
func (c *Connection) SetConnMaxLifetime(duration time.Duration) {
	c.db.SetConnMaxLifetime(duration)
}

//SetMaxIdleConns ...
func (c *Connection) SetMaxIdleConns(max int) {
	c.db.SetMaxIdleConns(max)
}

//SetMaxOpenConns ...
func (c *Connection) SetMaxOpenConns(max int) {
	c.db.SetMaxOpenConns(max)
}

//Begin ...
func (c *Connection) Begin() (*Tx, error) {
	sqlTx, e := c.db.Begin()
	if e != nil {
		return nil, e
	}
	return &Tx{driver: c.driver, sqlTx: sqlTx}, nil
}

//MustBegin ...
func (c *Connection) MustBegin() *Tx {
	tx, e := c.Begin()
	if e != nil {
		return nil
	}
	return tx
}

//MustExec ...
func (c *Connection) MustExec(query string, args ...interface{}) *Result {
	result, e := c.Exec(query, args...)
	if e != nil {
		panic(e)
	}
	return result
}

//Exec ...
func (c *Connection) Exec(query string, args ...interface{}) (*Result, error) {
	return exec(c.driver, c.db.Exec, query, args...)
}

//Get ...
func (c *Connection) Get(query string, args ...interface{}) (Resultset, error) {
	rs, e := c.Select(query, args...)
	if e != nil {
		return nil, e
	} else if rs == nil {
		return nil, nil
	}
	return rs[0], nil
}

//MustGet ...
func (c *Connection) MustGet(query string, args ...interface{}) Resultset {
	rs, e := c.Get(query, args...)
	if e != nil {
		panic(e)
	}
	return rs
}

//GetStruct ...
func (c *Connection) GetStruct(dest interface{}, query string, args ...interface{}) error {
	typeOf := reflect.TypeOf(dest)
	if typeOf.Kind() != reflect.Ptr {
		return errors.New(`dest is not a pointer`)
	}

	rs, e := c.Get(query, args...)
	if e != nil {
		return e
	}

	if rs == nil || len(rs) == 0 {
		return newSQLError(c.driver, errNotFound)
	}

	typeOf = typeOf.Elem()
	return assignStruct(dest, createFieldMap(typeOf), rs, typeOf)
}

//MustSelect ...
func (c *Connection) MustSelect(query string, args ...interface{}) []Resultset {
	rs, e := c.Select(query, args...)
	if e != nil {
		panic(e)
	}
	return rs
}

//Select ...
func (c *Connection) Select(query string, args ...interface{}) ([]Resultset, error) {
	return execQuery(c.driver, c.db.Query, query, args...)
}

//SelectStruct ...
func (c *Connection) SelectStruct(dest interface{}, query string, args ...interface{}) error {
	return execQueryStruct(c.driver, c.Select, dest, query, args...)
}

//MustInsert ...
func (c *Connection) MustInsert(tableName string, dataMap map[string]interface{}) *Result {
	result, e := c.Insert(tableName, dataMap)
	if e != nil {
		panic(e)
	}
	return result
}

//Insert ...
func (c *Connection) Insert(tableName string, dataMap map[string]interface{}) (*Result, error) {
	length := len(dataMap)
	fields := make([]string, length)
	values := make([]interface{}, length)
	idx := 0
	for name, value := range dataMap {
		fields[idx] = name
		values[idx] = value
		idx++
	}
	q := InsertInto(tableName, strings.Join(fields, `, `))
	return c.Exec(c.driver.StatementString(q), values...)
}

//EnumValues return enum values, parameter field using dot notation. Ex: profile.gender , returning ['male', 'female']
func (c *Connection) EnumValues(field string) ([]string, error) {
	cols := strings.Split(field, `.`)
	enum, e := c.Get(`SELECT column_type FROM information_schema.columns WHERE table_name = ?
		AND column_name = ?`, cols[0], cols[1])
	if e != nil {
		return nil, e
	}
	regexEnum := regexp.MustCompile(`'[a-zA-Z0-9]+'`)

	values := regexEnum.FindAllString(enum.String(`column_type`), -1)

	for i := 0; i < len(values); i++ {
		values[i] = strings.Trim(values[i], `'`)
	}
	return values, nil
}

//Close ...
func (c *Connection) Close() error {
	if c.db == nil {
		return nil
	}
	return c.db.Close()
}

func newConnection(driverName, hostname string, port int, username, password, name string) (*Connection, error) {
	if port < 0 || port > 65535 {
		return nil, fmt.Errorf(`invalid port %d`, port)
	}

	drv, e := getDriver(driverName)
	if e != nil {
		return nil, e
	}
	return &Connection{
		driver:   drv,
		hostname: hostname, port: port,
		username: username, password: password,
		name: name,
	}, nil

}
