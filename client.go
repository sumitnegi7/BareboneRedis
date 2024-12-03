package main

import (
	"fmt"
	"net"
)

func main() {

    conn, err := net.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer conn.Close()

    // Send data to the server
	data := []byte("Hello");
	_, err = conn.Write(data)

	if err != nil {
		fmt.Println("Error:", err)
	}
	
    // Read and process data from the server
}