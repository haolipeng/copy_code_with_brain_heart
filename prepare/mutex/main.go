package main

import (
	"fmt"
	"sync"
)

type Student struct {
	mu   sync.Mutex
	age  int
	name string
}

func (s Student) add() {
	//输出mutex的值
	fmt.Println("mutex:", s.mu)
}

func main() {
	var stu Student
	stu.mu.Lock()
	fmt.Println("origin mutex:", stu.mu)
	stu.add()
	stu.mu.Unlock()
}
