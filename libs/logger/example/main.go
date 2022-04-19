package main

import (
	"context"

	"bitbucket.org/efishery/go-efishery/libs/logger"
)

func main() {
	log := logger.Init(logger.DefaultLoggerOption)

	log.Info("test info")
	log.InfoWithContext(context.Background(), "test info")
}
