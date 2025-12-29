package event

import "bbs-go/internal/model"

// FollowEvent Follow
type FollowEvent struct {
	UserId  int64 `json:"userId"`
	OtherId int64 `json:"otherId"`
}

// UnfollowEvent Unfollow
type UnfollowEvent struct {
	UserId  int64 `json:"userId"`
	OtherId int64 `json:"otherId"`
}

type TopicCreatedEvent struct {
	Topic *model.Topic
}

type TopicApprovedEvent struct {
	UserID  int64
	TopicID int64
}

type TopicRejectedEvent struct {
	UserID  int64
	TopicID int64
	Reason  string
}

type TopicDeleteEvent struct {
	UserId       int64 `json:"userId"`
	TopicId      int64 `json:"topicId"`
	DeleteUserId int64 `json:"deleteUserId"`
}

type UserLikeEvent struct {
	UserId     int64  `json:"userId"`
	EntityId   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type UserUnLikeEvent struct {
	UserId     int64  `json:"userId"`
	EntityId   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type UserFavoriteEvent struct {
	UserId     int64  `json:"userId"`
	EntityId   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type UserUnfavoriteEvent struct {
	UserId     int64  `json:"userId"`
	EntityId   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type CommentCreatedEvent struct {
	Topic   *model.Topic
	Comment *model.Comment
}

type TopicRecommendedEvent struct {
	UserID int64
	Topic  *model.Topic
}

type TopicPinedEvent struct {
	UserID int64
	Topic  *model.Topic
}
