package main

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// IntSliceDeleteFunc1 最基本的删除——以int切片为例(不考虑顺序 O(1))
// 整体时间消耗固定在一个整体范围之内,不会有太大性能损耗
func IntSliceDeleteFunc1(inputSlice []int, deletePin int) []int {
	// 判断是否满足删除
	if deletePin < 0 || deletePin >= len(inputSlice) {
		panic("Out of Bound of inputSlice")
	}
	return append(inputSlice[:deletePin], inputSlice[deletePin+1:]...)
}

// IntSliceDeleteFunc2 最基本的删除——以int切片为例(考虑复制顺序牺牲性能 O(n))
// 该函数的性能由deletePin的位置以及Slice长度决定，具有突变性
func IntSliceDeleteFunc2(inputSlice []int, deletePin int) []int {
	// Use Copy to remain Sequence of Slice
	if deletePin < 0 || deletePin >= len(inputSlice) {
		panic("Out of Bound of inputSlice")
	}
	copy(inputSlice[deletePin:], inputSlice[deletePin+1:])
	return inputSlice[:len(inputSlice)-1]
}

// 老师给出的标准答案 - 实现的标准答案
// 用于测试性能
func IntSliceDeleteDemoFunc[T any](inputSlice []T, index int) ([]T, error) {
	var ErrIndexOutOfRange = errors.New("下标超出范围")
	length := len(inputSlice)
	if index < 0 || index >= length {
		return nil, fmt.Errorf("%w, 下标超出范围，长度 %d, 下标 %d", ErrIndexOutOfRange, length, index)
	}
	for i := index; i+1 < length; i++ {
		inputSlice[i] = inputSlice[i+1]
	}
	return inputSlice[:length-1], nil
}

// SliceDeleteGenericFunc 支持泛型 - 支持泛型的改造
func SliceDeleteGenericsFunc[T any](inputSlice []T, deletePin int) []T {
	if deletePin < 0 || deletePin >= len(inputSlice) {
		panic("Out of Bound of inputSlice")
	}
	return append(inputSlice[:deletePin], inputSlice[deletePin+1:]...)
}

// Slice Shrink
func SliceShrinkFunc[T any](inputSlice []T) ([]T, bool) {
	// Obtain Length and Capacity
	length, capacity := len(inputSlice), cap(inputSlice)
	newCapacity, changed := func(sliceLength int, sliceCapacity int) (_ int, _ bool) {
		if sliceCapacity < 64 {
			return sliceCapacity, false
		}
		if sliceCapacity >= 1024 && (sliceCapacity/sliceLength > 2) {
			return int(float64(sliceCapacity) * 0.615), true
		}
		if sliceCapacity < 1024 && (sliceCapacity/sliceLength >= 3) {
			return int(float64(sliceCapacity) * 0.525), true
		}
		// return default
		return sliceCapacity, false
	}(length, capacity)
	if !changed {
		return inputSlice, false
	} else {
		newSlice := make([]T, 0, newCapacity)
		newSlice = append(newSlice, inputSlice...)
		return newSlice, true
	}
}

var SliceLength int = 100000
var Seed = time.Now().UnixNano()

func BenchmarkIntSliceDeleteFunc1(b *testing.B) {
	b.StopTimer() // Stop Timer
	randomIntGenerator := rand.New(rand.NewSource(Seed))
	testIntSlice := make([]int, SliceLength)
	for i := 0; i < SliceLength; i++ {
		testIntSlice[i] = i
	}
	b.StartTimer() // Start Timer
	// Benchmark Test Part
	for i := 0; i < b.N; i++ {
		IntSliceDeleteFunc1(testIntSlice, randomIntGenerator.Intn(SliceLength-1))
	}
}

func BenchmarkIntSliceDeleteFunc2(b *testing.B) {
	b.StopTimer()
	randomIntGenerator := rand.New(rand.NewSource(Seed))
	testIntSlice := make([]int, SliceLength)
	for i := 0; i < SliceLength; i++ {
		testIntSlice[i] = i
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		IntSliceDeleteFunc2(testIntSlice, randomIntGenerator.Intn(SliceLength-1))
	}
}

func BenchmarkIntSliceDeleteDemoFunc(b *testing.B) {
	b.StopTimer()
	randomIntGenerator := rand.New(rand.NewSource(Seed))
	testIntSlice := make([]int, SliceLength)
	for i := 0; i < SliceLength; i++ {
		testIntSlice[i] = i
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IntSliceDeleteDemoFunc(testIntSlice, randomIntGenerator.Intn(SliceLength-1))
	}
}

func BenchmarkSliceDeleteGenericsFunc(b *testing.B) {
	b.StopTimer()
	randomIntGenerator := rand.New(rand.NewSource(Seed))
	testIntSlice := make([]int, SliceLength)
	for i := 0; i < SliceLength; i++ {
		testIntSlice[i] = i
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = SliceDeleteGenericsFunc(testIntSlice, randomIntGenerator.Intn(SliceLength-1))
	}
}

func BenchmarkSliceShrinkFunc(b *testing.B) {
	b.StopTimer()
	randomIntGenerator := rand.New(rand.NewSource(Seed))
	// 创建一个新的需要缩容的数组
	randomLength := randomIntGenerator.Intn(3000)
	testIntSlice := make([]int, randomLength, SliceLength)
	for i := 0; i < randomLength; i++ {
		testIntSlice[i] = i
	}
	b.StartTimer()
	newIntSlice, _ := SliceShrinkFunc(testIntSlice)
	b.StopTimer()
	fmt.Printf("Initialized Capacity:%d, After Shrinked Capacity is %d \n", SliceLength, cap(newIntSlice))
}

// go test -bench=. -benchmem -count 5
// 打印显示的结果
//(base) doheras@DoHerass-MBP first_codework % go test -bench=. -benchmem -count 5
//goos: darwin
//goarch: amd64
//pkg: GoLangProject/codework/week1
//cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
//BenchmarkIntSliceDeleteFunc1-12            86006             14335 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc1-12            83460             14263 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc1-12            83533             14344 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc1-12            83468             14300 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc1-12            84486             14211 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc2-12            84120             14233 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc2-12            84284             14260 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc2-12            84362             14308 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc2-12            83611             14341 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteFunc2-12            83426             14216 ns/op               0 B/op          0 allocs/op
//BenchmarkIntSliceDeleteDemoFunc-12         37184             32529 ns/op              16 B/op          1 allocs/op
//BenchmarkIntSliceDeleteDemoFunc-12         36165             32393 ns/op              16 B/op          1 allocs/op
//BenchmarkIntSliceDeleteDemoFunc-12         36816             32387 ns/op              16 B/op          1 allocs/op
//BenchmarkIntSliceDeleteDemoFunc-12         37149             32447 ns/op              16 B/op          1 allocs/op
//BenchmarkIntSliceDeleteDemoFunc-12         37016             32502 ns/op              16 B/op          1 allocs/op
//BenchmarkSliceDeleteGenericsFunc-12        83253             14223 ns/op               0 B/op          0 allocs/op
//BenchmarkSliceDeleteGenericsFunc-12        83122             14218 ns/op               0 B/op          0 allocs/op
//BenchmarkSliceDeleteGenericsFunc-12        83343             14015 ns/op               0 B/op          0 allocs/op
//BenchmarkSliceDeleteGenericsFunc-12        87692             13889 ns/op               0 B/op          0 allocs/op
//BenchmarkSliceDeleteGenericsFunc-12        86766             13852 ns/op               0 B/op          0 allocs/op
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//BenchmarkSliceShrinkFunc-12             Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//1000000000               0.0000190 ns/op               0 B/op          0 allocs/op
//BenchmarkSliceShrinkFunc-12             Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//1000000000               0.0000239 ns/op               0 B/op          0 allocs/op
//BenchmarkSliceShrinkFunc-12             Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//1000000000               0.0000296 ns/op               0 B/op          0 allocs/op
//BenchmarkSliceShrinkFunc-12             Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//1000000000               0.0000351 ns/op               0 B/op          0 allocs/op
//BenchmarkSliceShrinkFunc-12             Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//Initialized Capacity:100000, After Shrinked Capacity is 61500
//1000000000               0.0000334 ns/op               0 B/op          0 allocs/op
//PASS
//ok      GoLangProject/codework/week1    28.444s
