package server

import (
	"fmt"
	"net"
	"sync"
)

var (
	colors     = []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	colorMutex = &sync.Mutex{}
)

func AssignColor(input string) string {
	colorMutex.Lock()
	defer colorMutex.Unlock()
	n := len(Clients)
	if n+1 > len(colors) {
		return input
	}
	return fmt.Sprintf("%s%s\033[0m", colors[n], input)
}

func SendClientMessage(sender, message, timestamp string) {
	formattedMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, sender, message)
	resetter := "\r\033[K"
	Mutex.Lock()
	Messages = append(Messages, formattedMsg)
	fmt.Println(formattedMsg)
	for clientConn, clientName := range Clients {
		if clientName != sender {
			clientConn.Write([]byte(resetter+formattedMsg))
			clientConn.Write([]byte("\n[h1]enter your message: "))
		}
		// fmt.Println(clientConn.LocalAddr())
	}
	Mutex.Unlock()
}

func NotifyClients(name, action string, restrictConn net.Conn) {
	notification := "\r\033[K" + name + action
	Mutex.Lock()
	defer Mutex.Unlock()

	for clientConn := range Clients {
		clientConn.Write([]byte(notification))
		if clientConn != restrictConn {
			clientConn.Write([]byte("\r\033[K[h2|r]enter your message: "))
		}
	}
}

func PrintBanner(port string) {
	fmt.Printf("\033[1;32m\nTCPChat SERVER STATUS: 🟢 LIVE\n-> Server started on: localhost:%s\n-> Listening on the port: %s\033[0m\n", port, port)
}
