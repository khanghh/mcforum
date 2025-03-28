package admin

import (
	"bbs-go/internal/model"
	"bbs-go/internal/service"
	"strconv"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"
)

type UserFollowController struct {
	Ctx iris.Context
}

func (c *UserFollowController) GetBy(id int64) *web.JsonResult {
	t := service.UserFollowService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *UserFollowController) GetList() *web.JsonResult {
	list, paging := service.UserFollowService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *UserFollowController) PostCreate() *web.JsonResult {
	t := &model.UserFollow{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = service.UserFollowService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *UserFollowController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := service.UserFollowService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = service.UserFollowService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
