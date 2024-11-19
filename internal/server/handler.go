package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"TCPChat/internal/utils"
)

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

func SendClientMessage(sender, message, timestamp string) {
	formattedMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, sender, message)
	resetter := "\r\033[K"
	Mutex.Lock()
	Messages = append(Messages, formattedMsg)
	fmt.Println(formattedMsg)
	for clientConn, clientName := range Clients {
		if clientName != sender {
			clientConn.Write([]byte(resetter + formattedMsg))
			clientConn.Write([]byte("\n[h1]enter your message: "))
		}
	}
	Mutex.Unlock()
}

func HandleClient(clientConn net.Conn) {
	defer func() {
		Mutex.Lock()
		delete(Clients, clientConn)
		clientConn.Write([]byte("your disconnected"))
		Mutex.Unlock()
		clientConn.Close()
	}()
	clientConn.Write([]byte(utils.WelcomeMsg))

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
