package stmt

type whereStatement interface {
	where(WhereParam)
	orderBy(string)
	limit(int)
}

type WhereParam struct {
	Condition string
	Or        bool //OR, default: false = AND
}
