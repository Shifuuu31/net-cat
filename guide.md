Here's a **step-by-step guide** to successfully build a TCP Chat application similar to NetCat using Go:

* * *

### **1\. Setup Project Structure**

1.  Create a new directory for the project.
2.  Initialize a Go module:
    
    ```bash copy
    go mod init tcp-chat
    ```
    

* * *

### **2\. Define Objectives**

*   **Server**: Accept multiple client connections and broadcast messages.
*   **Client**: Connect to the server, send/receive messages.
*   Enforce name validation.
*   Maintain a message log to share with new clients.
*   Inform clients when others join or leave.
*   Limit to 10 simultaneous connections.

* * *

### **3\. Server Implementation**

*   Use the ```net``` package to create a TCP server.
*   Handle multiple clients concurrently with Goroutines.
*   Use a ```sync.Mutex``` or channels for safe access to shared resources (like client list and message log).

**Steps**:

1.  **Create a TCP Server**:
    
    ```go
    listener, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
    log.Println("Listening on port 8080")
    defer listener.Close()
    ```
    
2.  **Accept Client Connections**: Use a loop to accept incoming connections and start a Goroutine for each client.
    
    ```go
    for {
        conn, err := listener.Accept()     if err != nil {
        log.Println("Failed to accept connection:", err)
        continue
        }
        go handleClient(conn)
        }
    ```
    
3.  **Handle Client Communication**: Parse client messages and broadcast them to all other clients.
    
    *   Read/write using ```bufio```.
    *   Include timestamp, username, and message format.
4.  **Broadcast Messages**: Maintain a global list of active connections. Use a ```sync.Mutex``` to protect it.
    
5.  **Log Messages**: Store all messages in a slice or write them to a log file for new clients.
    

* * *

### **4\. Client Implementation**

*   Use ```net.Dial``` to connect to the server.
*   Prompt the user for their name, validate it, and send it to the server.
*   Use a Goroutine to continuously read messages from the server.
*   Send messages to the server in the main thread.

**Steps**:

1.  **Connect to Server**:
    
    ```go
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()
    ```
    
2.  **Prompt for Username**: Validate non-empty input and send to server.
    
    ```go
    fmt.Print("[ENTER YOUR NAME]: ")
    name := strings.TrimSpace()
    fmt.Fprintf(conn, name + "\n")
    ```
    
3.  **Read Messages**: Start a Goroutine to read messages from the server.
    
    ```go
    go func() {
        scanner := bufio.NewScanner(conn)
        for scanner.Scan() {
            fmt.Println(scanner.Text())
        }
    }()
    ```
    
4.  **Send Messages**: Use ```fmt.Fprintf``` to send user input to the server.
    

* * *

### **5\. Concurrency and Limits**

*   Use ```sync.Mutex``` or channels to manage:
    *   Shared message log.
    *   Active connections.
    *   Connection limit enforcement.

* * *

### **6\. Error Handling**

*   Ensure the server handles client disconnections gracefully.
*   Validate input on the server (e.g., reject empty usernames).
*   Log errors for debugging.

* * *

### **7\. Features to Implement**

*   **Username Broadcast**: Notify all clients when a new user joins or leaves.
*   **Message Formatting**: Prefix each message with the timestamp and username.
*   **Message History**: Send the chat log to new clients upon joining.
*   **Connection Limit**: Disconnect new clients if the limit (10) is reached.

* * *

### **8\. Additional Enhancements**

*   **Terminal UI**: Use the ```gocui``` package for a text-based UI.
*   **Logging**: Write chat logs to a file for persistence.
*   **Multiple Groups**: Handle multiple chat rooms (bonus).

* * *

### **9\. Testing**

1.  Create unit tests for:
    *   Server connection handling.
    *   Client connection and messaging.
2.  Use tools like ```telnet``` or ```nc``` to simulate clients:
    
    ```bash
    nc localhost 8080
    ```
    

* * *

### **10\. Run and Debug**

*   Start the server:
    
    ```bash
    go run server.go
    ```
    
*   Connect clients:
    
    ```bash
    go run client.go
    ```
    

* * *
