package router

import (
	"fmt"
	"gpsd-user-mgmt/src/config"
	"gpsd-user-mgmt/src/logger"
	"gpsd-user-mgmt/src/user"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func setupRouter(slogger *slog.Logger) *gin.Engine {
	router := gin.New()

	router.Use(logger.SlogMiddleware(slogger))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", user.List)
		v1.GET("/users/:id", user.Get)
		v1.POST("/users", user.Create)

		v1.PUT("/users/:id", user.Edit)

		v1.DELETE("/users/:id", user.Delete)

		v1.GET("/users/:id/incidents", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"API": "edit",
			})
		})

	}
	return router
}

func Run(config *config.Config, slogger *slog.Logger) bool {
	router := setupRouter(slogger)

	address := fmt.Sprintf(":%s", config.PORT)

	err := router.Run(address)
	if err != nil {
		slogger.Error("Unable to start server", "Error", err.Error())
	}
	slogger.Info("Starting server")

	return err == nil
}
