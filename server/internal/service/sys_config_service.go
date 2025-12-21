package service

import (
	"bbs-go/internal/model/constants"
	"errors"
	"log/slog"
	"strconv"

	"bbs-go/common/dates"
	"bbs-go/common/jsons"
	"bbs-go/common/strs"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

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

func (s *sysConfigService) GetSiteTitle() string {
	return cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
}

func (s *sysConfigService) GetSiteDescription() string {
	return cache.SysConfigCache.GetValue(constants.SysConfigSiteDescription)
}

func (s *sysConfigService) GetSiteKeywords() []string {
	str := cache.SysConfigCache.GetValue(constants.SysConfigSiteKeywords)
	var keywords []string
	if strs.IsNotBlank(str) {
		if err := jsons.Parse(str, &keywords); err != nil {
			slog.Error("Failed to get system config", "key", constants.SysConfigSiteKeywords, "error", err)
		}
	}
	return keywords
}

func (s *sysConfigService) GetSiteNofitication() string {
	return cache.SysConfigCache.GetValue(constants.SysConfigSiteNotification)
}

func (s *sysConfigService) GetMenuItems() []model.MenuItem {
	str := cache.SysConfigCache.GetValue(constants.SysConfigMenuItems)
	var menuItems []model.MenuItem
	if strs.IsNotBlank(str) {
		if err := jsons.Parse(str, &menuItems); err != nil {
			slog.Error("Failed to get system config", "key", constants.SysConfigMenuItems, "error", err)
		}
	}
	return menuItems
}

func (s *sysConfigService) GetRecommendTags() []string {
	str := cache.SysConfigCache.GetValue(constants.SysConfigRecommendTags)
	var recommendTags []string
	if strs.IsNotBlank(str) {
		if err := jsons.Parse(str, &recommendTags); err != nil {
			slog.Error("Failed to get system config", "key", constants.SysConfigRecommendTags, "error", err)
		}
	}
	return recommendTags
}

func (s *sysConfigService) GetTokenExpireDays() int {
	tokenExpireDaysStr := cache.SysConfigCache.GetValue(constants.SysConfigTokenExpireDays)
	tokenExpireDays, err := strconv.Atoi(tokenExpireDaysStr)
	if tokenExpireDays <= 0 || err != nil {
		return constants.DefaultTokenExpireDays
	}
	return tokenExpireDays
}

func (s *sysConfigService) GetLoginMethod() model.LoginMethod {
	str := cache.SysConfigCache.GetValue(constants.SysConfigLoginMethod)

	if strs.IsNotBlank(str) {
		var loginMethod model.LoginMethod
		err := jsons.Parse(str, &loginMethod)
		if err == nil {
			return loginMethod
		}
		slog.Error("Failed to get system config", "key", constants.SysConfigLoginMethod, "error", err)
	}
	return model.LoginMethod{
		Password: true,
		QQ:       true,
		Github:   true,
	}
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
			slog.Error("Failed to get system config", "key", constants.SysConfigSiteNavs, "error", err)
		}
	}
	return siteNavsArr
}

func (s *sysConfigService) GetUserObserveSeconds() int {
	userObserveSecondsStr := cache.SysConfigCache.GetValue(constants.SysConfigUserObserveSeconds)
	userObserveSeconds, _ := strconv.Atoi(userObserveSecondsStr)
	return userObserveSeconds
}

func (s *sysConfigService) GetModules() model.ModulesConfig {
	str := cache.SysConfigCache.GetValue(constants.SysConfigModules)

	if strs.IsNotBlank(str) {
		var modulesConfig model.ModulesConfig
		err := jsons.Parse(str, &modulesConfig)
		if err == nil {
			return modulesConfig
		}
		slog.Error("Failed to get system config", "key", constants.SysConfigModules, "error", err)
	}
	return model.ModulesConfig{
		Tweet:   true,
		Topic:   true,
		Article: true,
	}
}

// GetEmailWhitelist 邮箱白名单
func (s *sysConfigService) GetEmailWhitelist() []string {
	str := cache.SysConfigCache.GetValue(constants.SysConfigEmailWhitelist)
	var emailWhitelist []string
	if strs.IsNotBlank(str) {
		if err := jsons.Parse(str, &emailWhitelist); err != nil {
			slog.Error("Failed to get system config", "key", constants.SysConfigEmailWhitelist, "error", err)
		}
	}
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
