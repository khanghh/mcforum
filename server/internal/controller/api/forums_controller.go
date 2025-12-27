package api

import (
	"bbs-go/common/arrays"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type ForumsController struct {
	Ctx iris.Context
}

// func (c *ForumController) GetMenu() *web.JsonResult {
// 	sidebarMenu := []payload.ForumResponse{whatsNewForum, recommendedForum, followedForum}
// 	realNodes := payload.BuildForumList(service.ForumService.GetAll())
// 	sidebarMenu = append(sidebarMenu, realNodes...)
// 	return web.JsonData(sidebarMenu)
// }

// Get list of forums
func (c *ForumsController) Get() *web.JsonResult {
	nodes := payload.BuildForumList(service.ForumService.GetAll())
	return web.JsonData(nodes)
}

// Get topics by forum slug
func (c *ForumsController) GetBy(slug string) *web.JsonResult {
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
		pinnedTopics := service.TopicService.GetPinnedTopics(forum.ID, 3)
		temp = append(temp, pinnedTopics...)
	}
	topics, cursor, hasMore := service.TopicService.GetForumTopics(forum.ID, cursor)
	for _, topic := range topics {
		if !topic.Pinned {
			temp = append(temp, topic)
		}
	}
	list := arrays.Distinct(temp, func(t model.Topic) any { return t.ID })
	return web.JsonCursorData(payload.BuildSimpleTopics(list, user), cursor, hasMore)
}

func (c *ForumsController) GetByInfo(slug string) *web.JsonResult {
	forum := service.ForumService.FindOne(sqls.NewCnd().Eq("slug", slug))
	if forum == nil {
		return web.JsonError(errs.ErrForumNotFound)
	}
	return web.JsonData(payload.BuildForum(forum))
}
