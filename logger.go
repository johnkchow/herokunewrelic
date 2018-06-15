package main

import (
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

func init() {
	if os.Getenv("APP_ENV") == "production" {
		setupProductionLogger()
	} else {
		setupDevelopmentLogger()
	}
}
