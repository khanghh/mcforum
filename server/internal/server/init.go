package server

import (
	"bbs-go/internal/config"
	"bbs-go/internal/model"
	"bbs-go/internal/scheduler"
	"bbs-go/internal/search"
	"bbs-go/pkg/iplocator"
	"log/slog"
	"time"

	"bbs-go/sqls"

	_ "bbs-go/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Init() {
	initLogger()
	initDB()
	initCron()
	initIpLocator()
	initSearch()

}

func initDB() {
	conf := config.Instance().DB
	db, err := gorm.Open(mysql.Open(conf.Url), &gorm.Config{
		TranslateError: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.Instance().TablePrefix,
		},
	})

	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
		sqlDB.SetConnMaxIdleTime(time.Duration(conf.ConnMaxIdleTimeSeconds) * time.Second)
		sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeSeconds) * time.Second)
	}

	if err := db.AutoMigrate(model.Models...); nil != err {
		slog.Error(err.Error(), slog.Any("error", err))
	}

	sqls.SetDB(db)
}

func initCron() {
	if config.Instance().IsProd() {
		scheduler.Start()
	}
}

func initIpLocator() {
	iplocator.InitIpLocator(config.Instance().IpDataPath)
}

func initSearch() {
	search.Init(config.Instance().Search.IndexPath)
}
