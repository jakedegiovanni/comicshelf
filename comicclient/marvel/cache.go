package marvel

import "sync"

type Cache[T any] struct {
	mu    *sync.RWMutex
	cache map[string]T
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		mu:    &sync.RWMutex{},
		cache: make(map[string]T),
	}
}

func (c *Cache[T]) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.cache[key]
	return ok
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.cache[key]
	return val, ok
}

func (c *Cache[T]) Put(key string, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = val
}

func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}
