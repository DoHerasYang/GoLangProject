package main

import (
	"fmt"
	"sync"
)

func main() {
	var waitgroup sync.WaitGroup
	channel := make(chan string)
	waitgroup.Add(2)
	go multiprintstring(&waitgroup, channel, "GORoutine1\n")
	go multiprintstring(&waitgroup, channel, "GoRoutine2\n")
	go func() {
		waitgroup.Wait()
		close(channel)
	}()
	for obtain_string := range channel {
		fmt.Println(obtain_string)
	}
}

func multiprintstring(waitgroup *sync.WaitGroup, channel chan<- string, input_string string) {
	for i := 0; i < 10; i++ {
		channel <- input_string
		//fmt.Printf(input_string)
	}
	waitgroup.Done()
}
