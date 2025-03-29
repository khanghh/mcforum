package api

import (
	"bbs-go/common/arrays"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/locale"
	"strconv"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type ForumController struct {
	Ctx iris.Context
}

var (
	whatsNewForum = payload.ForumResponse{
		Name: locale.T("nav.whats-new"),
		Slug: "whats-new",
		Logo: "/images/new.png",
	}
	recommendedForum = payload.ForumResponse{
		Name: locale.T("nav.recommended"),
		Slug: "recommended",
		Logo: "/images/recommend.png",
	}
	followedForum = payload.ForumResponse{
		Name: locale.T("nav.followed"),
		Slug: "followed",
		Logo: "/images/feed.png",
	}
)

func (c *ForumController) GetMenu() *web.JsonResult {
	sidebarMenu := []payload.ForumResponse{whatsNewForum, recommendedForum, followedForum}
	realNodes := payload.BuildForumList(service.ForumService.GetAll())
	sidebarMenu = append(sidebarMenu, realNodes...)
	return web.JsonData(sidebarMenu)
}

// 节点
func (c *ForumController) GetList() *web.JsonResult {
	nodes := payload.BuildForumList(service.ForumService.GetAll())
	return web.JsonData(nodes)
}

func (c *ForumController) getWhatsNew() *web.JsonResult {
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
	return web.JsonCursorData(payload.BuildSimpleTopics(list, user), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *ForumController) getRecommended() *web.JsonResult {
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
	return web.JsonCursorData(payload.BuildSimpleTopics(list, user), strconv.FormatInt(cursor, 10), hasMore).
		SetProperty("forum", recommendedForum)
}

func (c *ForumController) getFollowed() *web.JsonResult {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = service.UserTokenService.GetCurrent(c.Ctx)
	)
	topics, cursor, hasMore := service.TopicService.GetFollowedTopics(user.Id, cursor)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}

// // 帖子列表
func (c *ForumController) GetBy(slug string) *web.JsonResult {
	if slug == whatsNewForum.Slug {
		return c.getWhatsNew()
	}
	if slug == recommendedForum.Slug {
		return c.getRecommended()
	}
	if slug == followedForum.Slug {
		return c.getFollowed()
	}

	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		user   = service.UserTokenService.GetCurrent(c.Ctx)
	)

	forum := service.ForumService.FindOne(sqls.NewCnd().Eq("slug", slug))
	if forum == nil {
		return web.JsonError(errs.ErrForumNotFound)
	}

	var temp []model.Topic
	if cursor <= 0 {
		pinnedTopics := service.TopicService.GetPinnedTopics(forum.Id, 3)
		temp = append(temp, pinnedTopics...)
	}
	topics, cursor, hasMore := service.TopicService.GetForumTopics(forum.Id, cursor)
	for _, topic := range topics {
		if !topic.Pinned {
			temp = append(temp, topic)
		}
	}
	list := arrays.Distinct(temp, func(t model.Topic) any { return t.Id })
	return web.JsonCursorData(payload.BuildSimpleTopics(list, user), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *ForumController) GetByInfo(slug string) *web.JsonResult {
	if slug == whatsNewForum.Slug {
		return web.JsonData(whatsNewForum)
	}
	if slug == recommendedForum.Slug {
		return web.JsonData(recommendedForum)
	}
	if slug == followedForum.Slug {
		return web.JsonData(followedForum)
	}

	forum := service.ForumService.FindOne(sqls.NewCnd().Eq("slug", slug))
	if forum == nil {
		return web.JsonError(errs.ErrForumNotFound)
	}
	return web.JsonData(payload.BuildForum(forum))
}
