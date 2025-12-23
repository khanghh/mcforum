package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/model/constants"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/cache"
	"bbs-go/internal/service"
)

type TagsController struct {
	Ctx iris.Context
}

// 标签详情
func (c *TagsController) GetBy(tagName string) *web.JsonResult {
	tag := cache.TagCache.Get(tagName)
	if tag == nil {
		return web.JsonError(errs.ErrBadRequest)
	}
	return web.JsonData(payload.BuildTag(tag))
}

// 标签列表
func (c *TagsController) GetTags() *web.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)
	tags, paging := service.TagService.FindPageByCnd(sqls.NewCnd().
		Eq("status", constants.StatusActive).
		Page(page, 200).Desc("id"))

	return web.JsonPageData(payload.BuildTags(tags), paging)
}

// 标签自动完成
func (c *TagsController) PostAutocomplete() *web.JsonResult {
	input := c.Ctx.FormValue("input")
	tags := service.TagService.Autocomplete(input)
	return web.JsonData(tags)
}
