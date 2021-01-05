package db

import (
	"regexp"
)

const (
	DriverMySQL     = `mysql`
	DriverSQLServer = `sqlserver`
)

var (
	regexStringColType *regexp.Regexp
)

func getRegex() *regexp.Regexp {
	if regexStringColType == nil {
		regexStringColType, _ = regexp.Compile(`(?i)^.*CHAR|.*TEXT$`)
	}
	return regexStringColType
}
