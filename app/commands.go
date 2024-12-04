package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SafeMap struct {
	data        map[string]string
	expiryTimes map[string]time.Time
	mutex       sync.RWMutex
}

func NewSafeMap()*SafeMap{
	safeMap := &SafeMap{
		data: make(map[string]string),
		expiryTimes: make(map[string]time.Time),
	}
	return safeMap
}


func(sf *SafeMap) Get(key string) (string, bool){
	fmt.Println(key,"key")
	sf.mutex.RLock()


	exp, exists := sf.expiryTimes[key]


	if exists && time.Now().After(exp){
		sf.mutex.RUnlock()
		sf.mutex.Lock() 
		delete(sf.data,key)
		delete(sf.expiryTimes,key)
		sf.mutex.Unlock()
		return "", false
	}

	val,exists := sf.data[key]
	fmt.Println("Var",sf.data)
	sf.mutex.RUnlock()
	if(!exists){
		return "", false 
	}

	return val,true
}


func(sf *SafeMap) TTL(key string) (time.Duration, bool){
	sf.mutex.RLock()

	exp, exists := sf.expiryTimes[key]

	if !exists {
		return 0, false
	}

	if (time.Now().After(exp)){
		sf.mutex.RUnlock()
		sf.mutex.Lock()
		delete(sf.data,key)
		delete(sf.expiryTimes,key)
		sf.mutex.Unlock()
		return 0 ,false
	}
	val := time.Until(exp)
	 sf.mutex.RUnlock()
	return val,true

}


// cleanupExpiredKeys periodically removes expired keys.
func (sf *SafeMap) cleanupExpiredKeys() {
	for {
		time.Sleep(1 * time.Second) // Check every second
		sf.mutex.Lock()           // Lock for writing
		now := time.Now()
		for key, expiration := range sf.expiryTimes {
			if now.After(expiration) {
				delete(sf.data, key)
				delete(sf.expiryTimes, key)
			}
		}
		sf.mutex.Unlock()
	}
}

func(sf *SafeMap) Set(key,val string ,ttl time.Duration){
	fmt.Println("ttl",ttl)
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	sf.data[key] = val
	 // Only set expiry if ttl > 0
	 if ttl > 0 {
        sf.expiryTimes[key] = time.Now().Add(ttl)
    } else {
        delete(sf.expiryTimes, key) // Ensure no expiry for this key
    }
}


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

func handleSet(conn net.Conn, args []string, store *SafeMap){
	response := "OK"
	fmt.Println("sad",len(args))
	if len(args)!=2 && len(args)!=4 {
		conn.Write([]byte("-ERR wrong number of arguments for 'SET' command\r\n"))
		return
	}
	
	if(len(args) == 2){
		store.Set(args[0],args[1],0)
		
	} else {
		if(strings.ToUpper(args[2])!="PX"){
			conn.Write([]byte("-ERR wrong argument type for 'SET' command\r\n"))
			return
		}
		ttl,err := strconv.Atoi(args[3])
		if err !=nil {
			conn.Write([]byte("-ERR wrong argument type for 'PX' command\r\n"))
			return
		}
		store.Set(args[0],args[1],time.Duration(ttl)*time.Millisecond)
	}

	conn.Write([]byte(fmt.Sprintf("+%s\r\n", response)))
}

func handleGet(conn net.Conn, args []string, store *SafeMap){
	if len(args)!=1 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ECHO' command\r\n"))
		return
	}
	  val, exists := store.Get(args[0])
	  if !exists {
        conn.Write([]byte("$-1\r\n")) 
        return
    }
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val),val)))
}


func handleTTL(conn net.Conn, args []string, store *SafeMap){
	if len(args)!=1 {
		conn.Write([]byte("-ERR wrong number of arguments for 'ECHO' command\r\n"))
		return
	}
	  val, exists := store.TTL(args[0])
	  if !exists {
        conn.Write([]byte("$-1\r\n")) 
        return
    }
	// Convert time.Duration to seconds
	ttlSeconds := int64(val.Seconds())
	value := fmt.Sprintf("%d", ttlSeconds)

	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
}



func handleCommand(conn net.Conn, command string, args []string, safeMap *SafeMap){
	command = strings.ToUpper(command)
	
	switch command {
	case  "PING" :
		 handlePing(conn,args) 
	case "ECHO":
		 handleEcho(conn,args)
	case "SET":
		 handleSet(conn,args,safeMap)
	case "GET":
		 handleGet(conn,args, safeMap)
	case "TTL":
		handleTTL(conn,args,safeMap)
default:
	conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command)))
	}
}