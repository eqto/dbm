package stmt

type whereStatement interface {
	where(WhereParam)
}

type WhereParam struct {
	Condition string
	Or        bool //OR, default: false = AND
}
