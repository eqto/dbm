package db

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var lastCn *Connection

//Connection ...
type Connection struct {
	db *sql.DB

	Driver   string
	Hostname string
	Port     uint16
	Username string
	Password string
	Name     string
}

//Connect ...
func (c *Connection) Connect() error {
	cnStr := ``
	switch c.Driver {
	case DriverMySQL:
		cnStr = fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local`,
			c.Username, c.Password,
			c.Hostname, c.Port,
			c.Name,
		)
	case DriverSQLServer:
		u := &url.URL{
			Scheme:   c.Driver,
			User:     url.UserPassword(c.Username, c.Password),
			Host:     fmt.Sprintf("%s:%d", c.Hostname, c.Port),
			RawQuery: c.Name,
		}
		cnStr = u.String()
	}
	if cnStr == `` {
		return errors.New(`Database driver not supported yet: ` + c.Driver)
	}
	db, e := sql.Open(c.Driver, cnStr)
	if e != nil {
		return e
	}
	if e := db.Ping(); e != nil {
		return e
	}
	c.db = db
	return nil
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
	tx, e := c.db.Begin()
	if e != nil {
		return nil, e
	}
	return &Tx{cn: c, tx: tx}, nil
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
func (c *Connection) MustExec(query string, params ...interface{}) *Result {
	result, e := c.Exec(query, params...)
	if e != nil {
		panic(e)
	}
	return result
}

//Exec ...
func (c *Connection) Exec(query string, params ...interface{}) (*Result, error) {
	tx, e := c.Begin()
	if e != nil {
		return nil, e
	}
	defer tx.Recover()
	return tx.Exec(query, params...)
}

//MustGet ...
func (c *Connection) MustGet(query string, params ...interface{}) Resultset {
	rs, e := c.Get(query, params...)
	if e != nil {
		panic(e)
	}
	return rs
}

//Get ...
func (c *Connection) Get(query string, params ...interface{}) (Resultset, error) {
	tx, e := c.Begin()
	if e != nil {
		return nil, e
	}
	defer tx.Recover()

	return tx.Get(query, params...)
}

//GetStruct ...
func (c *Connection) GetStruct(dest interface{}, query string, params ...interface{}) error {
	tx, e := c.Begin()
	if e != nil {
		return e
	}
	defer tx.Recover()

	return tx.GetStruct(dest, query, params...)
}

//MustSelect ...
func (c *Connection) MustSelect(query string, params ...interface{}) []Resultset {
	rs, e := c.Select(query, params...)
	if e != nil {
		panic(e)
	}
	return rs
}

//Select ...
func (c *Connection) Select(query string, params ...interface{}) ([]Resultset, error) {
	tx, e := c.Begin()
	if e != nil {
		return nil, e
	}
	defer tx.Recover()
	return tx.Select(query, params...)
}

//SelectStruct ...
func (c *Connection) SelectStruct(dest interface{}, query string, params ...interface{}) error {
	tx, e := c.Begin()
	if e != nil {
		return e
	}
	defer tx.Recover()
	return tx.SelectStruct(dest, query, params...)
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
	var names []string
	var questionMarks []string
	var values []interface{}

	for name, value := range dataMap {
		names = append(names, name)
		values = append(values, value)
		questionMarks = append(questionMarks, `?`)
	}
	var buffer bytes.Buffer
	buffer.WriteString(`INSERT INTO `)
	buffer.WriteString(tableName)
	buffer.WriteString(`(` + strings.Join(names, `, `) + `)`)
	buffer.WriteString(` VALUES(` + strings.Join(questionMarks, `, `) + `)`)

	return c.Exec(buffer.String(), values...)
}

//GetEnumValues ...
func (c *Connection) GetEnumValues(field string) ([]string, error) {
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

//Tx create new Tx when parameter tx is nil and the new Tx will have autocommit enabled. If parameter tx is not null then return tx from parameter
func (c *Connection) Tx(tx *Tx) *Tx {
	if tx == nil {
		return &Tx{cn: c}
	}
	return tx
}

func (c *Connection) wrapMsgErr(msg string) SQLError {
	return &sqlError{driver: c.Driver, msg: msg}
}

func (c *Connection) wrapErr(e error) SQLError {
	if e == nil {
		return nil
	}
	return &sqlError{driver: c.Driver, msg: e.Error()}
}

//Connect ...
func Connect(driver, host string, port int, username, password, name string) (*Connection, error) {
	cn, e := NewConnection(driver, host, port, username, password, name)
	if e != nil {
		return nil, e
	}
	if e := cn.Connect(); e != nil {
		return nil, e
	}
	lastCn = cn
	return cn, nil
}

//NewConnection ..
func NewConnection(driver, host string, port int, username, password, name string) (*Connection, error) {
	if port < 0 || port > 65535 {
		return nil, fmt.Errorf(`invalid port %d`, port)
	}
	return &Connection{Driver: driver, Hostname: host, Port: uint16(port), Username: username, Password: password, Name: name}, nil
}
