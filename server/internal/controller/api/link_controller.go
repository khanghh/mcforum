package api

import (
	"bbs-go/internal/model/constants"

	"bbs-go/sqls"
	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type LinkController struct {
	Ctx iris.Context
}

func (c *LinkController) GetBy(id int64) *web.JsonResult {
	link := service.LinkService.Get(id)
	if link == nil || link.Status == constants.StatusDeleted {
		return web.JsonErrorMsg("数据不存在")
	}
	return web.JsonData(c.buildLink(*link))
}

// 列表
func (c *LinkController) GetList() *web.JsonResult {
	links := service.LinkService.Find(sqls.NewCnd().
		Eq("status", constants.StatusOK).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return web.JsonData(itemList)
}

// 列表
func (c *LinkController) GetLinks() *web.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)

	links, paging := service.LinkService.FindPageByCnd(sqls.NewCnd().
		Eq("status", constants.StatusOK).Page(page, 20).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return web.JsonPageData(itemList, paging)
}

// 前10个链接
func (c *LinkController) GetToplinks() *web.JsonResult {
	links := service.LinkService.Find(sqls.NewCnd().
		Eq("status", constants.StatusOK).Limit(10).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return web.JsonData(itemList)
}

func (c *LinkController) buildLink(link model.Link) map[string]interface{} {
	return map[string]interface{}{
		"id":         link.Id,
		"linkId":     link.Id,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"logo":       link.Logo,
		"createTime": link.CreateTime,
	}
}
