package alpacamux

import (
	"net"
)

type UdpClient struct {
	Server     string // upstream IP:Port pair
	conn       *net.UDPConn
	serverAddr *net.UDPAddr
}

func (f *UdpClient) Listen() error {
	var err error
	if f.serverAddr, err = convertAddr(f.Server); err != nil {
		return err
	}

	addr := &net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("0.0.0.0"),
	}

	f.conn, err = net.ListenUDP("udp", addr)
	return err
}

func (f *UdpClient) Read(buf []byte) (int, error) {
	len, client, err := f.conn.ReadFromUDP(buf)
	log.Debug("udp client read from: %v", *client)
	return len, err
}

func (f *UdpClient) Write(buf []byte) error {
	_, err := f.conn.WriteToUDP(buf, f.serverAddr)
	log.Debug("udp client write to: %v", *f.serverAddr)
	return err
}
