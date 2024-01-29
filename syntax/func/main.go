package main

import (
	"math/rand"
	"reflect"
	"time"
	"unicode/utf8"
)

func MultipleParametersFunc(num int, parameters ...string) {
	print(reflect.TypeOf(parameters).String()) // string 类型的切片
	for _, i := range parameters {
		println(i)
		println(reflect.TypeOf(i).String())
	}
}

func ByteTestFunc() {
	buffer := make([]rune, 10)
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range buffer {
		buffer[i] = rune(randGenerator.Intn(95) + 32)
	}
	print(string(buffer))
}

func RuneTestFunc() {
	s := "你好Go语言"
	println(len(s))
	for _, r := range s {
		println(r)                          // 原始rune类型
		println(string(r))                  // 将rune转化成为string可以打印原始的
		println(reflect.TypeOf(r).String()) // Rune 类型 int32
	}
	println(utf8.RuneCountInString(s)) // 显示实际的长度、

}

func main() {
	MultipleParametersFunc(1, "a1", "b2", "c3", "d4")
	//ByteTestFunc()
	//RuneTestFunc()
}
