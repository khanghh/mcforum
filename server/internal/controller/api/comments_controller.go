package api

import (
	"bbs-go/common/utils"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/spam"

	"bbs-go/pkg/web"
	"bbs-go/pkg/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/service"
)

type CommentsController struct {
	Ctx iris.Context
}

// func (c *CommentController) GetClean() *web.JsonResult {
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

func (c *CommentsController) PostByReplies(commentId int64) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
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

	topic := service.TopicService.Get(parent.TopicId)
	if topic == nil || topic.Status != constants.StatusActive {
		return web.JsonError(errs.ErrTopicNotFound)
	}

	parentId := parent.Id
	if parent.ParentId != 0 {
		parentId = parent.ParentId
	}

	comment, err := service.CommentService.CreateComment(service.CreateCommentArgs{
		UserId:    user.Id,
		TopicId:   parent.TopicId,
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
	if err := service.UserLikeService.CommentLike(user.Id, commentId); err != nil {
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
	if err := service.UserLikeService.CommentUnLike(user.Id, commentId); err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}
