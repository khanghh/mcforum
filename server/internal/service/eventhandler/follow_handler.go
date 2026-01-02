package eventhandler

import (
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/service"
	"bbs-go/pkg/msg"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.FollowEvent{}), handleFollowEvent)
}

func handleFollowEvent(i interface{}) {
	e := i.(event.FollowEvent)

	// Add the user's topics to the feed
	// service.TopicService.ScanByUser(e.OtherID, func(topics []model.Topic) {
	// 	for _, topic := range topics {
	// 		if topic.Status != constants.StatusActive {
	// 			continue
	// 		}
	// 		_ = service.UserFeedService.Create(&model.UserFeed{
	// 			UserID:     e.UserID,
	// 			DataType:   constants.EntityTopic,
	// 			DataID:     topic.ID,
	// 			AuthorID:   topic.UserID,
	// 			CreateTime: topic.CreateTime,
	// 		})
	// 	}
	// })
	sendUserFollowNotification(&e)
}

func sendUserFollowNotification(e *event.FollowEvent) {
	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId: e.UserID,
		ToId:   e.OtherID,
		Type:   msg.TypeUserFollow,
		Title:  locale.T("message.title.started_following_you"),
	})
}
