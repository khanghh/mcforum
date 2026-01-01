package api

import (
	"bbs-go/internal/config"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/pkg/web"
	"log/slog"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/service"
)

type UploadController struct {
	Ctx iris.Context
}

func (c *UploadController) Post() (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return nil, errs.ErrUnauthorized
	}
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return nil, err
	}

	uploadURL := config.Instance().Uploader.SUpload.UploadURL
	if uploadURL == "" {
		return nil, errs.ErrInternalServer
	}

	file, header, err := c.Ctx.FormFile("avatar")
	if err != nil {
		return nil, errs.ErrBadRequest
	}

	result, err := service.UploadService.UploadStream(user, file, header.Filename)
	if err != nil {
		if errResp, ok := err.(*errs.ResponseError); ok {
			return nil, errResp
		}
		slog.Error("Upload error", slog.Any("err", err))
		return nil, errs.ErrInternalServer
	}

	resp := payload.UploadResponse{
		URL:       result.URL,
		Size:      result.Size,
		MimeType:  result.MimeType,
		FileName:  result.FileName,
		CreatedAt: result.CreatedAt,
	}
	return web.JsonData(resp), nil
}
