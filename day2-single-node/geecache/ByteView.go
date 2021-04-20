package geecache

type ByteView struct {
	b []byte //使用byte作为缓存值类型，就是
}

// Len 实现Value接口，必须实现
func (v ByteView) Len() int {
	return len(v.b)
}

//实现String()函数
func (v ByteView) String() string {
	return string(v.b)
}

// ByteSlice 实现底层数据data的一份拷贝
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

//cloneBytes实现底层数据的真正拷贝动作
func cloneBytes(b []byte) []byte {
	//new []byte
	newBytes := make([]byte, len(b))
	copy(newBytes, b)

	return newBytes
}
