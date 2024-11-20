package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"TCPChat/internal/utils"
)

func NotifyClients(name, action, timestamp string, restrictConn net.Conn) {
	notification := "\r\033[K" + name + action + "\n"
	Mutex.Lock()
	defer Mutex.Unlock()

	for clientConn, clientName := range Clients {
		_, err := clientConn.Write([]byte(notification))
		if err != nil {
			fmt.Errorf(err.Error())
			utils.Save("./internal/logs/logs.log", err.Error(), true)
		}
		if clientConn != restrictConn {
			_, err := clientConn.Write([]byte("\n[" + timestamp + "][" + clientName + "]: "))
			if err != nil {
				fmt.Errorf(err.Error())
				utils.Save("./internal/logs/logs.log", err.Error(), true)
			}
		}
	}
}

func SendClientMessage(name, message, timestamp string) {
	formattedMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, name, message)
	resetter := "\r\033[K"
	Mutex.Lock()
	utils.Save("./internal/db/archived_msgs", formattedMsg, false)
	for clientConn, clientName := range Clients {
		if clientName != name {
			_, err := clientConn.Write([]byte(resetter + formattedMsg))
			if err != nil {
				fmt.Errorf(err.Error())
				utils.Save("./internal/logs/logs.log", "[WARNING] User message failed to send: "+name+" -> Group: "+message, true)
			}
			_, err = clientConn.Write([]byte("\n[" + timestamp + "][" + clientName + "]: "))
			if err != nil {
				fmt.Errorf(err.Error())
				utils.Save("./internal/logs/logs.log", "[WARNING] User prompt failed to send.", true)
			}
		}
	}
	Mutex.Unlock()
}

func HandleClient(clientConn net.Conn) {
	defer func() {
		Mutex.Lock()
		name := Clients[clientConn]
		delete(Clients, clientConn)
		Mutex.Unlock()

		NotifyClients(name, " has left the chat...", time.Now().Format("2006-01-02 15:04:05"), nil)
		utils.Save("./internal/logs/logs.log", fmt.Sprintf("[INFO] Client '%s' disconnected.", name), true)
	}()

	_, err := clientConn.Write([]byte(utils.WelcomeMsg))
	if err != nil {
		utils.Save("./internal/logs/logs.log", fmt.Sprintf("[ERROR] Failed to send welcome message: %v", err), true)
		return
	}

	name := GetClientName(clientConn)
	if name == "" {
		utils.Save("./internal/logs/logs.log", "[INFO] Client disconnected without providing a valid name.", true)
		return
	}

	Mutex.Lock()
	Clients[clientConn] = name
	Mutex.Unlock()

	NotifyClients(name, " has joined the chat...", time.Now().Format("2006-01-02 15:04:05"), clientConn)
	BackupHistory(clientConn)

	_, err = clientConn.Write([]byte("\n[" + time.Now().Format("2006-01-02 15:04:05") + "][" + name + "]: "))
	if err != nil {
		utils.Save("./internal/logs/logs.log", fmt.Sprintf("[ERROR] Failed to prompt client '%s': %v", name, err), true)
		return
	}

	buf := make([]byte, 4096)
	for {
		n, err := clientConn.Read(buf)
		if err != nil {
			if err == io.EOF {
				utils.Save("./internal/logs/logs.log", fmt.Sprintf("[INFO] Client '%s' disconnected (EOF).", name), true)
			} else {
				utils.Save("./internal/logs/logs.log", fmt.Sprintf("[ERROR] Error reading from client '%s': %v", name, err), true)
			}
			break
		}
		message := strings.TrimSpace(string(buf[:n]))
		if message == "" {
			continue
		}
		utils.Save("./internal/logs/logs.log", fmt.Sprintf("[INFO] Received message from '%s': %s", name, message), true)
		SendClientMessage(name, message, time.Now().Format("2006-01-02 15:04:05"))
		_, err = clientConn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + name + "]: "))
		if err != nil {
			utils.Save("./internal/logs/logs.log", fmt.Sprintf("[ERROR] Failed to prompt client '%s': %v", name, err), true)
			break
		}
	}
}
