package payload

import (
	"bbs-go/internal/locale"
	"bbs-go/internal/model"
)

type ActivityResponse struct {
	Id         int64     `json:"id"`
	From       *UserInfo `json:"from"`   // message sender
	UserId     int64     `json:"userId"` // message receiver id
	Title      string    `json:"title"`
	Type       int       `json:"type"`
	DetailUrl  string    `json:"detailUrl"` // message detail url
	Status     int       `json:"status"`
	CreateTime int64     `json:"createTime"`
}

func BuildActivity(msg *model.Message) *ActivityResponse {
	if msg == nil {
		return nil
	}

	from := BuildUserInfoDefaultIfNull(msg.FromID)
	if msg.FromID <= 0 {
		from.Nickname = locale.T("system.user")
	}
	resp := &ActivityResponse{
		Id:         msg.ID,
		From:       from,
		UserId:     msg.UserID,
		Title:      msg.Title,
		Type:       msg.Type,
		DetailUrl:  getDetailUrl(msg.DetailUrl),
		Status:     msg.Status,
		CreateTime: msg.CreateTime,
	}
	return resp
}

// BuildActivities render activity list
func BuildActivities(messages []model.Message) []ActivityResponse {
	if len(messages) == 0 {
		return nil
	}
	var responses []ActivityResponse
	for _, message := range messages {
		responses = append(responses, *BuildActivity(&message))
	}
	return responses
}
