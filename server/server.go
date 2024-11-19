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
	fmt.Println(listener.Addr())
	fmt.Println(listener.Addr().Network())
	fmt.Println(listener.Addr().String())
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
			clientConn.Write([]byte("Max clients reached. Connection rejected.\n"))
			clientConn.Close()
			continue
		}

		go HandleClient(clientConn)
	}
}

func BackupHistory(clientConn net.Conn) {
	Mutex.Lock()
	defer Mutex.Unlock()
	clientConn.Write([]byte("\n"))
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

func IsClientNameExists(name string) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	for _, clientName := range Clients {
		if clientName[5:10] == name {
			fmt.Printf("\n[%s] | [%s]\n", name, clientName[5:10])
			return true
		}
	}
	return false
}

func GetClientName(clientConn net.Conn) string {
	clientConn.Write([]byte(welcomeMsg))
	buf := make([]byte, 256)

	for {
		clientConn.Write([]byte("\n[ENTER YOUR NAME]: "))
		n, err := clientConn.Read(buf)
		if err != nil {
			log.Printf("Error reading client name: %v", err)
			return ""
		}

		name := strings.TrimSpace(string(buf[:n]))
		if !IsValidName(name) {
			clientConn.Write([]byte("Invalid name. Names must be alphanumeric and non-empty. Please try again."))
		}else if IsClientNameExists(name) {
			clientConn.Write([]byte("\033[31m"+ name + "\033[0m already taken, try again"))
		}else {
			return AssignColor(name)
		}
		// if GetClientName(clientConn) == ""
		
	}
}

func HandleClient(clientConn net.Conn) {
	defer func() {
		Mutex.Lock()
		delete(Clients, clientConn)
		clientConn.Write([]byte("your disconnected"))
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

	NotifyClients(name, " has joined the chat...", clientConn)
	BackupHistory(clientConn)
	clientConn.Write([]byte("\n[s1]enter your message: "))

	buf := make([]byte, 4096)
	for {

		n, err := clientConn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from client %s: %v", name, err)
			}
			break
		}

		message := strings.TrimSpace(string(buf[:n]))
		if message != "" {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			SendClientMessage(name, message, timestamp)
		}
		clientConn.Write([]byte("[s2]enter your message: "))

	}

	NotifyClients(name, " has left the chat...", clientConn)
}
