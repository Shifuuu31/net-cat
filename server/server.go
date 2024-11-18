package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Clients    = make(map[net.Conn]string)
	Messages   = []string{}
	Mutex      = &sync.Mutex{}
	MaxClients = 4
)

const (
	welcomeMsg = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "
)

func StartupServer() net.Listener {
	argCount := len(os.Args[1:])
	if argCount > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}
	port := "8989"
	if argCount == 1 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	PrintBanner(port)
	return listener
}

func Run() {
	listener := StartupServer()
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		if len(Clients) >= MaxClients {
			fmt.Fprintf(connection, "Max Clients reached. New connections will be rejected.")
			connection.Close()
			continue
		}

		go HandleClient(connection)
	}
}

func HandleClient(connection net.Conn) {
	defer connection.Close()

	connection.Write([]byte(welcomeMsg))

	// fmt.Fprintf(connection, welcomeMsg)

	scanner := bufio.NewScanner(connection)
	scanner.Scan()

	name := AssignRandomColor(strings.TrimSpace(scanner.Text()))
	if name == "" {
		connection.Write([]byte("Name cannot be empty. Please cisConnectionect and try again.\n"))
		return
	}
	Mutex.Lock()
	Clients[connection] = name
	Mutex.Unlock()

	NotifyClients(name + " has joined the chat...")

	Mutex.Lock()
	for _, msg := range Messages {
		fmt.Fprintf(connection, "%s\n", msg)
	}
	Mutex.Unlock()

	for scanner.Scan() {
		message := scanner.Text()
		if message == "" {
			continue
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, name, message)

		Mutex.Lock()
		Messages = append(Messages, formattedMessage)
		Mutex.Unlock()

		NotifyClients(formattedMessage)
	}

	Mutex.Lock()
	delete(Clients, connection)
	Mutex.Unlock()

	NotifyClients(name + " has left the chat...")
}
