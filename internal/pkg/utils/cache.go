package utils

import (
	"sync"

	"github.com/silenceper/wechat/v2/cache"
)

var memCache *cache.Memory

func GetCache() *cache.Memory {
	var once sync.Once
	once.Do(func() {
		memCache = cache.NewMemory()
	})
	return memCache
}
