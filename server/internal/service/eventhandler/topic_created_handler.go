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
	event.RegHandler(reflect.TypeOf(event.TopicCreatedEvent{}), handleTopicCreatedEvent)
}

func notifiAdminReviewTopic(topic *model.Topic) {
	// admins := service.UserService.GetAdminUsers()
}

func handleTopicCreatedEvent(i interface{}) {
	e := i.(event.TopicCreatedEvent)

	// Points
	search.UpdateTopicIndex(e.Topic)
	service.UserService.IncrScoreForPostTopic(e.Topic)
	service.UserFollowService.ScanFollowers(e.Topic.UserID, func(fansId int64) {
		slog.With(slog.Any("topicId", e.Topic.ID), slog.Any("userId", e.Topic.UserID), slog.Any("fansId", fansId)).Info("Notify new topic created to followers")
		if err := service.UserFeedService.Create(&model.UserFeed{
			UserID:     fansId,
			DataID:     e.Topic.ID,
			DataType:   constants.EntityTopic,
			AuthorID:   e.Topic.UserID,
			CreateTime: e.Topic.CreateTime,
		}); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
			return
		}
		if e.Topic.Status == constants.StatusActive {
			service.MessageService.SendMsg(service.SendMessageArgs{
				FromId:    e.Topic.UserID,
				ToId:      fansId,
				Type:      msg.TypeFollowingUserCreateTopic,
				Title:     "message.title.posted_a_topic",
				DetailUrl: bbsurls.TopicUrl(e.Topic.Slug, e.Topic.ID),
				ExtraData: &msg.TopicRecommendExtraData{
					TopicId: e.Topic.ID,
				},
			})
		}
	})
}
