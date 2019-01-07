package main

import (
	"./protector"
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

func listenToClient(address string) {
	fmt.Println("Launching server...")
	// Listen to all interfaces
	ln, _ := net.Listen("tcp", address)
	// Accept connections
	conn, _ := ln.Accept()

	var prot *protector.Protector
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = message[:len(message) - 1] // Remove '\n' sign
		fmt.Println("Message Received: " + string(message))

		var clientSKey, clientHash string
		if strings.Contains(string(message), ":") {
			var keyAndHash = strings.Split(message, ":")
			clientSKey = keyAndHash[0]
			clientHash = keyAndHash[1]

			prot = protector.New(clientHash)
		} else {
			clientSKey = string(message)
		}

		var sKey = prot.GenerateNextSessionKey(clientSKey)
		conn.Write([]byte(sKey + "\n"))
	}
}

func dialToServer(port string) {
	fmt.Println("Dialing to " + port)
	conn, _ := net.Dial("tcp", port)


	var initialSKey = getSessionKey()
	var initialHash = getHashString()

	var sKey string
	var prot = protector.New(initialHash)
	for {
		var keyAndHash = initialSKey + ":" + initialHash
		if sKey == "" {
			fmt.Println("Text to send: " + keyAndHash)
			fmt.Fprintf(conn, keyAndHash + "\n")
			sKey = initialSKey
		} else {
			fmt.Println("Text to send: " + sKey)
			fmt.Fprintf(conn, sKey + "\n")
		}

		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = message[:len(message) - 1] // Remove '\n' sign
		fmt.Println("Message from server: " + message)

		sKey = prot.GenerateNextSessionKey(sKey)
		if message != sKey {
			panic("Keys don't match")
			break
		}
		fmt.Println("Keys are equal")

		sKey = prot.GenerateNextSessionKey(sKey)
		fmt.Fprintf(conn, sKey + "\n")
	}
}

func getSessionKey() string {
	var sb strings.Builder
	for i := 0; i < 10; i++ {
		var byte = byte('1' + rand.Intn(9))
		sb.WriteByte(byte)
	}

	return sb.String()
}

func getHashString() string {
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		var byte = byte('1' + rand.Intn(6))
		sb.WriteByte(byte)
	}

	return sb.String()
}

func validateIPv4(ip string) bool {
	if &ip == nil {
		return false
	}

	var splitValues = strings.Split(ip, ".")
	if len(splitValues) != 4 {
		return  false
	}

	for _, value := range splitValues {
		var integer, err = strconv.Atoi(value)
		if err != nil {
			return false
		}

		if (integer < 0) || (integer > math.MaxUint8) {
			return false
		}
	}

	return true
}

func main() {
	portPtr := flag.String("port", ":9000", "app ip*:port")
	instQuantityPtr := flag.Uint("inst", 1, "instances quantity")
	flag.Parse()

	var mode string
	var address = strings.Split(string(*portPtr), ":")
	if address[0] != "" {
		var ipValid = validateIPv4(address[0])
		if !ipValid {
			panic("IP address is not valid")
		}
		mode = "client"
	} else {
		mode = "server"
	}

	if mode == "server" {
		for i := 0; i < int(*instQuantityPtr); i++ {
			listenToClient(*portPtr)
		}
	}

	if mode == "client" {
		for i := 0; i < int(*instQuantityPtr); i++ {
			dialToServer(*portPtr)
		}
	}
}