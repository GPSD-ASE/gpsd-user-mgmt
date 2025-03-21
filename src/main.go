package main

import (
	"gpsd-user-mgmt/src/config"
	"gpsd-user-mgmt/src/db"
	"gpsd-user-mgmt/src/logger"
	"gpsd-user-mgmt/src/router"
	"os"
)

func main() {
	config.LoadConfig()
	slogger := logger.SetupLogger()
	slogger.Info("Loaded configs")

	ok := db.Connect()
	if !ok {
		slogger.Error("Failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()
	slogger.Info("Connected to database")

	_, ok = router.Run(slogger)
	if !ok {
		slogger.Error("Failed to start server")
		os.Exit(2)
	}
}
