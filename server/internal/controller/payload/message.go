package payload

import (
	"bbs-go/internal/locale"
	"bbs-go/internal/model"
	"bbs-go/pkg/bbsurls"
	"strings"
)

// 消息
type MessageResponse struct {
	Id           int64     `json:"id"`
	From         *UserInfo `json:"from"`    // 消息发送人
	UserId       int64     `json:"userId"`  // 消息接收人编号
	Title        string    `json:"title"`   // 标题
	Content      string    `json:"content"` // 消息内容
	QuoteContent string    `json:"quoteContent"`
	Type         int       `json:"type"`
	DetailUrl    string    `json:"detailUrl"` // 消息详情url
	ExtraData    string    `json:"extraData"`
	Status       int       `json:"status"`
	CreateTime   int64     `json:"createTime"`
}

func BuildMessage(msg *model.Message) *MessageResponse {
	if msg == nil {
		return nil
	}

	from := BuildUserInfoDefaultIfNull(msg.FromId)
	if msg.FromId <= 0 {
		from.Nickname = locale.T("system.user")
	}
	resp := &MessageResponse{
		Id:           msg.Id,
		From:         from,
		UserId:       msg.UserId,
		Title:        msg.Title,
		Content:      msg.Content,
		QuoteContent: msg.QuoteContent,
		Type:         msg.Type,
		DetailUrl:    getDetailUrl(msg.DetailUrl),
		ExtraData:    msg.ExtraData,
		Status:       msg.Status,
		CreateTime:   msg.CreateTime,
	}
	return resp
}

// BuildMessages 渲染消息列表
func BuildMessages(messages []model.Message) []MessageResponse {
	if len(messages) == 0 {
		return nil
	}
	var responses []MessageResponse
	for _, message := range messages {
		responses = append(responses, *BuildMessage(&message))
	}
	return responses
}

func getDetailUrl(urlOrPath string) string {
	if strings.HasPrefix(urlOrPath, "http://") || strings.HasPrefix(urlOrPath, "https://") {
		return urlOrPath
	}
	return bbsurls.AbsUrl(urlOrPath)
}
