package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Client struct {
	Name 	string
	Conn 	net.Conn
	Open	bool
}

func HandleClient(m *sync.Mutex, clients map[uuid.UUID]Client, clientId uuid.UUID, broadcast chan string) {
	currentClient := clients[clientId]

	defer func() {
		m.Lock()
		delete(clients, clientId) // from map, delete key
		m.Unlock()
		currentClient.Conn.Close()
	}()

	reader := bufio.NewReader(currentClient.Conn)

	name, err := reader.ReadString('\n')

	if err!=nil{
		if err == io.EOF {
			fmt.Printf("Client %s exited\n", currentClient.Conn.RemoteAddr())
			return
		}
		fmt.Println("Error reading from client:", err)
		return
	}

	// trimming to avoid leading whitespace
	currentClient.Name = strings.TrimSpace(name)
	
	m.Lock()
	clients[clientId] = currentClient // updating client with name
	m.Unlock()

	fmt.Fprintf(currentClient.Conn, "Hello, you are connected to %s\n", currentClient.Conn.LocalAddr())
	
	// sends a message to all clients that someone connected
	broadcast <- fmt.Sprintf(">>> %s entrou no chat\n", currentClient.Name)

	for {
		message, err := reader.ReadString('\n')

		if err!=nil{
			if err == io.EOF {
				fmt.Printf("Client %s exited\n", currentClient.Conn.RemoteAddr())
				broadcast <- fmt.Sprintf(">>> %s saiu do chat\n", currentClient.Name)
				return
			}
			fmt.Println("Error reading from client:", err)
			return
		}

		cleanedMsg := strings.TrimSpace(message)
		// broadcast to all
		broadcast <- fmt.Sprintf("[%s]: %s\n", currentClient.Name, cleanedMsg)

		fmt.Printf("Message received: %s from %s\n", cleanedMsg, currentClient.Conn.RemoteAddr())
	}
}
