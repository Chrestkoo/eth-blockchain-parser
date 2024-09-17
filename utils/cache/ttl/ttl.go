package ttl

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	expire time.Duration
	keys   map[K]*entry[V]
	sync.Mutex
}

type entry[V any] struct {
	value V
	t     time.Time
}

func newEntry[V any](val V, t time.Time) *entry[V] {
	return &entry[V]{
		value: val,
		t:     t,
	}
}

// New 创建定时过期的缓存
func New[K comparable, V any](expire time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		expire: expire,
		keys:   make(map[K]*entry[V]),
	}
}

// Put 插入
func (cache *Cache[K, V]) Put(key K, val V) {
	cache.Lock()
	defer cache.Unlock()
	cache.keys[key] = newEntry(val, time.Now())
}

// Get 获取，如果获取不到返回nil
func (cache *Cache[K, V]) Get(key K) (val V, ok bool) {
	cache.Lock()
	defer cache.Unlock()
	el, ok := cache.keys[key]
	if ok {
		if time.Since(el.t) > cache.expire {
			delete(cache.keys, key) //lazy delete
			ok = false
			return
		}
		return el.value, true
	}
	return
}

// Expire 重新設置key的过期时间
func (cache *Cache[K, V]) Expire(expire time.Duration) {
	cache.expire = expire
}
