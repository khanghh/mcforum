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
	event.RegHandler(reflect.TypeOf(event.TopicRecommendedEvent{}), handleTopicRecommend)
}

func handleTopicRecommend(i interface{}) {
	e, ok := i.(event.TopicRecommendedEvent)
	if !ok {
		return
	}

	if e.Topic.UserID == e.UserID {
		return
	}

	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:       e.UserID,
		ToId:         e.Topic.UserID,
		Title:        locale.T("message.title.topic_recommended"),
		QuoteContent: e.Topic.GetTitle(),
		DetailUrl:    bbsurls.TopicUrl(e.Topic.Slug, e.Topic.ID),
		Type:         msg.TypeTopicRecommend,
		ExtraData: &msg.TopicEventExtraData{
			TopicId: e.Topic.ID,
		},
	})
}
