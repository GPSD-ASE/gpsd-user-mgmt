package main

import (
	"gpsd-user-mgmt/src/config"
	"gpsd-user-mgmt/src/db"
	"gpsd-user-mgmt/src/logger"
	"gpsd-user-mgmt/src/router"
	"os"
)

func main() {
	slogger := logger.SetupLogger()
	config := config.Load()
	slogger.Info("Loaded configs")

	ok := db.Connect(config)
	if !ok {
		slogger.Error("Failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()
	slogger.Info("Connected to database")

	ok = router.Run(config, slogger)
	if !ok {
		slogger.Error("Failed to start server")
		os.Exit(2)
	}
}
