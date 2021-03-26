package lru

import (
	"container/list"
)

//核心结构体
//1.map 存储键和值的映射关系
//2.双向链表 保存所有缓存值
type Cache struct {
	//允许使用的最大内存
	maxBytes int64

	//当前已使用的内存
	nBytes int64

	//双向链表的节点
	ll *list.List //Package list implements a doubly linked list.

	// 保存对应关系
	cache map[string]*list.Element //链表中对应节点的指针

	// 回调函数
	OnEvicted func(key string, value Value)
}

func New(maxBytes int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	//1.从字典中查找对应的缓存节点是否存在
	if ele, ok := c.cache[key]; ok {
		//2.如果存在，则将对应节点移动到队首
		c.ll.PushFront(ele)

		//3.并返回查找到的值
		kv := ele.Value.(*entry)
		return kv.value, true
	}

	return nil, false
}

//删除接口
//即淘汰最近最少访问的数据
func (c *Cache) RemoveOldest() {
	//1.取队尾元素
	ele := c.ll.Back()

	//2.如存在，则从双向链表中删除
	if ele != nil {
		c.ll.Remove(ele)

		//3.从字典中删除该节点的对应关系
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)

		//4.更新当前占用的内存(减去 key + value)
		c.nBytes -= int64(len(kv.key) + kv.value.Len())

		//5.如果回调函数不为空，则在删除元素时调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	//1.从字典中查找对应的缓存节点是否存在
	if ele, ok := c.cache[key]; ok {
		//2.存在(更新场景)
		//2.1 将节点移动到队首
		c.ll.MoveToFront(ele)

		//2.2 计算当前占用内存(新值 减去 旧值的差值)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len() - kv.value.Len())

		//2.3 更新value值
		kv.value = value

		//2.4 2-3两步不能颠倒顺序
	} else {
		//3.不存在（新增场景）
		//3.1 构建新节点插入到队首
		newEle := c.ll.PushFront(&entry{key: key, value: value})

		//3.2 字典中建立key-value映射关系
		c.cache[key] = newEle

		//3.3 计算当前占用内存(key + value)
		c.nBytes += int64(len(key) + value.Len())
	}

	//4.当前占用内存超过maxBytes最大阈值时，启动淘汰策略
	if c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

//value是任意类型
type Value interface {
	Len() int
}

//链表的节点类型
type entry struct {
	key   string
	value Value
}
