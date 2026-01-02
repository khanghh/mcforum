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
	Env            string   // Environment
	BaseUrl        string   // base url
	Port           string   // Port
	IpDataPath     string   // IP data file
	AllowedOrigins []string // CORS whitelist
	Language       string
	TablePrefix    string

	// Log configuration
	Logger struct {
		Filename   string // Log file location
		MaxSize    int    // Max file size (in MB)
		MaxAge     int    // Max age to retain old files (days)
		MaxBackups int    // Max number of old files to keep
	}

	// Database configuration
	DB sqls.DbConfig

	Uploader struct {
		Enable string
		Local  struct {
			BaseURL string
			RootDir string
		}
		SUpload struct {
			UploadURL string
			Secret    string
		}
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
