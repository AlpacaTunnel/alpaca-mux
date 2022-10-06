package alpacamux

import (
	"net"
)

type UdpServer struct {
	Port       int
	conn       *net.UDPConn
	clientAddr net.UDPAddr
}

func (f *UdpServer) Listen() error {
	addr := &net.UDPAddr{
		Port: f.Port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", addr)
	f.conn = conn
	return err
}

func (f *UdpServer) Read(buf []byte) (int, error) {
	len, client, err := f.conn.ReadFromUDP(buf)
	if err != nil {
		return 0, err
	}
	f.clientAddr = *client
	log.Debug("udp server read from: %v", *client)

	return len, nil
}

func (f *UdpServer) Write(buf []byte) error {
	_, err := f.conn.WriteToUDP(buf, &f.clientAddr)
	log.Debug("udp server write to: %v", f.clientAddr)
	return err
}
