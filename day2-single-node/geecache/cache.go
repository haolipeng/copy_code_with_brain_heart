package geecache

import (
	"copy-group-cache/day2-single-node/geecache/lru"
	"sync"
)

//定义数据结构
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

//add函数 fix bug:c cache -> c *cache
//方法在传指针和传对象有啥区别呢?
func (c cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//lru exist or not
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

//get函数,返回获取结果 fix bug:c cache -> c *cache
func (c cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var value ByteView
	if c.lru == nil {
		return value, false
	}

	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), ok
	}

	return value, false
}
