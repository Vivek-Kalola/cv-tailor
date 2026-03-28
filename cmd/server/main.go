package main

import (
	"cv-tailoring/internal/api"
	"cv-tailoring/internal/db"

	"github.com/sirupsen/logrus"
)

func main() {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	database, err := db.NewDB(logger)
	if err != nil {
		logger.Fatal(err)
	}

	err = api.NewServer(8443, database, logger)
	if err != nil {
		logger.Fatal(err)
	}
}
