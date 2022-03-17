package ipcache

import "sync"

type IPCache interface {
	IP() string
	SetIP(string)
}

type ipCache struct {
	value string
	sync.RWMutex
}

func NewIPCache() IPCache {
	return newIPCache()
}

func newIPCache() *ipCache {
	return &ipCache{
		value:   "",
		RWMutex: sync.RWMutex{},
	}
}

func (i *ipCache) SetIP(value string) {
	i.RLock()
	defer i.RUnlock()
	i.value = value
}

func (i *ipCache) IP() string {
	i.Lock()
	defer i.Unlock()
	return i.value
}
