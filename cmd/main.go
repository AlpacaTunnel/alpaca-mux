package main

import (
	"flag"
	"fmt"
	"time"

	mux "github.com/AlpacaTunnel/alpaca-mux"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "array of strings"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func forward(fa mux.Forwarder, fb mux.Forwarder) {
	buf := make([]byte, mux.MAX_MTU)
	for {
		len, err := fa.Read(buf)
		if err != nil {
			fmt.Printf("error read: %v\n", err)
		}

		err = fb.Write(buf[:len])
		if err != nil {
			fmt.Printf("error write: %v\n", err)
		}
	}
}

func main() {
	client := flag.Bool("c", false, "start a client")
	server := flag.Bool("s", false, "start a server")
	listenPort := flag.Int("l", 1080, "listen port")
	var servers arrayFlags
	flag.Var(&servers, "u", "upstream server IP:Port")
	flag.Parse()

	if len(servers) == 0 {
		panic("empty upstream servers")
	}
	if len(servers) > 4 {
		panic("too many upstream servers")
	}

	var fa, fb mux.Forwarder

	if *client {
		fa = &mux.UdpServer{Port: *listenPort}
		fb = &mux.MuxClient{Servers: servers}
	}

	if *server {
		fa = &mux.MuxServer{Port: *listenPort}
		fb = &mux.UdpClient{Server: servers[0]}
	}

	if err := fa.Listen(); err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	if err := fb.Listen(); err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	go forward(fa, fb)
	go forward(fb, fa)

	for {
		time.Sleep(time.Second)
	}
}
