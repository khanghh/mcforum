package api

import (
	"bbs-go/web"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/service"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *web.JsonResult {
	config := service.SysConfigService.GetConfig()
	return web.JsonData(config)
}
