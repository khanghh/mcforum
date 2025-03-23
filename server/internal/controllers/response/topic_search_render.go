package response

import (
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/services"
)

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
		node := services.ForumService.Get(doc.NodeId)
		rsp.Forum = BuildForum(node)
	}

	tags := services.TopicService.GetTopicTags(doc.Id)
	rsp.Tags = BuildTags(tags)
	return rsp
}
