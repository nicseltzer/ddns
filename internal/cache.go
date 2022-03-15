package internal

import "sync"

type IPCache struct {
	value string
	sync.RWMutex
}

var cache IPCache = IPCache{}

func SetCachedIP(value string) {
	cache.RLock()
	defer cache.RUnlock()
	cache.value = value
}

func GetCachedIP() string {
	cache.Lock()
	defer cache.Unlock()
	return cache.value
}
