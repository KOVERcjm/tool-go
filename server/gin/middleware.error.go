package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kovercjm/tool-go/server/api"
	"google.golang.org/grpc/status"
	"net/http"
)

func ErrorFormatter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if err := ctx.Errors.ByType(gin.ErrorTypeBind).Last(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, &api.Error{
				Code:    http.StatusText(http.StatusBadRequest),
				Message: err.Error(),
			})
			return
		}

		if err := ctx.Errors.ByType(gin.ErrorTypePrivate).Last().Err; err != nil {
			ctx.Header("Content-Type", "application/json; charset=utf-8; error=true")

			httpStatus := http.StatusInternalServerError
			if errors.As(err, &api.Error{}) {
				httpStatus = http.StatusBadRequest
			}
			if grpcStatus, ok := status.FromError(err); ok {
				httpStatus = runtime.HTTPStatusFromCode(grpcStatus.Code())
			}
			ctx.AbortWithStatusJSON(httpStatus, &api.Error{
				Code:    http.StatusText(httpStatus),
				Message: err.Error(),
			})
			return
		}
	}
}
