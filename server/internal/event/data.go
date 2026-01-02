package event

import "bbs-go/internal/model"

// FollowEvent Follow
type FollowEvent struct {
	UserID  int64 `json:"userId"`
	OtherID int64 `json:"otherId"`
}

// UnfollowEvent Unfollow
type UnfollowEvent struct {
	UserID  int64 `json:"userId"`
	OtherID int64 `json:"otherId"`
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
	UserID       int64 `json:"userId"`
	TopicID      int64 `json:"topicId"`
	DeleteUserID int64 `json:"deleteUserId"`
}

type UserLikeEvent struct {
	UserID     int64  `json:"userId"`
	EntityID   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type UserUnLikeEvent struct {
	UserID     int64  `json:"userId"`
	EntityID   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type UserFavoriteEvent struct {
	UserID     int64  `json:"userId"`
	EntityID   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type UserUnfavoriteEvent struct {
	UserID     int64  `json:"userId"`
	EntityID   int64  `json:"entityId"`
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
