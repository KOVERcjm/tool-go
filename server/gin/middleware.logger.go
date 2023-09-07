package gin

import (
	"bytes"
	"reflect"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kovercjm/tool-go/logger"
)

func APILogging(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		logger.Debug("gin api call incoming", "time", start, "method", c.Request.Method,
			"path", path, "rest-path", c.FullPath(), "query", raw, "ip", c.ClientIP(), "user-agent", c.Request.UserAgent())

		// Process request
		c.Next()

		// Stop timer

		if raw != "" {
			path = path + "?" + raw
		}
		bw := &bodyWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = bw
		logger.Debug("gin api call finished", "timestamp", start, "method", c.Request.Method,
			"path", path, "query", raw, "latency", time.Since(start), "function", runtime.FuncForPC(reflect.ValueOf(c.Handler()).Pointer()).Name(),
			"status-code", c.Writer.Status(), "body", bw.body.String(), "bodySize", bw.ResponseWriter.Size())
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (bw bodyWriter) Write(b []byte) (int, error) {
	bw.body.Write(b)
	return bw.ResponseWriter.Write(b)
}
