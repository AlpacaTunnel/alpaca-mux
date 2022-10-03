package alpacamux

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type MuxClient struct {
	Servers     []string // upstream IP:Port pairs
	conn        *net.UDPConn
	serverAddrs []*net.UDPAddr
}

func (f *MuxClient) Listen() error {
	for _, s := range f.Servers {
		ipPort := strings.Split(s, ":")
		if len(ipPort) != 2 {
			return fmt.Errorf("wrong server IP:Port: %s", s)
		}
		ip := net.ParseIP(ipPort[0])
		port, err := strconv.Atoi(ipPort[1])
		if ip == nil || err != nil {
			return fmt.Errorf("wrong server IP:Port: %s", s)
		}
		f.serverAddrs = append(f.serverAddrs, &net.UDPAddr{
			IP:   ip,
			Port: port,
		})
	}

	addr := &net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", addr)
	f.conn = conn
	return err
}

func (f *MuxClient) Read(buf []byte) (int, error) {
	len, client, err := f.conn.ReadFromUDP(buf)
	fmt.Println("mux client read from", *client)
	return len, err
}

func (f *MuxClient) Write(buf []byte) error {
	var err error
	for _, addr := range f.serverAddrs {
		if addr.Port > 0 {
			fmt.Println("mux client write to", *addr)
			_, e := f.conn.WriteToUDP(buf, addr)
			if e != nil {
				err = e
			}
		}
	}
	return err
}
