package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/search"
	"bbs-go/internal/service"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"
)

type SearchController struct {
	Ctx iris.Context
}

func (c *SearchController) AnyReindex() *web.JsonResult {
	go service.TopicService.ScanDesc(func(topics []model.Topic) {
		for _, topic := range topics {
			if topic.Status != constants.StatusDeleted {
				search.UpdateTopicIndex(&topic)
			}
		}
	})
	return web.JsonSuccess()
}

func (c *SearchController) GetTopic() *web.JsonResult {
	var (
		cursor    = params.FormValueInt64Default(c.Ctx, "cursor", 1)
		keyword   = params.FormValue(c.Ctx, "keyword")
		nodeId    = params.FormValueInt64Default(c.Ctx, "nodeId", 0)
		timeRange = params.FormValueIntDefault(c.Ctx, "timeRange", 0)
		limit     = 20
	)
	list, _, err := search.SearchTopic(keyword, nodeId, timeRange, int(cursor), limit)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonCursorData(payload.BuildSearchTopics(list), cursor+1, len(list) >= limit)
}
