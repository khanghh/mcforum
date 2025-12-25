package payload

import (
	"bbs-go/common/utils"
	"bbs-go/internal/model/constants"
	"bbs-go/pkg/web/params"
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
			UserAgent:   utils.GetUserAgent(ctx.Request()),
			Ip:          utils.GetRequestIP(ctx.Request()),
		}
	}
	return *form
}

// CreateCommentForm post comment
type CreateCommentForm struct {
	QuoteId int64    `form:"quoteId"`
	Content string   `form:"content"`
	Images  []string `form:"images"`
}

func GetCreateCommentForm(ctx iris.Context) CreateCommentForm {
	contentType := ctx.GetHeader("Content-Type")

	if contentType == "application/json" {
		var form CreateCommentForm
		if err := ctx.ReadJSON(&form); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
		return form
	}

	form := CreateCommentForm{
		Content: strings.TrimSpace(params.FormValue(ctx, "content")),
		Images:  params.FormValueStringArray(ctx, "images"),
		QuoteId: params.FormValueInt64Default(ctx, "quoteId", 0),
	}

	return form
}
