package cache

import (
	"errors"
	"time"

	"bbs-go/sqls"

	"github.com/goburrow/cache"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

type sysConfigCache struct {
	cache cache.LoadingCache
}

var SysConfigCache = newSysConfigCache()

func newSysConfigCache() *sysConfigCache {
	return &sysConfigCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (cache.Value, error) {
				val := repository.SysConfigRepository.GetByKey(sqls.DB(), "")
				if val != nil {
					return val, nil
				}
				return nil, errors.New("数据不存在")
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *sysConfigCache) Get(key string) *model.SysConfig {
	val, err := c.cache.Get(key)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*model.SysConfig)
	}
	return nil
}

func (c *sysConfigCache) GetValue(key string) string {
	sysConfig := c.Get(key)
	if sysConfig == nil {
		return ""
	}
	return sysConfig.Value
}

func (c *sysConfigCache) Invalidate(key string) {
	c.cache.Invalidate(key)
}
