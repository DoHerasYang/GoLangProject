package main

func main() {
	var _ int = 123
	var _ = 123
}

const (
	value0 = iota<<2 + 1
	value1
)
