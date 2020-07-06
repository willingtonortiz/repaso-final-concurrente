package main

import (
	"fmt"
	"sync"
)

var counter int

func incrementCounter(mux *sync.Mutex) {
	mux.Lock()
	defer mux.Unlock()
	counter++
}

func printCounter(mux *sync.Mutex) {
	mux.Lock()
	defer mux.Unlock()
	fmt.Printf("Counter=%v\n", counter)
}

func main() {
	counter = 0
	mux := &sync.Mutex{}

	for i := 0; i < 1000; i++ {
		go incrementCounter(mux)
	}

	printCounter(mux)

}
