package main

import (
	"log/slog"
	"os"

	_ "bbs-go/internal/config"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/spf13/viper"
)

type Config struct {
	BaseURL        string // base url
	Port           string // Port
	RootDir        string // Upload root directory
	Secret         string // jwt secret
	AllowedOrigins []string
}

var config *Config

func init() {
	// Load configuration
	initConfig()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", slog.Any("err", err))
		panic(err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		slog.Error("failed to unmarshal config", slog.Any("err", err))
		panic(err)
	}
	config = cfg
}

func main() {

	app := iris.New()
	app.Logger().SetLevel("info")
	app.Use(logger.New())
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodOptions, iris.MethodHead, iris.MethodGet, iris.MethodPost, iris.MethodPut, iris.MethodPatch, iris.MethodDelete},
		AllowedHeaders:   []string{"*"},
	}))

	app.HandleDir("/files", iris.Dir(config.RootDir))

	if err := app.Listen(":"+config.Port,
		iris.WithConfiguration(iris.Configuration{
			TimeFormat: "2006-01-02 15:04:05",
			Charset:    "UTF-8",
		}),
	); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		os.Exit(-1)
	}
}
