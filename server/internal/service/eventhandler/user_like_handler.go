package eventhandler

import (
	"bbs-go/common/utils"
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/pkg/msg"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UserLikeEvent{}), handleUserLike)
	event.RegHandler(reflect.TypeOf(event.UserUnLikeEvent{}), handleUserUnLike)
}

func handleUserLike(i interface{}) {
	e := i.(event.UserLikeEvent)

	if e.EntityType == constants.EntityTopic {
		sendTopicLikedMsg(e.EntityId, e.UserId)
	} else if e.EntityType == constants.EntityComment {
		// TODO
	}
}

func handleUserUnLike(i interface{}) {
	e := i.(event.UserUnLikeEvent)
	if e.EntityType == constants.EntityTopic {
		// TODO
	}
}

// sendTopicLikedMsg sends a message to topic author when a topic is liked.
func sendTopicLikedMsg(topicId, likeUserId int64) {
	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return
	}
	if topic.UserId == likeUserId {
		return
	}
	var (
		from         = likeUserId
		to           = topic.UserId
		title        = locale.T("message.title.liked_your_topic")
		quoteContent = topic.GetTitle()
	)
	service.MessageService.SendMsg(from, to, msg.TypeTopicLike, title, "", quoteContent,
		&msg.TopicLikeExtraData{
			TopicId:    topicId,
			LikeUserId: likeUserId,
		})
}

// sendTopicLikedMsg sends a message to comment author when a comment is liked.
func sendCommentLikedMsg(commentId, likeUserId int64) {
	comment := service.CommentService.Get(commentId)
	if comment == nil || comment.Status != constants.StatusOK {
		return
	}
	if comment.UserId == likeUserId {
		return
	}
	var (
		from         = likeUserId
		to           = comment.UserId
		title        = locale.T("message.title.liked_your_comment")
		quoteContent = utils.GetSummaryHtml(comment.Content, constants.SummaryLen)
	)
	service.MessageService.SendMsg(from, to, msg.TypeTopicLike, title, "", quoteContent,
		&msg.TopicLikeExtraData{
			TopicId:    commentId,
			LikeUserId: likeUserId,
		})
}
