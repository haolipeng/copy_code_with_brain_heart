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

//结论：以值传递方式使用mutex，会复制其锁的状态，但是操作的并不是同一把锁，切记！！！

func (s Student) UnlockWithRef() {
	//输出mutex的值
	fmt.Printf("object: %p before mutex.UnlockWithRef() mutext:%v\n", &s, s.mu)
	s.mu.Unlock()
	fmt.Printf("object: %p after mutex.UnlockWithRef() mutext:%v\n", &s, s.mu)
}

func (s *Student) UnlockWithPointer() {
	//输出mutex的值
	fmt.Printf("object: %p before mutex.UnlockWithPointer() mutext:%v\n", &s, s.mu)
	s.mu.Unlock()
	fmt.Printf("object: %p after mutex.UnlockWithPointer() mutext:%v\n", &s, s.mu)
}

func main() {
	/*var stu Student
	stu.mu.Lock()

	//经过两次解锁后，看看是否为空
	fmt.Println("origin mutex:", &stu.mu)
	stu.UnlockWithRef()
	stu.UnlockWithRef()*/

	var stu2 Student
	stu2.mu.Lock()

	fmt.Println("origin mutex:", &stu2.mu)
	stu2.UnlockWithPointer()
	stu2.UnlockWithPointer() //调用两次unlock函数会触发如下错误，sync: unlock of unlocked mutex
}
