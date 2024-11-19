package server

import (
	"log"
	"net"
	"strings"

	"TCPChat/internal/utils"
)

func isAlphaNum(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}

func isValidName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !isAlphaNum(r) {
			return false
		}
	}
	return true
}

func isClientNameExists(name string) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	for _, clientName := range Clients {
		// clientName from idx 5 to 10 to remove assined color
		if clientName[5:10] == name {
			// fmt.Printf("\n[%s] | [%s]\n", name, clientName[5:10])
			return true
		}
	}
	return false
}

func GetClientName(clientConn net.Conn) string {
	buf := make([]byte, 256)

	for {
		clientConn.Write([]byte("\n[ENTER YOUR NAME]: "))
		n, err := clientConn.Read(buf)
		if err != nil {
			log.Printf("Error reading client name: %v", err)
			return ""
		}

		name := strings.TrimSpace(string(buf[:n]))
		if !isValidName(name) {
			clientConn.Write([]byte("Invalid name. Names must be alphanumeric and non-empty. Please try again."))
		} else if isClientNameExists(name) {
			clientConn.Write([]byte("\033[31m" + name + "\033[0m already taken, try again"))
		} else {
			return utils.AssignColor(name, len(Clients))
		}

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
