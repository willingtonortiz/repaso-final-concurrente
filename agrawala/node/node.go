package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

/*
[X] HACER QUE TODOS SE COMUNIQUEN
[ ] HACER QUE TODOS GENEREN UN NUMERO ALEATORIO Y SE LO COMPARTAN
*/

// Message ...
type Message struct {
	Command  string   `json:"command"`
	Hostname string   `json:"hostName"`
	List     []string `json:"list"`
	Number   int      `json:"number"`
}

// Variables globales
var end chan bool
var localPort string
var friendList []string

// Variables para el algoritmo
var hostNumbers map[string]int
var mutex = &sync.Mutex{}

/* ********** HELPER FUNCTIONS ********** */
func generateRandomNumber(max int) int {
	return rand.Intn(max)
}

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

		// Si se recibe un comando 'meet new friend'
		case "meet new friend":
			meetNewFriendHandler(message)

		// Si se recibe un comando 'finish'
		case "finish":
			finishHandler()

		case "start agrawala":
			startAgrawalaHandler(message)

		case "save number":
			saveNumberHandler(message)

		case "execute":
			executeHandler(message)
		}
	}
}

/* ********** HANDLERS ********** */

func executeHandler(message Message) {
	minNumber := hostNumbers[localPort]
	nextHost := ""
	posibleNumber := 9999999999

	for host, num := range hostNumbers {
		if minNumber < num {
			if num < posibleNumber {
				posibleNumber = num
				nextHost = host
			}
		}
	}

	if nextHost == "" {
		fmt.Println("Soy el ultimo")

		for _, friend := range friendList {
			fmt.Println(localPort, "sent finish to", friend)

			sendMessageToHost(friend, Message{Command: "finish"})
		}

		finishHandler()

		return
	}

	sendMessageToHost(nextHost, Message{Command: "execute"})
}

func saveNumberHandler(message Message) {
	mutex.Lock()

	hostNumbers[message.Hostname] = message.Number

	if len(hostNumbers) == len(friendList)+1 {

		minNumber := 999999
		prevNumber := 999999
		prevHost := ""
		nextHost := ""

		for host, num := range hostNumbers {
			if minNumber > num {
				prevHost = nextHost
				nextHost = host
				minNumber = num
			} else {
				if prevNumber > num {
					prevNumber = num
					prevHost = host
				}
			}
		}

		if minNumber == hostNumbers[localPort] {
			fmt.Println("Soy el primero", localPort, "->", prevHost)

			sendMessageToHost(prevHost, Message{Command: "execute"})
		}
	}

	mutex.Lock()
}

func startAgrawalaHandler(message Message) {
	randomNumber := rand.Int() % 9999
	hostNumbers[localPort] = randomNumber

	fmt.Println(localPort, "mi numero es", randomNumber)

	for _, friend := range friendList {
		fmt.Println(localPort, "sent", "save number", "to", friend)

		sendMessageToHost(friend, Message{Command: "save number", Number: randomNumber, Hostname: localPort})
	}

}

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

/* ********** MAIN ********** */

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	end = make(chan bool)
	localPort = os.Args[1]

	// Setup
	hostNumbers = make(map[string]int)

	go startServer()

	// Debo ser agregado a la red
	if len(os.Args) == 3 {
		knownPort := os.Args[2]
		friendList = append(friendList, knownPort)

		sendMessageToHost(knownPort, Message{Command: "hello", Hostname: localPort})
	}

	<-end
}
