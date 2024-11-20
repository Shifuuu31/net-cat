package server

import (
	"fmt"
	"net"
	"strings"

	"TCPChat/internal/utils"
)

func isAlphaNum(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}

func isValidName(name string) bool {
	if name == "" || len(name) < 3 {
		return false
	}

	for _, r := range name {
		if !isAlphaNum(r) {
			return false
		}
	}
	return true
}

func BackupHistory(clientConn net.Conn) {
	Mutex.Lock()
	defer Mutex.Unlock()
	archive := utils.Get("./internal/db/archived_msgs")
	_, err := clientConn.Write([]byte(archive))
	if err != nil {
		fmt.Errorf(err.Error())
		utils.Save("./internal/logs/logs.log", "[WARNING] failed to send archive to "+Clients[clientConn], true)
	}
}

func isClientNameExists(clientName string) bool {
	Mutex.Lock()
	defer Mutex.Unlock()
	if len(Clients) == 0 {
		return false
	}
	clients := strings.Split(utils.Get("./internal/db/clients"), "\n")
	for _, name := range clients {
		// name from idx 5 to 10 to remove assined color
		// fmt.Printf("\n[%s] | [%s]\n", clientName, name)
		if name == clientName {
			return true
		}
	}
	return false
}

func GetClientName(clientConn net.Conn) string {
	buf := make([]byte, 25)

	for {
		_, err := clientConn.Write([]byte("\n[ENTER YOUR NAME]: "))
		if err != nil {
			logErr := fmt.Errorf("[ERROR] failed prompting for client name: %v", err)
			utils.Save("./internal/logs/logs.log", logErr.Error(), true)
			return ""
		}

		n, readErr := clientConn.Read(buf)
		if readErr != nil {
			logErr := fmt.Errorf("[Error] failed reading client name: %v", readErr)
			utils.Save("./internal/logs/logs.log", logErr.Error(), true)
			return ""
		}

		name := strings.TrimSpace(string(buf[:n]))
		if !isValidName(name) {
			clientConn.Write([]byte("Invalid name. Alphanumeric characters only (min length: 3).\n"))
			continue
		}
		if isClientNameExists(name) {
			clientConn.Write([]byte("Name already in use. Try a different one.\n"))
			continue
		}

		utils.Save("./internal/db/clients", name, false)
		return utils.AssignColor(name, len(Clients))
	}
}
