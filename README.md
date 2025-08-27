# Chat Server com Go

## Básicos

- Objetivo: Aprender a configurar um servidor TCP simples que manipula um cliente.

O que é um servidor TCP?

Primeiramente, TCP (Transmission Control Protocol) é um protocolo de comunicação que permite a transmissão de dados entre dispositivos em uma rede. Um servidor TCP é um programa que escuta por conexões de clientes em uma porta específica e estabelece uma comunicação bidirecional com esses clientes.

Para criar um servidor TCP em Go, vamos utilizar o pacote `net`. Da documentação do pacote `net`, podemos extrair essa descrição: "Package net provides a portable interface for network I/O, including TCP/IP, UDP, domain name resolution, and Unix domain sockets."

Com isso em mente, vamos ao código.

Iniciamos criando um simples servidor TCP

```go
// a função Listen cria um servidor
ln, err := net.Listen("tcp", ":8080")
if err!=nil{
    fmt.Println("There was an error trying to create a TCP server: ", err)
    return
}
defer ln.Close()
```

`defer ln.Close()`: Isso garante que o servidor será fechado corretamente quando a função principal terminar sua execução.

Depois, precisamos aceitar conexões de clientes. Vamos começar com uma única conexão

```go
conn, err := ln.Accept()
if err!=nil{
    fmt.Println("There was an error trying to accept connections in the server: ", err)
    return
}
defer conn.Close()

fmt.Println("Client connected:", conn.RemoteAddr())
```

`ln.Accept()`: Esta função **bloqueia** até que um cliente se conecte ao servidor. Quando um cliente se conecta, ela retorna um objeto `net.Conn`, que representa a conexão com o cliente. Também retorna um erro, que deve ser tratado adequadamente.
`defer conn.Close()`: Isso garante que a conexão com o cliente será fechada corretamente quando a função principal terminar sua execução.

Agora vamos permitir que o cliente envie uma mensagem para o servidor

```go
status, err := bufio.NewReader(conn).ReadString('\n')
if err!=nil{
    fmt.Println("Error reading from client:", err)
    return
}

fmt.Println("Message received:", status)
```

`bufio.NewReader(conn).ReadString('\n')`: Esta linha cria um novo leitor de buffer para a conexão e lê uma string até encontrar um caractere de nova linha (`\n`). Isso é útil para ler mensagens delimitadas por novas linhas.
`status`: Esta variável armazena a mensagem recebida do cliente.

Com o servidor pronto, já podemos colocá-lo em execução. É importante lembrar que ele funciona apenas como servidor: para enviar mensagens, precisamos de um cliente que se conecte a ele.
Para testar, abra dois terminais:

- no primeiro, rode o servidor;
- no segundo, conecte-se a ele com o comando `nc localhost 8080`.

Esse comando estabelece a conexão com o servidor e mantém o terminal aguardando até que uma mensagem seja digitada. Assim que o cliente enviar algo, o servidor exibirá a mensagem no seu próprio terminal

## Múltiplos clientes

- Objetivo: Permitir que múltiplos clientes se conectem e se comuniquem com o servidor simultaneamente.

Para permitir que múltiplos clientes se conectem ao servidor simultaneamente, precisamos modificar o código para aceitar conexões em um loop infinito e lidar com cada conexão em uma goroutine separada. Isso permite que o servidor continue aceitando novas conexões enquanto processa as mensagens dos clientes existentes.

```go
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
```

E em `internal/client/client.go`:

```go
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
```

Note que envolvendo a seção do Reader temos um `for` que permite ler mensagens continuamente até que a conexão seja fechada.
E agora, já que temos múltiplos clientes, precisamos garantir que as mensagens sejam enviadas para todos os clientes conectados. Isso já está sendo feito na função `HandleClient`, onde usamos um loop para enviar a mensagem recebida para todos os outros clientes no slice de conexões.

## Transmissão de Mensagens

- Objetivo: Implementar um sistema de broadcast de mensagens usando channels em Go para melhor gerenciamento de concorrência.

Para implementar um sistema de broadcast eficiente, vamos usar um channel em Go, que é a forma idiomática de comunicação entre goroutines.

> Channels são um mecanismo de comunicação que permite que goroutines troquem dados de forma segura e sincronizada.

```go
// No arquivo cmd/server/main.go
var (
    clients   = make(map[uuid.UUID]client.Client) // clientes conectados
    broadcast = make(chan string)                 // canal para mensagens
    m         sync.Mutex                          // protege o mapa de clientes
)
```

A função `handleBroadcast` é executada em uma goroutine separada e é responsável por receber mensagens do canal e enviá-las para todos os clientes conectados:

```go
// handleBroadcast gerencia o envio de mensagens para todos os clientes
func handleBroadcast() {
    for {
        // Espera por mensagens no canal broadcast
        message := <-broadcast

        // Envia a mensagem para todos os clientes conectados
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
```

Quando um cliente envia uma mensagem, a função `HandleClient` a coloca no canal de broadcast:

```go
// No arquivo client.go
broadcast <- fmt.Sprintf("[%s]: %s\n", currentClient.Name, cleanedMsg)
```

Esta implementação tem várias vantagens:

1. **Concorrência segura**: O uso do mutex garante que não haja condições de corrida no acesso ao mapa de clientes.
2. **Eficiência**: O channel permite que a goroutine de broadcast processe mensagens quando elas estão disponíveis.
3. **Desacoplamento**: A lógica de processamento de mensagens está separada da lógica de comunicação com o cliente.

Para testar o sistema de broadcast, conecte vários clientes ao servidor usando telnet ou netcat:

```bash
nc localhost 8080
```

O primeiro texto enviado será considerado o nome do usuário, e todas as mensagens subsequentes serão transmitidas a todos os usuários conectados, com o formato `[Nome]: Mensagem`.
