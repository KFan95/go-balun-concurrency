package database

import (
	"database/internal/database/compute"
	queryPackage "database/internal/database/query"
	"database/internal/database/storage"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type DatabaseInterface interface {
	Execute(string) (*string, error)
}

type Database struct {
	parserLayer  compute.ParserInterface
	storageLayer storage.StorageInterface
	logger       *zap.Logger
}

func NewDatabase(parserLayer compute.ParserInterface, storageLayer storage.StorageInterface, logger *zap.Logger) (DatabaseInterface, error) {
	if parserLayer == nil {
		return nil, errors.New("parser is required")
	}

	if storageLayer == nil {
		return nil, errors.New("parser is required")
	}

	if logger == nil {
		return nil, errors.New("logger is required")
	}

	return &Database{
		parserLayer:  parserLayer,
		storageLayer: storageLayer,
		logger:       logger,
	}, nil
}

func (d *Database) Execute(rawQuery string) (*string, error) {
	query, err := d.parserLayer.Parse(rawQuery)

	if err != nil {
		return nil, err
	}

	switch query := query.(type) {
	case *queryPackage.Set:
		return d.executeSetQuery(query)

	case *queryPackage.Get:
		return d.executeGetQuery(query)

	case *queryPackage.Del:
		return d.executeDelQuery(query)
	}

	return nil, nil
}

func (d *Database) executeGetQuery(query *queryPackage.Get) (*string, error) {
	value, err := d.storageLayer.Get(query.Key)

	if err != nil {
		return nil, err
	}

	var response string
	if value != nil {
		response = fmt.Sprintf("[ok] %s", *value)
	} else {
		response = "[not found]"
	}

	return &response, nil
}

func (d *Database) executeSetQuery(query *queryPackage.Set) (*string, error) {
	err := d.storageLayer.Set(query.Key, query.Value)

	if err != nil {
		return nil, err
	}

	response := "[ok]"

	return &response, nil
}

func (d *Database) executeDelQuery(query *queryPackage.Del) (*string, error) {
	err := d.storageLayer.Del(query.Key)

	if err != nil {
		return nil, err
	}

	response := "[ok]"

	return &response, nil
}
