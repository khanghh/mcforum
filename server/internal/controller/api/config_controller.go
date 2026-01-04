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
	whatsNewMenuItem = model.MenuItem{
		Name:    locale.T("nav.whats-new"),
		URLPath: "/whats-new",
		LogoURL: "/icons/new.png",
	}
	recommendedMenuItem = model.MenuItem{
		Name:    locale.T("nav.recommended"),
		URLPath: "/recommended",
		LogoURL: "/icons/recommend.png",
	}
	followedMenuItem = model.MenuItem{
		Name:    locale.T("nav.followed"),
		URLPath: "/followed",
		LogoURL: "/icons/feed.png",
	}
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) Get() *web.JsonResult {
	configs := service.SysConfigService.GetAll()
	menuItems := []model.MenuItem{whatsNewMenuItem, recommendedMenuItem}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user != nil {
		menuItems = append(menuItems, followedMenuItem)
	}
	configResp := payload.BuildSysConfigResponse(configs)
	configResp.MenuItems = append(menuItems, configResp.MenuItems...)
	return web.JsonData(configResp)
}
