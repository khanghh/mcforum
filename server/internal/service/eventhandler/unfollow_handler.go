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

	// Clean up the user's feed
	service.UserFeedService.DeleteByUser(e.UserID, e.OtherID)
}
