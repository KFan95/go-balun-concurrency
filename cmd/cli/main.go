package main

import (
	"bufio"
	database "database/internal"
	"database/internal/database/compute"
	"database/internal/database/storage"
	"database/internal/database/storage/engine/memory"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger := createLogger()
	defer logger.Sync()

	logger.Info("start cli database")

	parserLayer := createParserLayer(logger)
	storageLayer := createStorageLayer(logger)

	db, err := database.NewDatabase(parserLayer, storageLayer, logger)

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("database cli > ")
		command, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		result, resultErr := db.Execute(command)

		if resultErr != nil {
			fmt.Printf("Error: %s\n", resultErr.Error())

			continue
		}

		fmt.Println(*result)
	}
}

func createLogger() *zap.Logger {
	logger, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}

	return logger
}

func createParserLayer(logger *zap.Logger) compute.ParserInterface {
	parserLayer, err := compute.NewParser(logger)

	if err != nil {
		panic(err)
	}

	return parserLayer
}

func createStorageLayer(logger *zap.Logger) storage.StorageInterface {
	engine, err := memory.NewMemoryEngine(logger)

	if err != nil {
		panic(err)
	}

	storageLayer, err := storage.NewStorage(engine, logger)

	if err != nil {
		panic(err)
	}

	return storageLayer
}
