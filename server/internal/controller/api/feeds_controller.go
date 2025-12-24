package api

import (
	"bbs-go/common/arrays"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model"
	"bbs-go/internal/service"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"
)

type FeedsController struct {
	Ctx iris.Context
}

func (c *FeedsController) GetWhatsNew() *web.JsonResult {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = service.UserTokenService.GetCurrent(c.Ctx)
	)

	var temp []model.Topic
	if cursor <= 0 {
		pinnedTopics := service.TopicService.GetPinnedTopics(0, 3)
		temp = append(temp, pinnedTopics...)
	}
	topics, cursor, hasMore := service.TopicService.GetNewestTopics(cursor)
	for _, topic := range topics {
		topic.Pinned = false // 正常列表不要渲染置顶
		temp = append(temp, topic)
	}
	list := arrays.Distinct(temp, func(t model.Topic) any { return t.Id })
	return web.JsonCursorData(payload.BuildSimpleTopics(list, user), cursor, hasMore)
}

func (c *FeedsController) GetRecommended() *web.JsonResult {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = service.UserTokenService.GetCurrent(c.Ctx)
	)
	var temp []model.Topic
	if cursor <= 0 {
		pinnedTopics := service.TopicService.GetPinnedTopics(0, 3)
		for _, topic := range pinnedTopics {
			if topic.Recommended {
				temp = append(temp, pinnedTopics...)
			}
		}
	}
	topics, cursor, hasMore := service.TopicService.GetRecommendedTopics(cursor)
	for _, topic := range topics {
		if !topic.Pinned {
			temp = append(temp, topic)
		}
	}
	list := arrays.Distinct(temp, func(t model.Topic) any { return t.Id })
	return web.JsonCursorData(payload.BuildSimpleTopics(list, user), cursor, hasMore).
		SetProperty("forum", "recommended")
}

func (c *FeedsController) GetFollowed() *web.JsonResult {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = service.UserTokenService.GetCurrent(c.Ctx)
	)
	topics, cursor, hasMore := service.TopicService.GetFollowedTopics(user.Id, cursor)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), cursor, hasMore)
}

func NewFeedController() *FeedsController {
	return &FeedsController{}
}
