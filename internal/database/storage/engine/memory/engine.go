package memory

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

type MemoryEngine struct {
	mutex sync.RWMutex

	data   map[string]string
	logger *zap.Logger
}

func NewMemoryEngine(logger *zap.Logger) (*MemoryEngine, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	return &MemoryEngine{
		data:   make(map[string]string),
		logger: logger,
	}, nil
}

func (e *MemoryEngine) Get(key string) (*string, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	value, ok := e.data[key]

	if !ok {
		return nil, nil
	}

	return &value, nil
}

func (e *MemoryEngine) Set(key string, value string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.data[key] = value

	return nil
}

func (e *MemoryEngine) Del(key string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.data, key)

	return nil
}
