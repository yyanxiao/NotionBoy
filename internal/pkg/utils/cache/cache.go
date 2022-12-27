package cache

import (
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	*gocache.Cache
}

var (
	client *Cache
	once   sync.Once
)

func DefaultClient() *Cache {
	if client == nil {
		once.Do(func() {
			c := gocache.New(1*time.Hour, 2*time.Hour)
			client = &Cache{
				Cache: c,
			}
		})
	}
	return client
}
