package main

import (
	"../lib"
	"context"
	"fmt"
	"net"
)

var listenConfig = net.ListenConfig{
	Control: lib.Control,
}

func main() {

LOOP:
	for {
		var cmd, arg1, arg2, arg3 string
		fmt.Scanln(&cmd, &arg1, &arg2, &arg3)
		switch cmd {
		case "server":
			cmdServer(arg1)
		case "send":
			cmdSend(arg1, arg2, arg3)
		case "exit":
			break LOOP
		default:
			fmt.Println("unknown command")
		}
	}
}

func cmdServer(port string) {
	go func() {

		listener, err := listenConfig.Listen(context.TODO(), "tcp", "0.0.0.0:"+port)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer listener.Close()

		data := make([]byte, 1024)
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(conn.RemoteAddr().String())

			n, err := conn.Read(data)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(data[:n]))
			}

			conn.Close()
		}

	}()
}

func cmdSend(portLocal, addrRemote, content string) {
	addrLocal, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+portLocal)
	checkErr(err)

	dialer := net.Dialer{
		Control:   lib.Control,
		LocalAddr: addrLocal,
	}

	conn, err := dialer.Dial("tcp", addrRemote)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(content))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
