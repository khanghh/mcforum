package api

import (
	"bbs-go/web"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/services"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *web.JsonResult {
	config := services.SysConfigService.GetConfig()
	return web.JsonData(config)
}
