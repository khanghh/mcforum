package api

import (
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/service"
	"bbs-go/web"

	"github.com/kataras/iris/v12"
)

type MeController struct {
	Ctx iris.Context
}

// POST /topics/{slugId}/favorite
func (c *MeController) AddFavoriteTopic(topicId int) (*web.JsonResult, int) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin), iris.StatusUnauthorized
	}
	err := service.FavoriteService.AddTopicFavorite(user.Id, int64(topicId))
	if err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}

// POST /topics/{slugId}/unfavorite
func (c *MeController) RemoveFavoriteTopic(topicId int) (*web.JsonResult, int) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin), iris.StatusUnauthorized
	}
	err := service.FavoriteService.RemoveTopicFavorite(user.Id, int64(topicId))
	if err != nil {
		return web.JsonError(err), iris.StatusInternalServerError
	}
	return web.JsonSuccess(), iris.StatusOK
}
