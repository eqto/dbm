package stmt

type Update struct {
	whereStatement
	table     string
	namevalue []string
	wheres    []WhereParam
	order     string
	count     int
}

func (u *Update) where(param WhereParam) {
	u.wheres = append(u.wheres, param)
}

// Set keyvalue pair to set
// Ex:
//
//	UPDATE books SET title = ?, publisher = ?
//	dbx.Update(`books`).Set(`title = ?`, `publisher = ?`).Set(`publisher = ?`)
func (u *Update) Set(namevalues ...string) *UpdateFields {
	u.namevalue = append(u.namevalue, namevalues...)
	return &UpdateFields{stmt: u}
}

func (u *Update) orderBy(orderBy string) {
	u.order = orderBy
}

func (u *Update) limit(count int) {
	u.count = count
}

type UpdateFields struct {
	stmt *Update
}

func (u *UpdateFields) Where(condition string) *UpdateWhere {
	u.stmt.wheres = []WhereParam{{condition, false}}
	return &UpdateWhere{u.stmt}
}
