package stmt

import (
	"regexp"
	"strings"
)

type Fields struct {
	names   []string
	aliases map[uint8]string
	values  map[uint8]string //used by Insert
}

func (f *Fields) Names() []string {
	return f.names
}

func (f *Fields) AliasByIndex(i int) string {
	if alias, ok := f.aliases[uint8(i)]; ok {
		return alias
	}
	return ``
}

func (f *Fields) ValueByIndex(i int) string {
	if value, ok := f.values[uint8(i)]; ok {
		return value
	}
	return ``
}

func parseFields(fields string) Fields {
	f := Fields{}
	split := strings.Split(strings.TrimSpace(fields), `,`)
	regexAs := regexp.MustCompile(`(?Uis)\s+AS\s+`)
	for _, str := range split {
		split := regexAs.Split(str, 2)
		if len(split) == 2 {
			f.names = append(f.names, strings.TrimSpace(split[0]))
			if f.aliases == nil {
				f.aliases = make(map[uint8]string)
			}
			f.aliases[uint8(len(f.names)-1)] = strings.TrimSpace(split[1])
		} else {
			f.names = append(f.names, strings.TrimSpace(str))
		}
	}
	return f
}
