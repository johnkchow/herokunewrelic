package main

import (
	"flag"
	"go.uber.org/zap"
	"os"
)

var logger *zap.SugaredLogger

func setupDevelopmentLogger() {
	log, _ := zap.NewDevelopment()
	logger = log.Sugar()
}

func setupProductionLogger() {
	log, _ := zap.NewProduction()
	logger = log.Sugar()
}

func setupNilLogger() {
	log := zap.NewNop()
	logger = log.Sugar()
}

func setupTestLogger() {
	log := zap.NewExample()
	logger = log.Sugar()
}

func init() {
	// TODO: actually honor the LOG_LEVEL env var...
	if os.Getenv("APP_ENV") == "production" {
		setupProductionLogger()
	} else if flag.Lookup("test.v") != nil {
		if os.Getenv("LOG_LEVEL") != "" {
			setupTestLogger()
		} else {
			setupNilLogger()
		}
	} else {
		setupDevelopmentLogger()
	}
}
