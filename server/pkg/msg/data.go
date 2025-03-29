package msg

// 消息状态
const (
	StatusUnread   = 0 // 消息未读
	StatusHaveRead = 1 // 消息已读
)

type Type int

// 消息类型
const (
	TypeTopicComment   Type = 0
	TypeCommentReply   Type = 1
	TypeTopicLike      Type = 2
	TypeTopicFavorite  Type = 3
	TypeTopicRecommend Type = 4
	TypeTopicDelete    Type = 5
	TypeUserFollow     Type = 6
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
