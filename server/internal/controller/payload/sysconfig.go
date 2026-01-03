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
// SysConfigResponse config response struct
type SysConfigResponse struct {
	SiteTitle                  string              `json:"siteTitle,omitempty"`
	SiteDescription            string              `json:"siteDescription,omitempty"`
	SiteKeywords               []string            `json:"siteKeywords,omitempty"`
	SiteNavs                   []model.ActionLink  `json:"siteNavs,omitempty"`
	SiteNotification           string              `json:"siteNotification,omitempty"`
	MenuItems                  []model.MenuItem    `json:"menuItems,omitempty"`
	RecommendTags              []string            `json:"recommendTags,omitempty"`
	UrlRedirect                bool                `json:"urlRedirect,omitempty"`
	ScoreConfig                model.ScoreConfig   `json:"scoreConfig,omitempty"`
	DefaultForumId             int64               `json:"defaultNodeId,omitempty"`
	ArticlePending             bool                `json:"articlePending,omitempty"`
	TopicCaptcha               bool                `json:"topicCaptcha,omitempty"`
	UserObserveSeconds         int                 `json:"userObserveSeconds,omitempty"`
	TokenExpireDays            int                 `json:"tokenExpireDays,omitempty"`
	LoginMethod                model.LoginMethod   `json:"loginMethod,omitempty"`
	CreateTopicEmailVerified   bool                `json:"createTopicEmailVerified,omitempty"`
	CreateArticleEmailVerified bool                `json:"createArticleEmailVerified,omitempty"`
	CreateCommentEmailVerified bool                `json:"createCommentEmailVerified,omitempty"`
	EnableHiddenContent        bool                `json:"enableHiddenContent,omitempty"`
	Modules                    model.ModulesConfig `json:"modules,omitempty"`
	EmailWhitelist             []string            `json:"emailWhitelist,omitempty"` // Email whitelist
}

func BuildSysConfigResponse(sysConfigs []model.SysConfig) *SysConfigResponse {
	var (
		siteTitle        = service.SysConfigService.GetSiteTitle()
		siteDescription  = service.SysConfigService.GetSiteDescription()
		siteKeywords     = service.SysConfigService.GetSiteKeywords()
		siteNotification = service.SysConfigService.GetSiteNofitication()
		menuItems        = service.SysConfigService.GetMenuItems()
		// recommendTags              = service.SysConfigService.GetRecommendTags()
		urlRedirect  = service.SysConfigService.IsEnabledUrlRedirectPage()
		pointsConfig = service.SysConfigService.GetPointConfig()
		// defaultForumIdStr          = service.SysConfigService.GetDefaultForumId()
		// topicCaptcha               = service.SysConfigService.IsEnabledTopicCaptcha()
		userObserveSecondsStr = cache.SysConfigCache.GetValue(constants.SysConfigUserObserveSeconds)
		siteNavs              = service.SysConfigService.GetSiteNavs()
		// loginMethod                = service.SysConfigService.GetLoginMethod()
		// createTopicEmailVerified   = service.SysConfigService.IsCreateTopicEmailVerified()
		// createArticleEmailVerified = service.SysConfigService.IsCreateArticleEmailVerified()
		// createCommentEmailVerified = service.SysConfigService.IsCreateCommentEmailVerified()
		// enableHideContent          = service.SysConfigService.IsEnableHideContent()
		// articlePending             = service.SysConfigService.IsArticlePending()
		tokenExpireDays = service.SysConfigService.GetTokenExpireDays()
		modules         = service.SysConfigService.GetModules()
		// emailWhiteList             = service.SysConfigService.GetEmailWhitelist()
	)

	var (
		userObserveSeconds = numbers.ToInt(userObserveSecondsStr)
	)

	return &SysConfigResponse{
		SiteTitle:          siteTitle,
		SiteDescription:    siteDescription,
		SiteKeywords:       siteKeywords,
		SiteNavs:           siteNavs,
		SiteNotification:   siteNotification,
		MenuItems:          menuItems,
		UrlRedirect:        urlRedirect,
		ScoreConfig:        pointsConfig,
		UserObserveSeconds: userObserveSeconds,
		TokenExpireDays:    tokenExpireDays,
		Modules:            modules,
	}
}
