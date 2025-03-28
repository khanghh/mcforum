package admin

import (
	"strconv"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type FavoriteController struct {
	Ctx iris.Context
}

func (c *FavoriteController) GetBy(id int64) *web.JsonResult {
	t := service.FavoriteService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *FavoriteController) GetList() *web.JsonResult {
	list, paging := service.FavoriteService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *FavoriteController) PostCreate() *web.JsonResult {
	t := &model.Favorite{}
	params.ReadForm(c.Ctx, t)

	err := service.FavoriteService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *FavoriteController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := service.FavoriteService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	params.ReadForm(c.Ctx, t)

	err = service.FavoriteService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
