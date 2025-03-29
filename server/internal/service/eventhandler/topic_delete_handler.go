package eventhandler

import (
	"bbs-go/internal/event"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicDeleteEvent{}), handleTopicDeleteEvent)
}

func handleTopicDeleteEvent(i interface{}) {
	e := i.(event.TopicDeleteEvent)

	// 处理userFeed
	service.UserFeedService.DeleteByDataId(e.TopicId, constants.EntityTopic)

	// 操作日志
	service.OperateLogService.AddOperateLog(e.DeleteUserId, constants.OpTypeDelete, constants.EntityTopic,
		e.TopicId, "", nil)
}

// sendTopicDeleteMsg 话题被删除消息
// func sendTopicDeleteMsg(topicId, deleteUserId int64) {
// 	topic := repository.TopicRepository.Get(sqls.DB(), topicId)
// 	if topic == nil {
// 		return
// 	}
// 	if topic.UserId == deleteUserId {
// 		return
// 	}

// 	service.MessageService.SendMsg(service.SendMessageArgs{
// 		FromId:       0,
// 		ToId:         topic.UserId,
// 		Type:         msg.TypeTopicDelete,
// 		Title:        "Topic deleted",
// 		QuoteContent: topic.GetTitle(),
// 		ExtraData: &msg.TopicDeleteExtraData{
// 			TopicId:      topicId,
// 			DeleteUserId: deleteUserId,
// 		},
// 	})
// }
