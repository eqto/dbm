package stmt

import "strings"

type InsertValues struct {
	insert *Insert
}

func (i *InsertValues) Values(values string) *Insert {
	i.insert.Values(values)
	return i.insert
}

type Insert struct {
	table  string
	fields Fields
	output string //used by sqlserver
}

func (i *Insert) Output(output string) *InsertValues {
	i.output = output
	return &InsertValues{i}
}

func (i *Insert) Values(values string) *Insert {
	split := strings.Split(values, `,`)
	idx := 0
	prev := ``
	str := ``
	for _, s := range split {
		if prev != `` {
			str = prev + `,` + s
		} else {
			str = s
		}
		if strings.Count(str, `(`) != strings.Count(str, `)`) {
			prev = str
			continue
		}
		prev = ``
		str = strings.TrimSpace(str)
		if str != `?` {
			if i.fields.values == nil {
				i.fields.values = make(map[uint8]string)
			}
			i.fields.values[uint8(idx)] = str
		}
		idx++
	}
	return i
}
