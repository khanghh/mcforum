package admin

import (
	"strconv"

	"bbs-go/common/dates"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type LinkController struct {
	Ctx iris.Context
}

func (c *LinkController) GetBy(id int64) *web.JsonResult {
	t := service.LinkService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *LinkController) GetList() *web.JsonResult {
	list, paging := service.LinkService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("status").LikeByReq("title").LikeByReq("url").PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *LinkController) PostCreate() *web.JsonResult {
	t := &model.Link{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}
	t.CreateTime = dates.NowTimestamp()

	err = service.LinkService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *LinkController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := service.LinkService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = service.LinkService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
