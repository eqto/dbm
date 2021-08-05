package stmt

type WhereParam struct {
	Condition string
	Or        bool //OR, default: false = AND
}
