package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

// Message ...
type Message struct {
	Command  string `json:"command"`
	Hostname string `json:"hostName"`
	Number   int    `json:"number"`
}

// Variables globales
var end chan bool
var localPort string
var nextPort string

// Variables para el algoritmo

// Ejecuta el servidor y acepta requests
func startServer() {
	fmt.Println("(", localPort, ")")
	ln, _ := net.Listen("tcp", "localhost:"+localPort)
	defer ln.Close()

	for {
		conn, _ := ln.Accept()
		go handleRequeset(conn)
	}
}

// Recibe cualquier request
func handleRequeset(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var message Message

	if err := dec.Decode(&message); err == nil {
		switch message.Command {

		case "start game":
			startGameHandler(message)

		case "process number":
			processNumberHandler(message)

		// Finaliza el servidor
		case "finish":
			finishHandler()
		}

	}
}

/* ********** HANDLERS ********** */

func startGameHandler(message Message) {
	randomNumber := generateRandomNumber(20)
	fmt.Println(localPort, "generated", randomNumber)
	sendMessageToHost(nextPort, Message{Command: "process number", Number: randomNumber})
}

func processNumberHandler(message Message) {
	if message.Number == 0 {
		fmt.Println("PERDI", localPort)
	} else {
		fmt.Println(localPort, "got", message.Number, " => sending", message.Number-1)
		sendMessageToHost(nextPort, Message{Command: "process number", Number: message.Number - 1})
	}
}

func finishHandler() {
	fmt.Println(localPort, "that's all folks")
	end <- true
}

func sendMessageToHost(port string, message Message) {
	// Se establece una conexion con el host remoto

	conn, _ := net.Dial("tcp", "localhost:"+port)
	defer conn.Close()

	enc := json.NewEncoder(conn)

	if err := enc.Encode(&message); err == nil {
		fmt.Println("Sending", "'"+message.Command+"'", "to", port)
	}
}

/* ********** HELPER FUNCTIONS ********** */

func generateRandomNumber(max int) int {
	return rand.Intn(max)
}

/* ********** MAIN ********** */

func main() {
	end = make(chan bool)
	localPort = os.Args[1]

	rand.Seed(time.Now().UTC().UnixNano())

	go startServer()

	// Debo ser agregado a la red
	if len(os.Args) == 3 {
		knownPort := os.Args[2]

		nextPort = knownPort
	}

	<-end
}
