package alpacamux

type MuxServer struct {
	Port int // server listen port
	peer MuxPeer
}

func (f *MuxServer) Listen() error {
	return f.peer.Init(f.Port, nil)
}

func (f *MuxServer) Read(buf []byte) (int, error) {
	return f.peer.Read(buf)
}

func (f *MuxServer) Write(buf []byte) error {
	return f.peer.Write(buf)
}
