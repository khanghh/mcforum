package api

import (
	"bbs-go/common/strs"
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

// Set background image
func (c *MeController) PostSetCoverPhoto() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	backgroundImage := params.FormValue(c.Ctx, "backgroundImage")
	if strs.IsBlank(backgroundImage) {
		return web.JsonErrorMsg("Please upload image")
	}
	if err := service.UserService.UpdateBackgroundImage(user.ID, backgroundImage); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *MeController) GetTopics() *web.JsonResult {
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	topics, cursor, hasMore := service.TopicService.GetUserTopics(user.ID, cursor)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), cursor, hasMore)
}

func (c *MeController) GetFollowers() *web.JsonResult {
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

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
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

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
	var count int64 = 0
	var messages []model.Message
	if user != nil {
		count = service.MessageService.GetUnReadCount(user.ID)
		messages = service.MessageService.Find(sqls.NewCnd().Eq("user_id", user.ID).
			Eq("status", msg.StatusUnread).Limit(3).Desc("id"))
	}
	return web.NewEmptyRspBuilder().Put("count", count).Put("messages", payload.BuildMessages(messages)).JsonResult()
}

// User messages
func (c *MeController) GetMessages() *web.JsonResult {
	user, err := service.UserTokenService.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(errs.NotLogin)
	}
	var (
		limit     = 20
		cursor, _ = params.GetInt64(c.Ctx, "cursor")
	)

	cnd := sqls.NewCnd().Eq("user_id", user.ID).Limit(limit).Desc("id")
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	list := service.MessageService.Find(cnd)

	var (
		nextCursor = cursor
		hasMore    = false
	)
	if len(list) > 0 {
		nextCursor = list[len(list)-1].ID
		hasMore = len(list) == limit
	}

	// Mark all as read
	service.MessageService.MarkRead(user.ID)

	return web.JsonCursorData(payload.BuildMessages(list), nextCursor, hasMore)
}

func (c *MeController) PutStatus() *web.JsonResult {
	var msgBody struct {
		Message string `json:"message"`
	}
	if err := c.Ctx.ReadJSON(&msgBody); err != nil {
		return web.JsonError(errs.ErrBadRequest)
	}
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	if len(msgBody.Message) > constants.StatusMessageMaxLength {
		return web.JsonError(errs.NewBadRequestError(locale.T("user.status_max_length_exceeded")))
	}
	if err := service.UserService.UpdateStatusMessage(user.ID, msgBody.Message); err != nil {
		return web.JsonError(errs.ErrInternalServer)
	}
	return web.JsonSuccess()
}
