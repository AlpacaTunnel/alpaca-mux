package alpacamux

import (
	"fmt"
	"net"
)

type MuxServer struct {
	Port        int
	conn        *net.UDPConn
	clientAddrs [4]*net.UDPAddr
}

func (f *MuxServer) Listen() error {
	f.clientAddrs = [4]*net.UDPAddr{
		{},
		{},
		{},
		{},
	}

	addr := &net.UDPAddr{
		Port: f.Port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", addr)
	f.conn = conn
	return err
}

func (f *MuxServer) Read(buf []byte) (int, error) {
	len, client, err := f.conn.ReadFromUDP(buf)
	fmt.Println("mux server read from", *client)
	*(f.clientAddrs[0]) = *client
	return len, err
}

func (f *MuxServer) Write(buf []byte) error {
	var err error
	for _, addr := range f.clientAddrs {
		if addr.Port > 0 {
			fmt.Println("mux server write to", *addr)
			_, e := f.conn.WriteToUDP(buf, addr)
			if e != nil {
				err = e
			}
		}
	}
	return err
}
