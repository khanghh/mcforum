package service

import (
	"bbs-go/internal/model/constants"
	"errors"
	"log/slog"
	"strconv"

	"bbs-go/common/dates"
	"bbs-go/common/jsons"
	"bbs-go/common/strs"
	"bbs-go/sqls"
	"bbs-go/web/params"

	"github.com/spf13/cast"
	"github.com/tidwall/gjson"

	"gorm.io/gorm"

	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var SysConfigService = newSysConfigService()

func newSysConfigService() *sysConfigService {
	return &sysConfigService{}
}

type sysConfigService struct {
}

func (s *sysConfigService) Get(id int64) *model.SysConfig {
	return repository.SysConfigRepository.Get(sqls.DB(), id)
}

func (s *sysConfigService) Take(where ...interface{}) *model.SysConfig {
	return repository.SysConfigRepository.Take(sqls.DB(), where...)
}

func (s *sysConfigService) Find(cnd *sqls.Cnd) []model.SysConfig {
	return repository.SysConfigRepository.Find(sqls.DB(), cnd)
}

func (s *sysConfigService) FindOne(cnd *sqls.Cnd) *model.SysConfig {
	return repository.SysConfigRepository.FindOne(sqls.DB(), cnd)
}

func (s *sysConfigService) FindPageByParams(params *params.QueryParams) (list []model.SysConfig, paging *sqls.Paging) {
	return repository.SysConfigRepository.FindPageByParams(sqls.DB(), params)
}

func (s *sysConfigService) FindPageByCnd(cnd *sqls.Cnd) (list []model.SysConfig, paging *sqls.Paging) {
	return repository.SysConfigRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *sysConfigService) GetAll() []model.SysConfig {
	return repository.SysConfigRepository.Find(sqls.DB(), sqls.NewCnd().Asc("id"))
}

func (s *sysConfigService) SetAll(configStr string) error {
	json := gjson.Parse(configStr)
	configs, ok := json.Value().(map[string]interface{})
	if !ok {
		return errors.New("配置数据格式错误")
	}
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		for k := range configs {
			v := json.Get(k).String()
			if err := s.setSingle(tx, k, v, "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

// Set 设置配置，如果配置不存在，那么创建
func (s *sysConfigService) Set(key, value, name, description string) error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.setSingle(tx, key, value, name, description); err != nil {
			return err
		}
		return nil
	})
}

func (s *sysConfigService) setSingle(db *gorm.DB, key, value, name, description string) error {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}
	sysConfig := repository.SysConfigRepository.GetByKey(db, key)
	if sysConfig == nil {
		sysConfig = &model.SysConfig{
			CreateTime: dates.NowTimestamp(),
		}
	}
	sysConfig.Key = key
	sysConfig.Value = value
	sysConfig.UpdateTime = dates.NowTimestamp()

	if strs.IsNotBlank(name) {
		sysConfig.Name = name
	}
	if strs.IsNotBlank(description) {
		sysConfig.Description = description
	}

	var err error
	if sysConfig.Id > 0 {
		err = repository.SysConfigRepository.Update(db, sysConfig)
	} else {
		err = repository.SysConfigRepository.Create(db, sysConfig)
	}
	if err != nil {
		return err
	} else {
		cache.SysConfigCache.Invalidate(key)
		return nil
	}
}

func (s *sysConfigService) GetTokenExpireDays() int {
	tokenExpireDaysStr := cache.SysConfigCache.GetValue(constants.SysConfigTokenExpireDays)
	tokenExpireDays, err := strconv.Atoi(tokenExpireDaysStr)
	if err != nil {
		tokenExpireDays = constants.DefaultTokenExpireDays
	}
	if tokenExpireDays <= 0 {
		tokenExpireDays = constants.DefaultTokenExpireDays
	}
	return tokenExpireDays
}

func (s *sysConfigService) GetLoginMethod() model.LoginMethod {
	loginMethodStr := cache.SysConfigCache.GetValue(constants.SysConfigLoginMethod)

	useDefault := true
	var loginMethod model.LoginMethod
	if strs.IsNotBlank(loginMethodStr) {
		if err := jsons.Parse(loginMethodStr, &loginMethod); err != nil {
			slog.Warn("登录方式数据错误", err)
		} else {
			useDefault = false
		}
	}
	if useDefault {
		loginMethod = model.LoginMethod{
			Password: true,
			QQ:       true,
			Github:   true,
		}
	}
	return loginMethod
}

func (s *sysConfigService) IsCreateTopicEmailVerified() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigCreateTopicEmailVerified)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsCreateArticleEmailVerified() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigCreateArticleEmailVerified)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsCreateCommentEmailVerified() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigCreateCommentEmailVerified)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsEnableHideContent() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigEnableHideContent)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsEnabledTopicCaptcha() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigTopicCaptcha)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsEnabledUrlRedirectPage() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigUrlRedirect)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsArticlePending() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigArticlePending)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) GetSiteNavs() []model.ActionLink {
	siteNavs := cache.SysConfigCache.GetValue(constants.SysConfigSiteNavs)
	var siteNavsArr []model.ActionLink
	if strs.IsNotBlank(siteNavs) {
		if err := jsons.Parse(siteNavs, &siteNavsArr); err != nil {
			slog.Warn("站点导航数据错误", slog.Any("err", err))
		}
	}
	return siteNavsArr
}

func (s *sysConfigService) GetModules() model.ModulesConfig {
	str := cache.SysConfigCache.GetValue(constants.SysConfigModules)

	useDefault := true
	var modulesConfig model.ModulesConfig
	if strs.IsNotBlank(str) {
		if err := jsons.Parse(str, &modulesConfig); err != nil {
			slog.Warn("启用模块配置错误", err)
		} else {
			useDefault = false
		}
	}
	if useDefault {
		modulesConfig = model.ModulesConfig{
			Tweet:   true,
			Topic:   true,
			Article: true,
		}
	}
	return modulesConfig
}

// GetEmailWhitelist 邮箱白名单
func (s *sysConfigService) GetEmailWhitelist() []string {
	str := cache.SysConfigCache.GetValue(constants.SysConfigEmailWhitelist)
	var emailWhitelist []string
	_ = jsons.Parse(str, &emailWhitelist)
	return emailWhitelist
}

func (s *sysConfigService) GetPointConfig() model.ScoreConfig {
	str := cache.SysConfigCache.GetValue(constants.SysConfigScoreConfig)
	var pointConfig model.ScoreConfig
	if err := jsons.Parse(str, &pointConfig); err != nil {
		slog.Error("Failed to get system config", "key", constants.SysConfigScoreConfig, "error", err)
	}
	return pointConfig
}

func (s *sysConfigService) GetDefaultForumId() int64 {
	str := cache.SysConfigCache.GetValue(constants.SysConfigDefaultForumId)
	forumId, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		slog.Error("Failed to get system config", "key", constants.SysConfigScoreConfig, "error", err)
	}
	return forumId
}

func (s *sysConfigService) GetConfig() []model.SysConfig {
	return s.GetAll()
	// var (
	// 	siteTitle                  = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
	// 	siteDescription            = cache.SysConfigCache.GetValue(constants.SysConfigSiteDescription)
	// 	siteKeywords               = cache.SysConfigCache.GetValue(constants.SysConfigSiteKeywords)
	// 	siteNotification           = cache.SysConfigCache.GetValue(constants.SysConfigSiteNotification)
	// 	recommendTags              = cache.SysConfigCache.GetValue(constants.SysConfigRecommendTags)
	// 	urlRedirect                = cache.SysConfigCache.GetValue(constants.SysConfigUrlRedirect)
	// 	scoreConfigStr             = cache.SysConfigCache.GetValue(constants.SysConfigScoreConfig)
	// 	defaultNodeIdStr           = cache.SysConfigCache.GetValue(constants.SysConfigDefaultNodeId)
	// 	topicCaptcha               = cache.SysConfigCache.GetValue(constants.SysConfigTopicCaptcha)
	// 	userObserveSecondsStr      = cache.SysConfigCache.GetValue(constants.SysConfigUserObserveSeconds)
	// 	siteNavs                   = s.GetSiteNavs()
	// 	loginMethod                = s.GetLoginMethod()
	// 	createTopicEmailVerified   = s.IsCreateTopicEmailVerified()
	// 	createArticleEmailVerified = s.IsCreateArticleEmailVerified()
	// 	createCommentEmailVerified = s.IsCreateCommentEmailVerified()
	// 	enableHideContent          = s.IsEnableHideContent()
	// 	articlePending             = s.IsArticlePending()
	// )

	// var siteKeywordsArr []string
	// if strs.IsNotBlank(siteKeywords) {
	// 	if err := jsons.Parse(siteKeywords, &siteKeywordsArr); err != nil {
	// 		slog.Warn("站点关键词数据错误", slog.Any("err", err))
	// 	}
	// }

	// var recommendTagsArr []string
	// if strs.IsNotBlank(recommendTags) {
	// 	if err := jsons.Parse(recommendTags, &recommendTagsArr); err != nil {
	// 		slog.Warn("推荐标签数据错误", slog.Any("err", err))
	// 	}
	// }

	// var scoreConfig models.ScoreConfig
	// if strs.IsNotBlank(scoreConfigStr) {
	// 	if err := jsons.Parse(scoreConfigStr, &scoreConfig); err != nil {
	// 		slog.Warn("积分配置错误", slog.Any("err", err))
	// 	}
	// }

	// var (
	// 	defaultNodeId      = numbers.ToInt64(defaultNodeIdStr)
	// 	userObserveSeconds = numbers.ToInt(userObserveSecondsStr)
	// )

	// return &models.SysConfig{
	// 	SiteTitle:                  siteTitle,
	// 	SiteDescription:            siteDescription,
	// 	SiteKeywords:               siteKeywordsArr,
	// 	SiteNavs:                   siteNavs,
	// 	SiteNotification:           siteNotification,
	// 	RecommendTags:              recommendTagsArr,
	// 	UrlRedirect:                strings.ToLower(urlRedirect) == "true",
	// 	ScoreConfig:                scoreConfig,
	// 	DefaultNodeId:              defaultNodeId,
	// 	ArticlePending:             articlePending,
	// 	TopicCaptcha:               strings.ToLower(topicCaptcha) == "true",
	// 	UserObserveSeconds:         userObserveSeconds,
	// 	TokenExpireDays:            s.GetTokenExpireDays(),
	// 	LoginMethod:                loginMethod,
	// 	CreateTopicEmailVerified:   createTopicEmailVerified,
	// 	CreateArticleEmailVerified: createArticleEmailVerified,
	// 	CreateCommentEmailVerified: createCommentEmailVerified,
	// 	EnableHideContent:          enableHideContent,
	// 	Modules:                    s.GetModules(),
	// 	EmailWhitelist:             s.GetEmailWhitelist(),
	// }
}

func (s *sysConfigService) GetStr(key, def string) (value string) {
	value = cache.SysConfigCache.GetValue(key)
	if strs.IsBlank(value) {
		value = def
	}
	return
}

func (s *sysConfigService) GetInt(key string, def int) (value int) {
	str := cache.SysConfigCache.GetValue(key)
	if strs.IsBlank(str) {
		value = def
		return
	}
	var err error
	if value, err = cast.ToIntE(str); err != nil {
		value = def
		slog.Warn("Get int config error, use default value", slog.Any("default", def), slog.Any("key", key), slog.Any("value", str))
		return
	}
	return
}
