package main

func ClosureFunc1() {
	for i := 0; i <= 5; i++ {
		defer func() {
			println(i)
		}()
	}
}

func ClosureFunc2() {
	for i := 0; i <= 5; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

func ClosureFunc3() {
	for i := 0; i <= 5; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}

func main() {
	println("ClosureFunc1:")
	ClosureFunc1()
	println("ClosureFunc2:")
	ClosureFunc2()
	println("ClosureFunc3:")
	ClosureFunc3()
}
