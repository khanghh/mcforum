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
	event.RegHandler(reflect.TypeOf(event.TopicRecommendedEvent{}), handleTopicRecommend)
}

func handleTopicRecommend(i interface{}) {
	e := i.(event.TopicRecommendedEvent)

	if e.Recommended {
		sendTopicRecommendMsg(e.TopicId)
	}
}

// sendTopicRecommendMsg 话题被设为推荐
func sendTopicRecommendMsg(topicId int64) {
	topic := service.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOK {
		return
	}
	var (
		from         int64 = 0
		to                 = topic.UserId
		title              = locale.T("errors.topic_recommended")
		quoteContent       = topic.GetTitle()
	)
	service.MessageService.SendMsg(from, to, msg.TypeTopicRecommend, title, "", quoteContent,
		&msg.TopicRecommendExtraData{
			TopicId: topicId,
		})
}
