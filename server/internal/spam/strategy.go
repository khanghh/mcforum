package spam

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model"
)

type Strategy interface {
	// Name Strategy name
	Name() string
	// CheckTopic Check topic
	CheckTopic(user *model.User, form payload.CreateTopicForm) error
	// CheckComment Check comment
	CheckComment(user *model.User, form payload.CreateCommentForm) error
}
