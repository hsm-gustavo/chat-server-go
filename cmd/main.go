package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// a função Listen cria um servidor
	ln, err := net.Listen("tcp", ":8080")
	if err!=nil{
		fmt.Println("There was an error trying to create a TCP server: ", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080")

	conn, err := ln.Accept()
	if err!=nil{
		fmt.Println("There was an error trying to accept connections in the server: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())
	
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err!=nil{
		fmt.Println("Error reading from client:", err)
		return
	}

	fmt.Println("Message received:", status)
		
}