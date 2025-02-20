package router

import (
	"fmt"
	"gpsd-user-mgmt/src/config"
	"gpsd-user-mgmt/src/logger"
	"gpsd-user-mgmt/src/user"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	router *gin.Engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.ServeHTTP(w, req)
}

func SetupRouter(slogger *slog.Logger) *Engine {
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
				"API": "Uninmplemented",
			})
		})

	}
	return &Engine{router: router}
}

func Run(config *config.Config, slogger *slog.Logger) (*Engine, bool) {
	router := SetupRouter(slogger)

	address := fmt.Sprintf(":%s", config.PORT)

	err := router.router.Run(address)
	if err != nil {
		slogger.Error("Unable to start server", "Error", err.Error())
		return nil, false
	}
	slogger.Info("Starting server")

	return router, true
}
