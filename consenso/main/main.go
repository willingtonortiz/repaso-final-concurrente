package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Message struct {
	Command  string
	Hostname string
	List     []string
}

func send(remote string, message Message) {
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()

	enc := json.NewEncoder(conn)
	if err := enc.Encode(&message); err == nil {
		fmt.Println("sending test consensus to", remote)
	}
}

func main() {
	for _, remote := range os.Args[1:] {
		send("localhost:"+remote, Message{Command: "test consensus"})
	}
}
