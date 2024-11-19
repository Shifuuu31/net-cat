package utils

import (
	"fmt"
	"sync"
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
