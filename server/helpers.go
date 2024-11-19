package server

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	colors         = []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	assignedColors = map[string]string{}
	colorMutex     = &sync.Mutex{}
)

// AssignRandomColor assigns a unique color to a client
func AssignRandomColor(input string) string {
	colorMutex.Lock()
	defer colorMutex.Unlock()

	if color, exists := assignedColors[input]; exists {
		return fmt.Sprintf("%s%s\033[0m", color, input)
	}

	rand.Seed(time.Now().UnixNano())
	color := colors[rand.Intn(len(colors))]
	assignedColors[input] = color
	return fmt.Sprintf("%s%s\033[0m", color, input)
}

// SendClientMessage broadcasts a message to all clients except the sender
func SendClientMessage(sender, message, timestamp string) {
	formattedMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, sender, message)

	Mutex.Lock()
	Messages = append(Messages, formattedMsg)
	for clientConn, clientName := range Clients {
		if clientName != sender {
			clientConn.Write([]byte(formattedMsg+"\n"))
		}
	}
	Mutex.Unlock()
}

// NotifyClients sends a notification (e.g., join/leave messages) to all clients
func NotifyClients(name, action string) {
	notification := fmt.Sprintf("\n%s %s\n", name, action)
	Mutex.Lock()
	defer Mutex.Unlock()

	for clientConn := range Clients {
		clientConn.Write([]byte(notification))
	}
}

// PrintBanner displays a startup banner
func PrintBanner(port string) {
	fmt.Printf("\033[1;32m\nTCPChat SERVER STATUS: 🟢 LIVE\n-> Server started on: localhost:%s\n-> Listening on the port: %s\033[0m\n", port, port)
}
