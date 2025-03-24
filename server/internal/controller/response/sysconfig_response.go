package response

import (
	"bbs-go/common/numbers"
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
)

func BuildSysConfigResponse(sysConfigs []model.SysConfig) *SysConfigResponse {
	var (
		siteTitle                  = service.SysConfigService.GetSiteTitle()
		siteDescription            = service.SysConfigService.GetSiteDescription()
		siteKeywords               = service.SysConfigService.GetSiteKeywords()
		siteNotification           = service.SysConfigService.GetSiteNofitication()
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
