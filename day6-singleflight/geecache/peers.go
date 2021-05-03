package geecache

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool) //根据key选择对应的节点
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error) //从对应group查找key对应的缓存值
}
