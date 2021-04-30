package lru

import (
	"testing"
)

type String string

func (str String) Len() int {
	return len(str)
}
func TestCache_Get(t *testing.T) {
	//1.创建lru对象
	cacheLru := New(5, nil)

	//2.添加元素到lru
	cacheLru.Add("key1", String("1234"))

	//3.验证是否加入成功
	if v, ok := cacheLru.Get("key1"); ok && v.(String) == "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}

	if _, ok := cacheLru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}
