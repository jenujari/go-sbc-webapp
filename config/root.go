package config

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)

func init() {
	// init log system
	logger = log.Default()
	logger.SetOutput(os.Stdout)
}

func GetLogger() *log.Logger {
	return logger
}
