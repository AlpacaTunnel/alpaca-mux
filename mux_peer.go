package alpacamux

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

const PATH_SIZE = 4

type MuxPeer struct {
	conn        *net.UDPConn
	dynamicPeer bool
	peerAddrs   [PATH_SIZE]*net.UDPAddr
	buffer      []byte
	timestamp   uint32
	sequence    uint32
	pktFilter   PktFilter
}

func (p *MuxPeer) initPeers(addrs []string) error {
	if addrs == nil {
		p.dynamicPeer = true
	}

	p.peerAddrs = [PATH_SIZE]*net.UDPAddr{
		{},
		{},
		{},
		{},
	}

	for i, s := range addrs {
		var err error
		if p.peerAddrs[i], err = convertAddr(s); err != nil {
			return err
		}
	}
	return nil
}

func (p *MuxPeer) Init(localPort int, peerAddrs []string) error {
	p.buffer = make([]byte, 2000)
	p.pktFilter.Init()

	if err := p.initPeers(peerAddrs); err != nil {
		return err
	}

	addr := &net.UDPAddr{
		Port: localPort,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", addr)
	p.conn = conn
	return err
}

func (p *MuxPeer) updateTimestampSeq() {
	now := uint32(time.Now().Unix())
	if now == p.timestamp {
		p.sequence += 1
	} else {
		p.timestamp = now
		p.sequence = 0
	}
}

func (p *MuxPeer) Write(buf []byte) error {
	p.updateTimestampSeq()
	bodyLen := uint16(len(buf))

	header := Header{
		Magic:     MAGIC,
		Length:    bodyLen,
		Timestamp: p.timestamp,
		Sequence:  p.sequence,
	}

	var err error
	for id, addr := range p.peerAddrs {
		if addr.Port == 0 {
			continue
		}

		log.Debug("mux peer write to path %v: %v", id, *addr)
		header.PathID = uint16(id)

		copy(p.buffer, header.ToNetwork())
		copy(p.buffer[HEADER_LEN:], buf)

		obfsLen := obfsLength(bodyLen)
		// feed random data for the obfs part
		rand.Read(p.buffer[HEADER_LEN+bodyLen : HEADER_LEN+obfsLen])

		_, e := p.conn.WriteToUDP(p.buffer[:HEADER_LEN+obfsLen], addr)
		if e != nil {
			err = e
		}
	}
	return err
}

func (p *MuxPeer) Read(buf []byte) (int, error) {
	len, client, err := p.conn.ReadFromUDP(buf)
	if len < HEADER_LEN {
		return 0, fmt.Errorf("invalid length")
	}

	header := Header{}
	header.FromNetwork(buf)
	log.Debug("mux peer read from: %v, path id: %v", *client, header.PathID)

	if header.Magic != MAGIC {
		return 0, fmt.Errorf("invalid magic")
	}
	if header.PathID > PATH_SIZE-1 {
		return 0, fmt.Errorf("invalid path ID")
	}

	if p.dynamicPeer {
		*(p.peerAddrs[header.PathID]) = *client
	}

	if !p.pktFilter.IsValid(header.Timestamp, header.Sequence) {
		log.Debug("Packet is filtered as invalid, drop it: %v:%v", header.Timestamp, header.Sequence)
		return 0, nil
	}

	copy(buf, buf[HEADER_LEN:HEADER_LEN+header.Length])
	return int(header.Length), err
}
