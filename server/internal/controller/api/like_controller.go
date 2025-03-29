package api

import (
	"bbs-go/internal/errs"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"

	"bbs-go/common/strs"
	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"
)

type LikeController struct {
	Ctx iris.Context
}

func (c *LikeController) PostLike() *web.JsonResult {
	var (
		entityType = params.FormValue(c.Ctx, "entityType")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
		user       = service.UserTokenService.GetCurrent(c.Ctx)
		err        error
	)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if entityType == constants.EntityTopic {
		err = service.UserLikeService.TopicLike(user.Id, entityId)
	} else if entityType == constants.EntityArticle {
		err = service.UserLikeService.ArticleLike(user.Id, entityId)
	} else if entityType == constants.EntityComment {
		err = service.UserLikeService.CommentLike(user.Id, entityId)
	}
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *LikeController) PostUnlike() *web.JsonResult {
	var (
		entityType = params.FormValue(c.Ctx, "entityType")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
		user       = service.UserTokenService.GetCurrent(c.Ctx)
		err        error
	)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if entityType == constants.EntityTopic {
		err = service.UserLikeService.TopicUnLike(user.Id, entityId)
	} else if entityType == constants.EntityArticle {
		err = service.UserLikeService.ArticleUnLike(user.Id, entityId)
	} else if entityType == constants.EntityComment {
		err = service.UserLikeService.CommentUnLike(user.Id, entityId)
	}
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *LikeController) GetLiked_ids() *web.JsonResult {
	var (
		user           = service.UserTokenService.GetCurrent(c.Ctx)
		entityType     = params.FormValue(c.Ctx, "entityType")
		entityIds      = params.FormValueInt64Array(c.Ctx, "entityIds")
		likedEntityIds []int64
	)
	if user != nil {
		likedEntityIds = service.UserLikeService.GetUserLikes(user.Id, entityType, entityIds)
	}
	return web.JsonData(likedEntityIds)
}

func (c *LikeController) GetLiked() *web.JsonResult {
	var (
		user       = service.UserTokenService.GetCurrent(c.Ctx)
		entityType = params.FormValue(c.Ctx, "entityType")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	if user == nil || strs.IsBlank(entityType) || entityId <= 0 {
		return web.JsonData(false)
	} else {
		liked := service.UserLikeService.IsLiked(user.Id, entityType, entityId)
		return web.JsonData(liked)
	}
}
