package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fmt.Println("Starting client...")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connected to localhost:8080")

	fmt.Print("Enter Your Name: ")
	var name string
	fmt.Scanln(&name)

	// Create scanner for user input
	scanner := bufio.NewScanner(os.Stdin)

	// Create scanner for reading messages from the server
	serverScanner := bufio.NewScanner(conn)

	// Goroutine for receiving messages from the server
	go func() {
		for serverScanner.Scan() {
			text := serverScanner.Text()
			message := Message{}
			err := json.Unmarshal([]byte(text), &message)
			if err != nil {
				fmt.Println("Error unmarshalling message:", err.Error())
				return
			}
			fmt.Println(message.Username + ": " + message.Message)
		}
	}()

	// Goroutine for sending messages to the server
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}
		message := Message{Username: name, Message: text}
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error marshalling message:", err.Error())
			return
		}
		_, err = conn.Write(append(jsonMessage, '\n'))
		if err != nil {
			fmt.Println("Error writing to stream:", err.Error())
			return
		}
	}
}
