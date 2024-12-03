package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

func parseInteger(reader *bufio.Reader)(string, error){
	data, err := reader.ReadString('\n')
	
	if err !=nil {
		return "", err
	}
	return strings.TrimSuffix(data,"\r\n"), nil
}

func parseSimpleString(reader *bufio.Reader)(string,error){
	data, err  := reader.ReadString('\n')
	if err !=nil {
		return "", err
	}
	 return strings.TrimSuffix(data,"\r\n"), nil
}

func parseArray(reader *bufio.Reader)([]interface{}, error){
	// data, err := reader.ReadString()
	// *2\r\n$4\r\nECHO\r\n$3\r\nsss\r\n
	data, err := reader.ReadString('\n')
	if err !=nil {
		return nil, err
	}

	len,err :=  strconv.Atoi(strings.TrimSuffix(data, "\r\n"))
	if err !=nil {
		return nil, err
	} 

	if len == -1 {
		return nil, nil
	} 

	elements := make([]interface{},len);

	for i:=0; i< len; i++ {
		element, err := ParseRESP(reader)
		if err!=nil {
			return nil, err
		}
		elements[i] = element
	}
	return elements,nil
}

func parseBulkString(reader *bufio.Reader)(interface{}, error){
		//4\r\nECHO\r\n
	str, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	} 

	length, err := strconv.Atoi(strings.TrimSuffix(str, "\r\n"))
	
	if err != nil {
		return "", err
	}

 	// Handle null bulk strings
	if length == -1 {
		return "", nil
	}

 	data := make([]byte, length)

	_, err = io.ReadFull(reader,data)
	if err != nil {
		return "", err
	}
	_, err = reader.Discard(2)
	if err != nil {
		return "", err
	}
	return string(data),nil
}



func parseError(reader *bufio.Reader) (error, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return errors.New(strings.TrimSuffix(line, "\r\n")), nil
}


func ParseRESP(reader *bufio.Reader)(interface{}, error){
	// Reading first byte to determine the type
	prefix,err := reader.ReadByte();

	if err != nil {
		return nil, err
	}

	switch prefix {
		case '+':
			return parseSimpleString(reader)
		case '-':
			return parseError(reader)
		case '*':
			return parseArray(reader)
		case '$':
			return parseBulkString(reader)
		case ':':
			return parseInteger(reader)
		default:
			return nil, errors.New("invalid RESP type prefix")
	}

}