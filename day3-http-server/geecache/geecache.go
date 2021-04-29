package geecache

import (
	"errors"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//全局变量
var (
	mu     sync.RWMutex              //读写锁，控制并发来获取缓存组
	groups = make(map[string]*Group) //缓存组列表，用于存储不同缓存组
)

type Group struct {
	name string //group name
	//在这块我出错了，getter变量到底应该填写什么类型呢？
	getter    Getter //getter function，缓存未命中时调用
	mainCache cache  //主存
}

// NewGroup 创建缓存组
func NewGroup(name string, cacheBytes int64, getter GetterFunc) *Group {
	//校验getterFunc是否为空
	if getter == nil {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()

	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	//添加到缓存组列表
	groups[name] = group

	return group
}

// GetGroup 根据名称来获取创建的缓存组，如果没有则返回空
func GetGroup(key string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	if g, exist := groups[key]; exist {
		return g
	}
	return nil
}

//Get 获取缓存组中key对应的缓存
func (g *Group) Get(key string) (ByteView, error) {
	//1.校验key
	bv := ByteView{}
	var err error

	if key == "" {
		//key是必须的，不能填空
		err = errors.New("key is required")
		return bv, err
	}

	//2.从mainCache中读取，成功则返回
	if bv, ok := g.mainCache.get(key); ok {
		return bv, nil
	}

	//3.失败，则从本地获取
	bv, err = g.load(key)
	return bv, err
}

//load 从远端节点或不同数据源获取数据
func (g *Group) load(key string) (ByteView, error) {
	//TODO:分布式场景下，会调用getFromPeer从其他节点获取，此处暂不实现
	//调用getLocally
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	//1.调用用户注册的Getter函数
	b, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	//2.封装成ByteView类型
	bv := ByteView{b: cloneBytes(b)}

	//3.将键值对回填到mainCache缓存中
	g.populateCache(key, bv)

	return bv, nil
}

//回填缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
