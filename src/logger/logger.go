package logger

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	})

	logger := slog.New(handler)
	return logger
}

func SlogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		var params map[string][]string = c.Request.URL.Query()

		logger.Info("Recieved payload",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("ip", c.ClientIP()),
			slog.Any("parameters", params),
		)
	}
}
