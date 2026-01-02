package eventhandler

import (
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/pkg/msg"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UserFavoriteEvent{}), handleUserFavorite)
}

func handleUserFavorite(i interface{}) {
	e := i.(event.UserFavoriteEvent)
	sendTopicFavoriteMsg(e.EntityID, e.UserID)
}

// sendTopicFavoriteMsg topic favorited
func sendTopicFavoriteMsg(topicId, favoriteUserId int64) {
	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusActive {
		return
	}
	if topic.UserID == favoriteUserId {
		return
	}
	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       favoriteUserId,
		ToId:         topic.UserID,
		Title:        locale.T("message.title.topic_favorited"),
		QuoteContent: topic.GetTitle(),
		Type:         msg.TypeTopicFavorite,
		ExtraData: &msg.TopicEventExtraData{
			TopicId: topicId,
			UserId:  favoriteUserId,
		},
	})
}
