package middleware

import (
	"bbs-go/internal/errs"
	"bbs-go/pkg/web"
	"log/slog"

	"github.com/kataras/iris/v12"
)

func ErrorHandler(ctx iris.Context, err error) {
	if err != nil {
		slog.Debug("Unhandled error occurred", "path", ctx.Path(), "error", err)
		if errs.IsDatabaseError(err) {
			ctx.JSON(web.JsonError(errs.ErrInternalServer))
			return
		}
		ctx.JSON(web.JsonError(err))
	}
}
