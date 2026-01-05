package api

import (
	"bbs-go/common/arrays"
	"bbs-go/internal/cache"
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
	currentRank := 0
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user != nil {
		currentRank = user.Role.Rank
	}
	forums := payload.BuildForumList(service.ForumService.GetAll(), currentRank)
	return web.JsonData(forums)
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
	return web.JsonData(payload.BuildForum(forum, 0))
}

func (c *ForumsController) GetStats() *web.JsonResult {
	stats, err := service.StatsService.GetForumStats()
	if err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}
	resp := payload.ForumStatsResponse{
		TotalTopics:  stats.TotalTopics,
		TotalPosts:   stats.TotalComments,
		TotalMembers: stats.TotalMembers,
		TotalVisits:  stats.TotalVisits,
		NewestMember: stats.NewestMember,
	}
	return web.JsonData(resp)
}

func (c *ForumsController) GetTopContributors() *web.JsonResult {
	users := cache.UserCache.GetScoreRank()
	respData := make([]payload.UserInfo, 0, len(users))
	respData = respData[:min(len(respData), 3)]
	for _, user := range users {
		item := payload.BuildUserInfo(&user)
		respData = append(respData, *item)
	}
	return web.JsonData(respData)
}
