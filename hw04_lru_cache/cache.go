package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	mutex    *sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*listItem),
		mutex:    &sync.Mutex{},
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mutex.Lock()
	if itm, ok := lc.items[key]; ok { // refresh
		refreshed := cacheItem{key, value}
		itm.Value = refreshed
		lc.queue.MoveToFront(lc.items[key])
		lc.mutex.Unlock()
		return true
	}
	// insert
	lc.items[key] = lc.queue.PushFront(cacheItem{key, value})
	if len(lc.items) > lc.capacity { // remove old record
		old := lc.queue.Back()
		delete(lc.items, old.Value.(cacheItem).key)
		lc.queue.Remove(old)
	}
	lc.mutex.Unlock()
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mutex.Lock()
	if itm, ok := lc.items[key]; ok { // return value
		lc.queue.MoveToFront(lc.items[key])
		lc.mutex.Unlock()
		return itm.Value.(cacheItem).value, true
	}
	lc.mutex.Unlock()
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mutex.Lock()
	lc.queue = NewList()
	lc.items = make(map[Key]*listItem)
	lc.mutex.Unlock()
}
