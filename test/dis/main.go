package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-discovery"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-protocol"
	"github.com/multiformats/go-multiaddr"
	"os"
	"strings"
)

type addrList []multiaddr.Multiaddr

func (list *addrList) String() string{
	strs := make([]string, len(*list))
	for i, addr := range *list {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (list *addrList)Set(value string) error{
	addr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*list = append(*list, addr)
	return nil
}

func handleStream(stream inet.Stream) {
	fmt.Println("in coming...")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

var isFirst = flag.Bool("f", false, "is first node")
var port = flag.String("p", "9527", "port")
var addrBootstrap = flag.String("b", "", "address of bootstrap")

func main() {

	ctx := context.Background()

	// 构造地址
	var listenAddresses addrList
	err := listenAddresses.Set("/ip4/0.0.0.0/tcp/"+*port)
	check(err)

	// 创建 host
	host, err := libp2p.New(ctx, libp2p.ListenAddrs(listenAddresses...))
	check(err)
	fmt.Println(host.Addrs())			// 地址
	fmt.Println(host.ID().Pretty())		// ID
	host.SetStreamHandler(protocol.ID("/ids/0.0.1"), handleStream)

	// 创建 dht client
	kademliaDHT, err := libp2pdht.New(ctx, host)

	// 启动 dht client
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Bootstrap...")

	if !*isFirst {

		fmt.Println("连接 Bootstrap Peers")
		ma, err := multiaddr.NewMultiaddr(*addrBootstrap)
		check(err)
		peerinfo, _ := peerstore.InfoFromP2pAddr(ma)
		if err := host.Connect(ctx, *peerinfo); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Connection established with bootstrap node:", *peerinfo)
		}

		fmt.Println("宣告自己")
		rendezvousString := "show me the money"
		routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
		discovery.Advertise(ctx, routingDiscovery, rendezvousString)

		fmt.Println("发现别人")
		peerChan, err := routingDiscovery.FindPeers(ctx, rendezvousString)
		check(err)
		for peer := range peerChan {
			if peer.ID == host.ID() {
				continue
			}
			fmt.Println("Connecting to:", peer)
			stream, err := host.NewStream(ctx, peer.ID, protocol.ID("/ids/0.0.1"))

			if err != nil {
				fmt.Println("Connection failed:", err)
				continue
			} else {
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

				go writeData(rw)
				go readData(rw)
			}
			fmt.Println("Connected to:", peer)
		}
	}

	select {}
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			fmt.Printf(str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func check(err error){
	if err != nil{
		fmt.Println(err)
	}
}