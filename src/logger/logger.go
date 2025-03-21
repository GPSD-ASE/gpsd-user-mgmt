package logger

import (
	"gpsd-user-mgmt/src/config"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupLogger() *slog.Logger {
	var lvl slog.Level
	switch config.USER_MGMT_ENV {
	case "PRODUCTION":
		lvl = slog.LevelInfo
	case "TEST":
		lvl = slog.LevelWarn
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     lvl,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
	return logger
}

func SlogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		var params map[string][]string = c.Request.URL.Query()

		logger.Info("Received payload",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("ip", c.ClientIP()),
			slog.Any("parameters", params),
		)
	}
}
