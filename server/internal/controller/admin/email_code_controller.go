package admin

import (
	"bbs-go/internal/model"
	"bbs-go/internal/service"
	"strconv"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"
)

type EmailCodeController struct {
	Ctx iris.Context
}

func (c *EmailCodeController) GetBy(id int64) *web.JsonResult {
	t := service.EmailCodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *EmailCodeController) GetList() *web.JsonResult {
	list, paging := service.EmailCodeService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *EmailCodeController) PostCreate() *web.JsonResult {
	t := &model.EmailCode{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = service.EmailCodeService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *EmailCodeController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := service.EmailCodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = service.EmailCodeService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
