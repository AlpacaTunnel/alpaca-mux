package main

import (
	"flag"
	"fmt"
	"math/rand"
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
	buf := make([]byte, 1500)
	for {
		len, err := fa.Read(buf)
		if err != nil {
			mux.Log.Error("error read: %v", err)
		}

		// duplicated packet
		if len == 0 {
			continue
		}

		err = fb.Write(buf[:len])
		if err != nil {
			mux.Log.Error("error write: %v", err)
		}
	}
}

func main() {
	client := flag.Bool("c", false, "start a client")
	server := flag.Bool("s", false, "start a server")
	listenPort := flag.Int("l", 1080, "listen port")
	debug := flag.Bool("d", false, "debug log")
	var servers arrayFlags
	flag.Var(&servers, "u", "upstream server IP:Port")
	flag.Parse()

	if len(servers) == 0 {
		panic("empty upstream servers")
	}
	if len(servers) > 4 {
		panic("too many upstream servers")
	}

	mux.Log.SetLevel("info")
	if *debug {
		mux.Log.SetLevel("debug")
	}
	rand.Seed(time.Now().UnixNano())

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
