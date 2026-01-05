package api

import (
	"bbs-go/common/strs"
	"bbs-go/internal/cache"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/locale"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/pkg/msg"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"
	"io"
	"log/slog"

	"github.com/kataras/iris/v12"
)

type MeController struct {
	Ctx iris.Context
}

func (c *MeController) Get() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	profile := payload.BuildUserProfile(user)
	return web.JsonData(profile)
}

func (c *MeController) Patch() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	form := payload.GetUpdateProfileForm(c.Ctx)
	if (len(form.Bio)) > constants.BioMaxLength {

		return web.JsonError(errs.NewBadRequestError(locale.T("user.bio_max_length_exceeded")))
	}

	if strs.IsBlank(form.Nickname) {
		return web.JsonError(errs.NewBadRequestError(locale.T("user.nickname_required")))
	}
	if len(form.Nickname) > constants.NicknameMaxLength {
		return web.JsonError(errs.NewBadRequestError(locale.T("user.nickname_max_length_exceeded")))
	}

	if len(form.Location) > constants.LocationMaxLength {
		return web.JsonError(errs.NewBadRequestError(locale.T("user.location_max_length_exceeded")))
	}

	columns := map[string]interface{}{
		"avatar":         form.Avatar,
		"nickname":       form.Nickname,
		"bio":            form.Bio,
		"location":       form.Location,
		"locked_profile": form.LockedProfile,
		"show_location":  form.ShowLocation,
		"email_notify":   form.EmailNotify,
	}
	if err := service.UserService.Updates(user.ID, columns); err != nil {
		slog.Error("Failed to update user profile:", "error", err, "userId", user.ID, "columns", columns)
		return web.JsonError(errs.ErrInternalServer)
	}

	cache.UserCache.Invalidate(user.ID)
	return web.JsonSuccess()
}

// PUT /api/users/me/favorites
func (c *MeController) PutFavoritesBy(topicId int64) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}

	err := service.FavoriteService.AddTopicFavorite(user.ID, topicId)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// DELETE /api/users/me/favorites/{topicId}
func (c *MeController) DeleteFavoritesBy(topicId int64) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}

	err := service.FavoriteService.RemoveTopicFavorite(user.ID, topicId)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// User favorites
func (c *MeController) GetFavorites() *web.JsonResult {
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

// Set avatar image
func (c *MeController) PutAvatar() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	applyNow := c.Ctx.URLParamBoolDefault("apply", false)
	file, _, err := c.Ctx.FormFile("file")
	if err != nil {
		return web.JsonError(errs.ErrBadRequest)
	}
	defer file.Close()

	lr := io.LimitReader(file, constants.MaxAvatarSizeBytes+1)
	data, err := io.ReadAll(lr)
	if err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}
	if len(data) > constants.MaxAvatarSizeBytes {
		return web.JsonError(errs.ErrAvatarTooLarge)
	}

	ret, err := service.UploadService.UploadAvatar(user, data)
	if err != nil {
		slog.Error("Upload avatar failed", "error", err)
		return web.JsonError(errs.ErrInternalServer)
	}

	if applyNow {
		if err := service.UserService.UpdateAvatar(user.ID, ret.URL); err != nil {
			slog.Error("Change user avatar failed", "error", err)
			return web.JsonError(errs.ErrInternalServer)
		}
	}

	cache.UserCache.Invalidate(user.ID)
	return web.JsonData(map[string]string{
		"avatar": ret.URL,
	})
}

func (c *MeController) DeleteAvatar() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	if err := service.UserService.UpdateAvatar(user.ID, ""); err != nil {
		slog.Error("Delete user avatar failed", "error", err)
		return web.JsonError(errs.ErrInternalServer)
	}
	cache.UserCache.Invalidate(user.ID)
	return web.JsonSuccess()
}

func (c *MeController) PutCover() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	file, header, err := c.Ctx.FormFile("file")
	if err != nil {
		return web.JsonError(errs.ErrBadRequest)
	}

	ret, err := service.UploadService.UploadStream(user, file, header.Filename)
	if err != nil {
		slog.Error("Upload cover photo failed", "error", err)
		return web.JsonError(errs.ErrInternalServer)
	}

	if err := service.UserService.UpdateCoverPhoto(user.ID, ret.URL); err != nil {
		slog.Error("Change user cover photo failed", "error", err)
		return web.JsonError(errs.ErrInternalServer)
	}

	cache.UserCache.Invalidate(user.ID)
	return web.JsonData(map[string]string{
		"coverImage": ret.URL,
	})
}

func (c *MeController) GetTopics() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	topics, cursor, hasMore := service.TopicService.GetUserTopics(user.ID, cursor)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), cursor, hasMore)
}

func (c *MeController) GetFollowers() *web.JsonResult {
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	followerIds, cursor, hasMore := service.UserFollowService.GetFollowers(currentUser.ID, cursor, 10)
	followedSet := service.UserFollowService.GetMutualFollowers(currentUser.ID, followerIds...)

	itemList := make([]payload.UserInfo, 0, len(followerIds))
	for _, id := range followerIds {
		item := payload.BuildUserInfoDefaultIfNull(id)
		item.IsFollowing = followedSet.Contains(id)
		itemList = append(itemList, *item)
	}
	return web.JsonCursorData(itemList, cursor, hasMore)
}

func (c *MeController) GetFollowing() *web.JsonResult {
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	followingIds, cursor, hasMore := service.UserFollowService.GetFollowing(currentUser.ID, cursor, 10)
	itemList := make([]payload.UserInfo, 0, len(followingIds))
	for _, id := range followingIds {
		item := payload.BuildUserInfoDefaultIfNull(id)
		item.IsFollowing = true
		itemList = append(itemList, *item)
	}
	return web.JsonCursorData(itemList, cursor, hasMore)
}

func (c *MeController) GetFollowingBy(userId int64) *web.JsonResult {
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	following := service.UserFollowService.IsFollowing(currentUser.ID, userId)
	return web.JsonData(map[string]bool{"following": following})
}

// Get last 3 unread messages
func (c *MeController) GetMessages_Recent() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	count := service.MessageService.GetUnReadCount(user.ID)
	messages := service.MessageService.Find(sqls.NewCnd().Eq("user_id", user.ID).
		Eq("status", msg.StatusUnread).Limit(3).Desc("id"))
	return web.NewEmptyRspBuilder().
		Put("count", count).
		Put("messages", payload.BuildMessages(messages)).
		JsonResult()
}

// User messages
func (c *MeController) GetMessages() *web.JsonResult {
	user, err := service.UserTokenService.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	messages, nextCursor, hasMore := service.MessageService.GetUserMessages(user.ID, cursor)

	// Mark all as read
	service.MessageService.MarkRead(user.ID)
	return web.JsonCursorData(payload.BuildMessages(messages), nextCursor, hasMore)
}

func (c *MeController) PutStatus() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	var msgBody struct {
		Message string `json:"message"`
	}
	if err := c.Ctx.ReadJSON(&msgBody); err != nil {
		return web.JsonError(errs.ErrBadRequest)
	}
	if len(msgBody.Message) > constants.StatusMessageMaxLength {
		return web.JsonError(errs.NewBadRequestError(locale.T("user.status_max_length_exceeded")))
	}
	if err := service.UserService.UpdateStatusMessage(user.ID, msgBody.Message); err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}
	cache.UserCache.Invalidate(user.ID)
	return web.JsonSuccess()
}
