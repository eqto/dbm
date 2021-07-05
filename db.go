package db

//Connect ...
func Connect(driver, host string, port int, username, password, name string) (*Connection, error) {
	cn, e := newConnection(driver, host, port, username, password, name)
	if e != nil {
		return nil, e
	}
	if e := cn.Connect(); e != nil {
		return nil, e
	}
	lastCn = cn
	return cn, nil
}
