package alpacamux

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

func convertAddr(addr string) (*net.UDPAddr, error) {
	ipPort := strings.Split(addr, ":")
	if len(ipPort) != 2 {
		return nil, fmt.Errorf("wrong address: %s", addr)
	}

	ip := net.ParseIP(ipPort[0])
	port, err := strconv.Atoi(ipPort[1])
	if ip == nil || err != nil {
		return nil, fmt.Errorf("wrong address: %s", addr)
	}
	return &net.UDPAddr{
		IP:   ip,
		Port: port,
	}, nil
}

func obfsLength(length uint16) uint16 {
	if length < 500 {
		length += uint16(rand.Intn(550)) + 250
	}
	return length
}
