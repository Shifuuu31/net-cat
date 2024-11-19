package utils

import "fmt"

const (
	WelcomeMsg = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'"
)

func PrintBanner(port string) {
	fmt.Printf("\033[1;32m\nTCPChat SERVER STATUS: 🟢 LIVE\n-> Server started on: localhost:%s\n-> Listening on the port: %s\033[0m\n", port, port)
}
