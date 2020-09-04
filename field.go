package db

type field struct {
	name, alias string
}

//String ...
func (f *field) String() string {
	if f.alias == `` {
		return f.name
	}
	return f.name + ` AS ` + f.alias
}
