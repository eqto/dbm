package stmt

import "strings"

type Insert struct {
	table  string
	fields Fields
}

func (i *Insert) Values(values string) *Insert {
	split := strings.Split(values, `,`)
	for idx, str := range split {
		str = strings.TrimSpace(str)
		if str != `?` {
			if i.fields.values == nil {
				i.fields.values = make(map[uint8]string)
			}
			i.fields.values[uint8(idx)] = str
		}
	}
	return i
}
