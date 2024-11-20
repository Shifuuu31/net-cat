package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"TCPChat/internal/utils"
)

func StartSupervisor(listener net.Listener) {
	files := []string{
		"./internal/db/clients",
		"./internal/db/archived_msgs",
	}
	Createdb(files)
	signalChan := make(chan os.Signal, 1)

	// os.Interrupt represents ctrl+c
	// syscall.SIGTERM sent by the os or other processes to request program termination
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan

		logErr := fmt.Errorf("\n[INFO] Server shutting down... ")
		utils.Save("internal/logs/logs/log", logErr.Error(), true)

		sendExitMessageToClients()

		removedb(files)

		if listener != nil {
			listener.Close()
		}
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
			utils.Save("./internal/logs/logs.log", logErr.Error(), true)
		}
		clientConn.Close()
	}
	fmt.Errorf("[INFO] Server sent shutdown message to all clients.")
	utils.Save("./internal/logs/logs.log", "[INFO] Server sent shutdown message to all clients.", true)
}

func Createdb(files []string) {
	if err := os.MkdirAll("./internal/db", 0o777); err != nil {
		fmt.Errorf("[ERROR] Failed to create dir: %v", err)
		utils.Save("./internal/logs/logs.log", "[ERROR] Failed to create dir: "+err.Error(), true)

	}
	for _, filename := range files {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o777)
		if err != nil {
			fmt.Println("[ERROR] Failed to open file '%s': %v\n", filename, err)
			utils.Save("./internal/logs/logs.log", "[ERROR] Failed to open file '"+filename+"': "+err.Error()+"\n", true)
			return
		}
		defer file.Close()
	}
}

func removedb(files []string) {
	err := os.RemoveAll("./internal/db/")
	if err != nil {
		logErr := fmt.Errorf("[ERROR] Failed to remove db", err)
		utils.Save("./internal/logs/logs.log", logErr.Error()+"\n", true)
	}
}
