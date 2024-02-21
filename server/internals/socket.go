package internals

import (
	"fmt"
	"net"
	"os"
)

func SocketInit() {

	listener, err := net.Listen("tcp", "localhost:8080")
	if(err != nil) {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening on localhost:8080")

	for {
		conn, err := listener.Accept()
		if(err != nil) {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}
		go HandleRequest(conn)
	}

}
