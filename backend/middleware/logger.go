package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a minimal request-logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				gin.DefaultErrorWriter.Write([]byte(e + "\n")) //nolint:errcheck
			}
		}

		_ = latency // silence unused-var warning while keeping the variable for future logging
		gin.DefaultWriter.Write([]byte( //nolint:errcheck
			time.Now().Format(time.RFC3339) + " " +
				clientIP + " " + method + " " + path + " " +
				strconv.Itoa(status) + "\n",
		))
	}
}
