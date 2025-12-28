package payload

import (
	"bbs-go/internal/model/constants"
	"bbs-go/pkg/web/params"
	"log/slog"
	"strings"

	"github.com/kataras/iris/v12"
)

type TopicPollOptions struct {
	Question        string   `json:"question"`
	Options         []string `json:"options"`
	DurationHours   uint64   `json:"duration"`
	MultiSelect     bool     `json:"multiSelect"`
	PublicResults   bool     `json:"publicResults"`
	AllowVoteChange bool     `json:"allowVoteChange"`
}

type CreateTopicForm struct {
	Type          constants.TopicType `json:"type"`
	ForumID       int64               `json:"forumId"`
	Title         string              `json:"title"`
	Content       string              `json:"content"`
	HiddenContent string              `json:"hiddenContent"`
	Tags          []string            `json:"tags"`
	Images        []string            `json:"images"`
	UserAgent     string              `json:"userAgent"`
	Ip            string              `json:"ip"`
}

func GetCreateTopicForm(ctx iris.Context) CreateTopicForm {
	contentType := ctx.GetHeader("Content-Type")

	var form CreateTopicForm
	if contentType == "application/json" {
		if err := ctx.ReadJSON(&form); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	}
	return form
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

type UpdateProfileForm struct {
	Avatar        string `form:"avatar"`
	Nickname      string `form:"nickname"`
	Bio           string `form:"bio"`
	Location      string `form:"location"`
	LockedProfile bool   `form:"lockedProfile"`
	ShowLocation  bool   `form:"showLocation"`
	EmailNotify   bool   `form:"emailNotify"`
}

func GetUpdateProfileForm(ctx iris.Context) UpdateProfileForm {
	contentType := ctx.GetHeader("Content-Type")

	if contentType == "application/json" {
		var form UpdateProfileForm
		if err := ctx.ReadJSON(&form); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
		return form
	}

	form := UpdateProfileForm{
		Avatar:        strings.TrimSpace(params.FormValue(ctx, "avatar")),
		Nickname:      strings.TrimSpace(params.FormValue(ctx, "nickname")),
		Bio:           strings.TrimSpace(params.FormValue(ctx, "bio")),
		Location:      strings.TrimSpace(params.FormValue(ctx, "location")),
		LockedProfile: params.FormValueBoolDefault(ctx, "lockedProfile", false),
		ShowLocation:  params.FormValueBoolDefault(ctx, "showLocation", false),
		EmailNotify:   params.FormValueBoolDefault(ctx, "emailNotify", false),
	}

	return form
}
