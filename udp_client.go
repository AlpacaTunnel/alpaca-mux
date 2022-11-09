package alpacamux

import (
	"encoding/binary"
	"fmt"
	"net"
)

type udpSessionS struct {
	conn       *net.UDPConn
	serverAddr *net.UDPAddr
}

type UdpClient struct {
	Servers  []string // upstream IP:Port pairs
	sessions []*udpSessionS
	pktIn    chan []byte
}

func (f *UdpClient) Listen() error {
	f.pktIn = make(chan []byte)

	for i, server := range f.Servers {
		var err error
		serverAddr, err := convertAddr(server)
		if err != nil {
			return err
		}

		conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0, IP: net.ParseIP("0.0.0.0")})
		if err != nil {
			return err
		}

		f.sessions = append(f.sessions, &udpSessionS{conn: conn, serverAddr: serverAddr})
		id := i
		go func() {
			for {
				collectPktS(f.pktIn, conn, uint16(id), f.sessions[id])
			}
		}()
	}

	return nil
}

func collectPktS(pktIn chan []byte, conn *net.UDPConn, id uint16, session *udpSessionS) {
	buf := make([]byte, 1500)
	binary.BigEndian.PutUint16(buf[0:2], id)

	len, client, err := conn.ReadFromUDP(buf[2:])
	if err != nil {
		log.Error("udp client read from: %v, session id: %d", *client, id)
	}
	log.Debug("udp client read from: %v, session id: %d", *client, id)

	pktIn <- buf[:2+len]
}

func (f *UdpClient) Read(buf []byte) (int, error) {
	data := <-f.pktIn
	log.Debug("udp client read from channel: %d", len(data))
	copy(buf, data)
	return len(data), nil
}

func (f *UdpClient) Write(buf []byte) error {
	if len(buf) < 3 {
		return fmt.Errorf("invalid data length: %d", len(buf))
	}

	id := binary.BigEndian.Uint16(buf[0:2])
	if int(id) >= len(f.sessions) {
		return fmt.Errorf("invalid session id: %d", id)
	}

	session := f.sessions[id]

	_, err := session.conn.WriteToUDP(buf[2:], session.serverAddr)
	log.Debug("udp client write to: %v, session id: %d", *session.serverAddr, id)
	return err
}
