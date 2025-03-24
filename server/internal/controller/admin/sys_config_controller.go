package admin

import (
	"strconv"

	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/controller/response"
	"bbs-go/internal/service"
)

type SysConfigController struct {
	Ctx iris.Context
}

func (c *SysConfigController) GetBy(id int64) *web.JsonResult {
	t := service.SysConfigService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *SysConfigController) GetList() *web.JsonResult {
	list, paging := service.SysConfigService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *SysConfigController) GetAll() *web.JsonResult {
	configs := service.SysConfigService.GetAll()
	return web.JsonData(response.BuildSysConfigResponse(configs))
}

func (c *SysConfigController) PostSave() *web.JsonResult {
	body, err := c.Ctx.GetBody()
	if err != nil {
		return web.JsonError(err)
	}
	if err := service.SysConfigService.SetAll(string(body)); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
