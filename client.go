package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

const MAX_MTU = 1500

var udpClient *net.UDPAddr

func local2remote(connServer *net.UDPConn, connClient *net.UDPConn, remoteAddr *net.UDPAddr) {
	buf := make([]byte, MAX_MTU)
	for {
		len, client, err := connServer.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("error recv: %v\n", err)
		}
		udpClient = client

		_, err = connClient.WriteToUDP(buf[:len], remoteAddr)
		if err != nil {
			fmt.Printf("error send: %v\n", err)
		}
	}
}

func remote2local(connClient *net.UDPConn, connServer *net.UDPConn) {
	buf := make([]byte, MAX_MTU)
	for {
		len, _, err := connClient.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("error recv: %v\n", err)
		}

		_, err = connClient.WriteToUDP(buf[:len], udpClient)
		if err != nil {
			fmt.Printf("error send: %v\n", err)
		}
	}
}

func getArgs() (net.UDPAddr, net.UDPAddr, net.UDPAddr) {
	localPort := flag.Int("l", 1080, "local port")
	remoteIp := flag.String("s", "", "remote server")
	remotePort := flag.Int("p", 8000, "remote server port")
	flag.Parse()

	localServer := net.UDPAddr{
		Port: *localPort,
		IP:   net.ParseIP("0.0.0.0"),
	}

	localClient := net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("0.0.0.0"),
	}

	remoteAddr := net.UDPAddr{
		Port: *remotePort,
		IP:   net.ParseIP(*remoteIp),
	}

	return localServer, localClient, remoteAddr
}

func main() {
	localServer, localClient, remoteAddr := getArgs()

	connServer, err := net.ListenUDP("udp", &localServer)
	if err != nil {
		fmt.Printf("error listen: %v\n", err)
		return
	}

	connClient, err := net.ListenUDP("udp", &localClient)
	if err != nil {
		fmt.Printf("error listen: %v\n", err)
		return
	}

	go local2remote(connServer, connClient, &remoteAddr)
	go remote2local(connClient, connServer)

	for {
		time.Sleep(time.Second)
	}
}
