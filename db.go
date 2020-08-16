package db

import (
	"database/sql"
	"fmt"
	"regexp"
)

var (
	lastCn             *Connection
	regexStringColType *regexp.Regexp
)

//NewConnection ...
func NewConnection(host string, port int, username, password, name string) (*Connection, error) {
	db, e := sql.Open(`mysql`,
		fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local`, username, password, host, port, name))

	if e != nil {
		return nil, e
	}
	if e := db.Ping(); e != nil {
		return nil, e
	}

	lastCn = &Connection{db: db, Hostname: host, Port: port, Username: username, Password: password, Name: name}
	return lastCn, nil
}

func getRegex() *regexp.Regexp {
	if regexStringColType == nil {
		regexStringColType, _ = regexp.Compile(`(?i)^.*CHAR|.*TEXT$`)
	}
	return regexStringColType
}
