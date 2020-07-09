package main

import (
	"fmt"
	"math/rand"
	"time"
)

func node(id int, in, out chan int) {

	var myNumber int = 999999

	for val := range in {
		// fmt.Println("Evaluando", val)

		if val < myNumber {
			out <- myNumber
			myNumber = val

			// fmt.Println(id, "actualizo", myNumber)
		} else {
			out <- val
		}

	}
	close(out)
	time.Sleep(time.Millisecond * 10 * time.Duration(id))
	fmt.Println(id, "=>", myNumber)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	size := 10

	channels := make([]chan int, size+1)

	for i := 0; i < size+1; i++ {
		channels[i] = make(chan int)
	}

	for i := 1; i < size+1; i++ {
		go node(i, channels[i-1], channels[i])
	}

	// 0 => 0 , 1
	// 1 => 1 , 2
	// 2 => 2 , 3

	// main -> () -> () -> () -> main

	for i := 0; i < size; i++ {
		randomNumber := rand.Int() % 99999
		channels[0] <- randomNumber
	}
	close(channels[0])

	for range channels[size] {
	}

	time.Sleep(time.Second * 2)
}
