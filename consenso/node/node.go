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
	Command  string   `json:"command"`
	Hostname string   `json:"hostName"`
	List     []string `json:"list"`
	Number   int      `json:"number"`
	Decision string   `json:"decision"`
}

// Variables globales
var end chan bool
var readyToListen chan bool
var localPort string
var friendList []string

// Variables para el algoritmo
var oneCounter int
var secondCounter int
var decisions map[string]string
var counter int

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

		case "hello":
			helloRequestHandler(conn, message)

		// Inicia el consenso
		case "test consensus":
			testConsensushandler(message)

		// Recibe una decision
		case "decision":
			decisionHandler(message)

		// Si se recibe un comando 'meet new friend'
		case "meet new friend":
			meetNewFriendHandler(message)

		// Si se recibe un comando 'finish'
		case "finish":
			finishHandler()
		}

	}
}

/* ********** HANDLERS ********** */

func helloRequestHandler(conn net.Conn, message Message) {
	// Le envio mi mista de amigos al nuevo host
	resp := Message{Command: "hey", Hostname: localPort, List: friendList}
	enc := json.NewEncoder(conn)

	// A cada amigo mio, le presento el nuevo host
	if err := enc.Encode(&resp); err == nil {

		for _, friend := range friendList {
			fmt.Println(localPort, "introduces", message.Hostname, "to", friend)

			// Presentar nuevo amigo a mis amigos
			sendMessageToHost(friend, Message{Command: "meet new friend", Hostname: message.Hostname})
		}
	}

	// Agrego al nuevo host
	friendList = append(friendList, message.Hostname)
	fmt.Println("Friend list updated:", friendList)
}

func testConsensushandler(message Message) {
	if rand.Intn(100)%2 == 0 {
		decisions[localPort] = "atacar"
	} else {
		decisions[localPort] = "retirada"
	}
	fmt.Println(localPort, "decidio", decisions[localPort])

	counter = 0

	for _, friend := range friendList {
		sendMessageToHost(friend, Message{Command: "decision", Hostname: localPort, Decision: decisions[localPort]})
	}

	readyToListen <- true
}

func decisionHandler(message Message) {
	<-readyToListen

	decisions[message.Hostname] = message.Decision
	counter++

	if counter == len(friendList) {
		attackCounter := 0
		fallCounter := 0

		for _, decision := range decisions {
			if decision == "atacar" {
				attackCounter++
			} else {
				fallCounter++
			}
		}

		if attackCounter > fallCounter {
			fmt.Println(localPort, "ATACAR!!!")
		} else {
			fmt.Println(localPort, "RETIRADA!!!")
		}

		end <- true
	} else {
		// Se debe aceptar el siguiente request
		readyToListen <- true
	}
}

func sendMessageToHost(port string, message Message) {
	// Se establece una conexion con el host remoto

	conn, _ := net.Dial("tcp", "localhost:"+port)
	defer conn.Close()

	enc := json.NewEncoder(conn)

	if err := enc.Encode(&message); err == nil {
		fmt.Println("Sending", "'"+message.Command+"'", "to", port)

		// Si el comando es hello, entonces se espera recibir una respuesta
		// con la lista de nodos en la red
		if message.Command == "hello" {
			dec := json.NewDecoder(conn)
			var response Message

			if err := dec.Decode(&response); err == nil {
				friendList = append(friendList, response.List...)

				fmt.Println("Receiving", response.List, "from", port)
			}
		}
	}
}

func meetNewFriendHandler(message Message) {
	// Agrega un nuevo amigo de un host conocido
	friendList = append(friendList, message.Hostname)
	fmt.Println("Friend list updated:", friendList)
}

func finishHandler() {
	fmt.Println(localPort, "that's all folks")
	end <- true
}

/* ********** HELPER FUNCTIONS ********** */

func generateRandomNumber(max int) int {
	return rand.Intn(max)
}

/* ********** MAIN ********** */

func main() {
	end = make(chan bool)
	readyToListen = make(chan bool)
	decisions = make(map[string]string)
	localPort = os.Args[1]

	rand.Seed(time.Now().UTC().UnixNano())

	go startServer()

	// Debo ser agregado a la red
	if len(os.Args) == 3 {
		knownPort := os.Args[2]
		friendList = append(friendList, knownPort)

		sendMessageToHost(knownPort, Message{Command: "hello", Hostname: localPort})
	}

	<-end
}
