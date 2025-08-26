package main

import (
	"fmt"
	"net"

	"github.com/hsm-gustavo/chat-server-go/internal/client"
)

func main() {
	var clients []net.Conn

	ln, err := net.Listen("tcp", ":8080")
	if err!=nil{
		fmt.Println("There was an error trying to create a TCP server: ", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		conn, err := ln.Accept()
		if err!=nil{
			fmt.Println("There was an error trying to accept connections in the server: ", err)
			return
		}

		clients = append(clients, conn)

		fmt.Println("Client connected:", conn.RemoteAddr())

		go client.HandleClient(conn, &clients)
	}
		
}