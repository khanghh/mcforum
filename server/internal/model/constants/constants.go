package constants

const (
	DefaultTokenExpireDays       = 7   // Default validity period for user login token
	SummaryLen                   = 256 // Summary length
	CommentSummaryLen            = 128 // Summary length
	UploadMaxM                   = 10
	UploadMaxBytes         int64 = 1024 * 1024 * 1024 * UploadMaxM
	CookieTokenKey               = "sid"
	BioMaxLength                 = 200
	NicknameMaxLength            = 50
	StatusMessageMaxLength       = 100
	ForumTitleMaxLength          = 50
	TopicTitleMaxLength          = 128
	LocationMaxLength            = 100
	MaxPendingReviewTopics       = 3 // Maximum number of pending review topics per user
)

// System config
const (
	SysConfigSiteTitle                  = "siteTitle"                  // Site title
	SysConfigSiteDescription            = "siteDescription"            // Site description
	SysConfigSiteKeywords               = "siteKeywords"               // Site keywords
	SysConfigSiteNavs                   = "siteNavs"                   // Site navigation
	SysConfigSiteNotification           = "siteNotification"           // Site notification
	SysConfigRecommendTags              = "recommendTags"              // Recommend tags
	SysConfigUrlRedirect                = "urlRedirect"                // Whether to enable link redirect
	SysConfigScoreConfig                = "scoreConfig"                // Score config
	SysConfigDefaultForumId             = "defaultNodeId"              // Default forum for posting
	SysConfigArticlePending             = "articlePending"             // Whether to enable article review
	SysConfigTopicCaptcha               = "topicCaptcha"               // Whether to enable topic captcha
	SysConfigUserObserveSeconds         = "userObserveSeconds"         // New user observation period
	SysConfigTokenExpireDays            = "tokenExpireDays"            // Login token valid days
	SysConfigLoginMethod                = "loginMethod"                // Login method
	SysConfigEnableHideContent          = "enableHideContent"          // Enable hide content feature
	SysConfigCreateTopicEmailVerified   = "createTopicEmailVerified"   // Create topic requires email verification
	SysConfigCreateArticleEmailVerified = "createArticleEmailVerified" // Create article requires email verification
	SysConfigCreateCommentEmailVerified = "createCommentEmailVerified" // Create comment requires email verification
	SysConfigModules                    = "modules"                    // Modules
	SysConfigEmailWhitelist             = "emailWhitelist"             // Email whitelist
	SysConfigMenuItems                  = "menuItems"                  // Menu config
)

// EntityType
const (
	EntityArticle = "article"
	EntityTopic   = "topic"
	EntityComment = "comment"
	EntityUser    = "user"
	EntityCheckIn = "checkIn"
)

// User roles
const (
	RoleOwner     = "owner"     // Owner
	RoleAdmin     = "admin"     // Admin
	RoleModerator = "moderator" // Moderator
)

// Action types
const (
	OpTypeCreate          = "create"
	OpTypeDelete          = "delete"
	OpTypeUpdate          = "update"
	OpTypeForbidden       = "forbidden"
	OpTypeRemoveForbidden = "removeForbidden"
)

// status
const (
	StatusActive  = 0
	StatusDeleted = 1
	StatusReview  = 2
)

// user type
const (
	UserTypeNormal   = 0 // normal user
	UserTypeEmployee = 1 // employee user
)

// role type
const (
	RoleTypeSystem = 0 // system role
	RoleTypeCustom = 1 // custom role
)

// content type
const (
	ContentTypeHtml     = "html"
	ContentTypeMarkdown = "markdown"
	ContentTypeText     = "text"
)

// third-party account type
const (
	ThirdAccountTypeGithub = "github"
	ThirdAccountTypeOSC    = "osc"
	ThirdAccountTypeQQ     = "qq"
)

// score operation type
const (
	ScoreTypeIncr = 0 // score +
	ScoreTypeDecr = 1 // score -
)

type TopicType int

const (
	TopicTypeTopic TopicType = 0
	TopicTypeTweet TopicType = 1
)

type LoginMethod string

const (
	LoginMethodQQ       LoginMethod = "qq"
	LoginMethodGithub   LoginMethod = "github"
	LoginMethodPassword LoginMethod = "password"
)

const (
	FollowStatusNONE   = 0
	FollowStatusFollow = 1
	FollowStatusBoth   = 2
)

const (
	NodeIdNewest    int64 = 0
	NodeIdRecommend int64 = -1
	NodeIdFollow    int64 = -2
)
const (
	ForumWhatsNew    = "whats-new"
	ForumRecommended = "recommended"
	ForumFeed        = "feed"
)

const (
	GenderMale   = "Male"
	GenderFemale = "Female"
)

// module
const (
	ModuleTweet   = "tweet"
	ModuleTopic   = "topic"
	ModuleArticle = "article"
)

const (
	ForbiddenWordTypeWord  = "word"
	ForbiddenWordTypeRegex = "regex"
)
