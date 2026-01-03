package model

import (
	"bbs-go/common/strs"
	"bbs-go/internal/model/constants"
)

// tag
type Tag struct {
	Model
	Name        string `gorm:"size:32;unique;not null" json:"name" form:"name"`
	Description string `gorm:"size:1024" json:"description" form:"description"`
	Status      int    `gorm:"type:int(11);index:idx_tag_status;not null" json:"status" form:"status"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`
}

// Article
type Article struct {
	Model
	CategoryId   int64  `gorm:"default:0;index:idx_category_id" json:"categoryId" form:"categoryId"` // Category ID
	UserId       int64  `gorm:"index:idx_article_user_id" json:"userId" form:"userId"`               // User ID
	Title        string `gorm:"size:128;not null;" json:"title" form:"title"`                        // Title
	Summary      string `gorm:"type:text" json:"summary" form:"summary"`                             // Summary
	Content      string `gorm:"type:longtext;not null;" json:"content" form:"content"`               // Content
	ContentType  string `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`     // Content type: markdown, html
	Cover        string `gorm:"type:text;" json:"cover" form:"cover"`                                // Cover image
	Status       int    `gorm:"type:int(11);index:idx_article_status" json:"status" form:"status"`   // Status
	SourceUrl    string `gorm:"type:text" json:"sourceUrl" form:"sourceUrl"`                         // Source URL
	ViewCount    int64  `gorm:"not null;" json:"viewCount" form:"viewCount"`                         // View count
	CommentCount int64  `gorm:"default:0" json:"commentCount" form:"commentCount"`                   // Comment count
	LikeCount    int64  `gorm:"default:0" json:"likeCount" form:"likeCount"`                         // Like count
	CreateTime   int64  `json:"createTime" form:"createTime"`                                        // Create time
	UpdateTime   int64  `json:"updateTime" form:"updateTime"`                                        // Update time
}

// Article tag
type ArticleTag struct {
	Model
	ArticleId  int64 `gorm:"not null;index:idx_article_id;" json:"articleId" form:"articleId"`  // Article ID
	TagId      int64 `gorm:"not null;index:idx_article_tag_tag_id;" json:"tagId" form:"tagId"`  // Tag ID
	Status     int64 `gorm:"not null;index:idx_article_tag_status" json:"status" form:"status"` // Status: normal, deleted
	CreateTime int64 `json:"createTime" form:"createTime"`                                      // Create time
}

// Comment
type Comment struct {
	Model
	UserID       int64      `gorm:"index:idx_comment_user_id;not null" json:"userId" form:"userId"`       // User ID
	TopicID      int64      `gorm:"index:idx_comment_topic_id;not null" json:"topicId" form:"topicId"`    // Topic ID
	ParentID     int64      `gorm:"index:idx_comment_parent_id;not null" json:"parentId" form:"parentId"` // Parent ID
	QuoteID      int64      `gorm:"not null" json:"quoteId" form:"quoteId"`                               // Quoted comment ID
	Content      string     `gorm:"type:text;not null" json:"content" form:"content"`                     // Content
	Images       []ImageDTO `gorm:"serializer:json" json:"images" form:"images"`                          // Image list
	ContentType  string     `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`      // Content type: markdown, html
	LikeCount    int64      `gorm:"not null;default:0" json:"likeCount" form:"likeCount"`                 // Like count
	CommentCount int64      `gorm:"not null;default:0" json:"commentCount" form:"commentCount"`           // Comment count
	UserAgent    string     `gorm:"size:1024" json:"userAgent" form:"userAgent"`                          // UserAgent
	IP           string     `gorm:"size:128" json:"ip" form:"ip"`                                         // IP
	IPLocation   string     `gorm:"size:64" json:"ipLocation" form:"ipLocation"`                          // IP location
	Status       int        `gorm:"type:int(11);index:idx_comment_status" json:"status" form:"status"`    // Status: 0: pending review, 1: approved, 2: rejected, 3: published
	CreateTime   int64      `json:"createTime" form:"createTime"`                                         // Create time
}

// Favorite
type Favorite struct {
	Model
	UserID     int64  `gorm:"index:idx_favorite_user_id;not null" json:"userId" form:"userId"`                     // User ID
	EntityType string `gorm:"index:idx_favorite_entity_type;size:32;not null" json:"entityType" form:"entityType"` // Entity type
	EntityID   int64  `gorm:"index:idx_favorite_entity_id;not null" json:"entityId" form:"entityId"`               // Entity ID
	CreateTime int64  `json:"createTime" form:"createTime"`                                                        // Create time
}

// Forum topic node
type Forum struct {
	Model
	Name        string `gorm:"size:32;unique" json:"name" form:"name"`                     // Name
	Slug        string `gorm:"size:32;unique" json:"slug" form:"slug"`                     // Slug
	Description string `gorm:"size:1024" json:"description" form:"description"`            // Description
	Logo        string `gorm:"size:1024" json:"logo" form:"logo"`                          // Logo
	SortNo      int    `gorm:"type:int(11);index:idx_sort_no" json:"sortNo" form:"sortNo"` // Sort number
	Status      int    `gorm:"type:int(11);not null" json:"status" form:"status"`          // Status
	CreateTime  int64  `json:"createTime" form:"createTime"`                               // Create time
}

// Topic
type Topic struct {
	Model
	Slug              string              `gorm:"size:128;not null:default:'untitled'" json:"slug" form:"slug"`                            // Name
	Type              constants.TopicType `gorm:"type:int(11);not null:default:0" json:"type" form:"type"`                                 // Type
	ForumId           int64               `gorm:"not null;index:idx_topic_forum_id;" json:"forumId" form:"forumId"`                        // Forum ID
	UserID            int64               `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"`                           // User ID
	Title             string              `gorm:"size:256" json:"title" form:"title"`                                                      // Title
	Content           string              `gorm:"type:longtext" json:"content" form:"content"`                                             // Content
	Images            []ImageDTO          `gorm:"serializer:json" json:"images" form:"images"`                                             // Image list
	HideContent       string              `gorm:"type:longtext" json:"hideContent" form:"hideContent"`                                     // Hidden content (visible after reply)
	Recommended       bool                `gorm:"not null" json:"recommended" form:"recommended"`                                          // Recommended
	RecommendedTime   int64               `gorm:"not null;index:idx_topic_recommended_time" json:"recommendedTime" form:"recommendedTime"` // Recommended time
	Pinned            bool                `gorm:"not null" json:"pinned" form:"pinned"`                                                    // Pinned
	PinnedTime        int64               `gorm:"not null;index:idx_topic_pinned_time" json:"pinnedTime" form:"pinnedTime"`                // Pinned time
	ViewCount         int64               `gorm:"not null" json:"viewCount" form:"viewCount"`                                              // View count
	CommentCount      int64               `gorm:"not null" json:"commentCount" form:"commentCount"`                                        // Comment count
	LikeCount         int64               `gorm:"not null" json:"likeCount" form:"likeCount"`                                              // Like count
	Status            int                 `gorm:"type:int(11);index:idx_topic_status;" json:"status" form:"status"`                        // Status: 0: normal, 1: deleted
	LastCommentTime   int64               `gorm:"index:idx_topic_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"`         // Last comment time
	LastCommentUserId int64               `json:"lastCommentUserId" form:"lastCommentUserId"`                                              // Last comment user ID
	UserAgent         string              `gorm:"size:1024" json:"userAgent" form:"userAgent"`                                             // UserAgent
	IP                string              `gorm:"size:128" json:"ip" form:"ip"`                                                            // IP
	IPLocation        string              `gorm:"size:64" json:"ipLocation" form:"ipLocation"`                                             // IP location
	CreateTime        int64               `gorm:"index:idx_topic_create_time" json:"createTime" form:"createTime"`                         // Create time
	ExtraData         string              `gorm:"type:text" json:"extraData" form:"extraData"`                                             // Extra data
}

// GetTitle gets the topic's title
func (t *Topic) GetTitle() string {
	if t.Type == constants.TopicTypeTweet {
		if strs.IsNotBlank(t.Content) {
			return t.Content
		} else {
			return "Share image"
		}
	} else {
		return t.Title
	}
}

// Topic tag
type TopicTag struct {
	Model
	TopicId           int64  `gorm:"not null;index:idx_topic_tag(topic_id,tag);" json:"topicId" form:"topicId"`
	Tag               string `gorm:"not null" json:"tag" form:"tag"` // Tag
	LastCommentTime   int64  `gorm:"index:idx_topic_last_comment_time" json:"lastCommentTime"`
	LastCommentUserId int64  `json:"lastCommentUserId"`            // Last comment user ID
	CreateTime        int64  `json:"createTime" form:"createTime"` // Create time
}

// Message
type Message struct {
	Model
	FromID       int64  `gorm:"not null" json:"fromId" form:"fromId"`                            // Sender ID
	UserID       int64  `gorm:"not null;index:idx_message_user_id;" json:"userId" form:"userId"` // User ID (recipient)
	Title        string `gorm:"size:1024" json:"title" form:"title"`                             // Title
	Content      string `gorm:"type:text;not null" json:"content" form:"content"`                // Content
	QuoteContent string `gorm:"type:text" json:"quoteContent" form:"quoteContent"`               // Quoted content
	Type         int    `gorm:"type:int(11);not null" json:"type" form:"type"`                   // Type
	DetailUrl    string `gorm:"type:text;not null" json:"detailUrl" form:"detailUrl"`            // Detail URL
	ExtraData    string `gorm:"type:text" json:"extraData" form:"extraData"`                     // Extra data
	Status       int    `gorm:"type:int(11);not null" json:"status" form:"status"`               // Status: 0: unread, 1: read
	CreateTime   int64  `json:"createTime" form:"createTime"`                                    // Create time
}

// System config
type SysConfig struct {
	Model
	Key         string `gorm:"not null;size:128;unique" json:"key" form:"key"` // Config key
	Value       string `gorm:"type:text" json:"value" form:"value"`            // Config value
	Name        string `gorm:"not null;size:64" json:"name" form:"name"`       // Config name
	Description string `gorm:"size:128" json:"description" form:"description"` // Config description
	CreateTime  int64  `gorm:"not null" json:"createTime" form:"createTime"`   // Create time
	UpdateTime  int64  `gorm:"not null" json:"updateTime" form:"updateTime"`   // Update time
}

// Link
type Link struct {
	Model
	URL        string `gorm:"not null;type:text" json:"url" form:"url"`          // URL
	Title      string `gorm:"not null;size:128" json:"title" form:"title"`       // Title
	Summary    string `gorm:"size:1024" json:"summary" form:"summary"`           // Site description
	Logo       string `gorm:"type:text" json:"logo" form:"logo"`                 // Logo
	Status     int    `gorm:"type:int(11);not null" json:"status" form:"status"` // Status
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`      // Create time
}

// Operation log
type OperateLog struct {
	Model
	UserID      int64  `gorm:"not null;index:idx_operate_log_user_id" json:"userId" form:"userId"`  // User ID
	OpType      string `gorm:"not null;index:idx_op_type;size:32" json:"opType" form:"opType"`      // Operation type
	DataType    string `gorm:"not null;index:idx_operate_log_data" json:"dataType" form:"dataType"` // Data type
	DataID      int64  `gorm:"not null;index:idx_operate_log_data" json:"dataId" form:"dataId" `    // Data ID
	Description string `gorm:"not null;size:1024" json:"description" form:"description"`            // Description
	IP          string `gorm:"size:128" json:"ip" form:"ip"`                                        // IP address
	UserAgent   string `gorm:"type:text" json:"userAgent" form:"userAgent"`                         // UserAgent
	Referer     string `gorm:"type:text" json:"referer" form:"referer"`                             // Referer
	CreateTime  int64  `json:"createTime" form:"createTime"`                                        // Create time
}

// Email code
type EmailCode struct {
	Model
	UserID     int64  `gorm:"not null;index:idx_user_score_log_user_id" json:"userId" form:"userId"` // User ID
	Email      string `gorm:"not null;size:128" json:"email" form:"email"`                           // Email
	Code       string `gorm:"not null;size:8" json:"code" form:"code"`                               // Code
	Token      string `gorm:"not null;size:32;unique" json:"token" form:"token"`                     // Code token
	Title      string `gorm:"size:1024" json:"title" form:"title"`                                   // Title
	Content    string `gorm:"type:text" json:"content" form:"content"`                               // Content
	Used       bool   `gorm:"not null" json:"used" form:"used"`                                      // Used
	CreateTime int64  `json:"createTime" form:"createTime"`                                          // Create time
}

// UserFollow fans follow
type UserFollow struct {
	Model
	UserID     int64 `gorm:"not null;uniqueIndex:idx_user_id" json:"userId"`           // User ID
	OtherID    int64 `gorm:"not null;uniqueIndex:idx_user_id" json:"otherId"`          // Other ID (followed user ID)
	Status     int   `gorm:"type:int(11);not null" json:"status"`                      // Follow status
	CreateTime int64 `gorm:"type:bigint;not null" json:"createTime" form:"createTime"` // Create time
}

// UserFeed user feed
type UserFeed struct {
	Model
	UserID     int64  `gorm:"not null;uniqueIndex:idx_data;index:idx_user_id;index:idx_search" json:"userId"`                   // User ID
	DataID     int64  `gorm:"not null;uniqueIndex:idx_data;index:idx_data_id" json:"dataId" form:"dataId"`                      // Data ID
	DataType   string `gorm:"not null;uniqueIndex:idx_data;index:idx_data_id;index:idx_search" json:"dataType" form:"dataType"` // Data type
	AuthorID   int64  `gorm:"not null;index:idx_user_id" json:"authorId" form:"authorId"`                                       // Author ID
	CreateTime int64  `gorm:"type:bigint;not null;index:idx_search" json:"createTime" form:"createTime"`                        // Data create time
}

// ForbiddenWord forbidden word
type ForbiddenWord struct {
	Model
	Type       string `gorm:"size:16" json:"type" form:"type"`       // Type: word/regex
	Word       string `gorm:"size:128" json:"word" form:"word"`      // Forbidden word
	Remark     string `gorm:"size:1024" json:"remark" form:"remark"` // Remark
	CreateTime int64  `json:"createTime" form:"createTime"`          // Report time
}
