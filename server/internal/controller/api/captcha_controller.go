package api

import (
	"bbs-go/pkg/bbsurls"
	"log/slog"

	"bbs-go/common/strs"
	"bbs-go/pkg/web"

	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
)

type CaptchaController struct {
	Ctx iris.Context
}

func (c *CaptchaController) GetRequest() *web.JsonResult {
	captchaId := c.Ctx.FormValue("captchaId")
	if strs.IsNotBlank(captchaId) { // reload
		if !captcha.Reload(captchaId) {
			// reload 失败，重新加载验证码
			captchaId = captcha.NewLen(4)
		}
	} else {
		captchaId = captcha.NewLen(4)
	}
	captchaUrl := bbsurls.AbsUrl("/api/captcha/show?captchaId=" + captchaId + "&r=" + strs.UUID())
	return web.NewEmptyRspBuilder().
		Put("captchaId", captchaId).
		Put("captchaUrl", captchaUrl).
		JsonResult()
}

func (c *CaptchaController) GetShow() {
	captchaId := c.Ctx.URLParam("captchaId")

	if captchaId == "" {
		c.Ctx.StatusCode(404)
		return
	}

	if !captcha.Reload(captchaId) {
		c.Ctx.StatusCode(404)
		return
	}

	c.Ctx.Header("Content-Type", "image/png")
	if err := captcha.WriteImage(c.Ctx.ResponseWriter(), captchaId, captcha.StdWidth, captcha.StdHeight); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *CaptchaController) GetVerify() *web.JsonResult {
	captchaId := c.Ctx.URLParam("captchaId")
	captchaCode := c.Ctx.URLParam("captchaCode")
	success := captcha.VerifyString(captchaId, captchaCode)
	return web.NewEmptyRspBuilder().Put("success", success).JsonResult()
}
