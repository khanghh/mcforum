package api

import (
	"bbs-go/internal/service"
	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"
)

type MeController struct {
	Ctx iris.Context
}

// PUT /api/me/favorites
func (c *MeController) PutFavorites() (*web.JsonResult, error) {
	topicId, err := params.FormValueInt64(c.Ctx, "topicId")
	if err != nil {
		return web.JsonError(service.ErrBadRequest), nil
	}

	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(service.ErrUnauthorized), nil
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
		return web.JsonError(service.ErrUnauthorized), nil
	}

	err := service.FavoriteService.RemoveTopicFavorite(user.Id, topicId)
	if err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}
