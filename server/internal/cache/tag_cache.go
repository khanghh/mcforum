package cache

import (
	"errors"
	"log/slog"
	"time"

	"github.com/goburrow/cache"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
	"bbs-go/sqls"
)

type tagCache struct {
	cache cache.LoadingCache // 标签缓存
}

var TagCache = newTagCache()

func newTagCache() *tagCache {
	return &tagCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (cache.Value, error) {
				val := repository.TagRepository.GetByName(sqls.DB(), key2String(key))
				if val != nil {
					return val, nil
				}
				return nil, errors.New("")
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *tagCache) Get(tagName string) *model.Tag {
	val, err := c.cache.Get(tagName)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return nil
	}
	if val != nil {
		return val.(*model.Tag)
	}
	return nil
}

func (c *tagCache) GetList(tagNames []string) (tags []model.Tag) {
	if len(tagNames) == 0 {
		return nil
	}
	for _, tagName := range tagNames {
		tag := c.Get(tagName)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}
	return
}

func (c *tagCache) Invalidate(tagName string) {
	c.cache.Invalidate(tagName)
}
