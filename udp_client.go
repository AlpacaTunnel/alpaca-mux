package alpacamux

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type UdpClient struct {
	Server     string // upstream IP:Port pair
	conn       *net.UDPConn
	serverAddr *net.UDPAddr
}

func (f *UdpClient) Listen() error {
	ipPort := strings.Split(f.Server, ":")
	if len(ipPort) != 2 {
		return fmt.Errorf("wrong server IP:Port: %s", f.Server)
	}
	ip := net.ParseIP(ipPort[0])
	port, err := strconv.Atoi(ipPort[1])
	if ip == nil || err != nil {
		return fmt.Errorf("wrong server IP:Port: %s", f.Server)
	}
	f.serverAddr = &net.UDPAddr{
		IP:   ip,
		Port: port,
	}

	addr := &net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", addr)
	f.conn = conn
	return err
}

func (f *UdpClient) Read(buf []byte) (int, error) {
	len, client, err := f.conn.ReadFromUDP(buf)
	fmt.Println("udp client read from", *client)
	return len, err
}

func (f *UdpClient) Write(buf []byte) error {
	_, err := f.conn.WriteToUDP(buf, f.serverAddr)
	fmt.Println("udp client write to", *f.serverAddr)
	return err
}
