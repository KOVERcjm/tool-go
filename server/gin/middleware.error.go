package gin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"

	"github.com/kovercjm/tool-go/server/api"
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

			apiError := &api.Error{
				HTTPStatus: http.StatusInternalServerError,
				Code:       http.StatusText(http.StatusInternalServerError),
				Message:    err.Error(),
			}
			if errors.As(err, apiError) {
				if apiError.HTTPStatus == 0 {
					apiError.HTTPStatus = http.StatusBadRequest
				}
				if apiError.Code == "" {
					apiError.Code = http.StatusText(apiError.HTTPStatus)
				}
				ctx.AbortWithStatusJSON(apiError.HTTPStatus, apiError)
				return
			}
			if grpcStatus, ok := status.FromError(err); ok {
				apiError.HTTPStatus = runtime.HTTPStatusFromCode(grpcStatus.Code())
				apiError.Code = http.StatusText(runtime.HTTPStatusFromCode(grpcStatus.Code()))
			}
			ctx.AbortWithStatusJSON(apiError.HTTPStatus, apiError)
			return
		}
	}
}
