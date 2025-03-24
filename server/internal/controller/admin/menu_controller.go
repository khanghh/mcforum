package admin

import (
	"bbs-go/internal/controller/response"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"strconv"

	"bbs-go/common/dates"
	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"
)

type MenuController struct {
	Ctx iris.Context
}

func (c *MenuController) GetBy(id int64) *web.JsonResult {
	t := service.MenuService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(response.BuildMenu(t))
}

func (c *MenuController) GetTree() *web.JsonResult {
	list := service.MenuService.Find(params.NewSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "status",
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(response.BuildMenuSimpleTree(0, list))
}

func (c *MenuController) GetList() *web.JsonResult {
	list := service.MenuService.Find(params.NewSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "status",
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(response.BuildMenuTree(0, list))
}

func (c *MenuController) PostCreate() *web.JsonResult {
	t := &model.Menu{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if t.SortNo <= 0 {
		t.SortNo = service.MenuService.GetNextSortNo(t.ParentId)
	}
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := service.MenuService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *MenuController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := service.MenuService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := service.MenuService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *MenuController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	for _, id := range ids {
		service.MenuService.Updates(id, map[string]interface{}{
			"status":      constants.StatusDeleted,
			"update_time": dates.NowTimestamp(),
		})
	}
	return web.JsonSuccess()
}

func (c *MenuController) GetUser_menus() *web.JsonResult {
	user, err := service.UserTokenService.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}
	list := service.MenuService.GetUserMenus(user)
	return web.JsonData(response.BuildMenuTree(0, list))
}

func (c *MenuController) PostUpdate_sort() *web.JsonResult {
	var ids []int64
	if err := c.Ctx.ReadJSON(&ids); err != nil {
		return web.JsonError(err)
	}
	if err := service.MenuService.UpdateSort(ids); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
