package payload

import (
	"bbs-go/common/base62"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/web/params"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kataras/iris/v12"
)

type CreateTopicForm struct {
	Type        constants.TopicType `json:"type"`
	CaptchaId   string              `json:"captchaId"`
	CaptchaCode string              `json:"captchaCode"`
	ForumId     int64               `json:"forumId"`
	Title       string              `json:"title"`
	Content     string              `json:"content"`
	HideContent string              `json:"hideContent"`
	Tags        []string            `json:"tags"`
	Images      []string            `json:"images"`
	UserAgent   string              `json:"userAgent"`
	Ip          string              `json:"ip"`
}

func GetCreateTopicForm(ctx iris.Context) CreateTopicForm {
	contentType := ctx.GetHeader("Content-Type")

	var form *CreateTopicForm
	if contentType == "application/json" {
		if err := ctx.ReadJSON(&form); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	} else {
		form = &CreateTopicForm{
			Type:        constants.TopicType(params.FormValueIntDefault(ctx, "type", int(constants.TopicTypeTopic))),
			CaptchaId:   params.FormValue(ctx, "captchaId"),
			CaptchaCode: params.FormValue(ctx, "captchaCode"),
			ForumId:     params.FormValueInt64Default(ctx, "nodeId", 0),
			Title:       strings.TrimSpace(params.FormValue(ctx, "title")),
			Content:     strings.TrimSpace(params.FormValue(ctx, "content")),
			HideContent: strings.TrimSpace(params.FormValue(ctx, "hideContent")),
			Tags:        params.FormValueStringArray(ctx, "tags"),
			Images:      params.FormValueStringArray(ctx, "images"),
			UserAgent:   common.GetUserAgent(ctx.Request()),
			Ip:          common.GetRequestIP(ctx.Request()),
		}
	}
	return *form
}

// CreateCommentForm 发表评论
type CreateCommentForm struct {
	QuoteId int64    `form:"quoteId"`
	Content string   `form:"content"`
	Images  []string `form:"images"`
}

func GetCreateCommentForm(ctx iris.Context) CreateCommentForm {
	fmt.Println("quoteIdStr", params.FormValue(ctx, "quoteId"))
	fmt.Println("quoteId", base62.Decode(params.FormValue(ctx, "quoteId")))
	form := CreateCommentForm{
		QuoteId: base62.Decode(params.FormValue(ctx, "quoteId")),
		Content: strings.TrimSpace(params.FormValue(ctx, "content")),
		Images:  params.FormValueStringArray(ctx, "images"),
	}
	return form
}
