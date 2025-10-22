package logging

import (
	"os"
	"sync"

	"go.uber.org/zap"
)

var (
	Base *zap.Logger
	once sync.Once
)

func Init() {
	once.Do(func() {
		var logger *zap.Logger
		var err error

		if os.Getenv("ENV") == "development" {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction()
		}

		if err != nil {
			panic("failed to initialize zap logger: " + err.Error())
		}
		Base = logger
	})
}

func GetBase() *zap.Logger {
	Init()
	return Base
}
