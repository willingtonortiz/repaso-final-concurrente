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

// var maxConsumer int

func fillProducer(items int, proCh chan bool) <-chan bool {
	await := make(chan bool)

	go func() {

		for i := 0; i < items; i++ {
			proCh <- true
		}

		close(await)
	}()

	return await
}

func producer(proCh, conCh chan bool) {
	for {

		if len(proCh) == 0 {
			randomValue := rand.Intn(5) + 5
			<-fillProducer(randomValue, proCh)

			fmt.Println("HACIENDO REFILL", len(proCh))
		}
	}

}

func consumer(proCh, conCh chan bool) {
	for {
		if rand.Int()%2 == 0 && len(conCh) > 0 {
			<-conCh

			fmt.Println("VENDI UN ITEM", len(conCh))
		}

		if len(conCh) < cap(conCh) {
			<-proCh

			conCh <- true
			fmt.Println("COMPRE")
		}
	}
}

func main() {

	// Setup
	INFINITO := 999999
	proCh := make(chan bool, INFINITO)
	conCh := make(chan bool, 10)

	<-fillProducer(10, proCh)

	fmt.Println("FILLED")

	go producer(proCh, conCh)
	go consumer(proCh, conCh)

	time.Sleep(time.Second * 100)
}
