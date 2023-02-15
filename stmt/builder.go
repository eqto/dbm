package stmt

type Builder struct {
}

func (Builder) Select(fields string) *SelectFields {
	f := parseFields(fields)
	return &SelectFields{value: f}
}

func (Builder) InsertInto(table, fields string) *Insert {
	f := parseFields(fields)
	return &Insert{table: table, fields: f}
}

func (Builder) Update(table string) *Update {
	return &Update{table: table}
}

func (Builder) DeleteFrom(table string) *Delete {
	return &Delete{table: table}
}

func Build() Builder {
	return Builder{}
}
