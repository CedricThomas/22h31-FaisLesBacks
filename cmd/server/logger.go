package main

import "github.com/sirupsen/logrus"

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	return logger
}
