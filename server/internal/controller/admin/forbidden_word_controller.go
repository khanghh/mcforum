package admin

import (
	"bbs-go/internal/model"
	"bbs-go/internal/service"
	"strconv"

	"bbs-go/common/dates"
	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"
)

type ForbiddenWordController struct {
	Ctx iris.Context
}

func (c *ForbiddenWordController) GetBy(id int64) *web.JsonResult {
	t := service.ForbiddenWordService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *ForbiddenWordController) GetList() *web.JsonResult {
	list, paging := service.ForbiddenWordService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("type").LikeByReq("word").EqByReq("status").PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c ForbiddenWordController) PostDelete() *web.JsonResult {
	id, _ := params.FormValueInt64(c.Ctx, "id")
	service.ForbiddenWordService.Delete(id)
	return web.JsonSuccess()
}

func (c *ForbiddenWordController) PostCreate() *web.JsonResult {
	t := &model.ForbiddenWord{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t.CreateTime = dates.NowTimestamp()
	err = service.ForbiddenWordService.Create(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *ForbiddenWordController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := service.ForbiddenWordService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = service.ForbiddenWordService.Update(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}
