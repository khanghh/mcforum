package payload

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/pkg/bbsurls"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/microcosm-cc/bluemonday"

	"bbs-go/common/strs"
	"bbs-go/common/urls"
	"bbs-go/pkg/web"

	"github.com/PuerkitoBio/goquery"

	"bbs-go/internal/service"
)

func xssProtection(htmlContent string) string {
	ugcProtection := bluemonday.UGCPolicy() // User generated content policy
	ugcProtection.AllowAttrs("class").OnElements("code")
	ugcProtection.AllowAttrs("start").OnElements("ol", "ul", "li")
	return ugcProtection.Sanitize(htmlContent)
}

// handleHtmlContent process html content
func handleHtmlContent(htmlContent string) string {
	htmlContent = xssProtection(htmlContent)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	doc.Find("a").Each(func(_ int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")

		if strs.IsBlank(href) {
			return
		}

		// not internal link
		if !bbsurls.IsInternalUrl(href) {
			selection.SetAttr("target", "_blank")
			selection.SetAttr("rel", "external nofollow") // mark external link, search engine crawler does not pass weight

			showRedirectPage := service.SysConfigService.IsEnabledUrlRedirectPage()
			if showRedirectPage { // enable non-internal link redirect
				newHref := urls.ParseUrl(bbsurls.AbsUrl("/redirect")).AddQuery("url", href).BuildStr()
				selection.SetAttr("href", newHref)
			}
		}

		// If a tag has no title, then set title
		title := selection.AttrOr("title", "")
		if len(title) == 0 {
			selection.SetAttr("title", selection.Text())
		}
	})

	// Process images
	doc.Find("img").Each(func(_ int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")

		// process third-party images
		if strings.Contains(src, "qpic.cn") {
			src = urls.ParseUrl("/api/img/proxy").AddQuery("url", src).BuildStr()
		}

		// process image style
		// // process lazyload
		// selection.SetAttr("data-src", src)
		// selection.RemoveAttr("src")

		selection.SetAttr("src", src)
	})

	if htmlStr, err := doc.Find("body").Html(); err == nil {
		return htmlStr
	}
	return htmlContent
}

/*
BuildLoginSuccess processes return data after successful login

Parameters:

	user - logged in user
	redirect - login source address, need to control jump to this address after successful login
*/
func BuildLoginSuccess(ctx iris.Context, user *model.User, redirect string) *web.JsonResult {
	token, err := service.UserTokenService.Generate(user.ID)
	if err != nil {
		return web.JsonError(err)
	}
	ctx.SetCookieKV(constants.CookieTokenKey, token, context.CookieHTTPOnly(true), context.CookieExpires(365*24*time.Hour))
	return web.NewEmptyRspBuilder().
		Put("token", token).
		Put("user", BuildUserProfile(user)).
		Put("redirect", redirect).JsonResult()
}
