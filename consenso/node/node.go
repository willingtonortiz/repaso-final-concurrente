package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
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

func generateRandomNumber(max int) int {
	return rand.Intn(max)
}

// Recibe cualquier request
func handleRequeset(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var message Message

	if err := dec.Decode(&message); err == nil {
		switch message.Command {
		case "generateNumber":
			randomNumber := generateRandomNumber(2)

			if randomNumber == 0 {
				oneCounter++
			} else {
				secondCounter++
			}

			fmt.Println("Decidí", randomNumber)

			for _, friend := range friendList {
				fmt.Println(localPort, "sending", randomNumber, "to", friend)

				sendMessageToHost(friend, Message{Command: "handleNumber", Number: randomNumber})
			}

		case "handleNumber":
			if message.Number == 0 {
				oneCounter++
			} else {
				secondCounter++
			}

			if oneCounter+secondCounter == len(friendList) {
				if oneCounter > secondCounter {
					fmt.Println("Gano el primero")
				} else if oneCounter < secondCounter {
					fmt.Println("Gano el segundo")
				} else {
					fmt.Println("Empate!")
				}
			}

		case "hello":
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

		// Inicia el consenso
		case "test consensus":
			if rand.Intn(100)%2 == 0 {
				decisions[localPort] = "atacar"
			} else {
				decisions[localPort] = "retirada"
			}
			fmt.Println(localPort, "decidio", decisions[localPort])

			counter = 0

			for _, friend := range friendList {
				sendMessageToHost(friend, Message{Command: "decision", Decision: decisions[localPort]})
			}

			readyToListen <- true

		// Recibe una decision
		case "decision":
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
			}

		// Si se recibe un comando 'meet new friend'
		case "meet new friend":
			// Agrega un nuevo amigo de un host conocido
			friendList = append(friendList, message.Hostname)
			fmt.Println("Friend list updated:", friendList)

		// Si se recibe un comando 'finish'
		case "finish":
			fmt.Println(localPort, "that's all folks")
			end <- true
		}

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

func main() {
	end = make(chan bool)
	localPort = os.Args[1]

	go startServer()

	// Debo ser agregado a la red
	if len(os.Args) == 3 {
		knownPort := os.Args[2]
		friendList = append(friendList, knownPort)

		sendMessageToHost(knownPort, Message{Command: "hello", Hostname: localPort})
	}

	<-end
}
