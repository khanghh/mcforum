package eventhandler

import (
	"bbs-go/internal/event"
	"bbs-go/internal/locale"
	"bbs-go/internal/service"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/msg"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicPinedEvent{}), handleTopicPinEvent)
}

func handleTopicPinEvent(i interface{}) {
	e, ok := i.(event.TopicPinedEvent)
	if !ok {
		return
	}

	if e.Topic.UserID == e.UserID {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       e.UserID,
		ToId:         e.Topic.UserID,
		Title:        locale.T("message.title.topic_pinned"),
		QuoteContent: e.Topic.GetTitle(),
		DetailUrl:    bbsurls.TopicUrl(e.Topic.Slug, e.Topic.ID),
		Type:         msg.TypeTopicPinned,
		ExtraData: &msg.TopicEventExtraData{
			TopicId: e.Topic.ID,
		},
	})
}
