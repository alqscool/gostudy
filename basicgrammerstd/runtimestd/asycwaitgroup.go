package main

import(
	"fmt"
	"sync"
)

var waitgroup sync.WaitGroup

func display(num int) {
	fmt.Println(num)
	defer waitgroup.Done()
}

func main()  {
	for i:=0; i<10; i++ {
		waitgroup.Add(1)
		go display(i)
	}

	waitgroup.Wait()
	fmt.Println("done!")
}

