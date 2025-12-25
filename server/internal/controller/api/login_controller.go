package api

import (
	"bbs-go/pkg/web"

	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"

	"bbs-go/internal/controller/payload"
	"bbs-go/internal/errs"
	"bbs-go/internal/service"
)

type LoginController struct {
	Ctx iris.Context
}

// Signup
func (c *LoginController) PostSignup() *web.JsonResult {
	var (
		captchaId   = c.Ctx.PostValueTrim("captchaId")
		captchaCode = c.Ctx.PostValueTrim("captchaCode")
		email       = c.Ctx.PostValueTrim("email")
		username    = c.Ctx.PostValueTrim("username")
		password    = c.Ctx.PostValueTrim("password")
		rePassword  = c.Ctx.PostValueTrim("rePassword")
		nickname    = c.Ctx.PostValueTrim("nickname")
		redirect    = c.Ctx.FormValue("redirect")
	)
	loginMethod := service.SysConfigService.GetLoginMethod()
	if !loginMethod.Password {
		return web.JsonErrorMsg("Username/password login/registration is disabled")
	}
	if !captcha.VerifyString(captchaId, captchaCode) {
		return web.JsonError(errs.CaptchaError)
	}
	user, err := service.UserService.SignUp(username, email, nickname, password, rePassword)
	if err != nil {
		return web.JsonError(err)
	}
	return payload.BuildLoginSuccess(c.Ctx, user, redirect)
}

// Signin (username/password)
func (c *LoginController) PostSignin() *web.JsonResult {
	var (
		username = c.Ctx.PostValueTrim("username")
		password = c.Ctx.PostValueTrim("password")
		redirect = c.Ctx.FormValue("redirect")
	)
	loginMethod := service.SysConfigService.GetLoginMethod()
	if !loginMethod.Password {
		return web.JsonErrorMsg("Username/password login/registration is disabled")
	}

	user, err := service.UserService.SignIn(username, password)
	if err != nil {
		return web.JsonError(err)
	}
	return payload.BuildLoginSuccess(c.Ctx, user, redirect)
}

// Signout
func (c *LoginController) GetSignout() *web.JsonResult {
	err := service.UserTokenService.Signout(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
