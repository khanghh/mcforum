package api

import (
	"bbs-go/internal/model/constants"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

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
		return web.JsonErrorMsg("Not found")
	}
	return web.JsonData(c.buildLink(*link))
}

// List
func (c *LinkController) GetList() *web.JsonResult {
	links := service.LinkService.Find(sqls.NewCnd().
		Eq("status", constants.StatusActive).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return web.JsonData(itemList)
}

// List
func (c *LinkController) GetLinks() *web.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)

	links, paging := service.LinkService.FindPageByCnd(sqls.NewCnd().
		Eq("status", constants.StatusActive).Page(page, 20).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return web.JsonPageData(itemList, paging)
}

// Top 10 links
func (c *LinkController) GetToplinks() *web.JsonResult {
	links := service.LinkService.Find(sqls.NewCnd().
		Eq("status", constants.StatusActive).Limit(10).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return web.JsonData(itemList)
}

func (c *LinkController) buildLink(link model.Link) map[string]interface{} {
	return map[string]interface{}{
		"id":         link.ID,
		"linkId":     link.ID,
		"url":        link.URL,
		"title":      link.Title,
		"summary":    link.Summary,
		"logo":       link.Logo,
		"createTime": link.CreateTime,
	}
}
