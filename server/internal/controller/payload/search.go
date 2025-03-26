package payload

import (
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/service"
)

type SearchTopicResponse struct {
	Id         int64          `json:"id"`
	User       *UserInfo      `json:"user"`
	Forum      *ForumResponse `json:"forum"`
	Tags       []string       `json:"tags"`
	Title      string         `json:"title"`
	Summary    string         `json:"summary"`
	CreateTime int64          `json:"createTime"`
}

func BuildSearchTopics(docs []search.TopicDocument) []SearchTopicResponse {
	var items []SearchTopicResponse
	for _, doc := range docs {
		items = append(items, BuildSearchTopic(doc))
	}
	return items
}

func BuildSearchTopic(doc search.TopicDocument) SearchTopicResponse {
	rsp := SearchTopicResponse{
		Id:         doc.Id,
		Title:      doc.Title,
		Summary:    doc.Content,
		CreateTime: doc.CreateTime,
		User:       BuildUserInfoDefaultIfNull(doc.UserId),
	}

	if doc.NodeId > 0 {
		node := service.ForumService.Get(doc.NodeId)
		rsp.Forum = BuildForum(node)
	}

	rsp.Tags = service.TopicService.GetTopicTags(doc.Id)
	return rsp
}
