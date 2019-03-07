package udp

import (
	"flag"
	"fmt"
	"net"
)

var strAddrServer = flag.String("a", "0.0.0.0:9527", "当前节点的地址和端口")
var listStrAddrClient = make([]string, 0, 20)
var connServer *net.UDPConn

func main() {
	flag.Parse()

	go server()

	mainLoop()
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
			// num, err := strconv.Atoi(arg1)
			// var strAddrClient string
			// if err != nil {
			//	strAddrClient = listStrAddrClient[num]
			//} else {
			//	strAddrClient := arg1
			//}
			cmdSend(arg1, arg2)
		case "exit":
			break LOOP
		default:
			fmt.Println("unknown command")
		}
	}
}

func cmdSend(strAddrClient, content string) {
	addrClient, err := net.ResolveUDPAddr("udp", strAddrClient)
	checkErr(err)
	_, err = connServer.WriteToUDP([]byte(content), addrClient)
	checkErr(err)
}

func cmdList() {
	for key, value := range listStrAddrClient {
		fmt.Printf("[%d] %s\n", key, value)
	}
}

func server() {
	addr, err := net.ResolveUDPAddr("udp", *strAddrServer)
	checkErr(err)
	connServer, err = net.ListenUDP("udp", addr)
	defer connServer.Close()
	fmt.Println("node is started...")
	data := make([]byte, 1024)
	for {
		n, addrClient, err := connServer.ReadFromUDP(data)
		checkErr(err)
		fmt.Printf("<%s> %s\n", addrClient, data[:n])
		addToClientList(addrClient.String())
		listClients()
	}
}

func listClients() {}
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
