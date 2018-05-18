package main

import (
	"fmt"
	"time"
	"sync"
)
var wg  = sync.WaitGroup{}

func say(s string) {
	wg.Add(1)
	defer wg.Done()
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
	wg.Wait()
}
