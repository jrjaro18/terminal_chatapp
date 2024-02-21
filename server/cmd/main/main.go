package main

import (
	"fmt"
	"github.com/jrjaro18/terminal_chatapp/server/internals"
)

func main() {
	fmt.Println("Starting server...")
	internals.SocketInit()
}