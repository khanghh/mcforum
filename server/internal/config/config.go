package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"bbs-go/common/jsons"
	"bbs-go/common/strs"
	"bbs-go/sqls"

	"github.com/spf13/viper"
)

var instance *Config

func Instance() *Config {
	return instance
}

func init() {
	env := os.Getenv("BBSGO_ENV")
	if strs.IsBlank(env) {
		env = "dev"
	}

	slog.Info("Load config", slog.String("ENV", env))

	viper.SetConfigName("bbs-go." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.bbs-go")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BBSGO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&instance); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	instance.Env = env
	slog.Info("Load config", slog.String("ENV", env), slog.String("config", jsons.ToJsonStr(instance)))
}

type Config struct {
	Env            string   // 环境
	BaseUrl        string   // base url
	Port           string   // 端口
	IpDataPath     string   // IP数据文件
	AllowedOrigins []string // 跨域白名单
	Language       string

	// 日志配置
	Logger struct {
		Filename   string // 日志文件的位置
		MaxSize    int    // 文件最大尺寸（以MB为单位）
		MaxAge     int    // 保留旧文件的最大天数
		MaxBackups int    // 保留的最大旧文件数量
	}

	// 数据库配置
	DB sqls.DbConfig

	// 阿里云oss配置
	Uploader struct {
		Enable    string
		AliyunOss struct {
			Host          string
			Bucket        string
			Endpoint      string
			AccessId      string
			AccessSecret  string
			StyleSplitter string
			StyleAvatar   string
			StylePreview  string
			StyleSmall    string
			StyleDetail   string
		}
		Local struct {
			Host string
			Path string
		}
	}

	// 百度SEO相关配置
	// 文档：https://ziyuan.baidu.com/college/courseinfo?id=267&page=2#h2_article_title14
	BaiduSEO struct {
		Site  string
		Token string
	}

	// 神马搜索SEO相关
	// 文档：https://zhanzhang.sm.cn/open/mip
	SmSEO struct {
		Site     string
		UserName string
		Token    string
	}

	// smtp
	Smtp struct {
		Host     string
		Port     string
		Username string
		Password string
		SSL      bool
	}

	Search struct {
		IndexPath string
	}
}

func (c *Config) IsProd() bool {
	e := strings.ToLower(c.Env)
	return e == "prod" || e == "production"
}
