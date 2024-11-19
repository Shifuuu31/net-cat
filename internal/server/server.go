package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"TCPChat/internal/utils"
)

var (
	Clients    = make(map[net.Conn]string)
	Messages   = []string{}
	Mutex      = &sync.Mutex{}
	MaxClients = 10
)

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

	utils.PrintBanner(port)
	return listener
}
