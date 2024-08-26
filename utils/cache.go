package utils

import (
	"sync"
	"time"
)

/*
	目前对内存并没有管理，后续完善。
	针对现有功能，希望使用时人工关注内存使用状态。
*/

type CacheMateData struct {
	Data       []byte
	Expiration int64 // Unix 时间戳，表示到期时间
}

type Cache struct {
	cache      map[any]CacheMateData
	cacheMutex sync.RWMutex // 读写锁
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[any]CacheMateData),
	}
}

// GetFromCache 从缓存中获取数据，如果数据过期则返回 false
func (c *Cache) GetFromCache(tileID any) ([]byte, bool) {
	c.cacheMutex.RLock()
	item, exists := c.cache[tileID]
	c.cacheMutex.RUnlock()

	if !exists {
		return nil, false
	}

	// 检查是否过期，如果 Expiration 为0则表示不过期
	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		// 过期删除
		c.DeleteFromCache(tileID)
		return nil, false
	}

	return item.Data, true
}

// SetToCache 将数据设置到缓存中，并指定过期时间，秒级
// 如果 duration 为0，表示该缓存项没有过期时间
func (c *Cache) SetToCache(tileID any, data []byte, duration ...time.Duration) {
	var expiration int64

	if len(duration) > 0 && duration[0] > 0 {
		expiration = time.Now().Add(duration[0]).Unix() // 计算到期时间
	} else {
		expiration = 0
	}

	c.cacheMutex.Lock()
	c.cache[tileID] = CacheMateData{
		Data:       data,
		Expiration: expiration,
	}
	c.cacheMutex.Unlock()
}

func (c *Cache) DeleteFromCache(tileID any) {
	c.cacheMutex.Lock()
	delete(c.cache, tileID)
	c.cacheMutex.Unlock()
}

func (c *Cache) CountFromCache() int {
	c.cacheMutex.RLock()
	cacheLen := len(c.cache)
	c.cacheMutex.RUnlock()
	return cacheLen
}

func (c *Cache) UpdateCache(tileID any, data []byte, duration ...time.Duration) bool {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	// 检查是否存在
	if _, exists := c.cache[tileID]; exists {
		// 更新数据
		var expiration int64
		if len(duration) > 0 && duration[0] > 0 {
			expiration = time.Now().Add(duration[0]).Unix() // 计算到期时间
		} else {
			expiration = 0
		}

		c.cache[tileID] = CacheMateData{
			Data:       data,
			Expiration: expiration,
		}
		return true
	}

	return false
}

func (c *Cache) Cleanup() {
	for {
		time.Sleep(1 * time.Minute) // 清理时间自定义
		c.cacheMutex.Lock()
		for key, item := range c.cache {
			if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
				delete(c.cache, key)
			}
		}
		c.cacheMutex.Unlock()
	}
}
