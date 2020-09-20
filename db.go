package db

import (
	"fmt"
	"regexp"
)

var (
	lastCn             *Connection
	regexStringColType *regexp.Regexp
)

//NewConnection ...
func NewConnection(host string, port int, username, password, name string) (*Connection, error) {
	cn, e := NewEmptyConnection(host, port, username, password, name)
	if e != nil {
		return nil, e
	}
	if e := cn.Connect(); e != nil {
		return nil, e
	}
	lastCn = cn
	return cn, nil
}

//NewEmptyConnection ...
func NewEmptyConnection(host string, port int, username, password, name string) (*Connection, error) {
	if port < 0 || port > 65535 {
		return nil, fmt.Errorf(`invalid port %d`, port)
	}
	lastCn = &Connection{Hostname: host, Port: uint16(port), Username: username, Password: password, Name: name}
	return lastCn, nil
}

func getRegex() *regexp.Regexp {
	if regexStringColType == nil {
		regexStringColType, _ = regexp.Compile(`(?i)^.*CHAR|.*TEXT$`)
	}
	return regexStringColType
}
