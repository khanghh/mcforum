package msg

// Message status
const (
	StatusUnread   = 0 // Message unread
	StatusHaveRead = 1 // Message read
)

type Type int

// Message type
const (
	TypeTopicComment             Type = 1
	TypeCommentReply             Type = 2
	TypeTopicLike                Type = 3
	TypeTopicFavorite            Type = 4
	TypeTopicRecommend           Type = 5
	TypeTopicPinned              Type = 6
	TypeTopicDelete              Type = 7
	TypeUserFollow               Type = 8
	TypeCommentLike              Type = 9
	TypeFollowingUserCreateTopic Type = 10
)

type TopicEventExtraData struct {
	TopicId int64 `json:"topicId,omitempty"`
	UserId  int64 `json:"userId,omitempty"`
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
