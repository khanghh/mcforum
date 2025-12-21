package model

import (
	"bbs-go/common/strs"
	"bbs-go/internal/model/constants"
)

// 标签
type Tag struct {
	Model
	Name        string `gorm:"size:32;unique;not null" json:"name" form:"name"`
	Description string `gorm:"size:1024" json:"description" form:"description"`
	Status      int    `gorm:"type:int(11);index:idx_tag_status;not null" json:"status" form:"status"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`
}

// 文章
type Article struct {
	Model
	CategoryId   int64  `gorm:"default:0;index:idx_category_id" json:"categoryId" form:"categoryId"` // 分类ID
	UserId       int64  `gorm:"index:idx_article_user_id" json:"userId" form:"userId"`               // 所属用户编号
	Title        string `gorm:"size:128;not null;" json:"title" form:"title"`                        // 标题
	Summary      string `gorm:"type:text" json:"summary" form:"summary"`                             // 摘要
	Content      string `gorm:"type:longtext;not null;" json:"content" form:"content"`               // 内容
	ContentType  string `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`     // 内容类型：markdown、html
	Cover        string `gorm:"type:text;" json:"cover" form:"cover"`                                // 封面图
	Status       int    `gorm:"type:int(11);index:idx_article_status" json:"status" form:"status"`   // 状态
	SourceUrl    string `gorm:"type:text" json:"sourceUrl" form:"sourceUrl"`                         // 原文链接
	ViewCount    int64  `gorm:"not null;" json:"viewCount" form:"viewCount"`                         // 查看数量
	CommentCount int64  `gorm:"default:0" json:"commentCount" form:"commentCount"`                   // 评论数量
	LikeCount    int64  `gorm:"default:0" json:"likeCount" form:"likeCount"`                         // 点赞数量
	CreateTime   int64  `json:"createTime" form:"createTime"`                                        // 创建时间
	UpdateTime   int64  `json:"updateTime" form:"updateTime"`                                        // 更新时间
}

// 文章标签
type ArticleTag struct {
	Model
	ArticleId  int64 `gorm:"not null;index:idx_article_id;" json:"articleId" form:"articleId"`  // 文章编号
	TagId      int64 `gorm:"not null;index:idx_article_tag_tag_id;" json:"tagId" form:"tagId"`  // 标签编号
	Status     int64 `gorm:"not null;index:idx_article_tag_status" json:"status" form:"status"` // 状态：正常、删除
	CreateTime int64 `json:"createTime" form:"createTime"`                                      // 创建时间
}

// 评论
type Comment struct {
	Model
	UserId       int64  `gorm:"index:idx_comment_user_id;not null" json:"userId" form:"userId"`       // 用户编号
	TopicId      int64  `gorm:"index:idx_comment_topic_id;not null" json:"topicId" form:"topicId"`    // 被评论实体类型
	ParentId     int64  `gorm:"index:idx_comment_parent_id;not null" json:"parentId" form:"parentId"` // 被评论实体编号
	QuoteId      int64  `gorm:"not null" json:"quoteId" form:"quoteId"`                               // 引用的评论编号
	Content      string `gorm:"type:text;not null" json:"content" form:"content"`                     // 内容
	ImageList    string `gorm:"type:longtext" json:"imageList" form:"imageList"`                      // 图片
	ContentType  string `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`      // 内容类型：markdown、html
	LikeCount    int64  `gorm:"not null;default:0" json:"likeCount" form:"likeCount"`                 // 点赞数量
	CommentCount int64  `gorm:"not null;default:0" json:"commentCount" form:"commentCount"`           // 评论数量
	UserAgent    string `gorm:"size:1024" json:"userAgent" form:"userAgent"`                          // UserAgent
	Ip           string `gorm:"size:128" json:"ip" form:"ip"`                                         // IP
	IpLocation   string `gorm:"size:64" json:"ipLocation" form:"ipLocation"`                          // IP属地
	Status       int    `gorm:"type:int(11);index:idx_comment_status" json:"status" form:"status"`    // 状态：0：待审核、1：审核通过、2：审核失败、3：已发布
	CreateTime   int64  `json:"createTime" form:"createTime"`                                         // 创建时间
}

// 收藏
type Favorite struct {
	Model
	UserId     int64  `gorm:"index:idx_favorite_user_id;not null" json:"userId" form:"userId"`                     // 用户编号
	EntityType string `gorm:"index:idx_favorite_entity_type;size:32;not null" json:"entityType" form:"entityType"` // 收藏实体类型
	EntityId   int64  `gorm:"index:idx_favorite_entity_id;not null" json:"entityId" form:"entityId"`               // 收藏实体编号
	CreateTime int64  `json:"createTime" form:"createTime"`                                                        // 创建时间
}

// Forum 话题节点
type Forum struct {
	Model
	Name        string `gorm:"size:32;unique" json:"name" form:"name"`                     // 名称
	Slug        string `gorm:"size:32;unique" json:"slug" form:"slug"`                     // 名称
	Description string `gorm:"size:1024" json:"description" form:"description"`            // 描述
	Logo        string `gorm:"size:1024" json:"logo" form:"logo"`                          // 图标
	SortNo      int    `gorm:"type:int(11);index:idx_sort_no" json:"sortNo" form:"sortNo"` // 排序编号
	Status      int    `gorm:"type:int(11);not null" json:"status" form:"status"`          // 状态
	CreateTime  int64  `json:"createTime" form:"createTime"`                               // 创建时间
}

// 话题节点
type Topic struct {
	Model
	Slug              string              `gorm:"size:128;not null:default:'untitled'" json:"slug" form:"slug"`                            // 名称
	Type              constants.TopicType `gorm:"type:int(11);not null:default:0" json:"type" form:"type"`                                 // 类型
	ForumId           int64               `gorm:"not null;index:idx_topic_forum_id;" json:"forumId" form:"forumId"`                        // 节点编号
	UserId            int64               `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"`                           // 用户
	Title             string              `gorm:"size:128" json:"title" form:"title"`                                                      // 标题
	Content           string              `gorm:"type:longtext" json:"content" form:"content"`                                             // 内容
	ImageList         string              `gorm:"type:longtext" json:"imageList" form:"imageList"`                                         // 图片
	HideContent       string              `gorm:"type:longtext" json:"hideContent" form:"hideContent"`                                     // 回复可见内容
	Recommended       bool                `gorm:"not null" json:"recommended" form:"recommended"`                                          // 是否推荐
	RecommendedTime   int64               `gorm:"not null;index:idx_topic_recommended_time" json:"recommendedTime" form:"recommendedTime"` // 推荐时间
	Pinned            bool                `gorm:"not null" json:"pinned" form:"pinned"`                                                    // 置顶
	PinnedTime        int64               `gorm:"not null;index:idx_topic_pinned_time" json:"pinnedTime" form:"pinnedTime"`                // 置顶时间
	ViewCount         int64               `gorm:"not null" json:"viewCount" form:"viewCount"`                                              // 查看数量
	CommentCount      int64               `gorm:"not null" json:"commentCount" form:"commentCount"`                                        // 跟帖数量
	LikeCount         int64               `gorm:"not null" json:"likeCount" form:"likeCount"`                                              // 点赞数量
	Status            int                 `gorm:"type:int(11);index:idx_topic_status;" json:"status" form:"status"`                        // 状态：0：正常、1：删除
	LastCommentTime   int64               `gorm:"index:idx_topic_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"`         // 最后回复时间
	LastCommentUserId int64               `json:"lastCommentUserId" form:"lastCommentUserId"`                                              // 最后回复用户
	UserAgent         string              `gorm:"size:1024" json:"userAgent" form:"userAgent"`                                             // UserAgent
	Ip                string              `gorm:"size:128" json:"ip" form:"ip"`                                                            // IP
	IpLocation        string              `gorm:"size:64" json:"ipLocation" form:"ipLocation"`                                             // IP属地
	CreateTime        int64               `gorm:"index:idx_topic_create_time" json:"createTime" form:"createTime"`                         // 创建时间
	ExtraData         string              `gorm:"type:text" json:"extraData" form:"extraData"`                                             // 扩展数据
}

// GetTitle 获取帖子的标题
func (t *Topic) GetTitle() string {
	if t.Type == constants.TopicTypeTweet {
		if strs.IsNotBlank(t.Content) {
			return t.Content
		} else {
			return "分享图片"
		}
	} else {
		return t.Title
	}
}

// 主题标签
type TopicTag struct {
	Model
	TopicId           int64  `gorm:"not null;index:idx_topic_tag(topic_id,tag);" json:"topicId" form:"topicId"`
	Tag               string `gorm:"not null" json:"tag" form:"tag"` // 标签编号
	LastCommentTime   int64  `gorm:"index:idx_topic_last_comment_time" json:"lastCommentTime"`
	LastCommentUserId int64  `json:"lastCommentUserId"`            // 最后回复用户
	CreateTime        int64  `json:"createTime" form:"createTime"` // 创建时间
}

// 消息
type Message struct {
	Model
	FromId       int64  `gorm:"not null" json:"fromId" form:"fromId"`                            // 消息发送人
	UserId       int64  `gorm:"not null;index:idx_message_user_id;" json:"userId" form:"userId"` // 用户编号(消息接收人)
	Title        string `gorm:"size:1024" json:"title" form:"title"`                             // 消息标题
	Content      string `gorm:"type:text;not null" json:"content" form:"content"`                // 消息内容
	QuoteContent string `gorm:"type:text" json:"quoteContent" form:"quoteContent"`               // 引用内容
	Type         int    `gorm:"type:int(11);not null" json:"type" form:"type"`                   // 消息类型
	DetailUrl    string `gorm:"type:text;not null" json:"detailUrl" form:"detailUrl"`            // Detail url
	ExtraData    string `gorm:"type:text" json:"extraData" form:"extraData"`                     // 扩展数据
	Status       int    `gorm:"type:int(11);not null" json:"status" form:"status"`               // 状态：0：未读、1：已读
	CreateTime   int64  `json:"createTime" form:"createTime"`                                    // 创建时间
}

// 系统配置
type SysConfig struct {
	Model
	Key         string `gorm:"not null;size:128;unique" json:"key" form:"key"` // 配置key
	Value       string `gorm:"type:text" json:"value" form:"value"`            // 配置值
	Name        string `gorm:"not null;size:64" json:"name" form:"name"`       // 配置名称
	Description string `gorm:"size:128" json:"description" form:"description"` // 配置描述
	CreateTime  int64  `gorm:"not null" json:"createTime" form:"createTime"`   // 创建时间
	UpdateTime  int64  `gorm:"not null" json:"updateTime" form:"updateTime"`   // 更新时间
}

// 友链
type Link struct {
	Model
	Url        string `gorm:"not null;type:text" json:"url" form:"url"`          // 链接
	Title      string `gorm:"not null;size:128" json:"title" form:"title"`       // 标题
	Summary    string `gorm:"size:1024" json:"summary" form:"summary"`           // 站点描述
	Logo       string `gorm:"type:text" json:"logo" form:"logo"`                 // LOGO
	Status     int    `gorm:"type:int(11);not null" json:"status" form:"status"` // 状态
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`      // 创建时间
}

// 操作日志
type OperateLog struct {
	Model
	UserId      int64  `gorm:"not null;index:idx_operate_log_user_id" json:"userId" form:"userId"`  // 用户编号
	OpType      string `gorm:"not null;index:idx_op_type;size:32" json:"opType" form:"opType"`      // 操作类型
	DataType    string `gorm:"not null;index:idx_operate_log_data" json:"dataType" form:"dataType"` // 数据类型
	DataId      int64  `gorm:"not null;index:idx_operate_log_data" json:"dataId" form:"dataId" `    // 数据编号
	Description string `gorm:"not null;size:1024" json:"description" form:"description"`            // 描述
	Ip          string `gorm:"size:128" json:"ip" form:"ip"`                                        // ip地址
	UserAgent   string `gorm:"type:text" json:"userAgent" form:"userAgent"`                         // UserAgent
	Referer     string `gorm:"type:text" json:"referer" form:"referer"`                             // Referer
	CreateTime  int64  `json:"createTime" form:"createTime"`                                        // 创建时间
}

// 邮箱验证码
type EmailCode struct {
	Model
	UserId     int64  `gorm:"not null;index:idx_user_score_log_user_id" json:"userId" form:"userId"` // 用户编号
	Email      string `gorm:"not null;size:128" json:"email" form:"email"`                           // 邮箱
	Code       string `gorm:"not null;size:8" json:"code" form:"code"`                               // 验证码
	Token      string `gorm:"not null;size:32;unique" json:"token" form:"token"`                     // 验证码token
	Title      string `gorm:"size:1024" json:"title" form:"title"`                                   // 标题
	Content    string `gorm:"type:text" json:"content" form:"content"`                               // 内容
	Used       bool   `gorm:"not null" json:"used" form:"used"`                                      // 是否使用
	CreateTime int64  `json:"createTime" form:"createTime"`                                          // 创建时间
}

// UserFollow 粉丝关注
type UserFollow struct {
	Model
	UserId     int64 `gorm:"not null;uniqueIndex:idx_user_id" json:"userId"`           // 用户编号
	OtherId    int64 `gorm:"not null;uniqueIndex:idx_user_id" json:"otherId"`          // 对方的ID（被关注用户编号）
	Status     int   `gorm:"type:int(11);not null" json:"status"`                      // 关注状态
	CreateTime int64 `gorm:"type:bigint;not null" json:"createTime" form:"createTime"` // 创建时间
}

// UserFeed 用户信息流
type UserFeed struct {
	Model
	UserId     int64  `gorm:"not null;uniqueIndex:idx_data;index:idx_user_id;index:idx_search" json:"userId"`                   // 用户编号
	DataId     int64  `gorm:"not null;uniqueIndex:idx_data;index:idx_data_id" json:"dataId" form:"dataId"`                      // 数据ID
	DataType   string `gorm:"not null;uniqueIndex:idx_data;index:idx_data_id;index:idx_search" json:"dataType" form:"dataType"` // 数据类型
	AuthorId   int64  `gorm:"not null;index:idx_user_id" json:"authorId" form:"authorId"`                                       // 作者编号
	CreateTime int64  `gorm:"type:bigint;not null;index:idx_search" json:"createTime" form:"createTime"`                        // 数据的创建时间
}

// ForbiddenWord 违禁词
type ForbiddenWord struct {
	Model
	Type       string `gorm:"size:16" json:"type" form:"type"`       // 类型：word/regex
	Word       string `gorm:"size:128" json:"word" form:"word"`      // 违禁词
	Remark     string `gorm:"size:1024" json:"remark" form:"remark"` // 备注
	CreateTime int64  `json:"createTime" form:"createTime"`          // 举报时间
}
