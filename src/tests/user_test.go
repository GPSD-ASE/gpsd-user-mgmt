package tests

import (
	"gpsd-user-mgmt/src/config"
	"gpsd-user-mgmt/src/db"
	"gpsd-user-mgmt/src/logger"
	"gpsd-user-mgmt/src/router"
	"testing"

	"github.com/gin-gonic/gin"
)

func startApp() *router.Engine {
	config := config.Load()
	slogger := logger.SetupLogger(config)
	ok := db.Connect(config)
	if !ok {
		slogger.Error("Failed to connect to database")
	}
	// db.CreateDatabase()

	gin.SetMode(gin.TestMode)
	r := router.SetupRouter(slogger)
	return r
}

func TestGetUser(t *testing.T) {
	r := startApp()

	t.Run("GetUsers", func(t *testing.T) {
		t.Run("Success", successGet(r))
		t.Run("Success - No Users", notFoundGet(r))
	})
}

func TestListUsers(t *testing.T) {
	r := startApp()

	t.Run("GetUsers", func(t *testing.T) {
		t.Run("Success", successList(r))
		t.Run("Success - Query Parameters", successListQuery(r))
		t.Run("Success - No Users", successListEmpty(r))
	})
}

func TestCreateUser(t *testing.T) {
	r := startApp()

	t.Run("GetUsers", func(t *testing.T) {
		t.Run("Success", successCreate(r))
	})
}

func TestUpdateUser(t *testing.T) {
	r := startApp()

	t.Run("GetUsers", func(t *testing.T) {
		t.Run("Success", successUpdate(r))
		t.Run("Success - No Users", notFoundUpdate(r))
	})
}

func TestDeleteUser(t *testing.T) {
	r := startApp()

	t.Run("GetUsers", func(t *testing.T) {
		t.Run("Success", successDelete(r))
		t.Run("Success - No Users", notFoundDelete(r))
	})
}
