package msg

// Message status
const (
	StatusUnread   = 0 // Message unread
	StatusHaveRead = 1 // Message read
)

type Type int

// Message type
const (
	TypeTopicComment             Type = 0
	TypeCommentReply             Type = 1
	TypeTopicLike                Type = 2
	TypeTopicFavorite            Type = 3
	TypeTopicRecommend           Type = 4
	TypeTopicDelete              Type = 5
	TypeUserFollow               Type = 6
	TypeCommentLike              Type = 7
	TypeFollowingUserCreateTopic Type = 8
)

type TopicLikeExtraData struct {
	TopicId int64 `json:"topicId,omitempty"`
	UserId  int64 `json:"userId,omitempty"`
}

type TopicFavoriteExtraData struct {
	TopicId int64 `json:"topicId,omitempty"`
	UserId  int64 `json:"userId,omitempty"`
}

type TopicRecommendExtraData struct {
	TopicId int64 `json:"topicId,omitempty"`
}

type CommentExtraData struct {
	TopicId  int64 `json:"topicId,omitempty"`
	ParentId int64 `json:"parentId,omitempty"`
	QuoteId  int64 `json:"quoteId,omitempty"`
}

type CommentLikeExtraData struct {
	UserId    int64 `json:"userId,omitempty"`
	TopicId   int64 `json:"topicId,omitempty"`
	CommentId int64 `json:"commentId,omitempty"`
}
