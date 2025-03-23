package admin

import (
	"bbs-go/internal/model"
	"bbs-go/internal/service"
	"strconv"

	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"
)

type UserFeedController struct {
	Ctx iris.Context
}

func (c *UserFeedController) GetBy(id int64) *web.JsonResult {
	t := service.UserFeedService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *UserFeedController) AnyList() *web.JsonResult {
	list, paging := service.UserFeedService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *UserFeedController) PostCreate() *web.JsonResult {
	t := &model.UserFeed{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = service.UserFeedService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *UserFeedController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := service.UserFeedService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = service.UserFeedService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
