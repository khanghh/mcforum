package spam

import "bbs-go/internal/model"

type Strategy interface {
	// Name 策略名称
	Name() string
	// CheckTopic 检查话题
	CheckTopic(user *model.User, form model.CreateTopicForm) error
	// CheckComment 检查评论
	CheckComment(user *model.User, form model.CreateCommentForm) error
}
