package cache

import (
	"errors"
	"time"

	"bbs-go/sqls"

	"github.com/goburrow/cache"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var UserTokenCache = newUserTokenCache()

type userTokenCache struct {
	cache cache.LoadingCache
}

func newUserTokenCache() *userTokenCache {
	return &userTokenCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (cache.Value, error) {
				val := repository.UserTokenRepository.GetByToken(sqls.DB(), key.(string))
				if val != nil {
					return val, nil
				}
				return nil, errors.New("cache miss")
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(60*time.Minute),
		),
	}
}

func (c *userTokenCache) Get(token string) *model.UserToken {
	if len(token) == 0 {
		return nil
	}
	val, err := c.cache.Get(token)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*model.UserToken)
	}
	return nil
}

func (c *userTokenCache) Invalidate(token string) {
	c.cache.Invalidate(token)
}
