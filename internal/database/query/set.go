package query

type Set struct {
	Key   string
	Value string
}

func (q *Set) GetType() string {
	return SetQuery
}
