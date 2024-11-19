package utils

import (
	"fmt"
	"sync"
)

const (
	WelcomeMsg = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'"
)

var (
	colors     = []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	colorMutex = &sync.Mutex{}
)

func AssignColor(input string, nclients int) string {
	colorMutex.Lock()
	defer colorMutex.Unlock()

	if nclients+1 > len(colors) {
		return input
	}
	return fmt.Sprintf("%s%s\033[0m", colors[nclients], input)
}

func PrintBanner(port string) {
	fmt.Printf("\033[1;32m\nTCPChat SERVER STATUS: 🟢 LIVE\n-> Server started on: localhost: %s\n-> Listening on the port: %s\033[0m\n", port, port)
}
