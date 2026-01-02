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
		sendTopicLikedNotification(e.EntityID, e.UserID)
	} else if e.EntityType == constants.EntityComment {
		sendCommentLikedNotification(e.EntityID, e.UserID)
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
	if topic.UserID == likeUserId {
		return
	}
	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       likeUserId,
		ToId:         topic.UserID,
		Title:        locale.T("message.title.liked_your_topic"),
		QuoteContent: topic.GetTitle(),
		DetailUrl:    bbsurls.TopicUrl(topic.Slug, topic.ID),
		Type:         msg.TypeTopicLike,
		ExtraData: &msg.TopicEventExtraData{
			TopicId: topic.ID,
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
	if comment.UserID == likedUserId {
		return
	}

	topic := service.TopicService.Get(comment.TopicID)
	if topic == nil || topic.Status != constants.StatusActive {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       likedUserId,
		ToId:         comment.UserID,
		Title:        locale.T("message.title.liked_your_comment"),
		QuoteContent: utils.GetSummaryHtml(comment.Content, constants.SummaryLen),
		DetailUrl:    bbsurls.TopicUrl(topic.Slug, topic.ID),
		Type:         msg.TypeCommentLike,
		ExtraData: &msg.CommentLikeExtraData{
			TopicId:   topic.ID,
			CommentId: commentId,
			UserId:    likedUserId,
		},
	})
}
