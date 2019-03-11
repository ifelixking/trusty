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
			cmdServer(arg1, false)
		case "serveronce":
			cmdServer(arg1, true)
		case "send":
			cmdSend(arg1, arg2, arg3)
		case "exit":
			break LOOP
		default:
			fmt.Println("unknown command")
		}
	}
}

func cmdServer(port string, once bool) {
	go func() {

		listener, err := listenConfig.Listen(context.TODO(), "tcp", "0.0.0.0:"+port)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer listener.Close()

		fmt.Println("server is started...")

		data := make([]byte, 1024)
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			addrRemote := conn.RemoteAddr().String()
			addrLocal := conn.LocalAddr().String()
			fmt.Println("coming remote: ", addrRemote, " local: ", addrLocal)

			n, err := conn.Read(data)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(data[:n]))
			}

			n, err = conn.Write([]byte("welcome"))
			checkErr(err)
			fmt.Println("sent welcome...")

			conn.Close()

			if once {
				break
			}

			// fmt.Println("send back...")
			// cmdSend(port, addrRemote, "welcome")
		}

		fmt.Println("server is exit...")

	}()
}

func cmdSend(strAddrLocal, strAddrRemote, content string) {
	addrLocal, err := net.ResolveTCPAddr("tcp", strAddrLocal)
	checkErr(err)

	dialer := net.Dialer{
		Control:   lib.Control,
		LocalAddr: addrLocal,
	}

	conn, err := dialer.Dial("tcp", strAddrRemote)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(content))
	checkErr(err)
	fmt.Println("msg is sent, and wait for welcome...")

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(data[:n]))
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
