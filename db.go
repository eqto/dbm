package db

import "log"

var (
	logE = log.Println
)

func Log(e func(...interface{})) {
	logE = e
}

//Connect ...
func Connect(driver, host string, port int, username, password, name string) (*Connection, error) {
	cn, e := newConnection(driver, host, port, username, password, name)
	if e != nil {
		return nil, e
	}
	if e := cn.Connect(); e != nil {
		return nil, e
	}
	return cn, nil
}

func Query(driverName string) *QueryBuilder {
	drv, e := getDriver(driverName)
	if e != nil {
		logE(e)
	}
	return &QueryBuilder{drv: drv}
}
