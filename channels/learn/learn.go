/*

Un proveedor inicialmente tiene 10 objetos de un stock ilimitado y desea darle objetos a un consumidor cuando este los requiera.
Cuando el proveedor se queda sin objetos y el consumidor requiera, entonces el proveedor hará un restock de entre 5 a 30 objetos
aleatoriamente y entregará 1 al consumidor.

El consumidor inicialmente no tiene objetos, y tiene una capacidad de 10 objetos como máximo. Cada iteración aleatoriamente
venderá 1 objeto y siempre pedirá objetos si la cantidad de objetos que tiene no es el límite máximo de su capacidad.

*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

var producerStock int
var consumerStock int

func reFill(proCh chan int) {

	producerStock = rand.Intn(25) + 5

}

func producer(proCh, conCh chan bool) {

}

func consumer(proCh, conCh chan bool) {
	for {
		if rand.Int()%2 == 0 && producerStock != 0 {
			producerStock--
			fmt.Println("VENDI UN OBJETO")
		}

		if producerStock < 10 {
			<-proCh
			producerStock++
		}
	}
}

func other(channel chan int) {
	fmt.Println("OTHER START")

	channel <- 10
	channel <- 20
	channel <- 30
	channel <- 40
	channel <- 50
	channel <- 60

	fmt.Println(<-channel)
	fmt.Println("OTHER FINISH")
}

func main() {

	// Setup
	// producerStock = 10
	// consumerStock = 0

	// proCh := make(chan bool)
	// conCh := make(chan bool)

	// producer(proCh, conCh)
	// consumer(proCh, conCh)

	channel := make(chan int, 10)

	go other(channel)

	fmt.Println("MAIN START")

	fmt.Println(<-channel)
	fmt.Println(<-channel)
	fmt.Println(<-channel)
	fmt.Println(<-channel)
	fmt.Println(<-channel)
	fmt.Println("TUGFA")

	fmt.Println("MAIN FINISH")
	time.Sleep(time.Second * 1)
}
