package cache

import (
	"errors"
	"log/slog"
	"time"

	"bbs-go/sqls"

	"github.com/goburrow/cache"

	"bbs-go/internal/model"
	"bbs-go/internal/repository"
)

type tagCache struct {
	cache cache.LoadingCache // 标签缓存
}

var TagCache = newTagCache()

func newTagCache() *tagCache {
	return &tagCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repository.TagRepository.Get(sqls.DB(), key2Int64(key))
				if value == nil {
					e = errors.New("数据不存在")
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *tagCache) Get(tagId int64) *model.Tag {
	val, err := c.cache.Get(tagId)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return nil
	}
	if val != nil {
		return val.(*model.Tag)
	}
	return nil
}

func (c *tagCache) GetList(tagIds []int64) (tags []model.Tag) {
	if len(tagIds) == 0 {
		return nil
	}
	for _, tagId := range tagIds {
		tag := c.Get(tagId)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}
	return
}

func (c *tagCache) Invalidate(tagId int64) {
	c.cache.Invalidate(tagId)
}
