package stmt

type Update struct {
	table     string
	namevalue []string
	wheres    []WhereParam
}

type UpdateFields struct {
	stmt *Update
}

//Set keyvalue pair to set
//Ex:
//	UPDATE books SET title = ?, publisher = ?
//	dbx.Update(`books`).Set(`title = ?`).Set(`publisher = ?`)
func (u *Update) Set(namevalue string) *UpdateFields {
	u.namevalue = append(u.namevalue, namevalue)
	return &UpdateFields{stmt: u}
}

func (u *UpdateFields) Set(keyvalue string) *UpdateFields {
	return u.stmt.Set(keyvalue)
}

func (u *UpdateFields) Where(condition string) *UpdateWhere {
	u.stmt.wheres = []WhereParam{{condition, false}}
	return &UpdateWhere{u.stmt}
}
