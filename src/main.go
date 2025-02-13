package main

import (
	"gpsd-user-mgmt/config"
	"gpsd-user-mgmt/db"
	"gpsd-user-mgmt/logger"
	"gpsd-user-mgmt/router"
	"os"
)

func main() {
	config := config.Load()
	slogger := logger.SetupLogger(config)
	slogger.Info("Loaded configs")

	ok := db.Connect(config)
	if !ok {
		slogger.Error("Failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()
	slogger.Info("Connected to database")

	_, ok = router.Run(config, slogger)
	if !ok {
		slogger.Error("Failed to start server")
		os.Exit(2)
	}
}
