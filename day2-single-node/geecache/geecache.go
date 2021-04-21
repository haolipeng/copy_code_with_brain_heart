package geecache

import "sync"

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
	name      string     //group name
	getter    GetterFunc //getter function，缓存未命中时调用
	mainCache cache      //主存
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

// Get 根据名称来获取缓存组
func Get(key string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	if g, exist := groups[key]; exist {
		return g
	}
	return nil
}
