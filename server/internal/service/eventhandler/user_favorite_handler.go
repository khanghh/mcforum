package eventhandler

import (
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/service"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UserFavoriteEvent{}), handleUserFavorite)
}

func handleUserFavorite(i interface{}) {
	e := i.(event.UserFavoriteEvent)

	if e.EntityType == constants.EntityTopic {
		sendTopicFavoriteMsg(e.EntityId, e.UserId)
	} else if e.EntityType == constants.EntityArticle {
		// TODO
	}
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
	var (
		from         = favoriteUserId
		to           = topic.UserId
		title        = "收藏了你的话题"
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	service.MessageService.SendMsg(from, to, msg.TypeTopicFavorite, title, "", quoteContent,
		&msg.TopicFavoriteExtraData{
			TopicId:        topicId,
			FavoriteUserId: favoriteUserId,
		})
}
