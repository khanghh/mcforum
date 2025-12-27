package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/model/constants"

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

// User details
func (c *UsersController) GetBy(username string) *web.JsonResult {
	user := cache.UserCache.GetByUsername(username)
	if user != nil && user.IsActive {
		return web.JsonData(payload.BuildUserDetail(user))
	}
	return web.JsonError(errs.ErrUserNotFound)
}

// User favorites
func (c *UsersController) GetFavorites() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	// User must be logged in
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}

	// Query list
	limit := 20
	var favorites []model.Favorite
	if cursor > 0 {
		favorites = service.FavoriteService.Find(sqls.NewCnd().Where("user_id = ? and id < ?",
			user.ID, cursor).Desc("id").Limit(20))
	} else {
		favorites = service.FavoriteService.Find(sqls.NewCnd().Where("user_id = ?", user.ID).Desc("id").Limit(limit))
	}

	hasMore := false
	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].ID
		hasMore = len(favorites) >= limit
	}

	return web.JsonCursorData(payload.BuildFavorites(favorites), cursor, hasMore)
}

// User score logs
func (c *UsersController) GetScore_logs() *web.JsonResult {
	user, err := service.UserTokenService.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}
	var (
		limit     = 20
		cursor, _ = params.GetInt64(c.Ctx, "cursor")
	)
	cnd := sqls.NewCnd().Eq("user_id", user.ID).Limit(limit).Desc("id")
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	list := service.UserScoreLogService.Find(cnd)

	var (
		nextCursor = cursor
		hasMore    = false
	)
	if len(list) > 0 {
		nextCursor = list[len(list)-1].ID
		hasMore = len(list) == limit
	}

	return web.JsonCursorData(list, nextCursor, hasMore)
}

// Score ranking
func (c *UsersController) GetScoreRank() *web.JsonResult {
	users := cache.UserCache.GetScoreRank()
	var results []*payload.UserInfo
	for _, user := range users {
		results = append(results, payload.BuildUserInfo(&user))
	}
	return web.JsonData(results)
}

// Ban
func (c *UsersController) PostForbidden() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("No permission")
	}
	var (
		userId = params.FormValueInt64Default(c.Ctx, "userId", 0)
		days   = params.FormValueIntDefault(c.Ctx, "days", 0)
		reason = params.FormValue(c.Ctx, "reason")
	)
	if userId < 0 {
		return web.JsonErrorMsg("Please provide: userId")
	}
	if days == -1 && !user.HasRole(constants.RoleOwner) {
		return web.JsonErrorMsg("No permanent ban permission")
	}
	if days == 0 {
		service.UserService.RemoveForbidden(user.ID, userId, c.Ctx.Request())
	} else {
		if err := service.UserService.Forbidden(user.ID, userId, days, reason, c.Ctx.Request()); err != nil {
			return web.JsonError(err)
		}
	}
	return web.JsonSuccess()
}

// Request email verification mail
func (c *UsersController) PostSend_verify_email() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if err := service.UserService.SendEmailVerifyEmail(user.ID); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// Get email verification code
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
	topics, cursor, hasMore := service.TopicService.GetUserTopics(user.ID, cursor)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), cursor, hasMore)
}

func (c *UsersController) GetByFollowers(username string) *web.JsonResult {
	user := cache.UserCache.GetByUsername(username)
	if user == nil {
		return web.JsonError(errs.ErrUserNotFound)
	}

	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	followerIds, cursor, hasMore := service.UserFollowService.GetFollowers(user.ID, cursor, 10)

	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	var followedSet *hashset.Set
	if currentUser != nil {
		followedSet = service.UserFollowService.GetMutualFollowers(currentUser.ID, followerIds...)
	} else {
		followedSet = hashset.New()
	}

	var itemList []*payload.UserInfo
	for _, id := range followerIds {
		item := payload.BuildUserInfoDefaultIfNull(id)
		item.IsFollowing = followedSet.Contains(id)
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
	userIds, cursor, hasMore := service.UserFollowService.GetFollowing(user.ID, cursor, 10)

	current := service.UserTokenService.GetCurrent(c.Ctx)
	var followedSet *hashset.Set
	if current != nil {
		if current.ID == user.ID {
			followedSet = hashset.New()
			for _, id := range userIds {
				followedSet.Add(id)
			}
		} else {
			followedSet = service.UserFollowService.GetMutualFollowers(current.ID, userIds...)
		}
	}

	var itemList []*payload.UserInfo
	for _, id := range userIds {
		item := payload.BuildUserInfoDefaultIfNull(id)
		item.IsFollowing = followedSet.Contains(id)
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

	err := service.UserFollowService.Follow(currentUser.ID, user.ID)
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

	err := service.UserFollowService.UnFollow(currentUser.ID, user.ID)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
