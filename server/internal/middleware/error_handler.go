package middleware

import (
	"bbs-go/internal/service"
	"bbs-go/web"
	"log/slog"

	"github.com/kataras/iris/v12"
)

func ErrorHandler(ctx iris.Context, err error) {
	if err != nil {
		slog.Debug("Unhandled error occurred", "path", ctx.Path(), "error", err)
		if service.IsDatabaseError(err) {
			ctx.JSON(web.JsonError(service.ErrInternalServer))
			return
		}
		ctx.JSON(web.JsonError(err))
	}
}
