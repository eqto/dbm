package db

import (
	"regexp"
)

var (
	lastCn             *Connection
	regexStringColType *regexp.Regexp
)

//NewConnection ...
func NewConnection(host string, port uint16, username, password, name string) (*Connection, error) {
	cn := NewEmptyConnection(host, port, username, password, name)
	if e := cn.Connect(); e != nil {
		return nil, e
	}
	lastCn = cn
	return cn, nil
}

//NewEmptyConnection ...
func NewEmptyConnection(host string, port uint16, username, password, name string) *Connection {
	lastCn = &Connection{Hostname: host, Port: port, Username: username, Password: password, Name: name}
	return lastCn
}

func getRegex() *regexp.Regexp {
	if regexStringColType == nil {
		regexStringColType, _ = regexp.Compile(`(?i)^.*CHAR|.*TEXT$`)
	}
	return regexStringColType
}
