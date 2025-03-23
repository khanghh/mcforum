package spam

import (
	"bbs-go/internal/model"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/service"

	"github.com/dchest/captcha"
)

type CaptchaStrategy struct{}

func (CaptchaStrategy) Name() string {
	return "CaptchaStrategy"
}

func (CaptchaStrategy) CheckTopic(user *model.User, form model.CreateTopicForm) error {
	if service.SysConfigService.IsEnabledTopicCaptcha() && !captcha.VerifyString(form.CaptchaId, form.CaptchaCode) {
		return errs.CaptchaError
	}
	return nil
}

func (CaptchaStrategy) CheckArticle(user *model.User, form model.CreateArticleForm) error {
	return nil
}

func (CaptchaStrategy) CheckComment(user *model.User, form model.CreateCommentForm) error {
	return nil
}
