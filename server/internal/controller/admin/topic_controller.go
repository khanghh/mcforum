package admin

import (
	"strconv"

	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/service"
)

type TopicController struct {
	Ctx iris.Context
}

func (c *TopicController) GetBy(id int64) *web.JsonResult {
	t := service.TopicService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *TopicController) GetList() *web.JsonResult {
	list, paging := service.TopicService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("id").EqByReq("user_id").EqByReq("status").EqByReq("recommend").LikeByReq("title").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, topic := range list {
		item := payload.BuildSimpleTopic(&topic)
		builder := web.NewRspBuilder(item)
		builder.Put("status", topic.Status)
		results = append(results, builder.Build())
	}

	return web.JsonData(&web.PageResult{Results: results, Page: paging})
}

// 推荐
func (c *TopicController) PostRecommend() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	err = service.TopicService.SetRecommended(id, true)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 取消推荐
func (c *TopicController) DeleteRecommend() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	err = service.TopicService.SetRecommended(id, false)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostDelete() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	err = service.TopicService.Delete(id, user.Id, c.Ctx.Request())
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostUndelete() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	err = service.TopicService.Undelete(id)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostAudit() *web.JsonResult {
	id := c.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return web.JsonErrorMsg("id is required")
	}
	err := service.TopicService.UpdateColumn(id, "status", constants.StatusOK)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
