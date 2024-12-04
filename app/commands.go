package main

import (
	"fmt"
	"net"
	"strings"
)

func handlePing(conn net.Conn,_ []string){
	response := "PONG"
	conn.Write([]byte(fmt.Sprintf("+%s\r\n",response)))
}

func handleEcho(conn net.Conn, args []string){
	if len(args)!=1 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ECHO' command\r\n"))
		return
	}
	
	arg := args[0]
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n",len(arg),arg)))
}

func handleSet(conn net.Conn, args []string, store map[string]string){
	response := "OK"
	if len(args)!=2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ECHO' command\r\n"))
		return
	}
	
	store[args[0]] = args[1]
	conn.Write([]byte(fmt.Sprintf("+%s\r\n", response)))
}

func handleGet(conn net.Conn, args []string, store map[string]string){
	if len(args)!=1 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ECHO' command\r\n"))
		return
	}
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(store[args[0]]),store[args[0]])))
}

func handleCommand(conn net.Conn, command string, args []string, store map[string]string){
	command = strings.ToUpper(command)
	
	switch command {
	case  "PING" :
		 handlePing(conn,args) 
	case "ECHO":
		 handleEcho(conn,args)
	case "SET":
		 handleSet(conn,args,store)
	case "GET":
		 handleGet(conn,args, store)
default:
	conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command)))
	}
}