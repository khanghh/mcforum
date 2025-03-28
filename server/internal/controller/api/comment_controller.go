package api

import (
	"bbs-go/common/base62"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/spam"
	"strconv"

	"bbs-go/web"
	"bbs-go/web/params"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/service"
)

type CommentController struct {
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

func (c *CommentController) PostByReplies(bas62Id string) *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	form := payload.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return web.JsonError(err)
	}

	commentId := base62.Decode(bas62Id)
	parent := service.CommentService.Get(commentId)
	if parent == nil {
		return web.JsonError(service.ErrCommentNotFound)
	}
	if parent.Status != constants.StatusOK {
		return web.JsonError(service.ErrCommentDeleted)
	}

	topic := service.TopicService.Get(parent.TopicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return web.JsonError(service.ErrTopicNotFound)
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
		UserAgent: common.GetUserAgent(c.Ctx.Request()),
		IPAddress: common.GetRequestIP(c.Ctx.Request()),
	})

	if err != nil {
		return web.JsonError(err)
	}

	return web.JsonData(payload.BuildComment(comment))
}

func (c *CommentController) GetByReplies(base62Id string) *web.JsonResult {
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	commentId := base62.Decode(base62Id)
	comments, cursor, hasMore := service.CommentService.GetReplies(commentId, cursor, 10)
	currentUser := service.UserTokenService.GetCurrent(c.Ctx)
	return web.JsonCursorData(payload.BuildComments(comments, currentUser, false, true), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) PostByReactions(base62Id string) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return nil, service.ErrForbidden
	}
	commentId := base62.Decode(base62Id)
	if err := service.UserLikeService.CommentLike(user.Id, commentId); err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}

func (c *CommentController) DeleteByReactionsBy(base62Id string, userId int64) (*web.JsonResult, error) {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil || user.Id != userId {
		return nil, service.ErrForbidden
	}

	commentId := base62.Decode(base62Id)
	if err := service.UserLikeService.CommentUnLike(userId, commentId); err != nil {
		return nil, err
	}
	return web.JsonSuccess(), nil
}
