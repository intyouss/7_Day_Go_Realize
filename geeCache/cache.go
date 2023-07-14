package geeCache

import (
	"geeCache/lru"
	"sync"
)

type cache struct {
	mutex      sync.RWMutex
	lru        *lru.Cache
	cacheBytes int64
}
