package api

import (
	"bbs-go/common/strs"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
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
	profile.Settings = payload.UserSettings{
		LockedProfile: user.LockedProfile,
		ShowLocation:  user.ShowLocation,
		EmailNotify:   user.EmailNotify,
	}
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
	if err := service.UserService.Updates(user.Id, columns); err != nil {
		slog.Error("Failed to update user profile:", "error", err, "userId", user.Id, "columns", columns)
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

	err := service.FavoriteService.AddTopicFavorite(user.Id, topicId)
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

	err := service.FavoriteService.RemoveTopicFavorite(user.Id, topicId)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// // POST /api/me/edit/{userId}
// func (c *MeController) PostEditBy(userId int64) *web.JsonResult {
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	if user == nil {
// 		return web.JsonError(errs.NotLogin)
// 	}
// 	if user.Id != userId {
// 		return web.JsonErrorMsg("No permission")
// 	}
// 	var (
// 		nickname    = strings.TrimSpace(params.FormValue(c.Ctx, "nickname"))
// 		homePage    = params.FormValue(c.Ctx, "homePage")
// 		description = params.FormValue(c.Ctx, "description")
// 		gender      = strings.TrimSpace(params.FormValue(c.Ctx, "gender"))
// 		birthdayStr = strings.TrimSpace(params.FormValue(c.Ctx, "birthday"))
// 		birthday    *time.Time
// 		err         error
// 	)

// 	if len(nickname) == 0 {
// 		return web.JsonErrorMsg("Nickname cannot be empty")
// 	}

// 	if strs.IsNotBlank(gender) {
// 		if gender != string(constants.GenderMale) && gender != string(constants.GenderFemale) {
// 			return web.JsonErrorMsg("Gender data error")
// 		}
// 	}
// 	if strs.IsNotBlank(birthdayStr) {
// 		*birthday, err = dates.Parse(birthdayStr, dates.FmtDate)
// 		if err != nil {
// 			return web.JsonError(err)
// 		}
// 	}

// 	if len(homePage) > 0 && validate.IsURL(homePage) != nil {
// 		return web.JsonErrorMsg("Homepage address error")
// 	}

// 	columns := map[string]interface{}{
// 		"nickname":    nickname,
// 		"home_page":   homePage,
// 		"description": description,
// 		"gender":      gender,
// 	}
// 	if birthday != nil {
// 		columns["birthday"] = birthday
// 	}
// 	err = service.UserService.Updates(user.Id, columns)
// 	if err != nil {
// 		return web.JsonError(err)
// 	}
// 	return web.JsonSuccess()
// }

// // Update password
// func (c *MeController) PostUpdatePassword() *web.JsonResult {
// 	user := service.UserTokenService.GetCurrent(c.Ctx)
// 	if user == nil {
// 		return web.JsonError(errs.NotLogin)
// 	}
// 	var (
// 		oldPassword = params.FormValue(c.Ctx, "oldPassword")
// 		password    = params.FormValue(c.Ctx, "password")
// 		rePassword  = params.FormValue(c.Ctx, "rePassword")
// 	)
// 	if err := service.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword); err != nil {
// 		return web.JsonError(err)
// 	}
// 	return web.JsonSuccess()
// }

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
	if err := service.UserService.UpdateBackgroundImage(user.Id, backgroundImage); err != nil {
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
	topics, cursor, hasMore := service.TopicService.GetUserTopics(user.Id, cursor)
	return web.JsonCursorData(payload.BuildSimpleTopics(topics, user), cursor, hasMore)
}

func (c *MeController) GetFollowers() *web.JsonResult {
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	followerIds, cursor, hasMore := service.UserFollowService.GetFollowers(currentUser.Id, cursor, 10)
	followedSet := service.UserFollowService.GetMutualFollowers(currentUser.Id, followerIds...)

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

	followingIds, cursor, hasMore := service.UserFollowService.GetFollowing(currentUser.Id, cursor, 10)
	itemList := make([]payload.UserInfo, 0, len(followingIds))
	for _, id := range followingIds {
		item := payload.BuildUserInfoDefaultIfNull(id)
		itemList = append(itemList, *item)
	}
	return web.JsonCursorData(itemList, cursor, hasMore)
}

func (c *MeController) GetFollowingBy(userId int64) *web.JsonResult {
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	following := service.UserFollowService.IsFollowing(currentUser.Id, userId)
	return web.JsonData(map[string]bool{"following": following})
}
