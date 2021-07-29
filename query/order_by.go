package query

type OrderBy struct {
	orders    []string
	stmt      interface{}
	limitFunc func(...int) interface{}
}

func (o *OrderBy) Limit(num ...int) interface{} {
	return o.limitFunc(num...)
}
