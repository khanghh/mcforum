package api

import (
	"bbs-go/internal/controllers/response"
	"bbs-go/internal/locale"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"strconv"

	"bbs-go/sqls"
	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type ForumController struct {
	Ctx iris.Context
}

func (c *ForumController) GetMenu() *web.JsonResult {
	forumList := []response.ForumResponse{
		{
			Name: locale.T("nav.news_feed"),
			Slug: "whats-new",
		},
		{
			Name: locale.T("nav.recommended"),
			Slug: "recommended",
		},
		{
			Name: locale.T("nav.feed"),
			Slug: "feed",
		},
	}
	realNodes := response.BuildForumList(services.ForumService.GetAll())
	forumList = append(forumList, realNodes...)
	return web.JsonData(forumList)
}

// 节点
func (c *ForumController) GetList() *web.JsonResult {
	nodes := response.BuildForumList(services.ForumService.GetAll())
	return web.JsonData(nodes)
}

// 节点信息
func (c *ForumController) GetNode() *web.JsonResult {
	nodeId := params.FormValueInt64Default(c.Ctx, "nodeId", 0)
	node := services.ForumService.Get(nodeId)
	return web.JsonData(response.BuildForum(node))
}

func (c *ForumController) GetWhatsNew() (*web.JsonResult, int) {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = services.UserTokenService.GetCurrent(c.Ctx)
	)

	var temp []models.Topic
	if cursor <= 0 {
		stickyTopics := services.TopicService.GetStickyTopics(0, 3)
		temp = append(temp, stickyTopics...)
	}
	topics, cursor, hasMore := services.TopicService.GetTopics(user, 1, cursor)
	for _, topic := range topics {
		topic.Sticky = false // 正常列表不要渲染置顶
		temp = append(temp, topic)
	}
	list := common.Distinct(temp, func(t models.Topic) any { return t.Id })
	return web.JsonCursorData(response.BuildSimpleTopics(list, user), strconv.FormatInt(cursor, 10), hasMore), iris.StatusOK
}

func (c *ForumController) GetRecommended() *web.JsonResult {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = services.UserTokenService.GetCurrent(c.Ctx)
	)
	var temp []models.Topic
	if cursor <= 0 {
		stickyTopics := services.TopicService.GetStickyTopics(0, 3)
		temp = append(temp, stickyTopics...)
	}
	topics, cursor, hasMore := services.TopicService.GetRecommendedTopics(cursor)
	for _, topic := range topics {
		if !topic.Sticky {
			temp = append(temp, topic)
		}
	}
	list := common.Distinct(temp, func(t models.Topic) any { return t.Id })
	return web.JsonCursorData(response.BuildSimpleTopics(list, user), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *ForumController) GetFollowed() (*web.JsonResult, int) {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = services.UserTokenService.GetCurrent(c.Ctx)
	)
	var temp []models.Topic
	if cursor <= 0 {
		pinnedTopics := services.TopicService.GetStickyTopics(0, 3)
		for _, topic := range pinnedTopics {
			if topic.Recommend {
				temp = append(temp, topic)
			}
		}
	}
	topics, cursor, hasMore := services.TopicService.GetFollowedAuthorsTopics(user.Id, cursor)
	for _, topic := range topics {
		if !topic.Sticky {
			temp = append(temp, topic)
		}
	}
	return web.JsonCursorData(response.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore), iris.StatusOK
}

// // 帖子列表
func (c *ForumController) GetBy(slug string) (*web.JsonResult, int) {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = services.UserTokenService.GetCurrent(c.Ctx)
	)

	forum := services.ForumService.FindOne(sqls.NewCnd().Eq("slug", slug))
	if forum == nil {
		return web.JsonError(web.NewError(iris.StatusNotFound, "not found")), iris.StatusNotFound
	}

	var temp []models.Topic
	if cursor <= 0 {
		stickyTopics := services.TopicService.GetStickyTopics(forum.Id, 3)
		temp = append(temp, stickyTopics...)
	}
	topics, cursor, hasMore := services.TopicService.GetTopics(user, forum.Id, cursor)
	for _, topic := range topics {
		if !topic.Sticky {
			temp = append(temp, topic)
		}
	}
	list := common.Distinct(temp, func(t models.Topic) any { return t.Id })
	return web.JsonCursorData(response.BuildSimpleTopics(list, user), strconv.FormatInt(cursor, 10), hasMore), iris.StatusOK
}

// 标签帖子列表
func (c *ForumController) GetTagTopics() *web.JsonResult {
	var (
		cursor     = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		tagId, err = params.FormValueInt64(c.Ctx, "tagId")
		user       = services.UserTokenService.GetCurrent(c.Ctx)
	)
	if err != nil {
		return web.JsonError(err)
	}
	topics, cursor, hasMore := services.TopicService.GetTagTopics(tagId, cursor)
	return web.JsonCursorData(response.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}

// 最新话题
func (c *ForumController) GetNewest() *web.JsonResult {
	topics := services.TopicService.Find(sqls.NewCnd().Eq("status", constants.StatusOK).Desc("id").Limit(6))
	return web.JsonData(response.BuildSimpleTopics(topics, nil))
}
