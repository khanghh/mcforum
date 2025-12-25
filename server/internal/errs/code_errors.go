package errs

import (
	"bbs-go/pkg/web"
)

var (
	NotLogin            = web.NewError(1, "Please login first")
	CaptchaError        = web.NewError(1000, "Invalid captcha")
	ForbiddenError      = web.NewError(1001, "Forbidden")
	UserDisabled        = web.NewError(1002, "Account disabled")
	InObservationPeriod = web.NewError(1003, "Account in observation period")
	EmailNotVerified    = web.NewError(1004, "Please verify your email")
)
