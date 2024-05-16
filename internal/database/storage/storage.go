package storage

import (
	"database/internal/database/storage/engine"
	"errors"

	"go.uber.org/zap"
)

type StorageInterface interface {
	Get(string) (*string, error)
	Set(string, string) error
	Del(string) error
}

type Storage struct {
	engine engine.EngineInterface
	logger *zap.Logger
}

func NewStorage(engine engine.EngineInterface, logger *zap.Logger) (StorageInterface, error) {
	if engine == nil {
		return nil, errors.New("engine is required")
	}

	if logger == nil {
		return nil, errors.New("logger is required")
	}

	return &Storage{
		engine: engine,
		logger: logger,
	}, nil
}

func (s *Storage) Get(key string) (*string, error) {
	return s.engine.Get(key)
}

func (s *Storage) Set(key string, value string) error {
	return s.engine.Set(key, value)
}

func (s *Storage) Del(key string) error {
	return s.engine.Del(key)
}
