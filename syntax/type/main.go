package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

type PersonInterface interface {
	PersonInitializerPointer(Name string, Age int)
	PersonInitializerWithOutPointer(Name string, Age int)
}

func (l *Person) PersonInitializerPointer(Name string, Age int) {
	l.Name = Name
	l.Age = Age
}

// 改不了原始的值
func (l *Person) PersonInitializerWithOutPointer(Name string, Age int) {
	l.Name = Name
	l.Age = Age
}

func main() {
	user1 := Person{}
	user1.PersonInitializerPointer("XiaoMing", 13)
	fmt.Printf("%+v \n", user1)

	user2 := &Person{}
	fmt.Printf("%+v \n", *user2)

	user3 := &Person{}
	// 发生了复制 user3 不是原来的
	user3.PersonInitializerWithOutPointer("XiaoFeng", 22)
	fmt.Printf("%+v \n", user3)

	user4 := &Person{}
	user4.PersonInitializerPointer("XiaoLiang", 18)
	fmt.Printf("%+v", *user4)

}
