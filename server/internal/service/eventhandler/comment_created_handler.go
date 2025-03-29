package eventhandler

import (
	"bbs-go/common/utils"
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/msg"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.CommentCreatedEvent{}), handleCommentCreate)
}

func handleCommentCreate(i interface{}) {
	e := i.(event.CommentCreatedEvent)

	// service.UserService.IncrCommentCount(userId)
	// service.UserService.IncrScoreForPostComment(comment)

	sendCommentTopicNotification(&e)
	sendReplyCommentNotification(&e)
	sendQuoteCommentNofitication(&e)
}

// sendCommentTopicNotification send message to notify topic owner of new comment created
func sendCommentTopicNotification(e *event.CommentCreatedEvent) {
	var (
		from = e.Comment.UserId
		to   = e.Topic.UserId
	)
	if from == to {
		return
	}

	// do not process replying to comment
	if e.Comment.ParentId != 0 {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       from,
		ToId:         to,
		Type:         msg.TypeTopicComment,
		Title:        locale.T("message.title.commented_your_topic"),
		QuoteContent: utils.GetSummaryHtml(e.Comment.Content, constants.CommentSummaryLen),
		DetailUrl:    bbsurls.TopicUrl(e.Topic.Slug, e.Topic.Id),
		ExtraData: msg.CommentExtraData{
			TopicId: e.Topic.Id,
		},
	})
}

func sendReplyCommentNotification(e *event.CommentCreatedEvent) {
	if e.Comment.ParentId == 0 {
		return
	}

	parentComment := service.CommentService.Get(e.Comment.ParentId)
	if parentComment == nil || parentComment.Status != constants.StatusOK {
		return
	}

	if e.Comment.UserId == parentComment.UserId {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       e.Comment.UserId,
		ToId:         parentComment.UserId,
		Title:        locale.T("message.title.replied_your_comment"),
		Type:         msg.TypeCommentReply,
		QuoteContent: utils.GetSummaryHtml(e.Comment.Content, constants.CommentSummaryLen),
		DetailUrl:    bbsurls.TopicUrl(e.Topic.Slug, e.Topic.Id),
		ExtraData: msg.CommentExtraData{
			TopicId:  e.Topic.Id,
			ParentId: e.Comment.ParentId,
			QuoteId:  e.Comment.QuoteId,
		},
	})

}

func sendQuoteCommentNofitication(e *event.CommentCreatedEvent) {
	if e.Comment.QuoteId == 0 {
		return
	}

	quoteComment := service.CommentService.Get(e.Comment.QuoteId)
	if quoteComment == nil || quoteComment.Status != constants.StatusOK {
		return
	}

	if e.Comment.UserId == quoteComment.UserId {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       e.Comment.UserId,
		ToId:         quoteComment.UserId,
		Title:        locale.T("message.title.replied_your_comment"),
		Type:         msg.TypeCommentReply,
		QuoteContent: utils.GetSummaryHtml(e.Comment.Content, constants.CommentSummaryLen),
		DetailUrl:    bbsurls.TopicUrl(e.Topic.Slug, e.Topic.Id),
		ExtraData: msg.CommentExtraData{
			TopicId:  e.Topic.Id,
			ParentId: e.Comment.ParentId,
			QuoteId:  e.Comment.QuoteId,
		},
	})
}
