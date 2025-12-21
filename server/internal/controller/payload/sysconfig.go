package payload

import (
	"bbs-go/common/numbers"
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
)

// SysConfigResponse
//
//	配置返回结构体
type SysConfigResponse struct {
	SiteTitle                  string              `json:"siteTitle"`
	SiteDescription            string              `json:"siteDescription"`
	SiteKeywords               []string            `json:"siteKeywords"`
	SiteNavs                   []model.ActionLink  `json:"siteNavs"`
	SiteNotification           string              `json:"siteNotification"`
	MenuItems                  []model.MenuItem    `json:"menuItems"`
	RecommendTags              []string            `json:"recommendTags"`
	UrlRedirect                bool                `json:"urlRedirect"`
	ScoreConfig                model.ScoreConfig   `json:"scoreConfig"`
	DefaultForumId             int64               `json:"defaultNodeId"`
	ArticlePending             bool                `json:"articlePending"`
	TopicCaptcha               bool                `json:"topicCaptcha"`
	UserObserveSeconds         int                 `json:"userObserveSeconds"`
	TokenExpireDays            int                 `json:"tokenExpireDays"`
	LoginMethod                model.LoginMethod   `json:"loginMethod"`
	CreateTopicEmailVerified   bool                `json:"createTopicEmailVerified"`
	CreateArticleEmailVerified bool                `json:"createArticleEmailVerified"`
	CreateCommentEmailVerified bool                `json:"createCommentEmailVerified"`
	EnableHideContent          bool                `json:"enableHideContent"`
	Modules                    model.ModulesConfig `json:"modules"`
	EmailWhitelist             []string            `json:"emailWhitelist"` // 邮箱白名单
}

func BuildSysConfigResponse(sysConfigs []model.SysConfig) *SysConfigResponse {
	var (
		siteTitle                  = service.SysConfigService.GetSiteTitle()
		siteDescription            = service.SysConfigService.GetSiteDescription()
		siteKeywords               = service.SysConfigService.GetSiteKeywords()
		siteNotification           = service.SysConfigService.GetSiteNofitication()
		menuItems                  = service.SysConfigService.GetMenuItems()
		recommendTags              = service.SysConfigService.GetRecommendTags()
		urlRedirect                = service.SysConfigService.IsEnabledUrlRedirectPage()
		pointsConfig               = service.SysConfigService.GetPointConfig()
		defaultForumIdStr          = service.SysConfigService.GetDefaultForumId()
		topicCaptcha               = service.SysConfigService.IsEnabledTopicCaptcha()
		userObserveSecondsStr      = cache.SysConfigCache.GetValue(constants.SysConfigUserObserveSeconds)
		siteNavs                   = service.SysConfigService.GetSiteNavs()
		loginMethod                = service.SysConfigService.GetLoginMethod()
		createTopicEmailVerified   = service.SysConfigService.IsCreateTopicEmailVerified()
		createArticleEmailVerified = service.SysConfigService.IsCreateArticleEmailVerified()
		createCommentEmailVerified = service.SysConfigService.IsCreateCommentEmailVerified()
		enableHideContent          = service.SysConfigService.IsEnableHideContent()
		articlePending             = service.SysConfigService.IsArticlePending()
		tokenExpireDays            = service.SysConfigService.GetTokenExpireDays()
		modules                    = service.SysConfigService.GetModules()
		emailWhiteList             = service.SysConfigService.GetEmailWhitelist()
	)

	var (
		userObserveSeconds = numbers.ToInt(userObserveSecondsStr)
	)

	return &SysConfigResponse{
		SiteTitle:                  siteTitle,
		SiteDescription:            siteDescription,
		SiteKeywords:               siteKeywords,
		SiteNavs:                   siteNavs,
		SiteNotification:           siteNotification,
		MenuItems:                  menuItems,
		RecommendTags:              recommendTags,
		UrlRedirect:                urlRedirect,
		ScoreConfig:                pointsConfig,
		DefaultForumId:             defaultForumIdStr,
		ArticlePending:             articlePending,
		TopicCaptcha:               topicCaptcha,
		UserObserveSeconds:         userObserveSeconds,
		TokenExpireDays:            tokenExpireDays,
		LoginMethod:                loginMethod,
		CreateTopicEmailVerified:   createTopicEmailVerified,
		CreateArticleEmailVerified: createArticleEmailVerified,
		CreateCommentEmailVerified: createCommentEmailVerified,
		EnableHideContent:          enableHideContent,
		Modules:                    modules,
		EmailWhitelist:             emailWhiteList,
	}
}
