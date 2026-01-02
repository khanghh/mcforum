package eventhandler

import (
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/model"
	"bbs-go/internal/service"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/msg"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicRejectedEvent{}), handleTopicRejectedEvent)
}

func handleTopicRejectedEvent(i interface{}) {
	e, ok := i.(event.TopicRejectedEvent)
	if !ok {
		return
	}

	topic := service.TopicService.Get(e.TopicID)
	if topic == nil {
		return
	}
	sendTopicRejectedNotification(&e, topic)
}

func sendTopicRejectedNotification(e *event.TopicRejectedEvent, topic *model.Topic) {
	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:    e.UserID,
		ToId:      topic.UserID,
		Type:      msg.TypeTopicRejected,
		Title:     locale.T("message.title.rejected_your_topic"),
		DetailUrl: bbsurls.TopicUrl(topic.Slug, topic.ID),
		ExtraData: &msg.TopicEventExtraData{
			TopicId: topic.ID,
			UserId:  topic.UserID,
		},
	})
}
