package query

import "strings"

type Field struct {
	Name  string
	Alias string
}

func parseFields(query string) []Field {
	fields := []Field{}
	split := strings.Split(query, `,`)
	for _, s := range split {
		split := strings.SplitN(s, ` AS `, 2)
		field := Field{}
		if len(split) == 2 {
			field.Name, field.Alias = strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
		} else {
			field.Name = strings.TrimSpace(s)
		}
		fields = append(fields, field)
	}
	return fields
}
