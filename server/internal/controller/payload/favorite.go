package payload

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/markdown"
)

// 收藏返回数据
type FavoriteResponse struct {
	Id         int64     `json:"id"`
	EntityType string    `json:"entityType"`
	EntityId   int64     `json:"entityId"`
	Deleted    bool      `json:"deleted"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	User       *UserInfo `json:"user"`
	Url        string    `json:"url"`
	CreateTime int64     `json:"createTime"`
}

func BuildFavorite(favorite *model.Favorite) *FavoriteResponse {
	rsp := &FavoriteResponse{}
	rsp.Id = favorite.Id
	rsp.EntityType = favorite.EntityType
	rsp.CreateTime = favorite.CreateTime

	topic := service.TopicService.Get(favorite.EntityId)
	if topic == nil || topic.Status != constants.StatusOK {
		rsp.Deleted = true
	} else {
		rsp.Url = bbsurls.AbsTopicUrl(topic.Slug, topic.Id)
		rsp.User = BuildUserInfoDefaultIfNull(topic.UserId)
		rsp.Title = topic.Title
		rsp.Content = markdown.GetSummary(topic.Content, constants.SummaryLen)
	}
	return rsp
}

func BuildFavorites(favorites []model.Favorite) []FavoriteResponse {
	if len(favorites) == 0 {
		return nil
	}
	var responses []FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, *BuildFavorite(&favorite))
	}
	return responses
}
