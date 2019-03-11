package main

import (
	"../lib"
	"context"
	"fmt"
	"net"
	"time"
)

var listenConfig = net.ListenConfig{
	Control: lib.Control,
}

func main() {
LOOP:
	for {
		var cmd, arg1, arg2, arg3 string
		fmt.Scanln(&cmd, &arg1, &arg2)
		switch cmd {
		case "send":
			cmdSend(arg1, arg2, arg3)
		case "server":
			cmdServer(arg1)
		case "callme":
			cmdCallme(arg1, arg2)
		case "exit":
			break LOOP
		default:
			fmt.Println("unknown command")
		}
	}
}

var bRepay bool = false

func cmdSend(strAddrLocal, strAddrRemote, content string) {
	addrLocal, err := net.ResolveTCPAddr("tcp", strAddrLocal)
	checkErr(err)
	send(addrLocal, strAddrRemote, content, true)
}

func send(addrLocal net.Addr, strAddrRemote, content string, autoClose bool) bool {
	dailer := net.Dialer{
		Control:   lib.Control,
		LocalAddr: addrLocal,
	}
	conn, err := dailer.Dial("tcp", strAddrRemote)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if autoClose {
		defer conn.Close()
	}

	_, err = conn.Write([]byte(content))
	if err != nil {
		fmt.Println(err)
		return false
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	checkErr(err)
	fmt.Println(string(data[:n]))

	return true
}

func cmdCallme(strAddrLocal, strAddrRemote string) {
	bRepay = false

	addrLocal, err := net.ResolveTCPAddr("tcp", strAddrLocal)
	checkErr(err)

	// for {
	send(addrLocal, strAddrRemote, "callme", false)
	fmt.Println("sent callme...")
	//		if bRepay {
	//			break
	//		}
	// }

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

			addrLocal := conn.LocalAddr()
			addrRemote := conn.RemoteAddr()
			fmt.Println("local:", addrLocal.String(), "remote:", addrRemote.String())

			n, err := conn.Read(data)
			content := string(data[:n])
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(content)
			}

			_, err = conn.Write([]byte("okok"))
			checkErr(err)

			err = conn.Close()
			checkErr(err)

			switch content {
			case "callme":
				for i := 0; i < 10; i++ {
					if send(addrLocal, addrRemote.String(), "reply", true) {
						fmt.Println("reply success")
					} else {
						fmt.Println("reply failed")
					}
					time.Sleep(time.Duration(100) * time.Millisecond)
				}
			case "reply":
				bRepay = true
			}

		}

	}()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
