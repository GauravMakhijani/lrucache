package app

import (
	"errors"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func (l *lruCache) GetCacheCapacity() GetCacheCapacityResponse {
	l.mu.Lock()
	defer l.mu.Unlock()

	return GetCacheCapacityResponse{
		Capacity: l.capacity,
	}
}

func (l *lruCache) GetKeyValue(input string) (*GetKeyValueResponse, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// check if key is present in cache
	if _, ok := l.HashMap[input]; !ok {
		log.Printf("key - %s, Not found", input)
		return nil, errors.New("key not found")
	}

	key := l.HashMap[input].Key
	val := l.HashMap[input].Value
	exp := l.HashMap[input].Expiry

	// remove the node from current position
	l.remove(input)

	// check if key is expired
	currTime := time.Now()
	if currTime.After(exp) {
		log.Printf("key - %s, expired", input)
		return nil, errors.New("key not found")
	}

	// add the node to head
	l.addToHead(*key, *val, exp.Unix())

	return &GetKeyValueResponse{
		Key:   *key,
		Value: *val,
	}, nil
}

func (l *lruCache) reset() {
	l.capacity = 0
	l.HashMap = make(map[string]*node)

	head := node{}
	tail := node{
		prev: &head,
	}
	head.next = &tail

	l.LL = doublyLL{
		head: &head,
		tail: &tail,
	}
}

func (l *lruCache) InitializeCache(input InitializeCacheInput) error {
	// check if capacity is valid
	if input.Capacity <= 0 {
		log.Println("capacity should be greater than 0")
		err := errors.New("capacity should be greater than 0")
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// reset the cache
	l.reset()

	// set the capacity
	l.capacity = input.Capacity

	// broadcast the cache update
	l.broadcastCacheUpdate()
	return nil
}

func (l *lruCache) Insert(input CacheItem) (*CacheItem, error) {

	// check if input is valid
	if err := input.Valid(); err != nil {
		log.Println("invalid input")
		return nil, err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// check if cache is initialized
	if l.capacity == 0 {
		log.Println("cache not initialized")
		err := errors.New("cache not initialized")
		return nil, err
	}

	log.Println("inserting key in cache - ", input.Key)

	// check if key is already present in cache
	if _, ok := l.HashMap[input.Key]; ok {
		log.Println("key already present in cache - ", input.Key)
		l.remove(input.Key)
		l.addToHead(input.Key, input.Value, input.Expiry)
		log.Println("key updated in cache - ", input.Key)
		return &CacheItem{
			Key:   input.Key,
			Value: input.Value,
		}, nil
	}

	log.Println("key not present in cache - ", input.Key)

	// if key is not present in cache
	// check if cache is full
	if len(l.HashMap) == l.capacity {
		log.Println("cache is full - removing the last node")
		// remove the last node
		l.removeLast()
	}

	log.Println("adding key in cache - ", input.Key)

	// add the new node
	l.addToHead(input.Key, input.Value, input.Expiry)

	log.Println("key added in cache - ", input.Key)

	return &CacheItem{
		Key:   input.Key,
		Value: input.Value,
	}, nil
}

func (l *lruCache) RemoveFromCache(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.remove(key)
}

func (l *lruCache) remove(key string) {

	// remove the node from current position
	existingNode := l.HashMap[key]
	prevNode := existingNode.prev
	nextNode := existingNode.next
	prevNode.next = nextNode

	nextNode.prev = prevNode
	existingNode = nil // will be taken care by garbage collector

	delete(l.HashMap, key)

	l.broadcastCacheUpdate()
}

func (l *lruCache) addToHead(key string, value string, expiry int64) {
	newNode := node{
		Key:    &key,
		Value:  &value,
		Expiry: time.Unix(expiry, 0),
	}
	nextNode := l.LL.head.next
	l.LL.head.next = &newNode
	newNode.prev = l.LL.head
	newNode.next = nextNode
	nextNode.prev = &newNode
	l.HashMap[key] = &newNode
	l.broadcastCacheUpdate()
}

func (l *lruCache) removeLast() {
	lastNode := l.LL.tail.prev
	prevNode := lastNode.prev
	prevNode.next = l.LL.tail

	l.LL.tail.prev = prevNode

	// remove the last node from hashmap
	delete(l.HashMap, *lastNode.Key)

	l.broadcastCacheUpdate()

	lastNode = nil // will be taken care by garbage collector
}

func (l *lruCache) StartExpirationCheckerWithCron() {
	// start a cron job to check for expirations
	cron := cron.New()
	cron.AddFunc("@every 1s", func() {
		log.Println("[CRON] Checking for expirations")
		l.checkExpirations()
	})
	cron.Start()
}

func (l *lruCache) checkExpirations() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for key, item := range l.HashMap {
		if now.After(item.Expiry) {
			l.remove(key)
		}
	}
}

func (l *lruCache) GetCacheState() []CacheItem {
	l.mu.Lock()
	defer l.mu.Unlock()

	var state []CacheItem
	for _, item := range l.HashMap {
		if time.Now().Before(item.Expiry) {
			cacheItem := CacheItem{
				Key:    *item.Key,
				Value:  *item.Value,
				Expiry: item.Expiry.Unix(),
			}
			state = append(state, cacheItem)
		}
	}

	return state

}

func (l *lruCache) broadcastCacheUpdate() {
	BroadcastChannel <- struct{}{}
}

func ProvideNewCache() *lruCache {
	hashmap := make(map[string]*node)
	head := node{}
	tail := node{
		prev: &head,
	}
	head.next = &tail
	ll := doublyLL{
		head: &head,
		tail: &tail,
	}

	BroadcastChannel = make(chan struct{})

	l := &lruCache{
		capacity: 0,
		HashMap:  hashmap,
		LL:       ll,
	}
	l.StartExpirationCheckerWithCron()
	return l
}
