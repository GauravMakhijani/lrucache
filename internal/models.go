package app

import (
	"errors"
	"sync"
	"time"
)

type Cache interface {
	Insert(input CacheItem) (*CacheItem, error)
	InitializeCache(input InitializeCacheInput) error
	GetKeyValue(input string) (*GetKeyValueResponse, error)
	GetCacheState() []CacheItem
	RemoveFromCache(key string)
	GetCacheCapacity() GetCacheCapacityResponse
}

var BroadcastChannel chan struct{}

type node struct {
	Key    *string
	Value  *string
	Expiry time.Time
	next   *node
	prev   *node
}

type doublyLL struct {
	head *node
	tail *node
}

type lruCache struct {
	capacity int
	HashMap  map[string]*node
	LL       doublyLL
	mu       sync.Mutex
}

type CacheItem struct {
	Key    string
	Value  string
	Expiry int64
}

func (i CacheItem) Valid() error {
	if i.Key == "" || i.Value == "" {
		return errors.New("key or value is empty")
	}
	// expiry should be valid timestamp
	if i.Expiry < time.Now().Unix() {
		return errors.New("expiry should be greater than current time")
	}

	return nil
}

type InitializeCacheInput struct {
	Capacity int `json:"capacity"`
}

type GetKeyValueResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetCacheCapacityResponse struct {
	Capacity int `json:"capacity"`
}
