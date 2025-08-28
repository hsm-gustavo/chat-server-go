package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/hsm-gustavo/chat-server-go/internal/client"
)

var (
	clients   = make(map[uuid.UUID]client.Client) // connected clients
	broadcast = make(chan string)		// channel for messages
	m		  sync.Mutex				// protect client map
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	
	if err!=nil{
		fmt.Println("There was an error trying to create a TCP server: ", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080")
	
	// starts a goroutine to manage broadcast messages
	go handleBroadcast()

	for {
		conn, err := ln.Accept()
		if err!=nil{
			fmt.Println("There was an error trying to accept connections in the server: ", err)
			return
		}

		clientId := uuid.New()
		newClient := client.Client{Name: "", Conn: conn, Open: true}

		m.Lock()
		clients[clientId] = newClient
		m.Unlock()

		fmt.Println("Client connected:", conn.RemoteAddr())

		go client.HandleClient(&m, clients, clientId, broadcast)
	}
}

// handles sending messages to all clients
func handleBroadcast() {
	for {
		// wait for messages in broadcast channel
		message := <-broadcast
		
		// sends to all connected clients
		m.Lock()
		for _, client := range clients {
			if client.Open {
				_, err := fmt.Fprint(client.Conn, message)
				if err != nil {
					client.Open = false
					client.Conn.Close()
				}
			}
		}
		m.Unlock()
	}
}