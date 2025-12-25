package admin

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"strconv"

	"bbs-go/common/dates"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"github.com/kataras/iris/v12"
)

type RoleController struct {
	Ctx iris.Context
}

func (c *RoleController) GetBy(id int64) *web.JsonResult {
	t := service.RoleService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *RoleController) GetList() *web.JsonResult {
	list, paging := service.RoleService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
		},
		params.QueryFilter{
			ParamName: "status",
		},
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "code",
			Op:        params.Like,
		},
	).Asc("sort_no").Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *RoleController) GetAll_roles() *web.JsonResult {
	list := service.RoleService.Find(sqls.NewCnd().Eq("status", constants.StatusActive).Asc("sort_no").Desc("id"))
	return web.JsonData(list)
}

func (c *RoleController) PostCreate() *web.JsonResult {
	t := &model.Role{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if service.RoleService.GetByCode(t.Code) != nil {
		return web.JsonErrorMsg("Role code already exists")
	}

	t.SortNo = service.RoleService.GetNextSortNo()
	t.Type = constants.RoleTypeCustom
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := service.RoleService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *RoleController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := service.RoleService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if t.Type == constants.RoleTypeSystem {
		return web.JsonErrorMsg("System roles cannot be edited")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if exists := service.RoleService.GetByCode(t.Code); exists != nil && exists.Id != t.Id {
		return web.JsonErrorMsg("Role code already exists")
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := service.RoleService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *RoleController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	for _, id := range ids {
		service.RoleService.Updates(id, map[string]interface{}{
			"status":      constants.StatusDeleted,
			"update_time": dates.NowTimestamp(),
		})
	}
	return web.JsonSuccess()
}

func (s *RoleController) GetRoles() *web.JsonResult {
	roles := service.RoleService.Find(sqls.NewCnd().Eq("status", constants.StatusActive).Asc("sort_no").Desc("id"))
	return web.JsonData(roles)
}

func (c *RoleController) GetRole_menu_ids() *web.JsonResult {
	roleId, _ := params.GetInt64(c.Ctx, "roleId")
	menuIds := service.RoleMenuService.GetMenuIdsByRole(roleId)
	return web.JsonData(menuIds)
}

func (c *RoleController) PostSave_role_menus() *web.JsonResult {
	roleId, _ := params.GetInt64(c.Ctx, "roleId")
	menuIds := params.GetInt64Arr(c.Ctx, "menuIds")
	if err := service.RoleMenuService.SaveRoleMenus(roleId, menuIds); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *RoleController) PostUpdate_sort() *web.JsonResult {
	var ids []int64
	if err := c.Ctx.ReadJSON(&ids); err != nil {
		return web.JsonError(err)
	}
	if err := service.RoleService.UpdateSort(ids); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
