package utils

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const WelcomeMsg = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'"

var (
	colors     = []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	colorMutex = &sync.Mutex{}
)

func PrintBanner(port string) {
	fmt.Printf("\033[1;32m\nTCPChat SERVER STATUS: 🟢 LIVE\n-> Server started on: localhost: %s\n-> Listening on the port: %s\033[0m\n", port, port)
}

func Get(filename string) string {
	readed, err := os.ReadFile(filename)
	if err != nil {
		fmt.Errorf("[ERROR] Failed to open file: %v", err)
		Save("./internal/logs/logs.log", "[ERROR] Failed to open file: "+err.Error(), true)
		// return true to siulate existent of user to avoid conflict
		// and keep new clients away while file connat be opened
		return ""
	}
	return string(readed)
}


func Save(filename, content string, islog bool) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		fmt.Errorf("[ERROR] Failed to open file '%s': %v\n", filename, err)
		Save("./internal/logs/logs.log", "[ERROR] Failed to open file '"+filename+"': "+err.Error()+"\n", true)
		return
	}
	defer file.Close()

	if islog {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		content = fmt.Sprintf("[%s] %s", timestamp, content)
	}

	_, err = file.WriteString(content + "\n")
	if err != nil {
		fmt.Errorf("[ERROR] Failed to write to file '%s': %v\n", filename, err)
		Save("./internal/logs/logs.log", "[ERROR] Failed to write to file '"+filename+"': "+err.Error()+"\n", true)
		return
	}
}

func AssignColor(input string, nclients int) string {
	colorMutex.Lock()
	defer colorMutex.Unlock()

	if nclients+1 > len(colors) {
		return input
	}
	return fmt.Sprintf("%s%s\033[0m", colors[nclients], input)
}
