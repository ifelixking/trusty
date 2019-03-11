package main

import (
	"../lib"
	"context"
	"fmt"
	"net"
	"strings"
)

var listenConfig = net.ListenConfig{
	Control: lib.Control,
}

func main() {
LOOP:
	for {
		var cmd, arg1, arg2 string
		fmt.Scanln(&cmd, &arg1, &arg2)

		switch cmd {
		case "server":
			cmdServer(arg1)
		case "send":
			cmdSend(arg1, arg2)
		case "exit":
			break LOOP
		default:
			fmt.Println("unknown command")
		}
	}
}

func cmdServer(port string) {
	addr := "0.0.0.0:" + port
	listener, err := listenConfig.Listen(context.TODO(), "tcp", addr)
	checkErr(err)
	fmt.Println("started...")
	for {
		conn, err := listener.Accept()
		checkErr(err)

		var data = make([]byte, 1024)
		n, err := conn.Read(data)
		checkErr(err)

		fmt.Println(string(data[:n]))

		conn.Close()
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
