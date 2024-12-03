package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	li, err := net.Listen("tcp",":6379")
	if err != nil {
		fmt.Println("error")
	}

	defer li.Close()

	fmt.Println("Server in listening a port 6379")

	for {
		conn, err := li.Accept();
		if err != nil {
			fmt.Println("error")
			continue
		}		
			go handleClient(conn)
	}
}

func handleClient(conn net.Conn){
	defer conn.Close()


	reader := bufio.NewReader(conn)

	for {
		line,err := reader.ReadString('\n')
		if err !=nil {
			fmt.Println("Error :", err)
			return 
		} 

		fmt.Printf("Recieved: %s\n", line)
	}

	// Write data to the client
}