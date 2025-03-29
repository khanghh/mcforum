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
	sendTopicFavoriteMsg(e.EntityId, e.UserId)
}

// sendTopicFavoriteMsg 话题被收藏
func sendTopicFavoriteMsg(topicId, favoriteUserId int64) {
	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return
	}
	if topic.UserId == favoriteUserId {
		return
	}
	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       favoriteUserId,
		ToId:         topic.UserId,
		Title:        locale.T("message.title.topic_favorited"),
		QuoteContent: topic.GetTitle(),
		ExtraData: &msg.TopicFavoriteExtraData{
			TopicId: topicId,
			UserId:  favoriteUserId,
		},
	})
}
