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

	sendCommentTopicNotification(&e)
	sendReplyCommentNotification(&e)
	sendQuoteCommentNofitication(&e)
	service.UserService.IncreaseCommentCount(e.Comment.UserID)
	service.UserService.IncrScoreForPostComment(e.Comment)
}

// sendCommentTopicNotification send message to notify topic owner of new comment created
func sendCommentTopicNotification(e *event.CommentCreatedEvent) {
	var (
		from = e.Comment.UserID
		to   = e.Topic.UserID
	)
	if from == to {
		return
	}

	// do not process replying to comment
	if e.Comment.ParentID != 0 {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       from,
		ToId:         to,
		Type:         msg.TypeTopicComment,
		Title:        locale.T("message.title.commented_your_topic"),
		QuoteContent: utils.GetSummaryHtml(e.Comment.Content, constants.CommentSummaryLen),
		DetailUrl:    bbsurls.TopicUrl(e.Topic.Slug, e.Topic.ID),
		ExtraData: msg.CommentExtraData{
			TopicId: e.Topic.ID,
		},
	})
}

func sendReplyCommentNotification(e *event.CommentCreatedEvent) {
	if e.Comment.ParentID == 0 {
		return
	}

	parentComment := service.CommentService.Get(e.Comment.ParentID)
	if parentComment == nil || parentComment.Status != constants.StatusActive {
		return
	}

	if e.Comment.UserID == parentComment.UserID {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       e.Comment.UserID,
		ToId:         parentComment.UserID,
		Title:        locale.T("message.title.replied_your_comment"),
		Type:         msg.TypeCommentReply,
		QuoteContent: utils.GetSummaryHtml(e.Comment.Content, constants.CommentSummaryLen),
		DetailUrl:    bbsurls.TopicUrl(e.Topic.Slug, e.Topic.ID),
		ExtraData: msg.CommentExtraData{
			TopicId:  e.Topic.ID,
			ParentId: e.Comment.ParentID,
			QuoteId:  e.Comment.QuoteID,
		},
	})

}

func sendQuoteCommentNofitication(e *event.CommentCreatedEvent) {
	if e.Comment.QuoteID == 0 {
		return
	}

	quoteComment := service.CommentService.Get(e.Comment.QuoteID)
	if quoteComment == nil || quoteComment.Status != constants.StatusActive {
		return
	}

	if e.Comment.UserID == quoteComment.UserID {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       e.Comment.UserID,
		ToId:         quoteComment.UserID,
		Title:        locale.T("message.title.replied_your_comment"),
		Type:         msg.TypeCommentReply,
		QuoteContent: utils.GetSummaryHtml(e.Comment.Content, constants.CommentSummaryLen),
		DetailUrl:    bbsurls.TopicUrl(e.Topic.Slug, e.Topic.ID),
		ExtraData: msg.CommentExtraData{
			TopicId:  e.Topic.ID,
			ParentId: e.Comment.ParentID,
			QuoteId:  e.Comment.QuoteID,
		},
	})
}
