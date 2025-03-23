package admin

import (
	"bbs-go/internal/service"
	"strconv"

	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"
)

type OperateLogController struct {
	Ctx iris.Context
}

func (c *OperateLogController) GetBy(id int64) *web.JsonResult {
	t := service.OperateLogService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *OperateLogController) AnyList() *web.JsonResult {
	list, paging := service.OperateLogService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("user_id").EqByReq("op_type").EqByReq("data_type").EqByReq("data_id").
		PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}
