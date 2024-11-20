package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"TCPChat/internal/utils"
)

func StartSupervisor(shutdownChan chan struct{}, listener net.Listener, filesToRemove []string) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan

		fmt.Println("\n[INFO] Server shutting down...")

		sendExitMessageToClients()

		removeFiles(filesToRemove)

		if listener != nil {
			listener.Close()
		}
		close(shutdownChan)
		os.Exit(0)
	}()
}

func sendExitMessageToClients() {
	Mutex.Lock()
	defer Mutex.Unlock()

	exitMessage := "\n\033[31m[SERVER] The server is shutting down. Goodbye!\033[0m\n"
	for clientConn := range Clients {
		_, err := clientConn.Write([]byte(exitMessage))
		if err != nil {
			logErr := fmt.Errorf("[ERROR] Failed to send shutdown message to client: %v", err)
			utils.Save("./internal/logs/logs.log", logErr.Error() , true)
		}
		clientConn.Close()
	}
	fmt.Errorf("[INFO] Server sent shutdown message to all clients.")
	utils.Save("./internal/logs/logs.log", "[INFO] Server sent shutdown message to all clients.", true )
}

func removeFiles(files []string) {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			logErr := fmt.Errorf("[ERROR] Failed to remove file '%s': %v", file, err)
			utils.Save("./internal/logs/logs.log",logErr.Error(), true)
		} else {
			logErr := fmt.Errorf("[INFO] Successfully removed file: %s", file)
			utils.Save("./internal/logs/logs.log", logErr.Error(), true)
		}
	}
}
