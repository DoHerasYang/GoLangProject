package main

import (
	"math"
	"unicode/utf8"
)

func main() {
	var a_num int = 456
	var b_num int = 123
	var c_float float64 = 1.23
	println(a_num * b_num)
	println(a_num - b_num)
	print(a_num & b_num)
	print(float64(a_num) + c_float)
	a_num++
	b_num--
	math.Abs(float64(a_num))
	print(utf8.RuneCountInString("Hello你好"))
	Byte()
}

func Byte() {
	var a = "hello"
	var b = []byte(a)
	var str2 = string(b)
	print(str2)
}
