package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type clientType map[string]bool

var clients = clientType{}

type conns map[string]*net.TCPConn

var openConnections = conns{}

func (c clientType) keys(filter string) string {
	output := []string{}
	for key := range c {
		if key != filter {
			output = append(output, key)
		}
	}
	return strings.Join(output, ",")
}

// Server --
func Server() {
	fmt.Println("Running Server")
	localAddress := ":9595"
	if len(os.Args) > 2 {
		localAddress = os.Args[2]
	}

	addr, _ := net.ResolveTCPAddr("tcp", localAddress)
	listener, _ := net.ListenTCP("tcp", addr)
	for {
		fmt.Println("Waiting for connection")
		conn, _ := listener.AcceptTCP()
		fmt.Println("accepted: ", conn.RemoteAddr())
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		fmt.Println(err)
		fmt.Println(string(buffer[:n]))

		openConnections[conn.RemoteAddr().String()] = conn
		incoming := string(buffer[0:n])
		fmt.Println("[INCOMING]", incoming)
		if incoming != "register" {
			continue
		}
		clients[conn.RemoteAddr().String()] = true

		for connection := range openConnections {
			if connection != conn.RemoteAddr().String() {
				conn.Write([]byte(connection))
				openConnections[connection].Write([]byte(conn.RemoteAddr().String()))
				fmt.Println(conn.RemoteAddr())
				fmt.Printf("[INFO] Responded to %s with %s\n", conn.RemoteAddr(), connection)
				fmt.Printf("[INFO] Responded to %s with %s\n", connection, conn.RemoteAddr())

			}
		}

	}
}
