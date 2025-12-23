package cache

import (
	"bbs-go/internal/model/constants"
	"errors"
	"time"

	"bbs-go/sqls"

	"github.com/goburrow/cache"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

var (
	topicRecommendCacheKey = "recommend_topics_cache"
)

var TopicCache = newTopicCache()

type topicCache struct {
	recommendCache cache.LoadingCache
}

func newTopicCache() *topicCache {
	return &topicCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (cache.Value, error) {
				val := repository.TopicRepository.Find(sqls.DB(),
					sqls.NewCnd().Eq("status", constants.StatusActive).Desc("id").Limit(50))
				if val != nil {
					return val, nil
				}
				return nil, errors.New("not found")
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (c *topicCache) GetRecommendTopics() []model.Topic {
	val, err := c.recommendCache.Get(topicRecommendCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Topic)
	}
	return nil
}

func (c *topicCache) InvalidateRecommend() {
	c.recommendCache.Invalidate(topicRecommendCacheKey)
}
