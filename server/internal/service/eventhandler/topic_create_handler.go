package eventhandler

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/seo"
	"bbs-go/internal/service"
	"log/slog"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicCreateEvent{}), handleTopicCreateEvent)
}

func handleTopicCreateEvent(i interface{}) {
	e := i.(event.TopicCreateEvent)

	service.UserFollowService.ScanFans(e.UserId, func(fansId int64) {
		slog.With(slog.Any("topicId", e.TopicId), slog.Any("userId", e.UserId), slog.Any("fansId", fansId)).Info("用户关注，处理帖子")
		if err := service.UserFeedService.Create(&model.UserFeed{
			UserId:     fansId,
			DataId:     e.TopicId,
			DataType:   constants.EntityTopic,
			AuthorId:   e.UserId,
			CreateTime: e.CreateTime,
		}); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	})

	// 百度链接推送
	seo.Push(bbsurls.TopicUrl(e.TopicId))
}
