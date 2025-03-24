package api

import (
	"bbs-go/web"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/controller/response"
	"bbs-go/internal/service"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *web.JsonResult {
	configs := service.SysConfigService.GetAll()
	return web.JsonData(response.BuildSysConfigResponse(configs))
}
