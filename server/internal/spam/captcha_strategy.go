package spam

import (
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/model"
)

type CaptchaStrategy struct{}

func (CaptchaStrategy) Name() string {
	return "CaptchaStrategy"
}

func (CaptchaStrategy) CheckTopic(user *model.User, form payload.CreateTopicForm) error {
	// if service.SysConfigService.IsEnabledTopicCaptcha() && !captcha.VerifyString(form.CaptchaId, form.CaptchaCode) {
	// 	return errs.CaptchaError
	// }
	return nil
}

func (CaptchaStrategy) CheckComment(user *model.User, form payload.CreateCommentForm) error {
	return nil
}
