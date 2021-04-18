package main

import (
	"container/list"
	"fmt"
)

type entry struct {
	key   string
	value string
}

func main() {
	//创建链表
	ll := list.New()

	//向链表添加元素
	ll.PushBack(entry{
		key:   "haolipeng",
		value: "32",
	})
	ll.PushBack(entry{
		key:   "zhouyang",
		value: "33",
	})

	//遍历链表中所有元素
	for e := ll.Front(); e != nil; e = e.Next() {
		kv := e.Value.(entry)
		fmt.Printf("entry key:%s value:%s\n", kv.key, kv.value)
	}
}
