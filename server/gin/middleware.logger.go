package gin

import (
	"bytes"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kovercjm/tool-go/logger"
)

func APILogging(logger logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		logger.Debug("gin api call incoming", "timestamp", start, "method", ctx.Request.Method,
			"path", path, "query", raw, "ip", ctx.ClientIP(), "user-agent", ctx.Request.UserAgent())

		// Process request
		ctx.Next()

		if raw != "" {
			path = path + "?" + raw
		}
		bw := &bodyWriter{ResponseWriter: ctx.Writer, body: &bytes.Buffer{}}
		ctx.Writer = bw
		funcName := runtime.FuncForPC(reflect.ValueOf(ctx.Handler()).Pointer()).Name()
		pieces := strings.Split(funcName, ".")
		funcName = pieces[len(pieces)-1]
		// Stop timer
		logger.Debug("gin api call finished", "timestamp", start, "method", ctx.Request.Method,
			"path", path, "query", raw, "latency", time.Since(start), "function", funcName,
			"status-code", ctx.Writer.Status(), "body", bw.body.String(), "bodySize", bw.ResponseWriter.Size())
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
