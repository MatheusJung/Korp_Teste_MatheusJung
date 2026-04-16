package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latencia", time.Since(start)),
			zap.String("ip", c.ClientIP()),
		)
	}
}

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recuperado", zap.Any("erro", err))
				c.JSON(http.StatusInternalServerError, gin.H{
					"erro":      "erro interno inesperado",
					"timestamp": time.Now().UTC(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
