首先我先介绍一下录制此视频课程的初衷，我在阅读并实践《七天用Go从零实现系列》时，遇到了一些困惑，而且我发现博客留言区也有不少小伙伴有迫切的想看视频教程的心理，所以我做了这份视频教程。

视频教程是基于geektutu.com博客上的《七天用Go从零实现系列》来节选录制的，希望能帮忙到一起学习进步的小伙伴。

# 一、简介

LRU简写是(Least Recently Used)，意义为最近最少访问，LRU算法的实现很简单，维护一个队列，队列中存储全部的缓存值，每次把最近访问的元素重新放置到队首。

这里使用go标准库实现的双向链表list.List，双向链表作为队列，队首队尾是相对的，这里约定front是队首。

队列的头部是**经常被访问**的数据，队列的尾部是**最近最少访问**的数据，如果所有元素的内存超过设置的阈值时，将队尾的数据元素淘汰即可。

# 二、基本功能需求

1、频繁查找某个元素在缓存中是否存在

2、淘汰最近访问最少的元素

3、新增/更新元素

# 三、实现思路

基于功能需求的3条需求，既要求增加元素和删除元素的效率，也要求查找元素的效率，所以我们选择将链表和字典相结合

## 3、1 核心数据结构

1、map 存储缓存名称和缓存值的映射关系，其增删操作的复杂度都是O（1）

2、将所有缓存值都存到双向链表中，在队尾新增一条记录以及删除一条记录的复杂度均为O(1)

**值类型是interface类型**

为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。

```go
// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}
```

**链表的节点类型**

键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队尾节点时，需要用 key 从字典中删除对应的映射。

```go
type entry struct {
	key   string
	value Value
}
```

**缓存Cache结构体**

这里使用go标准库实现的双向链表list.List

字典的定义是 map[string] *list.Element，键是字符串，值是双向链表中对应节点的指针。
maxBytes 是允许使用的最大内存，nbytes 是当前已使用的内存，

OnEvicted 是某条记录被移除时的回调函数，可以为 nil。



```go
// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}
```

为方便实例化Cache，我们实现New()函数

```go
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
```

## 3、2 查找元素

从字典中找到对应的缓存节点是否存在：

- 如果键对应的链表节点存在，则将对应节点移动到队尾，并返回查找到的值。

## 3、3 删除元素

移除最近最少访问的节点（即队尾的元素）

步骤：

1、取队尾元素，从双向链表中删除

2、从字典中删除该节点的映射关系

3、更新当前所用的内存c.nbytes

4、如果回调函数OnEvicted不为空，则在删除元素时调用此回调函数。

## 3、4 新增/修改

1、如果键对应的链表节点存在，则更新对应节点的值，并将该节点移动至队首。

不存在则是新增场景，在队首添加新节点&entry{key,value},并在字典中添加key和节点的映射关系

2、更新c.nbytes，如果超过了设定的最大值c.maxbytes，则移除最少访问的节点

更新操作c.nbytes：更新操作,需要在原计数上 + 新值与旧值的差值（可正可负）