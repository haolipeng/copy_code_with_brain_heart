package consistenthash

import (
	"sort"
	"strconv"
)

// Hash hash函数类型
type Hash func([]byte) uint32

type Map struct {
	//hash函数,默认为crc32哈希
	hash Hash

	//虚拟节点倍数
	replies int

	//哈希环
	keys []int //哈希环中的元素是已经排序的

	//虚拟节点和真实节点的映射关系
	hashMap map[int]string
}

func New(replies int, hash Hash) *Map {
	return &Map{replies: replies, hash: hash}
}

//Add 可一次性添加多个缓存服务器
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		//虚拟节点的名称为序号 + key名称
		for i := 0; i < m.replies; i++ {
			//计算虚拟节点的哈希值,进行类型转换
			virNodeHash := int(m.hash([]byte(strconv.Itoa(i) + key)))

			//添加到哈希环中
			m.keys = append(m.keys, virNodeHash)

			//维护虚拟节点到真实节点之间的映射关系
			m.hashMap[virNodeHash] = key
		}
	}

	//对哈希环进行排序
	sort.Ints(m.keys)
}

// Get 从缓存中获取指定名称的缓存数据
func (m *Map) Get(key string) {

}
