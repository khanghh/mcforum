package cache

import (
	"errors"
	"time"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"

	"bbs-go/common/dates"
	"bbs-go/sqls"

	"github.com/goburrow/cache"
)

type userCache struct {
	cache            cache.LoadingCache
	idCache          cache.Cache
	scoreRankCache   cache.LoadingCache
	checkInRankCache cache.LoadingCache
}

var UserCache = newUserCache()

func newUserCache() *userCache {
	return &userCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (cache.Value, error) {
				val := repository.UserRepository.Get(sqls.DB(), key2Int64(key))
				if val != nil {
					return val, nil
				}
				return nil, errors.New("not found")
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		idCache: cache.New(
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		scoreRankCache: cache.NewLoadingCache(
			func(key cache.Key) (cache.Value, error) {
				val := repository.UserRepository.Find(sqls.DB(), sqls.NewCnd().Desc("score").Limit(10))
				if val != nil {
					return val, nil
				}
				return nil, errors.New("not found")
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(10*time.Minute),
		),
		checkInRankCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				today := dates.GetDay(time.Now())
				value = repository.CheckInRepository.Find(sqls.DB(),
					sqls.NewCnd().Eq("latest_day_name", today).Asc("update_time").Limit(10))
				return
			},
			cache.WithMaximumSize(10),
			cache.WithExpireAfterAccess(1*time.Hour),
		),
	}
}

func (c *userCache) Get(userId int64) *model.User {
	if userId <= 0 {
		return nil
	}
	val, err := c.cache.Get(userId)
	if err != nil {
		return nil
	}
	return val.(*model.User)
}

func (c *userCache) GetByUsername(username string) *model.User {
	id, _ := c.idCache.GetIfPresent(username)
	if id != nil {
		return c.Get(id.(int64))
	}
	user := repository.UserRepository.GetByUsername(sqls.DB(), username)
	if user != nil {
		c.cache.Put(user.ID, user)
		c.idCache.Put(username, user.ID)
	}
	return user
}

func (c *userCache) Invalidate(userId int64) {
	c.cache.Invalidate(userId)
}

func (c *userCache) GetScoreRank() []model.User {
	val, err := c.scoreRankCache.Get("data")
	if err != nil {
		return nil
	}
	return val.([]model.User)
}

func (c *userCache) GetCheckInRank() []model.CheckIn {
	today := dates.GetDay(time.Now())
	val, err := c.checkInRankCache.Get(today)
	if err != nil {
		return nil
	}
	return val.([]model.CheckIn)
}

func (c *userCache) RefreshCheckInRank() {
	c.checkInRankCache.Refresh(dates.GetDay(time.Now()))
}
