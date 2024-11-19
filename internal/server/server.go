package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"TCPChat/internal/utils"
)

var (
	Clients    = make(map[net.Conn]string)
	Messages   = []string{}
	Mutex      = &sync.Mutex{}
	MaxClients = 10
	portHelp   = `
Allowed Ports
-Privileged Ports (0–1023):

	-> Typically require root/administrator permissions to bind.
	-> Reserved for well-known services (e.g., HTTP on 80, HTTPS on 443, SSH on 22).
	-> Avoid unless your application is implementing a standard service.


-Registered Ports (1024–49151):

	-> Suitable for most server applications.
	-> These ports are not restricted and are widely used for custom applications (e.g., 8080 for web development, 3306 for MySQL).


-Dynamic/Private Ports (49152–65535):

	-> Primarily used for client-side ephemeral connections.
	-> Servers generally avoid binding to these ports unless required by specific use cases.
`
)

func Run() {
	listener := startupServer()
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

func getPort() string {
	if len(os.Args) == 2 {
		pt, err := strconv.Atoi(os.Args[1])
		if err != nil || pt < 1024 || pt > 49151 {
			fmt.Println("Invalid port")
			fmt.Println(portHelp)
			os.Exit(0)
		}
		return os.Args[1]
	}
	return "8989"
}

func startupServer() net.Listener {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}

	port := getPort()

	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	utils.PrintBanner(port)
	return listener
}
