package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func HandleClient(conn net.Conn, clients *[]net.Conn) {
	defer conn.Close()

	fmt.Fprintf(conn, "Hello, you are connected to %s\n", conn.LocalAddr())

	for {
		status, err := bufio.NewReader(conn).ReadString('\n')

		if err!=nil{
			fmt.Println("Error reading from client:", err)
			return
		}

		// broadcast to all
		for _, client := range *clients {
			if client != conn {
				fmt.Fprintf(client, "Connection %s said: %s\n", conn.RemoteAddr(), status)
			}
		}

		fmt.Fprintf(os.Stdin, "Message received: %s from %s\n", status, conn.RemoteAddr())
	}
}