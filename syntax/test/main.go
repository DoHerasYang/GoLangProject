package main

import (
	"fmt"
)

type Integer int

type Number interface {
	~int | // 你需要加上|符号
		int | int32 | int64
}

func Sum[T Number](vals []T) T {
	var res T
	for _, v := range vals {
		res = res + v
	}
	return res
}

func main() {
	res := Sum[Integer]([]Integer{1, 2, 3, 4, 5})
	fmt.Println(res)
}
