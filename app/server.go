package main

import (
	"bufio"
	"fmt"
	"net"
)


var store map[string]string =  make(map[string]string)

func main() {
	li, err := net.Listen("tcp",":6382")
	if err != nil {
		fmt.Println("error")
	}

	defer li.Close()

	fmt.Println("Server in listening a port 6382")

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
		req, err := ParseRESP(reader)
		fmt.Println("Client  :", req)
		if err !=nil {
			fmt.Println("Client Error :", err)
			conn.Write([]byte("-ERR invalid RESP format\r\n"))
			return 
		} 
		array, ok := req.([]interface{})
		if !ok || len(array) < 1{
			conn.Write([]byte("-ERR invalid resp array \r\n"))
			continue
		}

		command, ok := array[0].(string)
		if !ok {
			conn.Write([]byte("-ERR command must be a string\r\n"))
			continue
		}

		args := array[1:]
		stringArgs := make([]string, len(args))
		for i, arg := range args {
			strArg, ok := arg.(string)
			if !ok {
				fmt.Printf("Error: argument at index %d is not a string\n", i)
				return
			}
			stringArgs[i]= strArg
		}
		fmt.Println("Client 12 :", command, stringArgs)
		handleCommand(conn, command, stringArgs,store)

	}
}