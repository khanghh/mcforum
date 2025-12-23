package cache

import (
	"github.com/goburrow/cache"
)

func key2Int64(key cache.Key) int64 {
	return key.(int64)
}

func key2String(key cache.Key) string {
	return key.(string)
}
