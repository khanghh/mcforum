package api

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/spam"
	"strconv"

	"bbs-go/sqls"
	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"

	"bbs-go/internal/service"
)

type CommentController struct {
	Ctx iris.Context
}

func (c *CommentController) GetClean() *web.JsonResult {
	go func() {
		p, _ := ants.NewPool(10)
		service.CommentService.Scan(func(comments []model.Comment) {
			var ids []int64
			for _, comment := range comments {
				if comment.ContentType == constants.ContentTypeHtml {
					ids = append(ids, comment.Id)
				}
			}
			if len(ids) > 0 {
				p.Submit(func() {
					sqls.DB().Delete(&model.Comment{}, "id in ?", ids)
					logrus.Info("清理评论:", ids)
				})
			}
		})
	}()
	return web.JsonSuccess()
}

func (c *CommentController) GetComments() *web.JsonResult {
	var (
		err        error
		cursor     int64
		entityType string
		entityId   int64
	)
	cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)

	if entityType, err = params.FormValueRequired(c.Ctx, "entityType"); err != nil {
		return web.JsonError(err)
	}
	if entityId, err = params.FormValueInt64(c.Ctx, "entityId"); err != nil {
		return web.JsonError(err)
	}
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := service.CommentService.GetComments(entityType, entityId, cursor)
	return web.JsonCursorData(payload.BuildComments(comments, currentUser, true, false), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) GetReplies() *web.JsonResult {
	var (
		cursor    = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		commentId = params.FormValueInt64Default(c.Ctx, "commentId", 0)
	)
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := service.CommentService.GetReplies(commentId, cursor, 10)
	return web.JsonCursorData(payload.BuildComments(comments, currentUser, false, true), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) PostCreate() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	form := model.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return web.JsonError(err)
	}

	comment, err := service.CommentService.Publish(user.Id, form)
	if err != nil {
		return web.JsonError(err)
	}

	return web.JsonData(payload.BuildComment(comment))
}
