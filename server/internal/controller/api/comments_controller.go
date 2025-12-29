package api

import (
	"bbs-go/common/utils"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/spam"
	"log/slog"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/service"
)

type CommentsController struct {
	Ctx iris.Context
}

func (c *CommentsController) PostByReplies(commentId int64) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}

	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	form := payload.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return web.JsonError(err)
	}

	parent := service.CommentService.Get(commentId)
	if parent == nil {
		return web.JsonError(errs.ErrCommentNotFound)
	}
	if parent.Status != constants.StatusActive {
		return web.JsonError(errs.ErrCommentDeleted)
	}

	topic := service.TopicService.Get(parent.TopicID)
	if topic == nil {
		return web.JsonError(errs.ErrTopicNotFound)
	}
	if err := canViewTopic(user, topic); err != nil {
		return web.JsonError(err)
	}

	parentId := parent.ID
	if parent.ParentID != 0 {
		parentId = parent.ParentID
	}

	comment, err := service.CommentService.CreateComment(service.CreateCommentArgs{
		UserId:    user.ID,
		TopicId:   parent.TopicID,
		ParentId:  parentId,
		QuoteId:   form.QuoteId,
		Content:   form.Content,
		Images:    form.Images,
		UserAgent: utils.GetUserAgent(c.Ctx.Request()),
		IPAddress: utils.GetRequestIP(c.Ctx.Request()),
	})

	if err != nil {
		return web.JsonError(err)
	}

	return web.JsonData(payload.BuildComment(comment))
}

func (c *CommentsController) GetByReplies(commentId int64) *web.JsonResult {
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	comments, cursor, hasMore := service.CommentService.GetReplies(commentId, cursor, 10)
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	return web.JsonCursorData(payload.BuildComments(comments, currentUser, false, true), cursor, hasMore)
}

// POST /api/comments/:id/reactions
func (c *CommentsController) PostByReactions(commentId int64) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return nil, errs.ErrUnauthorized
	}
	if err := service.UserLikeService.CommentLike(user.ID, commentId); err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

// DELETE /api/comments/:id/reactions
func (c *CommentsController) DeleteByReactions(commentId int64) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return nil, errs.ErrUnauthorized
	}
	if err := service.UserLikeService.CommentUnLike(user.ID, commentId); err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

func (c *CommentsController) DeleteBy(commentId int64) *web.JsonResult {
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	if currentUser == nil {
		return web.JsonError(errs.ErrUnauthorized)
	}
	comment := service.CommentService.Get(commentId)
	if comment == nil {
		return web.JsonError(errs.ErrCommentNotFound)
	}

	if comment.UserID != currentUser.ID && !currentUser.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonError(errs.ErrForbidden)
	}
	if err := service.CommentService.Delete(comment.ID); err != nil {
		slog.Error("Delete comment failed", "error", err, "commentId", commentId)
		return web.JsonError(errs.ErrInternalServer)
	}
	return web.JsonSuccess()
}

// func (c *CommentsController) GetClean() *web.JsonResult {
// 	go func() {
// 		p, _ := ants.NewPool(10)
// 		service.CommentService.Scan(func(comments []model.Comment) {
// 			var ids []int64
// 			for _, comment := range comments {
// 				if comment.ContentType == constants.ContentTypeHtml {
// 					ids = append(ids, comment.Id)
// 				}
// 			}
// 			if len(ids) > 0 {
// 				p.Submit(func() {
// 					sqls.DB().Delete(&model.Comment{}, "id in ?", ids)
// 					slog.Info("Comments cleaned up", "ids", ids)
// 				})
// 			}
// 		})
// 	}()
// 	return web.JsonSuccess()
// }
