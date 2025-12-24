package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/model/constants"
	"bbs-go/pkg/msg"

	"bbs-go/common/strs"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/iris/v12"

	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

type UsersController struct {
	Ctx iris.Context
}

// 用户详情
func (c *UsersController) GetBy(username string) *web.JsonResult {
	user := cache.UserCache.GetByUsername(username)
	if user != nil && user.IsActive {
		return web.JsonData(payload.BuildUserDetail(user))
	}
	return web.JsonError(errs.ErrUserNotFound)
}

// 用户收藏
func (c *UsersController) GetFavorites() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	// 用户必须登录
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}

	// 查询列表
	limit := 20
	var favorites []model.Favorite
	if cursor > 0 {
		favorites = service.FavoriteService.Find(sqls.NewCnd().Where("user_id = ? and id < ?",
			user.Id, cursor).Desc("id").Limit(20))
	} else {
		favorites = service.FavoriteService.Find(sqls.NewCnd().Where("user_id = ?", user.Id).Desc("id").Limit(limit))
	}

	hasMore := false
	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].Id
		hasMore = len(favorites) >= limit
	}

	return web.JsonCursorData(payload.BuildFavorites(favorites), cursor, hasMore)
}

// 获取最近3条未读消息
func (c *UsersController) GetMsgrecent() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	var count int64 = 0
	var messages []model.Message
	if user != nil {
		count = service.MessageService.GetUnReadCount(user.Id)
		messages = service.MessageService.Find(sqls.NewCnd().Eq("user_id", user.Id).
			Eq("status", msg.StatusUnread).Limit(3).Desc("id"))
	}
	return web.NewEmptyRspBuilder().Put("count", count).Put("messages", payload.BuildMessages(messages)).JsonResult()
}

// 用户消息
func (c *UsersController) GetMessages() *web.JsonResult {
	user, err := service.UserTokenService.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(errs.NotLogin)
	}
	var (
		limit     = 20
		cursor, _ = params.GetInt64(c.Ctx, "cursor")
	)

	cnd := sqls.NewCnd().Eq("user_id", user.Id).Limit(limit).Desc("id")
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	list := service.MessageService.Find(cnd)

	var (
		nextCursor = cursor
		hasMore    = false
	)
	if len(list) > 0 {
		nextCursor = list[len(list)-1].Id
		hasMore = len(list) == limit
	}

	// 全部标记为已读
	service.MessageService.MarkRead(user.Id)

	return web.JsonCursorData(payload.BuildMessages(list), nextCursor, hasMore)
}

// 用户积分记录
func (c *UsersController) GetScore_logs() *web.JsonResult {
	user, err := service.UserTokenService.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}
	var (
		limit     = 20
		cursor, _ = params.GetInt64(c.Ctx, "cursor")
	)
	cnd := sqls.NewCnd().Eq("user_id", user.Id).Limit(limit).Desc("id")
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	list := service.UserScoreLogService.Find(cnd)

	var (
		nextCursor = cursor
		hasMore    = false
	)
	if len(list) > 0 {
		nextCursor = list[len(list)-1].Id
		hasMore = len(list) == limit
	}

	return web.JsonCursorData(list, nextCursor, hasMore)
}

// 积分排行
func (c *UsersController) GetScoreRank() *web.JsonResult {
	users := cache.UserCache.GetScoreRank()
	var results []*payload.UserInfo
	for _, user := range users {
		results = append(results, payload.BuildUserInfo(&user))
	}
	return web.JsonData(results)
}

// 禁言
func (c *UsersController) PostForbidden() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("无权限")
	}
	var (
		userId = params.FormValueInt64Default(c.Ctx, "userId", 0)
		days   = params.FormValueIntDefault(c.Ctx, "days", 0)
		reason = params.FormValue(c.Ctx, "reason")
	)
	if userId < 0 {
		return web.JsonErrorMsg("请传入：userId")
	}
	if days == -1 && !user.HasRole(constants.RoleOwner) {
		return web.JsonErrorMsg("无永久禁言权限")
	}
	if days == 0 {
		service.UserService.RemoveForbidden(user.Id, userId, c.Ctx.Request())
	} else {
		if err := service.UserService.Forbidden(user.Id, userId, days, reason, c.Ctx.Request()); err != nil {
			return web.JsonError(err)
		}
	}
	return web.JsonSuccess()
}

// PostEmailVerify 请求邮箱验证邮件
func (c *UsersController) PostSend_verify_email() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if err := service.UserService.SendEmailVerifyEmail(user.Id); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// PostVerify_email 获取邮箱验证码
func (c *UsersController) PostVerify_email() *web.JsonResult {
	token := params.FormValue(c.Ctx, "token")
	if strs.IsBlank(token) {
		return web.JsonErrorMsg("Illegal request")
	}
	var (
		email string
		err   error
	)
	if email, err = service.UserService.VerifyEmail(token); err != nil {
		return web.JsonError(err)
	}
	return web.NewEmptyRspBuilder().Put("email", email).JsonResult()
}

func (c *UsersController) GetByTopics(username string) *web.JsonResult {
	user := cache.UserCache.GetByUsername(username)
	if user == nil {
		return web.JsonError(errs.ErrUserNotFound)
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	topics, cursor, hasMore := service.TopicService.GetUserTopics(user.Id, cursor)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), cursor, hasMore)
}

func (c *UsersController) GetByFollowers(username string) *web.JsonResult {
	user := cache.UserCache.GetByUsername(username)
	if user == nil {
		return web.JsonError(errs.ErrUserNotFound)
	}

	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	followerIds, cursor, hasMore := service.UserFollowService.GetFollowers(user.Id, cursor, 10)

	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	var followedSet *hashset.Set
	if currentUser != nil {
		followedSet = service.UserFollowService.GetMutualFollowers(currentUser.Id, followerIds...)
	}

	var itemList []*payload.UserInfo
	for _, id := range followerIds {
		item := payload.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return web.JsonCursorData(itemList, cursor, hasMore)
}

func (c *UsersController) GetByFollowing(username string) *web.JsonResult {
	user := cache.UserCache.GetByUsername(username)
	if user == nil {
		return web.JsonError(errs.ErrUserNotFound)
	}

	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	userIds, cursor, hasMore := service.UserFollowService.GetFollowing(user.Id, cursor, 10)

	current := service.UserTokenService.GetCurrent(c.Ctx)
	var followedSet *hashset.Set
	if current != nil {
		if current.Id == user.Id {
			followedSet = hashset.New()
			for _, id := range userIds {
				followedSet.Add(id)
			}
		} else {
			followedSet = service.UserFollowService.GetMutualFollowers(current.Id, userIds...)
		}
	}

	var itemList []*payload.UserInfo
	for _, id := range userIds {
		item := payload.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return web.JsonCursorData(itemList, cursor, hasMore)
}

func (c *UsersController) PostByFollow(username string) *web.JsonResult {
	user := cache.UserCache.GetByUsername(username)
	if user == nil {
		return web.JsonError(errs.ErrUserNotFound)
	}

	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	err := service.UserFollowService.Follow(currentUser.Id, user.Id)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *UsersController) DeleteByFollow(username string) *web.JsonResult {
	user := cache.UserCache.GetByUsername(username)
	if user == nil {
		return web.JsonError(errs.ErrUserNotFound)
	}

	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	err := service.UserFollowService.UnFollow(currentUser.Id, user.Id)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
