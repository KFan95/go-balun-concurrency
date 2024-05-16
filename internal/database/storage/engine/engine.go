package engine

type EngineInterface interface {
	Get(string) (*string, error)
	Set(string, string) error
	Del(string) error
}
