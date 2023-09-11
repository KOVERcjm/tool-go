package gin

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequestToken(ignore ...interface{}) gin.HandlerFunc {
	ignoreFuncs := map[string]struct{}{}
	for _, handler := range ignore {
		funcName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		pieces := strings.Split(funcName, ".")
		funcName = pieces[len(pieces)-1]
		ignoreFuncs[funcName] = struct{}{}
	}
	return func(ctx *gin.Context) {
		funcName := runtime.FuncForPC(reflect.ValueOf(ctx.Handler()).Pointer()).Name()
		pieces := strings.Split(funcName, ".")
		funcName = pieces[len(pieces)-1]

		if strings.Contains(funcName, "func") {
			// Return 404 when there are no matching handlers
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if _, ok := ignoreFuncs[funcName]; ok {
			ctx.Next()
			return
		}

		token := strings.TrimSpace(ctx.GetHeader("Authorization"))
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			cookieToken, err := ctx.Cookie("token") // looking for token in cookie with key 'token'
			if err != nil {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			token = cookieToken
		}

		ctx.Set("token", token)
		ctx.Next()
	}
}
