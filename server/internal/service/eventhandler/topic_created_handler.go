package eventhandler

import (
	"bbs-go/internal/event"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/search"
	"bbs-go/internal/service"
	"log/slog"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicCreatedEvent{}), handleTopicCreatedEvent)
}

func handleTopicCreatedEvent(i interface{}) {
	e := i.(event.TopicCreatedEvent)

	// 积分
	search.UpdateTopicIndex(e.Topic)
	service.UserService.IncrScoreForPostTopic(e.Topic)
	service.UserFollowService.ScanFollowers(e.Topic.UserId, func(fansId int64) {
		slog.With(slog.Any("topicId", e.Topic.Id), slog.Any("userId", e.Topic.UserId), slog.Any("fansId", fansId)).Info("Notify new topic created to followers")
		if err := service.UserFeedService.Create(&model.UserFeed{
			UserId:     fansId,
			DataId:     e.Topic.Id,
			DataType:   constants.EntityTopic,
			AuthorId:   e.Topic.UserId,
			CreateTime: e.Topic.CreateTime,
		}); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	})
}
