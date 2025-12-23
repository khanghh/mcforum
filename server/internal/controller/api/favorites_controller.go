package api

import (
	"errors"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/errs"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
)

type FavoritesController struct {
	Ctx iris.Context
}

func (c *FavoritesController) PostAdd() *web.JsonResult {
	var (
		user       = service.UserTokenService.GetCurrent(c.Ctx)
		entityType = params.FormValue(c.Ctx, "entityType")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	var err error
	if entityType == constants.EntityTopic {
		err = service.FavoriteService.AddTopicFavorite(user.Id, entityId)
	} else if entityType == constants.EntityArticle {
		err = service.FavoriteService.AddArticleFavorite(user.Id, entityId)
	} else {
		err = errors.New("unsupproted")
	}

	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 取消收藏
func (c *FavoritesController) PostDelete() *web.JsonResult {
	var (
		user       = service.UserTokenService.GetCurrent(c.Ctx)
		entityType = params.FormValue(c.Ctx, "entityType")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	tmp := service.FavoriteService.GetBy(user.Id, entityType, entityId)
	if tmp != nil {
		service.FavoriteService.Delete(tmp.Id)
	}
	return web.JsonSuccess()
}
