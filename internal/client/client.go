package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func HandleClient(conn net.Conn, clients *[]net.Conn) {
	defer conn.Close()

	fmt.Fprintf(conn, "Hello, you are connected to %s\n", conn.LocalAddr())

	for {
		status, err := bufio.NewReader(conn).ReadString('\n')

		if err!=nil{
			if err == io.EOF {
				fmt.Printf("Client %s exited\n", conn.RemoteAddr())
				return
			}
			fmt.Println("Error reading from client:", err)
			return
		}

		cleanedMsg := strings.ReplaceAll(status, "\n","")
		// broadcast to all
		broadcastMessage(clients, conn, cleanedMsg)

		fmt.Fprintf(os.Stdin, "Message received: %s from %s\n", status, conn.RemoteAddr())
	}
}

func broadcastMessage(clients *[]net.Conn, currentConn net.Conn, message string) {
	for _, client := range *clients {
		if client != currentConn {
			fmt.Fprintf(client, "Connection %s said: %s\n", currentConn.RemoteAddr(), message)
		}
	}
}