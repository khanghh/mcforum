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

	// 将该用户下的帖子添加到信息流
	// service.TopicService.ScanByUser(e.OtherId, func(topics []model.Topic) {
	// 	for _, topic := range topics {
	// 		if topic.Status != constants.StatusOK {
	// 			continue
	// 		}
	// 		_ = service.UserFeedService.Create(&model.UserFeed{
	// 			UserId:     e.UserId,
	// 			DataType:   constants.EntityTopic,
	// 			DataId:     topic.Id,
	// 			AuthorId:   topic.UserId,
	// 			CreateTime: topic.CreateTime,
	// 		})
	// 	}
	// })
	sendUserFollowNotification(&e)
}

func sendUserFollowNotification(e *event.FollowEvent) {
	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId: e.UserId,
		ToId:   e.OtherId,
		Type:   msg.TypeUserFollow,
		Title:  locale.T("message.title.started_following_you"),
	})
}
