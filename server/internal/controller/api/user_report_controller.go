package api

import (
	"bbs-go/internal/model"
	"bbs-go/internal/service"

	"bbs-go/common/dates"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"
)

type UserReportController struct {
	Ctx iris.Context
}

func (c *UserReportController) PostSubmit() *web.JsonResult {
	var (
		dataId, _ = params.FormValueInt64(c.Ctx, "dataId")
		dataType  = params.FormValue(c.Ctx, "dataId")
		reason    = params.FormValue(c.Ctx, "reason")
	)
	report := &model.UserReport{
		DataID:     dataId,
		DataType:   dataType,
		Reason:     reason,
		CreateTime: dates.NowTimestamp(),
	}

	if user := service.UserTokenService.GetCurrent(c.Ctx); user != nil {
		report.UserID = user.ID
	}
	service.UserReportService.Create(report)
	return web.JsonSuccess()
}
