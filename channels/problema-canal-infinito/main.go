package main

import (
	"fmt"
	"time"
)

func main() {

	channel := make(chan int, 20)

	go func() {
		channel <- 10
		fmt.Println("GAGAGA")

		channel <- 10
		fmt.Println("GAGAGA")
	}()

	time.Sleep(time.Second * 3)
	<-channel
	time.Sleep(time.Second * 3)
	<-channel
	time.Sleep(time.Second * 1)
}
