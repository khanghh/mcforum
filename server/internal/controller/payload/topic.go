package payload

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"bbs-go/pkg/markdown"
	"fmt"
	"html"

	"bbs-go/common/arrays"
	"bbs-go/common/base62"
	"bbs-go/common/strs"
)

type TopicEditResponse struct {
	Id            int64               `json:"id"`
	Slug          string              `json:"slug"`
	Type          constants.TopicType `json:"type"`
	User          *UserInfo           `json:"user"`
	Forum         *ForumResponse      `json:"forum"`
	Tags          []string            `json:"tags"`
	Title         string              `json:"title"`
	Summary       string              `json:"summary"`
	Content       string              `json:"content"`
	HiddenContent string              `json:"hiddenContent"`
	ImageList     []ImageInfo         `json:"imageList"`
	ViewCount     int64               `json:"viewCount"`
	CreateTime    int64               `json:"createTime"`
	Status        int                 `json:"status"`
}

// topic list response entity
type TopicResponse struct {
	Id              int64               `json:"id"`
	Slug            string              `json:"slug"`
	Type            constants.TopicType `json:"type"`
	User            *UserInfo           `json:"user"`
	Forum           *ForumResponse      `json:"forum"`
	Tags            []string            `json:"tags"`
	Title           string              `json:"title"`
	Summary         string              `json:"summary"`
	Content         string              `json:"content"`
	ImageList       []ImageInfo         `json:"imageList"`
	LastCommentTime int64               `json:"lastCommentTime"`
	ViewCount       int64               `json:"viewCount"`
	CommentCount    int64               `json:"commentCount"`
	LikeCount       int64               `json:"likeCount"`
	Liked           bool                `json:"liked"`
	CreateTime      int64               `json:"createTime"`
	Recommended     bool                `json:"recommended"`
	RecommendedTime int64               `json:"recommendedTime"`
	Pinned          bool                `json:"pinned"`
	PinnedTime      int64               `json:"pinnedTime"`
	Status          int                 `json:"status"`
	Favorited       bool                `json:"favorited"`
	IpLocation      string              `json:"ipLocation"`
}

func BuildTopic(topic *model.Topic, currentUser *model.User) *TopicResponse {
	resp := _buildTopic(topic, true)
	if currentUser != nil {
		resp.Liked = service.UserLikeService.IsLiked(currentUser.ID, constants.EntityTopic, topic.ID)
		resp.Favorited = service.FavoriteService.IsFavorited(currentUser.ID, constants.EntityTopic, topic.ID)
	}
	return resp
}

func BuildSimpleTopic(topic *model.Topic) *TopicResponse {
	buildContent := topic.Type == constants.TopicTypeTweet // render content when tweet
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
			topicIds = append(topicIds, topic.ID)
		}
		likedTopicIds = service.UserLikeService.GetUserLikes(currentUser.ID, constants.EntityTopic, topicIds)
	}

	var responses []TopicResponse
	for _, topic := range topics {
		item := BuildSimpleTopic(&topic)
		item.Liked = arrays.Contains(likedTopicIds, topic.ID)
		responses = append(responses, *item)
	}
	return responses
}

func _buildTopic(topic *model.Topic, isBriefContent bool) *TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &TopicResponse{}

	rsp.Id = topic.ID
	rsp.Type = topic.Type
	rsp.Slug = fmt.Sprintf("%s.%s", topic.Slug, base62.Encode(topic.ID))
	rsp.Title = topic.Title
	rsp.User = BuildUserInfoDefaultIfNull(topic.UserID)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount
	rsp.Recommended = topic.Recommended
	rsp.RecommendedTime = topic.RecommendedTime
	rsp.Pinned = topic.Pinned
	rsp.PinnedTime = topic.PinnedTime
	rsp.Status = topic.Status
	rsp.IpLocation = topic.IPLocation

	// build content
	if isBriefContent {
		if topic.Type == constants.TopicTypeTopic {
			content := markdown.ToHTML(topic.Content)
			rsp.Content = handleHtmlContent(content)
		} else {
			rsp.Content = html.EscapeString(topic.Content)
		}
	} else {
		rsp.Summary = markdown.GetSummary(topic.Content, constants.SummaryLen)
	}

	if topic.Type == constants.TopicTypeTweet {
		if strs.IsBlank(topic.Content) {
			rsp.Content = "share image"
		} else {
			rsp.Content = html.EscapeString(topic.Content)
		}
		rsp.ImageList = BuildImageList(topic.ImageList)
	}

	if topic.ForumId > 0 {
		node := service.ForumService.Get(topic.ForumId)
		rsp.Forum = BuildForum(node)
	}

	rsp.Tags = service.TopicService.GetTopicTags(topic.ID)

	return rsp
}

func BuildTopicEdit(topic *model.Topic) *TopicEditResponse {
	if topic == nil {
		return nil
	}

	rsp := &TopicEditResponse{}

	rsp.Id = topic.ID
	rsp.Type = topic.Type
	rsp.Slug = fmt.Sprintf("%s.%s", topic.Slug, base62.Encode(topic.ID))
	rsp.Title = topic.Title
	rsp.User = BuildUserInfoDefaultIfNull(topic.UserID)
	rsp.CreateTime = topic.CreateTime
	rsp.Status = topic.Status
	// build content
	if topic.Type == constants.TopicTypeTopic {
		rsp.Content = topic.Content
		rsp.HiddenContent = topic.HideContent
		rsp.ImageList = BuildImageList(topic.ImageList)
	}

	if topic.ForumId > 0 {
		node := service.ForumService.Get(topic.ForumId)
		rsp.Forum = BuildForum(node)
	}

	rsp.Tags = service.TopicService.GetTopicTags(topic.ID)

	return rsp
}
