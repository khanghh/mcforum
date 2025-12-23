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
	event.RegHandler(reflect.TypeOf(event.UserLikeEvent{}), handleUserLike)
	// event.RegHandler(reflect.TypeOf(event.UserUnLikeEvent{}), handleUserUnLike)
}

func handleUserLike(i interface{}) {
	e := i.(event.UserLikeEvent)

	if e.EntityType == constants.EntityTopic {
		sendTopicLikedNotification(e.EntityId, e.UserId)
	} else if e.EntityType == constants.EntityComment {
		sendCommentLikedNotification(e.EntityId, e.UserId)
	}
}

// func handleUserUnLike(i interface{}) {
// 	e := i.(event.UserUnLikeEvent)
// 	if e.EntityType == constants.EntityTopic {
// 		// TODO
// 	}
// }

// sendTopicLikedNotification sends a message to topic author when a topic is liked.
func sendTopicLikedNotification(topicId, likeUserId int64) {
	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusActive {
		return
	}
	if topic.UserId == likeUserId {
		return
	}
	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       likeUserId,
		ToId:         topic.UserId,
		Title:        locale.T("message.title.liked_your_topic"),
		QuoteContent: topic.GetTitle(),
		DetailUrl:    bbsurls.TopicUrl(topic.Slug, topic.Id),
		ExtraData: &msg.TopicLikeExtraData{
			TopicId: topic.Id,
			UserId:  likeUserId,
		},
	})
}

// sendTopicLikedMsg sends a message to comment author when a comment is liked.
func sendCommentLikedNotification(commentId, likedUserId int64) {
	comment := service.CommentService.Get(commentId)
	if comment == nil || comment.Status != constants.StatusActive {
		return
	}
	if comment.UserId == likedUserId {
		return
	}

	topic := service.TopicService.Get(comment.TopicId)
	if topic == nil || topic.Status != constants.StatusActive {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       likedUserId,
		ToId:         comment.UserId,
		Title:        locale.T("message.title.liked_your_comment"),
		QuoteContent: utils.GetSummaryHtml(comment.Content, constants.SummaryLen),
		DetailUrl:    bbsurls.TopicUrl(topic.Slug, topic.Id),
		ExtraData: &msg.CommentLikeExtraData{
			TopicId:   topic.Id,
			CommentId: commentId,
			UserId:    likedUserId,
		},
	})
}
