package admin

import (
	"bbs-go/internal/controller/response"
	"bbs-go/internal/model/constants"
	"strconv"
	"strings"

	"bbs-go/common/dates"
	"bbs-go/sqls"
	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type TagController struct {
	Ctx iris.Context
}

func (c *TagController) GetBy(id int64) *web.JsonResult {
	t := service.TagService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *TagController) GetList() *web.JsonResult {
	list, paging := service.TagService.FindPageByParams(params.NewQueryParams(c.Ctx).
		LikeByReq("id").
		LikeByReq("name").
		EqByReq("status").
		PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *TagController) PostCreate() *web.JsonResult {
	t := &model.Tag{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	if len(t.Name) == 0 {
		return web.JsonErrorMsg("name is required")
	}
	if service.TagService.GetByName(t.Name) != nil {
		return web.JsonErrorMsg("标签「" + t.Name + "」已存在")
	}

	t.Status = constants.StatusOK
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()

	err = service.TagService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *TagController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := service.TagService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	if len(t.Name) == 0 {
		return web.JsonErrorMsg("name is required")
	}
	if tmp := service.TagService.GetByName(t.Name); tmp != nil && tmp.Id != id {
		return web.JsonErrorMsg("标签「" + t.Name + "」已存在")
	}

	t.UpdateTime = dates.NowTimestamp()
	err = service.TagService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

// 自动完成
func (c *TagController) GetAutocomplete() *web.JsonResult {
	keyword := strings.TrimSpace(c.Ctx.URLParam("keyword"))
	var tags []model.Tag
	if len(keyword) > 0 {
		tags = service.TagService.Find(sqls.NewCnd().Starting("name", keyword).Desc("id"))
	} else {
		tags = service.TagService.Find(sqls.NewCnd().Desc("id").Limit(10))
	}
	return web.JsonData(response.BuildTags(tags))
}

// 根据标签编号批量获取
func (c *TagController) GetTags() *web.JsonResult {
	tagIds := params.FormValueInt64Array(c.Ctx, "tagIds")
	var tags *[]response.TagResponse
	if len(tagIds) > 0 {
		tagArr := service.TagService.Find(sqls.NewCnd().In("id", tagIds))
		if len(tagArr) > 0 {
			tags = response.BuildTags(tagArr)
		}
	}
	return web.JsonData(tags)
}
