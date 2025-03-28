package api

import (
	"bbs-go/pkg/web"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/controller/payload"
	"bbs-go/internal/service"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *web.JsonResult {
	configs := service.SysConfigService.GetAll()
	return web.JsonData(payload.BuildSysConfigResponse(configs))
}
