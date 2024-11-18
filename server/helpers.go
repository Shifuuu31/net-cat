package server

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	colors         = []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	assignedColors = map[string]string{}
)

func AssignRandomColor(input string) string {
	if color, exists := assignedColors[input]; exists {
		return fmt.Sprintf("%s%s\033[0m", color, input)
	}
	if len(colors) == 0 {
		return input
	}
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(colors))
	color := colors[idx]
	colors = append(colors[:idx], colors[idx+1:]...)
	assignedColors[input] = color
	return fmt.Sprintf("%s%s\033[0m", color, input)
}

func NotifyClients(message string) {
	Mutex.Lock()
	defer Mutex.Unlock()

	for client := range Clients {
		fmt.Fprintf(client, "%s\n", message)
	}
}


func PrintBanner(port string) {
	// banner := 
	fmt.Printf(`
` + "\033[1;32m" + `                                           ` + "\033[0m" + `
` + "\033[1;32m" + `       ` + "\033[1;33m" + `TCPChat SERVER STATUS: ` + "\033[1;32m" + `🟢 LIVE       ` + "\033[1;32m" + "\033[0m" + `
` + "\033[1;32m" + `                                           ` + "\033[0m" + `
` + "\033[1;32m" + `     Server started on: ` + "\033[1;34m" + `localhost:` + port + `     ` + "\033[1;32m" + "\033[0m" + `
` + "\033[1;32m" + `     Listening on the port: ` + port + `           ` + "\033[1;32m" + "\033[0m" + `
` + "\033[1;32m" + `                                           ` + "\033[0m")
}