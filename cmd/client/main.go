package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run cmd/client/main.go <seu-nome>")
		return
	}

	name := os.Args[1]

	// connects to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	// sends name to the server
	fmt.Fprintf(conn, "%s\n", name)

	// goroutine to read server messages
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Desconectado do servidor")
				os.Exit(0)
			}
			fmt.Print(message)
		}
	}()

	// reads user messages and sends to the server
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Conectado ao servidor. Digite suas mensagens:")
	for scanner.Scan() {
		message := scanner.Text()
		if strings.ToLower(strings.TrimSpace(message)) == "/sair" {
			fmt.Println("Saindo...")
			break
		}
		fmt.Fprintf(conn, "%s\n", message)
	}
}
