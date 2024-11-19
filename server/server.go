package server

import (
	"fmt"
	"io"
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
	MaxClients = 10               
)

const (
	welcomeMsg = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'"
)

func StartupServer() net.Listener {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}

	port := "8989"
	if len(os.Args) == 2 {
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
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		if len(Clients) >= MaxClients {
			clientConn.Write([]byte("\nMax clients reached. Connection rejected.\n"))
			clientConn.Close()
			continue
		}

		go HandleClient(clientConn)
	}
}

func BackupHistory(clientConn net.Conn) {
	Mutex.Lock()
	defer Mutex.Unlock()
	for _, msg := range Messages {
		clientConn.Write([]byte(msg + "\n"))
	}
}

func IsAlphaNum(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}

func IsValidName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !IsAlphaNum(r) {
			return false
		}
	}
	return true
}

func GetClientName(clientConn net.Conn) string {
	clientConn.Write([]byte(welcomeMsg))
	buf := make([]byte, 256)

	for {
		clientConn.Write([]byte("\n[ENTER YOUR NAME]: "))
		n, err := clientConn.Read(buf)
		if err != nil {
			log.Printf("Error reading client name: %v\n", err)
			return ""
		}

		name := strings.TrimSpace(string(buf[:n]))
		if IsValidName(name) {
			return AssignRandomColor(name)
		}
		clientConn.Write([]byte("\nInvalid name. Names must be alphanumeric and non-empty. Please try again.\n"))
	}
}

func HandleClient(clientConn net.Conn) {
	defer func() {
		Mutex.Lock()
		delete(Clients, clientConn)
		Mutex.Unlock()
		clientConn.Close()
	}()

	name := GetClientName(clientConn)
	if name == "" {
		return
	}

	Mutex.Lock()
	Clients[clientConn] = name
	Mutex.Unlock()

	NotifyClients(name, "has joined the chat...")
	BackupHistory(clientConn)

	buf := make([]byte, 4096)
	for {

		n, err := clientConn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from client %s: %v\n", name, err)
			}
			break
		}

		clientConn.Write([]byte("\nenter your message: "))
		message := strings.TrimSpace(string(buf[:n]))
		if message != "" {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			SendClientMessage(name, message, timestamp)
		}
	}

	NotifyClients(name, "has left the chat...")
}
