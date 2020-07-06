package main

import (
	"fmt"
	"time"
)

func simpleChannel() {
	c := make(chan int)

	go func(ch chan<- int, x int) {
		time.Sleep(time.Second)
		ch <- x * x
	}(c, 3)

	done := make(chan int32)

	go func(ch <-chan int) {

		n := <-ch
		fmt.Println(n)

		time.Sleep(time.Second)

		done <- 10
	}(c)

	<-done

	fmt.Println("BYE")
}

func bufferedChannel() {
	c := make(chan int, 2)
	c <- 3
	c <- 5

	close(c)

	fmt.Println(len(c), cap(c))

	x, ok := <-c
	fmt.Println(x, ok)
	fmt.Println(len(c), cap(c))

	x, ok = <-c
	fmt.Println(x, ok)
	fmt.Println(len(c), cap(c))

	x, ok = <-c
	fmt.Println(x, ok)

	x, ok = <-c
	fmt.Println(x, ok)
	fmt.Println(len(c), cap(c))

	// No se puede cerrar un canal ya cerrado
	// close(c)

	// No se puede enviar un valor a un canal cerrado
	// c <- 7
}

func forChannel() {
	ch := make(chan int)

	go func(n int, c chan int) {
		x, y := 0, 1

		for i := 0; i < n; i++ {
			c <- x
			x, y = y, x+y
			time.Sleep(time.Second)
		}

		close(c)
	}(5, ch)

	for i := range ch {
		fmt.Println(i)
	}
}

func main() {
	// simpleChannel()
	// bufferedChannel()
	forChannel()

}
