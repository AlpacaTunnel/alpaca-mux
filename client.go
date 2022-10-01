package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const MAX_MTU = 1500

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func udp2mux(udpServer *net.UDPConn, udpClientAddr *net.UDPAddr, muxClient *net.UDPConn, remoteAddrs []*net.UDPAddr) {
	buf := make([]byte, MAX_MTU)
	for {
		len, client, err := udpServer.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("error recv: %v\n", err)
		}
		*udpClientAddr = *client

		for _, addr := range remoteAddrs {
			_, err = muxClient.WriteToUDP(buf[:len], addr)
			if err != nil {
				fmt.Printf("error send: %v\n", err)
			}
		}
	}
}

func mux2udp(muxClient *net.UDPConn, udpServer *net.UDPConn, udpClientAddr *net.UDPAddr) {
	buf := make([]byte, MAX_MTU)
	for {
		len, _, err := muxClient.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("error recv: %v\n", err)
		}

		_, err = udpServer.WriteToUDP(buf[:len], udpClientAddr)
		if err != nil {
			fmt.Printf("error send: %v\n", err)
		}
	}
}

func getArgs() (*net.UDPAddr, *net.UDPAddr, []*net.UDPAddr) {
	localPort := flag.Int("l", 1080, "local port")
	var servers arrayFlags
	flag.Var(&servers, "s", "Remote server IP:Port")
	flag.Parse()

	localServer := &net.UDPAddr{
		Port: *localPort,
		IP:   net.ParseIP("0.0.0.0"),
	}

	localClient := &net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("0.0.0.0"),
	}

	var remoteAddrs []*net.UDPAddr

	for _, s := range servers {
		ipPort := strings.Split(s, ":")
		ip := net.ParseIP(ipPort[0])
		port, err := strconv.Atoi(ipPort[1])
		if ip == nil || err != nil {
			panic("wrong IP:Port")
		}
		remoteAddrs = append(remoteAddrs, &net.UDPAddr{
			IP:   ip,
			Port: port,
		})
	}

	return localServer, localClient, remoteAddrs
}

func main() {
	udpServerAddr, muxClientAddr, muxServerAddrs := getArgs()
	var udpClientAddr net.UDPAddr

	udpServer, err := net.ListenUDP("udp", udpServerAddr)
	if err != nil {
		fmt.Printf("error listen: %v\n", err)
		return
	}

	muxClient, err := net.ListenUDP("udp", muxClientAddr)
	if err != nil {
		fmt.Printf("error listen: %v\n", err)
		return
	}

	go udp2mux(udpServer, &udpClientAddr, muxClient, muxServerAddrs)
	go mux2udp(muxClient, udpServer, &udpClientAddr)

	for {
		time.Sleep(time.Second)
	}
}
