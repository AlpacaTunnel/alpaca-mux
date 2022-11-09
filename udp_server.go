package alpacamux

import (
	"encoding/binary"
	"fmt"
	"net"
)

type udpSessionC struct {
	conn       *net.UDPConn
	clientAddr net.UDPAddr
}

type UdpServer struct {
	Ports    []int
	sessions []*udpSessionC
	pktIn    chan []byte
}

func (f *UdpServer) Listen() error {
	f.pktIn = make(chan []byte)
	for i, port := range f.Ports {
		addr := &net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("0.0.0.0"),
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			return err
		}

		f.sessions = append(f.sessions, &udpSessionC{conn: conn})
		id := i
		go func() {
			for {
				collectPktC(f.pktIn, conn, uint16(id), f.sessions[id])
			}
		}()
	}
	return nil
}

func collectPktC(pktIn chan []byte, conn *net.UDPConn, id uint16, session *udpSessionC) {
	buf := make([]byte, 1500)
	binary.BigEndian.PutUint16(buf[0:2], id)

	len, client, err := conn.ReadFromUDP(buf[2:])
	if err != nil {
		log.Error("udp server read from: %v, session id: %d", *client, id)
	}
	log.Debug("udp server read from: %v, session id: %d", *client, id)

	session.clientAddr = *client
	pktIn <- buf[:2+len]
}

func (f *UdpServer) Read(buf []byte) (int, error) {
	data := <-f.pktIn
	log.Debug("udp server read from channel: %d", len(data))
	copy(buf, data)
	return len(data), nil
}

func (f *UdpServer) Write(buf []byte) error {
	if len(buf) < 3 {
		return fmt.Errorf("invalid data length: %d", len(buf))
	}

	id := binary.BigEndian.Uint16(buf[0:2])
	if int(id) >= len(f.sessions) {
		return fmt.Errorf("invalid session id: %d", id)
	}

	session := f.sessions[id]

	_, err := session.conn.WriteToUDP(buf[2:], &session.clientAddr)
	log.Debug("udp server write to: %v, session id: %d", session.clientAddr, id)
	return err
}
