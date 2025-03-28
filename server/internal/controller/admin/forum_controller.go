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

type ForumController struct {
	Ctx iris.Context
}

func (c *ForumController) GetBy(id int64) *web.JsonResult {
	t := service.ForumService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *ForumController) GetList() *web.JsonResult {
	list, paging := service.ForumService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *ForumController) PostCreate() *web.JsonResult {
	t := &model.Forum{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}
	t.CreateTime = dates.NowTimestamp()
	err = service.ForumService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *ForumController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := service.ForumService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	err = service.ForumService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *ForumController) GetNodes() *web.JsonResult {
	list := service.ForumService.GetAll()
	return web.JsonData(list)
}
