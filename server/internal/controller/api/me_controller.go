package api

import (
	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/internal/validate"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/iris/v12"
)

type MeController struct {
	Ctx iris.Context
}

func (c *MeController) Get() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user != nil {
		return web.JsonData(payload.BuildUserProfile(user))
	}
	return web.JsonSuccess()
}

// PUT /api/me/favorites
func (c *MeController) PutFavorites() (*web.JsonResult, error) {
	topicId, err := params.FormValueInt64(c.Ctx, "topicId")
	if err != nil {
		return web.JsonError(errs.ErrBadRequest), nil
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized), nil
	}

	err = service.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// DELETE /api/me/favorites/{topicId}
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

// POST /api/me/edit/{userId}
func (c *MeController) PostEditBy(userId int64) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if user.Id != userId {
		return web.JsonErrorMsg("无权限")
	}
	var (
		nickname    = strings.TrimSpace(params.FormValue(c.Ctx, "nickname"))
		homePage    = params.FormValue(c.Ctx, "homePage")
		description = params.FormValue(c.Ctx, "description")
		gender      = strings.TrimSpace(params.FormValue(c.Ctx, "gender"))
		birthdayStr = strings.TrimSpace(params.FormValue(c.Ctx, "birthday"))
		birthday    *time.Time
		err         error
	)

	if len(nickname) == 0 {
		return web.JsonErrorMsg("昵称不能为空")
	}

	if strs.IsNotBlank(gender) {
		if gender != string(constants.GenderMale) && gender != string(constants.GenderFemale) {
			return web.JsonErrorMsg("性别数据错误")
		}
	}
	if strs.IsNotBlank(birthdayStr) {
		*birthday, err = dates.Parse(birthdayStr, dates.FmtDate)
		if err != nil {
			return web.JsonError(err)
		}
	}

	if len(homePage) > 0 && validate.IsURL(homePage) != nil {
		return web.JsonErrorMsg("个人主页地址错误")
	}

	columns := map[string]interface{}{
		"nickname":    nickname,
		"home_page":   homePage,
		"description": description,
		"gender":      gender,
	}
	if birthday != nil {
		columns["birthday"] = birthday
	}
	err = service.UserService.Updates(user.Id, columns)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 修改头像
func (c *MeController) PostUpdateAvatar() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	avatar := strings.TrimSpace(params.FormValue(c.Ctx, "avatar"))
	if len(avatar) == 0 {
		return web.JsonErrorMsg("头像不能为空")
	}
	err := service.UserService.UpdateAvatar(user.Id, avatar)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *MeController) PostUpdateNickname() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	nickname := strings.TrimSpace(params.FormValue(c.Ctx, "nickname"))
	if len(nickname) == 0 {
		return web.JsonErrorMsg("Nickname cannot be empty")
	}
	err := service.UserService.UpdateNickname(user.Id, nickname)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *MeController) PostUpdateBio() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	description := strings.TrimSpace(params.FormValue(c.Ctx, "description"))
	err := service.UserService.UpdateDescription(user.Id, description)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *MeController) PostUpdateGender() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	gender := strings.TrimSpace(params.FormValue(c.Ctx, "gender"))
	err := service.UserService.UpdateGender(user.Id, gender)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *MeController) PostUpdateBirthday() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	birthday := strings.TrimSpace(params.FormValue(c.Ctx, "birthday"))
	err := service.UserService.UpdateBirthday(user.Id, birthday)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

// // 修改密码
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

// 设置背景图
func (c *MeController) PostSet_background_image() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	backgroundImage := params.FormValue(c.Ctx, "backgroundImage")
	if strs.IsBlank(backgroundImage) {
		return web.JsonErrorMsg("请上传图片")
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

	var followedSet *hashset.Set
	if currentUser != nil {
		followedSet = service.UserFollowService.GetMutualFollowers(currentUser.Id, followerIds...)
	}

	itemList := make([]payload.UserInfo, 0, len(followerIds))
	for _, id := range followerIds {
		item := payload.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
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
