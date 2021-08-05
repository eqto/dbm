package stmt

const (
	InnerJoin = iota
	LeftJoin
	RightJoin
)

type Table struct {
	Name   string
	Alias  string
	Join   uint8
	JoinOn string
}

type Tables struct {
	names   []string
	aliases map[uint8]string
	joinOn  []string
	join    map[uint8]uint8 //key: indexOf names, value = const of InnerJoin/LeftJoin/RightJoin
}

func (t *Tables) Names() []string {
	return t.names
}

func (t *Tables) TableByIndex(i int) Table {
	table := Table{}
	if i > len(t.names)-1 {
		return table
	}
	table.Name = t.names[i]
	if alias, ok := t.aliases[uint8(i)]; ok {
		table.Alias = alias
	}
	if kind, ok := t.join[uint8(i)]; ok {
		table.Join = kind
	}

	if i > 0 && i < len(t.joinOn)+1 {
		table.JoinOn = t.joinOn[i-1]
	}
	return table
}

func (t *Tables) add(name, alias string, joinKind int, joinOn string) {
	t.names = append(t.names, name)
	if alias != `` {
		if t.aliases == nil {
			t.aliases = make(map[uint8]string)
		}
		t.aliases[uint8(len(t.names)-1)] = alias
	}
	if joinOn != `` {
		t.joinOn = append(t.joinOn, joinOn)
	}
	if joinKind != InnerJoin {
		if t.join == nil {
			t.join = make(map[uint8]uint8)
		}
		t.join[uint8(len(t.names)-1)] = uint8(joinKind)
	}
}
