package spam

import (
	"bbs-go/internal/model"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/service"
	"fmt"
)

type EmailVerifyStrategy struct{}

func (EmailVerifyStrategy) Name() string {
	return "EmailVerifyStrategy"
}

func (EmailVerifyStrategy) CheckTopic(user *model.User, form model.CreateTopicForm) error {
	fmt.Println("EmailVerified", user.EmailVerified)
	if service.SysConfigService.IsCreateTopicEmailVerified() && !user.EmailVerified {
		return errs.EmailNotVerified
	}
	return nil
}

func (EmailVerifyStrategy) CheckArticle(user *model.User, form model.CreateArticleForm) error {
	if service.SysConfigService.IsCreateArticleEmailVerified() && !user.EmailVerified {
		return errs.EmailNotVerified
	}
	return nil
}

func (EmailVerifyStrategy) CheckComment(user *model.User, form model.CreateCommentForm) error {
	if service.SysConfigService.IsCreateCommentEmailVerified() && !user.EmailVerified {
		return errs.EmailNotVerified
	}
	return nil
}
