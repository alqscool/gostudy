package main

import (
	"fmt"
	"runtime"
	"time"
)

func say(s string) {
	for i:=0; i<5; i++ {
		time.Sleep(100*time.Millisecond)
		fmt.Println(s)
	}
}

func main()  {
	defer fmt.Println("cpus:", runtime.NumCPU())
	fmt.Println("goroot:", runtime.GOROOT())
	fmt.Println("archive:",runtime.GOOS)

	go say("world")
	say("hello")
}


