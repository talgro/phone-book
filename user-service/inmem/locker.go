package inmem

import (
	"context"
	"sync"
)

type lockCache struct {
	mutex sync.RWMutex
	data  map[string]interface{}
}

func NewLockCache() *lockCache {
	return &lockCache{
		data: make(map[string]interface{}),
	}
}

func (c *lockCache) Lock(_ context.Context, key string) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.data[key]; !ok {
		c.data[key] = struct{}{}
		return true, nil
	}

	return false, nil
}

func (c *lockCache) Unlock(_ context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)

	return nil
}
