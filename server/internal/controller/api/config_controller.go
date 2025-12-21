package api

import (
	"bbs-go/pkg/web"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/controller/payload"
	"bbs-go/internal/locale"
	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

var (
	defaultMeuItems = []model.MenuItem{
		{
			Name:    locale.T("nav.whats-new"),
			URLPath: "/whats-new",
			LogoURL: "/icons/new.png",
		},
		{
			Name:    locale.T("nav.recommended"),
			URLPath: "/recommended",
			LogoURL: "/icons/recommend.png",
		},
		{
			Name:    locale.T("nav.followed"),
			URLPath: "/followed",
			LogoURL: "/icons/feed.png",
		},
	}
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) Get() *web.JsonResult {
	configs := service.SysConfigService.GetAll()
	configResp := payload.BuildSysConfigResponse(configs)
	configResp.MenuItems = append(defaultMeuItems, configResp.MenuItems...)
	return web.JsonData(configResp)
}
