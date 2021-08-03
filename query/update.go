package query

import "strings"

type UpdateStmt struct {
	table     Table
	condition *Condition
}

func (u *UpdateStmt) Set(query string) *UpdateStmt {
	split := strings.Split(query, `,`)
	for _, str := range split {
		u.set(str)
	}
	return u
}

func (u *UpdateStmt) Where(condition string) *Condition {
	if u.condition == nil {
		u.condition = &Condition{stmt: u}
	}
	u.condition.conditions = append(u.condition.conditions, condition)
	return u.condition
}

func (u *UpdateStmt) set(query string) {
	split := strings.SplitN(query, `=`, 2)
	field := Field{}
	if len(split) == 2 {
		field.Name, field.Placeholder = strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
	} else {
		field.Name = strings.TrimSpace(query)
	}
	u.table.Fields = append(u.table.Fields, field)
}

func NewUpdate(table string) *UpdateStmt {
	return &UpdateStmt{table: Table{Name: table}}
}
