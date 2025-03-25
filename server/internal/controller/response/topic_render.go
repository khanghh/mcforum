package response

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/service"
	"fmt"
	"html"

	"bbs-go/common/arrays"
	"bbs-go/common/strs"
)

func BuildTopic(topic *model.Topic, currentUser *model.User) *TopicResponse {
	resp := _buildTopic(topic, true)
	if currentUser != nil {
		resp.Liked = service.UserLikeService.Exists(currentUser.Id, constants.EntityTopic, topic.Id)
		resp.Favorited = service.FavoriteService.IsFavorited(currentUser.Id, constants.EntityTopic, topic.Id)
	}
	return resp
}

func BuildSimpleTopic(topic *model.Topic) *TopicResponse {
	buildContent := topic.Type == constants.TopicTypeTweet // 动态时渲染内容
	return _buildTopic(topic, buildContent)
}

func BuildSimpleTopics(topics []model.Topic, currentUser *model.User) []TopicResponse {
	if len(topics) == 0 {
		return nil
	}

	var likedTopicIds []int64
	if currentUser != nil {
		var topicIds []int64
		for _, topic := range topics {
			topicIds = append(topicIds, topic.Id)
		}
		likedTopicIds = service.UserLikeService.IsLiked(currentUser.Id, constants.EntityTopic, topicIds)
	}

	var responses []TopicResponse
	for _, topic := range topics {
		item := BuildSimpleTopic(&topic)
		item.Liked = arrays.Contains(topic.Id, likedTopicIds)
		responses = append(responses, *item)
	}
	return responses
}

func _buildTopic(topic *model.Topic, isBriefContent bool) *TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &TopicResponse{}

	rsp.Id = topic.Id
	rsp.Type = topic.Type
	rsp.Slug = fmt.Sprintf("%s.%d", topic.Slug, topic.Id)
	rsp.Title = topic.Title
	rsp.User = BuildUserInfoDefaultIfNull(topic.UserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount
	rsp.Recommend = topic.Recommended
	rsp.RecommendTime = topic.RecommendedTime
	rsp.Pinned = topic.Pinned
	rsp.PinnedTime = topic.PinnedTime
	rsp.Status = topic.Status
	rsp.IpLocation = topic.IpLocation

	// 构建内容
	if isBriefContent {
		if topic.Type == constants.TopicTypeTopic {
			content := markdown.ToHTML(topic.Content)
			rsp.Content = handleHtmlContent(content)
		} else {
			rsp.Content = html.EscapeString(topic.Content)
		}
	} else {
		rsp.Summary = markdown.GetSummary(topic.Content, 128)
	}

	if topic.Type == constants.TopicTypeTweet {
		if strs.IsBlank(topic.Content) {
			rsp.Content = "分享图片"
		} else {
			rsp.Content = html.EscapeString(topic.Content)
		}
		rsp.ImageList = BuildImageList(topic.ImageList)
	}

	if topic.ForumId > 0 {
		node := service.ForumService.Get(topic.ForumId)
		rsp.Forum = BuildForum(node)
	}

	tags := service.TopicService.GetTopicTags(topic.Id)
	rsp.Tags = BuildTags(tags)

	return rsp
}
