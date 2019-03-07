package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
)

var strAddrMain = flag.String("m", "0.0.0.0:9527", "主连接的端口号")
var strAddrComm = flag.String("h", "0.0.0.0:9528", "通信用的端口号")
var listStrAddrClient = make([]string, 0, 20)

func main() {

	go server()
	go serverHole()

	mainLoop()
}

func server() {
	addr, err := net.ResolveTCPAddr("tcp4", *strAddrMain)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", addr)
	checkErr(err)
	defer listener.Close()
	fmt.Println("node started...")
	for {
		conn, err := listener.Accept()
		checkErr(err)
		onMainAccept(conn)
	}
}

// B login 10.1.165.105:9527
// C login 10.1.165.105:9527

// B send 10.1.165.105:9527

func serverHole() {
	addr, err := net.ResolveTCPAddr("tcp4", *strAddrComm)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", addr)
	checkErr(err)
	defer listener.Close()
	fmt.Println("node hole started...")
	for {
		conn, err := listener.Accept()
		checkErr(err)
		onMainAccept(conn)
	}
}

func mainLoop() {
LOOP:
	for {
		var cmd, arg1, arg2 string
		_, err := fmt.Scanln(&cmd, &arg1, &arg2)
		checkErr(err)

		switch cmd {
		case "list":
			cmdList()
		case "login":
			cmdSend(arg1, "login")
		case "send":
			cmdSend(arg1, arg2)
		case "exit":
			break LOOP
		default:
			fmt.Println("unknown command")
		}
	}
}

func cmdSend(strAddrClient, content string) {
	addrComm, err := net.ResolveTCPAddr("tcp", *strAddrComm)
	checkErr(err)
	addrRemote, err := net.ResolveTCPAddr("tcp", strAddrClient)
	checkErr(err)
	var conn *net.TCPConn
	for i := 0; i < 10; i++ {
		conn, err = net.DialTCP("tcp", addrComm, addrRemote)
		if err == nil {
			break
		}
		fmt.Println("retry")
		time.Sleep(time.Duration(100) * time.Millisecond)
		conn = nil
	}
	if conn == nil {
		return
	}
	fmt.Println("ok")
	defer conn.Close()
	_, err = conn.Write([]byte(content))
	checkErr(err)
}

func onMainAccept(conn net.Conn) {
	var data = make([]byte, 1024)
	n, err := conn.Read(data)
	checkErr(err)
	defer conn.Close()
	var args = strings.Split(string(data[:n]), " ")
	switch args[0] {
	case "login":
		addrRemote := conn.RemoteAddr().String()
		addToClientList(addrRemote)
		fmt.Println(addrRemote)
		_, err = conn.Write([]byte("msg welcome"))
	case "msg":
		fmt.Println(args[1])
	}
	checkErr(err)
}

func cmdList() {
	for key, value := range listStrAddrClient {
		fmt.Printf("[%d] %s\n", key, value)
	}
}

func addToClientList(strAddrClient string) {
	for _, value := range listStrAddrClient {
		if value == strAddrClient {
			return
		}
	}
	listStrAddrClient = append(listStrAddrClient, strAddrClient)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
