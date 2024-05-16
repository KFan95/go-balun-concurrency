package query

type Get struct {
	Key string
}

func (q *Get) GetType() string {
	return GetQuery
}
