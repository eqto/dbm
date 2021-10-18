package dbm

type SQLStatement struct {
	statement string
}

func SQL(statement string) SQLStatement {
	return SQLStatement{statement}
}
