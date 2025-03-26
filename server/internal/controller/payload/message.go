package payload

import (
	"bbs-go/internal/locale"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/msg"

	"github.com/tidwall/gjson"
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
	detailUrl := getMessageDetailUrl(msg)
	resp := &MessageResponse{
		Id:           msg.Id,
		From:         from,
		UserId:       msg.UserId,
		Title:        msg.Title,
		Content:      msg.Content,
		QuoteContent: msg.QuoteContent,
		Type:         msg.Type,
		DetailUrl:    detailUrl,
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

// getMessageDetailUrl 查看消息详情链接地址
func getMessageDetailUrl(t *model.Message) string {
	msgType := msg.Type(t.Type)
	if msgType == msg.TypeTopicComment || msgType == msg.TypeArticleComment {
		entityType := gjson.Get(t.ExtraData, "entityType")
		entityId := gjson.Get(t.ExtraData, "entityId")
		if entityType.String() == constants.EntityArticle {
			return bbsurls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			return bbsurls.TopicUrl(entityId.Int())
		}
	} else if msgType == msg.TypeCommentReply {
		entityType := gjson.Get(t.ExtraData, "rootEntityType")
		entityId := gjson.Get(t.ExtraData, "rootEntityId")

		if entityType.String() == constants.EntityArticle {
			return bbsurls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			return bbsurls.TopicUrl(entityId.Int())
		}
	} else if msgType == msg.TypeTopicLike ||
		msgType == msg.TypeTopicFavorite ||
		msgType == msg.TypeTopicRecommend {
		topicId := gjson.Get(t.ExtraData, "topicId")
		if topicId.Exists() && topicId.Int() > 0 {
			return bbsurls.TopicUrl(topicId.Int())
		}
	}
	return bbsurls.AbsUrl("/user/messages")
}
