package eventhandler

import (
	"bbs-go/internal/event"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/search"
	"bbs-go/internal/service"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/msg"
	"log/slog"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicApprovedEvent{}), handleTopicApprovedEvent)
}

func handleTopicApprovedEvent(i interface{}) {
	e, ok := i.(event.TopicApprovedEvent)
	if !ok {
		return
	}

	topic := service.TopicService.Get(e.TopicID)
	if topic == nil {
		return
	}

	// Points
	search.UpdateTopicIndex(topic)
	service.UserService.IncrScoreForPostTopic(topic)
	sendTopicApprovedNotification(&e, topic)
	service.UserFollowService.ScanFollowers(topic.UserID, func(fansId int64) {
		slog.With(slog.Any("topicId", topic.ID), slog.Any("userId", topic.UserID), slog.Any("fansId", fansId)).Info("Notify new topic created to followers")
		if err := service.UserFeedService.Create(&model.UserFeed{
			UserID:     fansId,
			DataID:     topic.ID,
			DataType:   constants.EntityTopic,
			AuthorID:   topic.UserID,
			CreateTime: topic.CreateTime,
		}); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
			return
		}
		if topic.Status == constants.StatusActive {
			service.MessageService.SendMsg(service.SendMessageArgs{
				FromId:    topic.UserID,
				ToId:      fansId,
				Type:      msg.TypeFollowingUserCreateTopic,
				Title:     topic.Title,
				DetailUrl: bbsurls.TopicUrl(topic.Slug, topic.ID),
				ExtraData: &msg.TopicEventExtraData{
					TopicId: topic.ID,
					UserId:  topic.UserID,
				},
			})
		}
	})
}

func sendTopicApprovedNotification(e *event.TopicApprovedEvent, topic *model.Topic) {
	service.MessageService.SendMsg(service.SendMessageArgs{
		FromId:    e.UserID,
		ToId:      topic.UserID,
		Type:      msg.TypeTopicApproved,
		Title:     topic.Title,
		DetailUrl: bbsurls.TopicUrl(topic.Slug, topic.ID),
		ExtraData: &msg.TopicEventExtraData{
			TopicId: topic.ID,
			UserId:  topic.UserID,
		},
	})
}
