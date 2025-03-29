package eventhandler

import (
	"bbs-go/internal/event"
	"bbs-go/internal/service"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.UnfollowEvent{}), handleUnFollowEvent)
}

func handleUnFollowEvent(i interface{}) {
	e := i.(event.UnfollowEvent)

	// 清理该用户下的信息流
	service.UserFeedService.DeleteByUser(e.UserId, e.OtherId)
}
