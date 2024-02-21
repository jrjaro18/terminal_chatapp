package internals

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net"
    "sync"
)

type Message struct {
    Username string `json:"username"`
    Message  string `json:"message"`
}

var (
    connections = make(map[net.Conn]struct{})
    mu          sync.Mutex
)

func HandleRequest(conn net.Conn) {
    mu.Lock()
    connections[conn] = struct{}{}
    mu.Unlock()

    defer func() {
        mu.Lock()
        delete(connections, conn)
        mu.Unlock()
        conn.Close()
    }()

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        text := scanner.Text()
        message := Message{}
        if err := json.Unmarshal([]byte(text), &message); err != nil {
            fmt.Println("Error unmarshalling message:", err.Error())
            return
        }
        fmt.Println(message)

        if message.Message == "exit" {
            break
        }

        broadcast(message)
    }
}

func broadcast(message Message) {
    encodedMessage, err := json.Marshal(message)
    if err != nil {
        fmt.Println("Error marshalling message:", err.Error())
        return
    }

    mu.Lock()
    defer mu.Unlock()

    for conn := range connections {
        if _, err := conn.Write(append(encodedMessage, '\n')); err != nil {
            fmt.Println("Error writing to connection:", err.Error())
            // Close the connection if there's an error
            conn.Close()
            delete(connections, conn)
        }
    }
}
