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

func NotifyClients(name, action, timestamp string, restrictConn net.Conn) {
	notification := "\r\033[K" + name + action
	Mutex.Lock()
	defer Mutex.Unlock()

	for clientConn := range Clients {
		clientConn.Write([]byte(notification))
		if clientConn != restrictConn {
			clientConn.Write([]byte("\r\033[K[" + timestamp + "][" + name + "]: "))
		}
	}
}

func SendClientMessage(name, message, timestamp string) {
	formattedMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, name, message)
	resetter := "\r\033[K"
	Mutex.Lock()
	Messages = append(Messages, formattedMsg)
	fmt.Println(formattedMsg)
	for clientConn, clientName := range Clients {
		if clientName != name {
			clientConn.Write([]byte(resetter + formattedMsg))
			clientConn.Write([]byte("\n["+timestamp+"]["+name+"]: "))
		}
	}
	Mutex.Unlock()
}

func HandleClient(clientConn net.Conn) {
	defer func() {
		Mutex.Lock()
		delete(Clients, clientConn)
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
	// timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	NotifyClients(name, " has joined the chat...", time.Now().Format("2006-01-02 15:04:05"), clientConn)
	BackupHistory(clientConn)
	clientConn.Write([]byte("\n["+time.Now().Format("2006-01-02 15:04:05")+"]["+name+"]: "))

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
			SendClientMessage(name, message, time.Now().Format("2006-01-02 15:04:05"))
		}
		clientConn.Write([]byte("["+time.Now().Format("2006-01-02 15:04:05")+"]["+name+"]: "))

	}

	NotifyClients(name, " has left the chat...", time.Now().Format("2006-01-02 15:04:05"), clientConn)
}
